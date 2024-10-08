package api

import (
	"LineBotCreator/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateChannelInfo(c *gin.Context, db *gorm.DB) {
	var req database.ChannelRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	channel := database.LineBotChannelSetting{
		ChannelSecretKey:   req.ChannelSecretKey,
		ChannelAccessToken: req.ChannelAccessToken,
	}
	if err := db.Create(&channel).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create channel"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Channel": channel})
}
