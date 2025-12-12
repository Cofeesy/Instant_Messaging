package service

import (
	"context"
	"encoding/json"
	"fmt"
	"ZustChat/global"
	"ZustChat/model"
	"ZustChat/model/request"
	"ZustChat/utils"
	"time"

	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

func GetSingleHistoryMsg(singleReq request.SingleHistoryMsgReq) ([]*model.Message, error) {
	ctx := context.Background()

	key := " "
	if singleReq.UserId < singleReq.TargetId {
		key = fmt.Sprintf("chat:private:%d:%d", singleReq.UserId, singleReq.TargetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", singleReq.TargetId, singleReq.UserId)
	}

	// 准备Redis查询参数 (ZRevRangeByScore)
	// 我们要查 Score < Cursor 的数据，也就是比这个时间更早的消息
	max := "+inf" // 默认查最新的
	if singleReq.Cursor > 0 {
		// "(" 表示开区间，即不包含 Cursor 这条本身，防止重复
		max = fmt.Sprintf("(%d", singleReq.Cursor)
	}

	// 从 Redis 倒序查询
	// 结果是 []string，里面存的是 json
	redisVals, err := utils.RDB.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    max,
		Count:  int64(singleReq.Limit),
		Offset: 0,
	}).Result()

	if err != nil {
		// 如果 Redis 挂了，记录日志，不直接返回错，尝试去 DB 查兜底
		global.Logger.Warn("Redis查询失败,降级查DB",
			zap.Uint("user_id", singleReq.UserId),
			zap.Uint("target_id", singleReq.TargetId),
			zap.String("key", key),
			zap.Error(err),
		)
		redisVals = []string{}
	}

	msgs := make([]*model.Message, 0)
	for _, v := range redisVals {
		var msg model.Message
		if err := json.Unmarshal([]byte(v), &msg); err == nil {
			msgs = append(msgs, &msg)
		}
	}

	// 如果 Redis 里拿到的数量少于 Limit (说明 Redis 数据没了，或者发生了断层)
	if len(msgs) < singleReq.Limit {
		need := singleReq.Limit - len(msgs)

		// 计算MySQL查询的起始时间 (nextCursor)
		var nextCursor int64

		if len(msgs) > 0 {
			// A:Redis里有几条，但不够。接着 Redis 最后一条往下查
			nextCursor = int64(msgs[len(msgs)-1].CreateTime)
		} else {
			// B:Redis 完全没数据 (过期了或冷门会话)
			nextCursor = singleReq.Cursor
			// 如果是第一次加载 (Cursor=0) 且 Redis 为空，就从当前时间开始查
			if nextCursor <= 0 {
				nextCursor = time.Now().UnixMilli()
			}
		}

		// 去 MySQL 查剩下的
		var dbMsgs []*model.Message

		// GORM 查询：
		// 1. 双方 ID 匹配 (注意 A->B 和 B->A 都是同一个会话)
		// 2. create_time < nextCursor
		// 3. 倒序 + 限制条数
		err := utils.DB.Where(
			"((user_id = ? AND target_id = ?) OR (user_id = ? AND target_id = ?)) AND create_time < ?",
			singleReq.UserId, singleReq.TargetId, singleReq.TargetId, singleReq.UserId, nextCursor,
		).
			Order("create_time DESC").
			Limit(need).
			Find(&dbMsgs).Error

		if err != nil {
			global.Logger.Error("私聊历史消息DB查询失败",
				zap.Uint("user_id", singleReq.UserId),
				zap.Uint("target_id", singleReq.TargetId),
				zap.Error(err),
			)
			// DB 也错了，返回当前已有的
			return msgs, nil
		}

		// 合并结果: Redis在前，MySQL在后
		msgs = append(msgs, dbMsgs...)
	}

	return msgs, nil
}

