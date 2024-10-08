package utils

import (
	"LineBotCreator/database"
	"errors"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"gorm.io/gorm"
)

func GetLineBotChannelSetting(db *gorm.DB) (*database.LineBotChannelSetting, error) {
	var setting database.LineBotChannelSetting
	if err := db.First(&setting).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("no channel settings found")
		}
		return nil, err
	}
	return &setting, nil
}

func ConnectLineBot(c *gin.Context, db *gorm.DB) *linebot.Client {
	var setting database.LineBotChannelSetting
	if err := db.First(&setting).Error; err != nil {
		c.HTML(http.StatusOK, "channel.html", nil)
		return nil
	}
	bot, err := linebot.New(setting.ChannelSecretKey, setting.ChannelAccessToken)
	if err != nil {
		log.Fatal(err)
		return nil
	}
	return bot
}
