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
		case linebot.EventTypeFollow:
			nextNode, err := addFriendHandler(event, db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			gotoNextNode(nextNode, db, bot, event.Source.UserID)
		case linebot.EventTypeJoin:
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				nextNode, err := checkMessageCondition(event, db, message)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				gotoNextNode(nextNode, db, bot, event.Source.UserID)

			default:
				log.Printf("不支援的訊息類型: %T", message)
			}
		default:
			log.Printf("不支援的事件類型: %s", event.Type)
		}
	}
	c.Status(http.StatusOK)
}