func GetGroupHistoryMessages(groupReq *request.GroupHistoryMsgReq) ([]*model.Message, error) {
	ctx := context.Background()

	key := fmt.Sprintf("chat:group:%d", groupReq.GroupId)

	// 准备Redis查询参数 (ZRevRangeByScore)
	// 要查 Score < Cursor 的数据，也就是比这个时间更早的消息
	max := "+inf" // 默认查最新的
	if groupReq.Cursor > 0 {
		// "(" 表示开区间，即不包含 Cursor 这条本身，防止重复
		max = fmt.Sprintf("(%d", groupReq.Cursor)
	}

	// 从 Redis 倒序查询
	// 结果是 []string，里面存的是 json
	redisVals, err := utils.RDB.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    max,
		Count:  int64(groupReq.Limit),
		Offset: 0,
	}).Result()

	if err != nil {
		// 如果 Redis 挂了，记录日志，不直接返回错，尝试去DB兜底
		global.Logger.Warn("Redis查询失败,降级查DB",
			zap.Uint("group_id", groupReq.GroupId),
			zap.String("key", key),
			zap.Error(err),
		)
		redisVals = []string{}
	}

	msgs := make([]*model.Message, 0)
	for _, v := range redisVals {
		var msg model.Message
		if err := json.Unmarshal([]byte(v), &msg); err == nil {
			msgs = append(msgs, &msg)
		}
	}

	// 如果 Redis 里拿到的数量少于 Limit (说明 Redis 数据没了，或者发生了断层)
	if len(msgs) < groupReq.Limit {
		need := groupReq.Limit - len(msgs)

		// 计算MySQL查询的起始时间 (nextCursor)
		var nextCursor int64

		if len(msgs) > 0 {
			// A:Redis里有几条，但不够。接着 Redis 最后一条往下查
			nextCursor = int64(msgs[len(msgs)-1].CreateTime)
		} else {
			// B:Redis 完全没数据 (过期了或冷门会话)
			nextCursor = groupReq.Cursor
			// 如果是第一次加载 (Cursor=0) 且 Redis 为空，就从当前时间开始查
			if nextCursor <= 0 {
				nextCursor = time.Now().UnixMilli()
			}
		}

		// 去 MySQL 查剩下的
		var dbMsgs []*model.Message

		// GORM 查询：
		// 1. 双方 ID 匹配 (注意 A->B 和 B->A 都是同一个会话)
		// 2. create_time < nextCursor
		// 3. 倒序 + 限制条数
		err := utils.DB.Where("target_id = ? AND type = 2 AND create_time < ?", groupReq.GroupId, nextCursor).
			Order("create_time DESC").
			Limit(need).
			Find(&dbMsgs).Error

		if err != nil {
			global.Logger.Error("群聊历史消息DB查询失败",
				zap.Uint("group_id", groupReq.GroupId),
				zap.Error(err),
			)
			// DB 也错了，返回当前已有的
			return msgs, nil
		}

		// 合并结果: Redis在前，MySQL在后
		msgs = append(msgs, dbMsgs...)
	}

	return msgs, nil
}

func GetAiHistoryMessages(aiReq *request.AiHistoryMsgReq) ([]*model.Message, error) {
	ctx := context.Background()

	key := " "

	if aiReq.UserId < aiReq.TargetId {
		key = fmt.Sprintf("chat:private:%d:%d", aiReq.UserId, aiReq.TargetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", aiReq.TargetId, aiReq.UserId)
	}

	max := "+inf"
	if aiReq.Cursor > 0 {
		max = fmt.Sprintf("(%d", aiReq.Cursor)
	}

	redisVals, err := utils.RDB.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    max,
		Count:  int64(aiReq.Limit),
		Offset: 0,
	}).Result()

	if err != nil {
		global.Logger.Warn("Redis查询失败,降级查DB",
			zap.Uint("user_id", aiReq.UserId),
			zap.Uint("target_id", aiReq.TargetId),
			zap.String("key", key),
			zap.Error(err),
		)
		redisVals = []string{}
	}

	msgs := make([]*model.Message, 0)
	for _, v := range redisVals {
		var msg model.Message
		if err := json.Unmarshal([]byte(v), &msg); err == nil {
			msgs = append(msgs, &msg)
		}
	}

	if len(msgs) < aiReq.Limit {
		need := aiReq.Limit - len(msgs)

		// 计算MySQL查询的起始时间 (nextCursor)
		var nextCursor int64

		if len(msgs) > 0 {
			// A:Redis里有几条，但不够。接着 Redis 最后一条往下查
			nextCursor = int64(msgs[len(msgs)-1].CreateTime)
		} else {
			// B:Redis 完全没数据 (过期了或冷门会话)
			nextCursor = aiReq.Cursor
			// 如果是第一次加载 (Cursor=0) 且 Redis 为空，就从当前时间开始查
			if nextCursor <= 0 {
				nextCursor = time.Now().UnixMilli()
			}
		}

		// 去 MySQL 查剩下的
		var dbMsgs []*model.Message

		// GORM 查询：
		// 1. 双方 ID 匹配 (注意 A->B 和 B->A 都是同一个会话)
		// 2. create_time < nextCursor
		// 3. 倒序 + 限制条数
		err := utils.DB.Where(
			"((user_id = ? AND target_id = ?) OR (user_id = ? AND target_id = ?)) AND create_time < ?",
			aiReq.UserId, aiReq.TargetId, aiReq.TargetId, aiReq.UserId, nextCursor,
		).
			Order("create_time DESC").
			Limit(need).
			Find(&dbMsgs).Error

		if err != nil {
			global.Logger.Error("AI历史消息DB查询失败",
				zap.Uint("user_id", aiReq.UserId),
				zap.Uint("target_id", aiReq.TargetId),
				zap.Error(err),
			)
			// DB 也错了，返回当前已有的
			return msgs, nil
		}

		// 合并结果: Redis在前，MySQL在后
		msgs = append(msgs, dbMsgs...)
	}

	return msgs, nil
}
