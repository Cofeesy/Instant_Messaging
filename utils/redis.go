package utils

import (
	"log"

	"github.com/redis/go-redis/v9"

	"go.uber.org/zap"
)

var RDB *redis.Client

func InitRedis() {
	var (
		err            error
		addr, password string
		dbNumber       int
	)

	sec, err := Cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	addr = sec.Key("ADDR").String()
	password = sec.Key("PASSWORD").String()
	dbNumber, err = sec.Key("DB").Int()
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       dbNumber, // use default DB
	})

	// 如果 logger 已初始化，使用 logger，否则使用 log
	if Logger != nil {
		Logger.Info("Redis初始化成功",
			zap.String("addr", addr),
			zap.Int("db", dbNumber),
		)
	} else {
		log.Printf("Redis已经加载: %s, DB: %d", addr, dbNumber)
	}
}
