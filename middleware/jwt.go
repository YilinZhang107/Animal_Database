/*
* @Author: Oatmeal107
* @Date:   2023/6/16 17:42
 */

package middleware

import (
	"Animal_database/utils"
	"github.com/gin-gonic/gin"
	"time"
)

// JWT token验证中间件
func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
		} else {
			claims, err := utils.VerifyToken(token)
			if err != nil {
				code = utils.ErrorAuthCheckTokenFail
			} else if time.Now().Unix() > claims.ExpiresAt {
				code = utils.ErrorAuthCheckTokenTimeout
			}
		}
		if code != 200 {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    utils.GetMsg(code),
			})
			c.Abort() //调用Abort以确保剩余的处理程序不会被调用。
			return
		}
		c.Next()
	}
}
