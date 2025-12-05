package utils

import (
	"time"
)

//初始化定时器
// CleanConnection
// DelayHeartbeat = 3 
// HeartbeatHz = 30
func InitTimer(f TimerFunc) {
	Timer(time.Duration(DelayHeartbeat)*time.Second, time.Duration(HeartbeatHz)*time.Second, f, "")
}

type TimerFunc func(interface{}) bool


/**
delay  首次延迟
tick  间隔
fun  定时执行的方法
param  方法的参数 --> CleanConnection
**/
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		// 延迟几秒启动，保证“系统准备好”再开始定时检查
		time.Sleep(delay)
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for range ticker.C {
			if !fun(param) {
			// if fun(param) == false {
				return
			}
		}
	}()
}
