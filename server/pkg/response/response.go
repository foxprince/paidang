package response

import "github.com/gin-gonic/gin"

// 统一返回体：{ code, msg, data }
func OK(c *gin.Context, data any) {
	c.JSON(200, gin.H{"code": 0, "msg": "ok", "data": data})
}

func Err(c *gin.Context, code int, msg string) {
	c.JSON(200, gin.H{"code": code, "msg": msg, "data": nil})
}

func NotImpl(c *gin.Context) {
	Err(c, 50001, "未实现")
}
