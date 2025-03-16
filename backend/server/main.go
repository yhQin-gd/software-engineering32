package main

import (
	"cmd/server/config"
	"cmd/server/handle/agent/install"
	"cmd/server/handle/server/monitor" // 引入 monitor 包
	"cmd/server/handle/user/login"
	"cmd/server/middlewire"
	"cmd/server/middlewire/cors"
	db "cmd/server/model/init"
	"fmt"
	"log"
	"os"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

func main() {
	go monitor.CheckServerStatus() //读取DBConfig.yaml文件
	config, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("加载配置失败: %v", err)
	}

	fmt.Println("-----------------------")
	//设置数据库连接的环境变量
	os.Setenv("DB_USER", config.DB.User)
	os.Setenv("DB_PASSWORD", config.DB.Password)
	os.Setenv("DB_HOST", config.DB.Host)
	os.Setenv("DB_PORT", config.DB.Port)
	os.Setenv("DB_NAME", config.DB.Name)
	// OSS服务
	os.Setenv("OSS_REGION", config.OSS.OSS_REGION)
	os.Setenv("OSS_ACCESS_KEY_ID", config.OSS.OSS_ACCESS_KEY_ID)
	os.Setenv("OSS_ACCESS_KEY_SECRET", config.OSS.OSS_ACCESS_KEY_SECRET)
	os.Setenv("OSS_BUCKET", config.OSS.OSS_BUCKET)
	// 邮箱服务
	os.Setenv("EMAIL_NAME", config.Email.Name)
	os.Setenv("EMAIL_PASSWORD", config.Email.Password)
	os.Setenv("BASE_URL", config.Email.Url)
	os.Setenv("SMTP_SERVER_HOST", config.SMTPServer.Host)
	os.Setenv("SMTP_SERVER_PORT", config.SMTPServer.Port)
	// Redis服务
	os.Setenv("REDIS_ADDR", config.Redis.Addr)
	os.Setenv("REDIS_PASSWORD", config.Redis.Password)
	os.Setenv("REDIS_DB", config.Redis.DB)

	// fmt.Println(os.Getenv("DB_USER"))
	// fmt.Println(os.Getenv("DB_PASSWORD"))
	// fmt.Println(os.Getenv("DB_HOST"))
	// fmt.Println(os.Getenv("DB_PORT"))
	// fmt.Println(os.Getenv("DB_NAME"))

	router := gin.Default()
	router.Use(cors.CORSMiddleware())
	// 连接数据库
	if err := db.ConnectDatabase(); err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	// 初始化数据库
	if err := db.InitDB(); err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	// 初始化数据库数据
	if err := db.InitDBData(); err != nil {
		log.Fatalf("Failed to initialize data: %v", err)
	}
	// 初始化redis
	// if err := db.InitRedis(); err!= nil {
	// 	log.Fatalf("Failed to connect to redis: %v", err)
	// }

	router.Static("/static", "./static")

	// 注册 Swagger 路由
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

	router.POST("/agent/register", login.Register)
	router.POST("/agent/login", login.Login)
	// 需要 JWT 认证的路由
	auth := router.Group("/agent", middlewire.JWTAuthMiddleware())
	{
		router.POST("/reset_password", login.ResetPassword)
		auth.POST("/request_reset_password", login.RequestResetPassword)
		auth.POST("/install", install.InstallAgent)
		auth.POST("/addSystemInfo", monitor.ReceiveAndStoreSystemMetrics)
		auth.GET("/list", monitor.ListAgent)
		router.GET("/monitor/:hostname", monitor.GetAgentInfo)
	}

	router.Run("0.0.0.0:8080")
}
