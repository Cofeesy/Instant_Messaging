package service

import(
	"context"
	"encoding/json"
	"fmt"
	"gin_chat/models"
	"gin_chat/models/system"
	"gin_chat/utils"
)

// 读取start-end的redis数据，并返回给前端
func GetSingleHistoryMsg(singleRedisPayload system.SingleRedisPayload) ([]*models.Message, error) {
	ctx := context.Background()

	key := " "
	if singleRedisPayload.UserId < singleRedisPayload.TargetId {
		key = fmt.Sprintf("chat:private:%d:%d", singleRedisPayload.UserId, singleRedisPayload.TargetId)
	} else {
		key = fmt.Sprintf("chat:private:%d:%d", singleRedisPayload.TargetId, singleRedisPayload.UserId)
	}

	// 前端传来的 Start/End 是相对于最新消息的偏移（Start=0, End=9 表示最近 10 条）
	// 需要把它们转换为 Redis 有序集合的正序索引：
	// 假设总数为 total，则想要的区间为 [total-1-End, total-1-Start]
	// 如果 End == -1，则返回全部（ZRange 0 -1）
	var stringmsgs []string
	total, err := utils.RDB.ZCard(ctx, key).Result()
	if err != nil {
		return nil, err
	}

	if singleRedisPayload.End == -1 {
		// 返回所有消息（正序）
		stringmsgs, err = utils.RDB.ZRange(ctx, key, 0, -1).Result()
		if err != nil {
			return nil, err
		}
	} else {
		// 计算正序的 start/end
		if total == 0 {
			return []*models.Message{}, nil
		}
		// 计算索引
		s := int64(total) - 1 - singleRedisPayload.End
		e := int64(total) - 1 - singleRedisPayload.Start
		if s < 0 {
			s = 0
		}
		if e < 0 {
			// 没有可返回的消息
			return []*models.Message{}, nil
		}
		stringmsgs, err = utils.RDB.ZRange(ctx, key, s, e).Result()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	msgs := make([]*models.Message, 0)
	for _, v := range stringmsgs {
		stringmsg := []byte(v)
		var msg models.Message
		json.Unmarshal(stringmsg, &msg)
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

func GetGroupHistoryMessages(groupRedis *system.GroupRedisPayload) ([]*models.Message, error) {
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
			return []*models.Message{}, nil
		}
		s := int64(total) - 1 - groupRedis.End
		e := int64(total) - 1 - groupRedis.Start
		if s < 0 {
			s = 0
		}
		if e < 0 {
			return []*models.Message{}, nil
		}
		stringmsgs, err = utils.RDB.ZRange(ctx, key, s, e).Result()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	msgs := make([]*models.Message, 0)
	for _, v := range stringmsgs {
		stringmsg := []byte(v)
		var msg models.Message
		json.Unmarshal(stringmsg, &msg)
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}

func GetAiHistoryMessages(AiRedisMsgPayload *system.AiRedisMsgPayload) ([]*models.Message, error) {
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
			return []*models.Message{}, nil
		}
		s := int64(total) - 1 - AiRedisMsgPayload.End
		e := int64(total) - 1 - AiRedisMsgPayload.Start
		if s < 0 {
			s = 0
		}
		if e < 0 {
			return []*models.Message{}, nil
		}
		stringmsgs, err = utils.RDB.ZRange(ctx, key, s, e).Result()
		if err != nil {
			return nil, err
		}
	}
	if err != nil {
		return nil, err
	}

	msgs := make([]*models.Message, 0)
	for _, v := range stringmsgs {
		var msg models.Message
		err := json.Unmarshal([]byte(v), &msg)
		if err != nil {
			fmt.Println("解析消息失败:", err)
			continue
		}
		msgs = append(msgs, &msg)
	}

	return msgs, nil
}
