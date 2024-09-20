package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"gorm.io/gorm"
)

func CallbackHandler(c *gin.Context, bot *linebot.Client, db *gorm.DB) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request"})
		}
		return
	}
	for _, event := range events {
		switch event.Type {

		default:
			log.Printf("不支援的事件類型: %s", event.Type)
		}
	}
	c.Status(http.StatusOK)
}
