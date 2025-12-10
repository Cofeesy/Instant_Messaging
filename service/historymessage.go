package service

import(
	"context"
	"encoding/json"
	"fmt"
	"gin_chat/model"
	"gin_chat/model/request"
	"gin_chat/utils"
	"time"
	"github.com/redis/go-redis/v9"
)

func GetSingleHistoryMsg(req request.SingleHistoryMsgReq) ([]*model.Message, error) {
	ctx := context.Background()

	key := " "
	if req.UserId < req.TargetId {
		key = fmt.Sprintf("chat:private:%d:%d", req.UserId, req.TargetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", req.TargetId, req.UserId)
	}

	// 准备Redis查询参数 (ZRevRangeByScore)
	// 我们要查 Score < Cursor 的数据，也就是比这个时间更早的消息
	max := "+inf" // 默认查最新的
	if req.Cursor > 0 {
		// "(" 表示开区间，即不包含 Cursor 这条本身，防止重复
		max = fmt.Sprintf("(%d", req.Cursor) 
	}

	// 从 Redis 倒序查询
	// 结果是 []string，里面存的是 json
	redisVals, err := utils.RDB.ZRevRangeByScore(ctx, key, &redis.ZRangeBy{
		Min:    "-inf",
		Max:    max,
		Count:  int64(req.Limit),
		Offset: 0,
	}).Result()

	if err != nil {
		// 如果 Redis 挂了，记录日志，不直接返回错，尝试去 DB 查兜底
		fmt.Println("Redis查询失败,降级查DB:", err)
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
	if len(msgs) < req.Limit {
		need := req.Limit - len(msgs)
		
		// 计算MySQL查询的起始时间 (nextCursor)
		var nextCursor int64
		
		if len(msgs) > 0 {
			// A:Redis 里有几条，但不够。接着 Redis 最后一条往下查
			nextCursor = int64(msgs[len(msgs)-1].CreateTime)
		} else {
			// B:Redis 完全没数据 (过期了或冷门会话)
			nextCursor = req.Cursor
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
			req.UserId, req.TargetId, req.TargetId, req.UserId, nextCursor, 
		).
		Order("create_time DESC").
		Limit(need).
		Find(&dbMsgs).Error
			
		if err != nil {
			fmt.Println("DB Query Error:", err)
			// DB 也错了，返回当前已有的
			return msgs, nil 
		} 

		// 合并结果: Redis在前，MySQL在后
		msgs = append(msgs, dbMsgs...)
	}


	return msgs, nil
}

func GetGroupHistoryMessages(groupRedis *request.GroupRedisPayload) ([]*model.Message, error) {
	ctx := context.Background()

	key := fmt.Sprintf("chat:group:%d", groupRedis.GroupId)

	var stringmsgs []string
	total, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if groupRedis.End == -1 {
		stringmsgs, err = utils.RDB.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
	} else {
		if total == 0 {
			return []*model.Message{}, nil
		}
		s := int64(total) - 1 - groupRedis.End
		e := int64(total) - 1 - groupRedis.Start
		if s < 0 {
			s = 0
		}
		if e < 0 {
			return []*model.Message{}, nil
		}
		stringmsgs, err = utils.RDB.ZRange(ctx, key, s, e).Result()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	msgs := make([]*model.Message, 0)
	for _, v := range stringmsgs {
		stringmsg := []byte(v)
		var msg model.Message
		json.Unmarshal(stringmsg, &msg)
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

func GetAiHistoryMessages(AiRedisMsgPayload *request.AiRedisMsgPayload) ([]*model.Message, error) {
	ctx := context.Background()

	key := fmt.Sprintf("aichat:%d:%d", 0, AiRedisMsgPayload.UserId)

	var stringmsgs []string
	total, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if AiRedisMsgPayload.End == -1 {
		stringmsgs, err = utils.RDB.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
	} else {
		if total == 0 {
			return []*model.Message{}, nil
		}
		s := int64(total) - 1 - AiRedisMsgPayload.End
		e := int64(total) - 1 - AiRedisMsgPayload.Start
		if s < 0 {
			s = 0
		}
		if e < 0 {
			return []*model.Message{}, nil
		}
		stringmsgs, err = utils.RDB.ZRange(ctx, key, s, e).Result()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	msgs := make([]*model.Message, 0)
	for _, v := range stringmsgs {
		var msg model.Message
		err := json.Unmarshal([]byte(v), &msg)
		if err != nil {
			fmt.Println("解析消息失败:", err)
			continue
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}
