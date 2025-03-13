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
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"log"
	"os"

	"github.com/gin-gonic/gin"
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
<<<<<<< HEAD
	fmt.Println(os.Getenv("DB_USER"))
	fmt.Println(os.Getenv("DB_PASSWORD"))
	fmt.Println(os.Getenv("DB_HOST"))
	fmt.Println(os.Getenv("DB_PORT"))
	fmt.Println(os.Getenv("DB_NAME"))
=======
	// fmt.Println(os.Getenv("DB_USER"))
	// fmt.Println(os.Getenv("DB_PASSWORD"))
	// fmt.Println(os.Getenv("DB_HOST"))
	// fmt.Println(os.Getenv("DB_PORT"))
	// fmt.Println(os.Getenv("DB_NAME"))
>>>>>>> fe9f630c60fdee46caca79950757ff2b94ca695d

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
	if err := db.InitDBData(); err!= nil {
		log.Fatalf("Failed to initialize data: %v", err)
	}

	// 注册 Swagger 路由
<<<<<<< HEAD
	router.GET("/swagger/*any", func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json"))(c)
	})
=======
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, ginSwagger.URL("/docs/swagger.json")))

>>>>>>> fe9f630c60fdee46caca79950757ff2b94ca695d
	router.POST("/agent/register", login.Register)
	router.POST("/agent/login", login.Login)
	// 需要 JWT 认证的路由
	auth := router.Group("/agent", middlewire.JWTAuthMiddleware())
	{
		auth.POST("/install", install.InstallAgent)
		auth.POST("/addSystemInfo", monitor.ReceiveAndStoreSystemMetrics)
		auth.GET("/list", monitor.ListAgent)
		router.GET("/monitor/:hostname", monitor.GetAgentInfo)
	}
	router.Run("0.0.0.0:8080")
}
