package main

import (
	"fmt"
	"gin_chat/router"
	"gin_chat/service"
	"gin_chat/utils"
	"net/http"
)

func main() {
	utils.InitializeSystem() 
	utils.InitTimer(service.CleanConnection) 

	r := router.InitRouter()
	s := &http.Server{
		Addr:           fmt.Sprintf(":%d", utils.HTTPPort),
		Handler:        r,
		ReadTimeout:    utils.ReadTimeout,
		WriteTimeout:   utils.WriteTimeout,
		MaxHeaderBytes: 1 << 20,
	}
	
	s.ListenAndServe()
}

