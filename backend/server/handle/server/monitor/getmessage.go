package monitor

import (
	"cmd/server/middlewire"
	"cmd/server/model"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)

// 接收数据
func GetMessage(c *gin.Context) {
	// 初始化数据库
	db, err := model.InitDB()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	} else {
		fmt.Println("Init Successfully")
	}
	defer db.Close()

	// 解析请求数据
	var requestData struct {
		CPUInfo  []model.CPUInfo     `json:"cpu_info"`
		HostInfo model.HostInfo      `json:"host_info"`
		MemInfo  model.MemoryInfo    `json:"mem_info"`
		ProInfo  []model.ProcessInfo `json:"pro_info"`
		NetInfo  model.NetworkInfo   `json:"net_info"`
	}
	if err := c.ShouldBindJSON(&requestData); err != nil {
		s := fmt.Sprintf("Invalid JSON data: %s", err)
		c.JSON(http.StatusBadRequest, gin.H{"error": s})
		return
	}
	tokenh := requestData.HostInfo.Token
	var tokens string
	querySQL := `
	SELECT token (SELECT 1 FROM hostandtoken WHERE host_name = $1)
	FROM hostandtoken WHERE hostname = $1`

	err = db.QueryRow(querySQL, requestData.HostInfo.Hostname).Scan(&tokens)
	if err == sql.ErrNoRows {
		fmt.Println("No matching token.")
	}
	if len(tokenh) != 16 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Wrong token length"})
		return
	} else if tokenh != tokens {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unequal string"})
		return
	}

	// 从 Authorization Header 中提取 JWT token
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header is missing"})
		return
	}
	// 提取 token（格式为 "Bearer <token>"）
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == authHeader {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token format"})
		return
	}

	// 解析 token 并验证
	claims := &middlewire.Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return middlewire.JwtKey, nil
	})
	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
		return
	}
	// 更新心跳时间和状态为在线
	updateSQL := `
    UPDATE hostandtoken 
    SET last_heartbeat = NOW(), status = 'online' 
    WHERE host_name = $1`
	_, err = db.Exec(updateSQL, requestData.HostInfo.Hostname)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update heartbeat and status"})
		return
	}
	// 从解析的 token 中获取 username存入数据库
	username := claims.Username
	// 将数据插入数据库
	err = model.InsertSystemInfo(db, requestData.CPUInfo, requestData.MemInfo, requestData.HostInfo, requestData.ProInfo, requestData.NetInfo, username)
	if err != nil {
		s := fmt.Sprintf("Failed to insert system info: %s", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": s})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"status": "System information inserted successfully"})
}
