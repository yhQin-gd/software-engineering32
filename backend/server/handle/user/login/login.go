package login

import (
	"cmd/server/middlewire"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

type User struct {
	ID         int
	Name       string
	Email      string
	Password   string
	IsVerified bool
}

// 初始化数据库并创建用户表
func InitDB() (*sql.DB, error) {
	// 连接 PostgreSQL 数据库，替换连接信息
	connStr := "host=192.168.31.251 port=5432 user=postgres password=cCyjKKMyweCer8f3 dbname=monitor sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("连接数据库时出错: %v", err)
	}

	// 检查数据库连接
	if err = db.Ping(); err != nil {
		return nil, fmt.Errorf("检查连接时出错: %v", err)
	}
	return db, err
}

// 注册
func Register(c *gin.Context) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("数据库连接错误: %v", err)
	}

	// 定义用于接收 JSON 数据的结构体
	var input struct {
		Name     string `json:"name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}

	// 数据验证
	if len(input.Name) == 0 || len(input.Email) == 0 || len(input.Password) < 6 {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "参数不完整或密码长度不足"})
		return
	}

	// 检查邮箱是否存在
	var user User
	err = db.QueryRow("SELECT id FROM users WHERE email = $1", input.Email).Scan(&user.ID)
	if err == nil {
		c.JSON(http.StatusUnprocessableEntity, gin.H{"message": "用户已存在"})
		return
	}
	// 创建用户
	_, err = db.Exec("INSERT INTO users (name, email, password, is_verified) VALUES ($1, $2, $3, $4)", input.Name, input.Email, input.Password, true)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "用户创建失败"})
		return
	}

	// 返回结果，包括加密后的密码
	c.JSON(http.StatusOK, gin.H{
		"message": "注册成功",
	})
}

// 登录
func Login(c *gin.Context) {
	db, err := InitDB()
	if err != nil {
		log.Fatalf("数据库连接错误: %v", err)
	}
	defer db.Close()

	// 定义用于接收 JSON 数据的结构体
	var input struct {
		Name     string `json:"name"`
		Password string `json:"password"`
	}

	// 解析 JSON 数据
	if err := c.BindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "请求数据格式错误"})
		return
	}

	// 查找用户
	var user User
	err = db.QueryRow("SELECT id, password FROM users WHERE name = $1", input.Name).Scan(&user.ID, &user.Password)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
		return
	}

	// 验证密码
	if user.Password != input.Password {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "密码错误"})
		return
	}
	// 生成 JWT
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := &middlewire.Claims{
		Username: input.Name,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(middlewire.JwtKey)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"message": "生成 token 错误"})
		return
	}
	// 登录成功
	c.JSON(http.StatusOK, gin.H{
		"message": "登录成功",
		"token":   tokenString,
	})
}
