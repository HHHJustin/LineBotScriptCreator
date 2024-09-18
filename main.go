package main

import (
	"LineBotCreator/api"
	db "LineBotCreator/database"
	"LineBotCreator/utils"

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
	bot := utils.ConnectLineBot()
	db := db.Connect()
	router := gin.Default()
	router.POST("/callback", func(c *gin.Context) {
		api.CallbackHandler(c, bot, db)
	})

	// Node operator
	nodeRouter := router.Group("/nodes")
	nodeRouter.POST("/create", func(c *gin.Context) {
		api.CreateNodeHandler(c, db)
	})
	nodeRouter.POST("/previous", func(c *gin.Context) {
		api.UpdateNodePreviousHandler(c, db)
	})
	nodeRouter.POST("/next", func(c *gin.Context) {
		api.UpdateNodeNextHandler(c, db)
	})
	nodeRouter.POST("/title", func(c *gin.Context) {
		api.UpdateNodeTitleHandler(c, db)
	})
	nodeRouter.POST("/delete", func(c *gin.Context) {
		api.DeleteNodeHandler(c, db)
	})
	// Swagger ui
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	if err := router.Run(":8080"); err != nil {
		panic(err)
	}
}
