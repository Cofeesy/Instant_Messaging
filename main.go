package main

import (
	"fmt"
	"gin_chat/router"
	"gin_chat/utils/setting"
	"net/http"
)

func main() {
	r := router.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", setting.HTTPPort),
		Handler:        r,
		ReadTimeout:    setting.ReadTimeout,
		WriteTimeout:   setting.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	
	s.ListenAndServe()
}
