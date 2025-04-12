package middle

import (
	"github.com/gin-gonic/gin"
)

func isLoggedIn(c *gin.Context) bool {
	userId := c.GetString("userId")
	return userId != " "

}
