package api

import (
	"LineBotCreator/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateLinkHandler(c *gin.Context, db *gorm.DB) {
	var req database.LinkCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var fromNode database.Node
	if err := db.Where("id = ?", req.FromNodeID).First(&fromNode).Error; err == nil {
		fromNode.NextNode = req.ToNodeID
		if err := db.Save(&fromNode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fromNode"})
			return
		}
	}

	var toNode database.Node
	if err := db.Where("id = ?", req.ToNodeID).First(&toNode).Error; err == nil {
		toNode.PreviousNode = req.FromNodeID
		if err := db.Save(&toNode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update toNode"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link created successfully"})
}

func DeleteLinkHandler(c *gin.Context, db *gorm.DB) {
	var req database.LinkDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var fromNode database.Node
	if err := db.Where("id = ?", req.FromNodeID).First(&fromNode).Error; err == nil {
		fromNode.NextNode = 0
		if err := db.Save(&fromNode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update fromNode"})
			return
		}
	}

	var toNode database.Node
	if err := db.Where("id = ?", req.ToNodeID).First(&toNode).Error; err == nil {
		toNode.PreviousNode = 0
		if err := db.Save(&toNode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update toNode"})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{"message": "Link deleted successfully"})
}
