package main

import (
	"goPush/controllers"
	"goPush/lib/config"
	"goPush/lib/log"
	"goPush/lib/push"
	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/pprof"
)

func main() {
	//初始化配置文件
	config.Instance = new(config.Config)
	config.Instance.InitConfig("conf/app.ini")
	runMode := config.Instance.Get("default", "runMode")
	port := config.Instance.Get("default", "port")

	//初始化logger
	log.Instance = new(log.Logger)
	log.Instance.InitLogger("logs")
	if runMode == "dev" {
		gin.SetMode(gin.DebugMode)
		log.Instance.SetLevel(log.DEBUG)
		log.Instance.SetConsole(true)
	} else {
		gin.SetMode(gin.ReleaseMode)
		log.Instance.SetLevel(log.INFO)
		log.Instance.SetConsole(false)
	}

	//初始化ios client
	push.IosInstance = new(push.Ios)
	push.IosInstance.InitClient()

	//初始化路由
	r := gin.Default()
	//debug cpu 内存等
	pprof.Register(r, &pprof.Options{
		RoutePrefix: "debug/pprof",
	})
	pushController := new(controllers.PushController)
	pushGroup := r.Group("/api")
	pushGroup.POST("/push", pushController.Push)
	r.Run(":" + port)
}
