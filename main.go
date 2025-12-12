package main

import (
	"fmt"
	"ZustChat/global"
	"ZustChat/router"
	"ZustChat/service"
	"ZustChat/utils"
	"net/http"

	"go.uber.org/zap"
)

func main() {
	utils.InitializeSystem()

	// 设置全局 logger
	global.Logger = utils.Logger

	// 确保程序退出时同步日志缓冲区
	defer utils.SyncLogger()

	utils.InitTimer(service.CleanConnection)

	r := router.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", utils.HTTPPort),
		Handler:        r,
		ReadTimeout:    utils.ReadTimeout,
		WriteTimeout:   utils.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}

	global.Logger.Info("服务器启动",
		zap.String("addr", s.Addr),
		zap.String("mode", utils.RunMode),
	)

	if err := s.ListenAndServe(); err != nil {
		global.Logger.Fatal("服务器启动失败", zap.Error(err))
	}
}
