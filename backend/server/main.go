package main

import (
	"cmd/server/handle/agent/install"
	"cmd/server/handle/server/monitor" // 引入 monitor 包
	"cmd/server/handle/user/login"
	"cmd/server/middlewire"
	"cmd/server/model"

	"github.com/gin-gonic/gin"
)

func main() {
	go monitor.CheckServerStatus()
	router := gin.Default()
	model.InitDB()
	router.POST("/agent/register", login.Register)
	router.POST("/agent/login", login.Login)
	// 需要 JWT 认证的路由
	auth := router.Group("/agent", middlewire.JWTAuthMiddleware())
	{
		auth.POST("/install", install.InstallAgent)
		auth.POST("/system_info", monitor.GetMessage)
		auth.GET("/list", monitor.ListAgent)
		router.GET("/:hostname", monitor.GetAgentInfo)
	}
	router.Run("0.0.0.0:8080")
}
