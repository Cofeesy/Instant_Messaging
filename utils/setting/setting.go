package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg      *ini.File
	RunMode  string
	HTTPPort uint64

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	JwtSecret string
	GeminiKey string

	DelayHeartbeat   uint64
	HeartbeatHz      uint64
	HeartbeatMaxTime uint64
	RedisOnlineTime  uint64
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini:%v'", err)
	}

	LoadAPP()
	LoadServer()
	LoadTimer()
}

func LoadAPP() {
	app := Cfg.Section("app")
	JwtSecret = app.Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
	GeminiKey = app.Key("GEMINI_KEY").MustString("")
	RunMode = app.Key("RUN_MODES").MustString("debug")
}


func LoadServer() {
	server, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	// 为什么要加这个？

	// 含义：从客户端发出请求开始，到服务器读取完整的请求体为止的最长等待时间。
	// 目的：防止恶意客户端缓慢发送数据（比如 “Slowloris 攻击”），拖住连接占用资源
	ReadTimeout = time.Duration(server.Key("READ_TIMEOUT").MustInt(60)) * time.Second

	// 含义：从服务器开始写响应，到写完响应的最大允许时间。
	// 目的：防止客户端读得太慢、网络异常等情况造成服务器阻塞
	WriteTimeout = time.Duration(server.Key("WRITE_TIMEOUT").MustInt(60)) * time.Second

	// HTTPPort, err := HTTPPort = server.Key("HTTP_PORT").Int()
	// 上面这种方法指定返回int,但要处理错误，即读取到的值可能不是int
	// 下面这种方法不用处理错误，即如果不是int则使用设置的默认值
	HTTPPort = server.Key("HTTP_PORT").MustUint64(8000)
}

func LoadTimer() {
	Timer := Cfg.Section("timer")
	DelayHeartbeat = Timer.Key("DelayHeartbeat").MustUint64(3)
	HeartbeatHz = Timer.Key("HeartbeatHz").MustUint64(30)
	HeartbeatMaxTime = Timer.Key("HeartbeatMaxTime").MustUint64(30000)
	RedisOnlineTime = Timer.Key("RedisOnlineTime").MustUint64(40)
}


