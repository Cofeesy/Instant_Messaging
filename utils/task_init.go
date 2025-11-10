package utils

import (
	"gin_chat/utils/setting"
	"time"
)

//初始化定时器
func InitTimer(f TimerFunc) {
	Timer(time.Duration(setting.DelayHeartbeat)*time.Second, time.Duration(setting.HeartbeatHz)*time.Second, f, "")
}

type TimerFunc func(interface{}) bool


/**
delay  首次延迟
tick  间隔
fun  定时执行的方法
param  方法的参数
**/
func Timer(delay, tick time.Duration, fun TimerFunc, param interface{}) {
	go func() {
		if fun == nil {
			return
		}
		time.Sleep(delay)
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		for range ticker.C {
			if fun(param) == false {
				return
			}
		}

		// for {
		// 	select {
		// 	case <-t.C:
		// 		if fun(param) == false {
		// 			return
		// 		}
		// 		t.Reset(tick)
		// 	}
		// }
	}()
}
