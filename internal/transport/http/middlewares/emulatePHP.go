package middlewares

import "github.com/gin-gonic/gin"

// EmulatePHP makes things more funny
func EmulatePHP() func(c *gin.Context) {
	return func(c *gin.Context) {
		c.Header("X-Powered-By", "PHP/5.6.14") //https://schd.io/ENp
	}
}
