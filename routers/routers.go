package routers

import (
	"03.threeGames/logger"
	"03.threeGames/middleWares"
	"github.com/gin-gonic/gin"
	"github.com/mattn/go-colorable"
	"github.com/spf13/viper"
	"net/http"
)

func SetUp() *gin.Engine {
	gin.ForceConsoleColor()
	gin.DefaultWriter = colorable.NewColorableStdout()

	var (
		mode = viper.GetString("app.mode")
	)
	if mode == gin.ReleaseMode {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(logger.GinLogger(), logger.GinRecovery(true), middleWares.Cors())

	routers(r)

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "All being OK!")
	})
	r.NoRoute(func(c *gin.Context) {
		c.JSON(http.StatusNotFound, gin.H{
			"msg": "404",
		})
	})
	return r
}

func routers(r *gin.Engine) *gin.Engine  {
	r.GET("/get", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"msg": "hello",
			"data": "none",
			"code": 200,
		})
	})

	return r
}
