package utils

import (
	"fmt"
	"gin_chat/utils/setting"
	"github.com/redis/go-redis/v9"
	"log"
)

var rdb *redis.Client

func init() {
	var (
		err            error
		addr, password string
		dbNumber             int
	)

	sec, err := setting.Cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	addr = sec.Key("ADDR").String()
	password = sec.Key("PASSWORD").String()
	dbNumber,err = sec.Key("DB").Int()
	if err!=nil{
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       dbNumber,        // use default DB
	})
	fmt.Println("redis已经加载")
}
