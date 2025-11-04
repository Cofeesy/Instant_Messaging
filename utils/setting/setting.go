package setting

import (
	"github.com/go-ini/ini"
	"log"
	"time"
)

var (
	Cfg      *ini.File
	RunMode  string
	HTTPPort int

	ReadTimeout  time.Duration
	WriteTimeout time.Duration

	JwtSecret string
)

func init() {
	var err error
	Cfg, err = ini.Load("conf/app.ini")
	if err != nil {
		log.Fatalf("Fail to parse 'conf/app.ini:%v'", err)
	}

	LoadAPP()
	LoadBase()
	LoadServer()

}

func LoadAPP(){
	JwtSecret = Cfg.Section("app").Key("JWT_SECRET").MustString("!@)*#)!@U#@*!@!)")
}

func LoadBase() {
	RunMode = Cfg.Section("").Key("RUN_MODES").MustString("debug")
	// 为什么要加这个？

	// 含义：从客户端发出请求开始，到服务器读取完整的请求体为止的最长等待时间。
	// 目的：防止恶意客户端缓慢发送数据（比如 “Slowloris 攻击”），拖住连接占用资源
	ReadTimeout = time.Duration(Cfg.Section("server").Key("READ_TIMEOUT").MustInt(60)) * time.Second

	// 含义：从服务器开始写响应，到写完响应的最大允许时间。
	// 目的：防止客户端读得太慢、网络异常等情况造成服务器阻塞
	WriteTimeout = time.Duration(Cfg.Section("server").Key("WRITE_TIMEOUT").MustInt(60)) * time.Second
}

func LoadServer() {
	sec, err := Cfg.GetSection("server")
	if err != nil {
		log.Fatalf("Fail to get section 'server': %v", err)
	}
	HTTPPort = sec.Key("HTTP_PORT").MustInt(8000)

}
