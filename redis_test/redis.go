package main


// import (
//     "context"
//     "fmt"

//     "github.com/redis/go-redis/v9"
// )

// var ctx = context.Background()

// func main() {
//     rdb := redis.NewClient(&redis.Options{
//         Addr:     "localhost:6379",
//         Password: "", // no password set
//         DB:       0,  // use default DB
//     })

//     err := rdb.Set(ctx, "key", "value", 0).Err()
//     if err != nil {
//         panic(err)
//     }

//     val, err := rdb.Get(ctx, "key").Result()
//     if err != nil {
//         panic(err)
//     }
//     fmt.Println("key", val)

//     val2, err := rdb.Get(ctx, "key2").Result()
//     if err == redis.Nil {
//         fmt.Println("key2 does not exist")
//     } else if err != nil {
//         panic(err)
//     } else {
//         fmt.Println("key2", val2)
//     }
//     // Output: key value
//     // key2 does not exist
// }



import (
	"fmt"
	"gin_chat/utils/setting"
	"github.com/redis/go-redis/v9"
	"log"
	"context"
)

var ctx = context.Background()
var rdb *redis.Client

func main() {
	var (
		err            error
		addr, password string
		db             int
	)

	sec, err := setting.Cfg.GetSection("redis")
	if err != nil {
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	addr = sec.Key("ADDR").String()
	password = sec.Key("PASSWORD").String()
	db, err = sec.Key("DB").Int()
	if err!=nil{
		log.Fatalf("Fail to get section 'redis': %v", err)
	}

	rdb = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,        // use default DB
	})

	fmt.Println("redis已经加载")

	err = rdb.Set(ctx, "key", "value", 0).Err()
    if err != nil {
        panic(err)
    }

    val, err := rdb.Get(ctx, "key").Result()
    if err != nil {
        panic(err)
    }
    fmt.Println("key", val)

    val2, err := rdb.Get(ctx, "key2").Result()
    if err == redis.Nil {
        fmt.Println("key2 does not exist")
    } else if err != nil {
        panic(err)
    } else {
        fmt.Println("key2", val2)
    }
}