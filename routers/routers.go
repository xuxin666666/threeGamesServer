package routers

import (
	"03.threeGames/controller"
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

	auth := r.Group("/auth")
	{
		auth.POST("/register", controller.AuthRegister)
		auth.POST("", controller.AuthLogin)
		auth.Use(middleWares.JWTAuthMiddleware())
		auth.GET("/auto", controller.AuthLoginAuto)
		auth.POST("/changeAvatar", controller.AuthChangeAvatar)
	}
	r.Use(middleWares.JWTAuthMiddleware())

	tetris := r.Group("/tetris")
	{
		tetris.GET("", controller.GetTetrisScores)
		tetris.POST("/postScores", controller.PostTetrisScores)
		tetris.POST("/postScores/list", controller.PostTetrisScoresList)
	}

	mineSweep := r.Group("/mineSweep")
	{
		mineSweep.GET("", controller.GetMineSweepScores)
		mineSweep.POST("/postScores", controller.PostMineSweepScores)
		mineSweep.POST("/postScores/list", controller.PostMineSweepScoresList)
	}

	community := r.Group("/community")
	{
		community.POST("", controller.CommunityGetPageList)
		community.POST("/createPage", controller.CommunityCreatePage)
		community.GET("/:page_id", controller.CommunityGetPageDetail)
		community.POST("/:page_id", controller.CommunityModifyPage)
		community.POST("/:page_id/approve", controller.CommunityPageApprove)

		community.POST("/:page_id/comments", controller.CommunityGetComments)
		community.POST("/:page_id/comments/insert", controller.CommunityAddComment)
		community.POST("/:page_id/comments/reply", controller.CommunityAddReply)
		community.POST("/:page_id/comments/approve", controller.CommunityCommentApprove)
	}

	return r
}
