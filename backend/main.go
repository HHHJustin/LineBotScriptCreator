package main

import (
	"LineBotCreator/api"
	db "LineBotCreator/database"
	"LineBotCreator/utils"
	"net/http"

	_ "LineBotCreator/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title Line Bot API
// @version 1.0
// @description This is a sample server for Line Bot API.
// @BasePath /
func main() {
	db := db.Connect()
	router := gin.Default()

	templatesPath := "templates/*" // 使用相對路徑
	staticPath := "static"

	router.LoadHTMLGlob(templatesPath)
	router.Static("/assets", staticPath)
	router.POST("/callback", func(c *gin.Context) {
		bot := utils.ConnectLineBot(c, db)
		if bot == nil {
			return
		}
		api.CallbackHandler(c, bot, db)
	})
	// LineBot Channel
	channelRouter := router.Group("/channel")
	channelRouter.GET("/read", func(c *gin.Context) {
		c.HTML(http.StatusOK, "channel.html", nil)
	})
	channelRouter.POST("/create", func(c *gin.Context) {
		api.CreateChannelInfo(c, db)
	})
	// Node
	nodeRouter := router.Group("/nodes")
	nodeRouter.GET("/get", func(c *gin.Context) {
		api.ReadNodeHandler(c, db)
	})
	nodeRouter.GET("/read", func(c *gin.Context) {
		c.HTML(http.StatusOK, "nodes.html", nil)
	})
	nodeRouter.GET("/type", func(c *gin.Context) {
		api.GetNodeTypeHandler(c, db)
	})
	nodeRouter.GET("/get/:nodeID/:nodeType", func(c *gin.Context) {
		api.EditPageHandler(c, db)
	})
	nodeRouter.POST("/updatelocation", func(c *gin.Context) {
		api.UpdateLocationHandler(c, db)
	})
	nodeRouter.POST("/title", func(c *gin.Context) {
		api.UpdateNodeTitleHandler(c, db)
	})
	nodeRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteNodeHandler(c, db)
	})

	// New Node
	newNodeRouter := nodeRouter.Group("/create").Use(func(c *gin.Context) {
		api.CreateNodeMiddleware(c)
	})
	newNodeRouter.POST("/next", func(c *gin.Context) {
		api.CreateNextNodeHandler(c, db)
	})
	newNodeRouter.POST("/previous", func(c *gin.Context) {
		api.CreatePreviousNodeHandler(c, db)
	})
	newNodeRouter.POST("/firststep", func(c *gin.Context) {
		api.CreateFirstStepHandler(c, db)
	})
	newNodeRouter.POST("/branch", func(c *gin.Context) {
		api.CreateBranchHandler(c, db)
	})
	newNodeRouter.POST("/keywordDecision", func(c *gin.Context) {
		api.CreateKWDecisionHandler(c, db)
	})

	// Link
	linkRouter := router.Group("/links")
	linkRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteLinkHandler(c, db)
	})
	linkRouter.POST("/create", func(c *gin.Context) {
		api.CreateLinkHandler(c, db)
	})

	// First Step
	firstStepRouter := router.Group("/firststep")
	firstStepRouter.GET("/read", func(c *gin.Context) {
		api.FirstStepPageHandler(c, db)
	})
	firstStepRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteFirstStepHandler(c, db)
	})

	// Message
	messageRouter := router.Group("/messages")
	messageRouter.POST("/create", func(c *gin.Context) {
		api.CreateMessageHandler(c, db)
	})
	messageRouter.POST("/update", func(c *gin.Context) {
		api.UpdateMessageHandler(c, db)
	})
	messageRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteMessageHandler(c, db)
	})
	messageRouter.POST("/updateorder", func(c *gin.Context) {
		api.UpdateMessageOrderHandler(c, db)
	})

	// Keyword Decision
	keywordDecisionRouter := router.Group("/keywordDecisions")
	keywordDecisionRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteKWDecisionHandler(c, db)
	})
	keywordDecisionRouter.POST("/update", func(c *gin.Context) {
		api.UpdateKWDecisionHandler(c, db)
	})

	// QuickReply
	quickReplyRouter := router.Group("/quickreplies")
	quickReplyRouter.POST("/create", func(c *gin.Context) {
		api.CreateQuickReplyHandler(c, db)
	})
	quickReplyRouter.POST("/update", func(c *gin.Context) {
		api.UpdateQuickReplyHandler(c, db)
	})
	quickReplyRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteQuickReplyHandler(c, db)
	})
	// Swagger ui
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
