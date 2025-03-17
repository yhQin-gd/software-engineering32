package info

import (
	m_init "cmd/server/model/init"
	u "cmd/server/model/user"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func GetUserInfo(c *gin.Context) {
	username, exists := c.Get("username")
	if !exists {
		log.Printf("用户还未登录")
		c.JSON(401, gin.H{
			"message": "用户未登录",
		})
	}

	var user u.User
	err := m_init.DB.Where("name = ?", username.(string)).First(&user).Error
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "用户不存在"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "获取用户信息成功",
		"user":    user,
	})
}
