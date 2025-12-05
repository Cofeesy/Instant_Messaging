package utils

import (
	"fmt"
	"github.com/redis/go-redis/v9"
	"log"
)

var RDB *redis.Client

func InitRedis() {
	var (
		err            error
		addr, password string
		dbNumber             int
	)

	sec, err := Cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	addr = sec.Key("ADDR").String()
	password = sec.Key("PASSWORD").String()
	dbNumber,err = sec.Key("DB").Int()
	if err!=nil{
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	RDB = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       dbNumber,        // use default DB
	})
	fmt.Println("redis已经加载")
}
