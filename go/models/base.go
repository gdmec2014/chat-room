package models

import (
	"chat-room/go/helper"
	"github.com/astaxie/beego"
	"net"
	"strconv"
	"time"
)

var (
	App *app = new(app)
)

type app struct {
	Runmode string
	IsDebug bool
	ApiUrl  string
	RunTime int64
}

func init() {
	//app
	// App Config Init
	App.Runmode = beego.AppConfig.String("runmode")
	if beego.AppConfig.String("isdebug") == "true" {
		App.IsDebug = true
	} else {
		App.IsDebug = false
	}
	App.ApiUrl = "http://" + GetIp() + ":" + strconv.Itoa(beego.BConfig.Listen.HTTPPort)
	RunTime := helper.StringDateFormatInt(beego.AppConfig.String("runtime"))
	if RunTime == 0 {
		RunTime = time.Now().Unix()
	}
	App.RunTime = RunTime
}

//获取本地IP
func GetIp() string {
	if App.IsDebug {
		return "127.0.0.1"
	}
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "127.0.0.1"
	}
	for _, address := range addrs {
		// 检查ip地址判断是否回环地址
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "127.0.0.1"
}
