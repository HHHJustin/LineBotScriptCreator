package api

import (
	"LineBotCreator/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateMessage(c *gin.Context, db *gorm.DB) {
	var node database.Node
	nodeId := c.PostForm("nodeId")
	nodeIdInt, err := strconv.Atoi(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodePreviousId"})
		return
	}
	if err := db.Where("id = ?", nodeIdInt).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	if node.Type != "Message" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is wrong type"})
		return
	}
	messageType := c.PostForm("messageType")
	messageContent := c.PostForm("messageContent")
	newMessage := database.Message{
		Type:    messageType,
		Content: messageContent,
		NodeID:  nodeIdInt,
	}
	if err := db.Create(&newMessage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Create Message fail"})
		return
	}
	node.Range = append(node.Range, newMessage.MessageID)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": newMessage})
}

func UpdateMessage(c *gin.Context, db *gorm.DB) {
	var message database.Message
	messageId := c.PostForm("messageId")
	messageIdInt, err := strconv.Atoi(messageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodeId"})
		return
	}
	if err := db.Where("id = ?", messageIdInt).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	newMessageContent := c.PostForm("newMessageContent")
	message.Content = newMessageContent
	if err := db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": message})
}

func DeleteMessage(c *gin.Context, db *gorm.DB) {
	var node database.Node
	nodeId := c.PostForm("nodeId")
	nodeIdInt, err := strconv.Atoi(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodePreviousId"})
		return
	}
	if err := db.Where("id = ?", nodeIdInt).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	var message database.Message
	messageId := c.PostForm("messageId")
	messageIdInt, err := strconv.Atoi(messageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messsageId"})
		return
	}
	if err := db.Where("id = ?", messageIdInt).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is not exist"})
		return
	}
	removeValue(node.Range, messageIdInt)
	if err := db.Delete(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}

func removeValue(arr []int, value int) []int {
	var result []int
	for _, v := range arr {
		if v != value {
			result = append(result, v)
		}
	}
	return result
}

// FlexMessage最後研究
