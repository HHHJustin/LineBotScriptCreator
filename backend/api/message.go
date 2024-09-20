package api

import (
	"LineBotCreator/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateMessage godoc
// @Summary Create a new message for a specific node
// @Description Create a new message associated with a given node
// @Tags messages
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param messageType formData string true "Message Type"
// @Param messageContent formData string true "Message Content"
// @Success 200 {object} map[string]interface{} "Created message"
// @Failure 400 {object} map[string]interface{} "Invalid nodeId format"
// @Failure 400 {object} map[string]interface{} "Node does not exist"
// @Failure 400 {object} map[string]interface{} "Node is wrong type"
// @Failure 500 {object} map[string]interface{} "Create Message fail"
// @Failure 500 {object} map[string]interface{} "Failed to update node"
// @Router /messages/create [post]
func CreateMessageHandler(c *gin.Context, db *gorm.DB) {
	var node database.Node
	nodeId := c.PostForm("nodeId")
	nodeIdInt, err := strconv.Atoi(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodeId"})
		return
	}
	if err := db.Where("id = ?", nodeIdInt).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	if node.Type != "message" {
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

func ReadMessageHandler(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, ".html", nil)
}

// UpdateMessage godoc
// @Summary Update a message by ID
// @Description Update the content of a specific message using its ID
// @Tags messages
// @Accept x-www-form-urlencoded
// @Produce json
// @Param messageId formData int true "Message ID"
// @Param newMessageContent formData string true "New Message Content"
// @Success 200 {object} map[string]interface{} "Updated message"
// @Failure 400 {object} map[string]interface{} "Invalid messageId"
// @Failure 400 {object} map[string]interface{} "Node does not exist"
// @Failure 400 {object} map[string]interface{} "Message ID does not exist in node range"
// @Failure 500 {object} map[string]interface{} "Failed to update message"
// @Router /messages/update [post]
func UpdateMessageHandler(c *gin.Context, db *gorm.DB) {
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
	if node.Type != "message" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is wrong type"})
		return
	}
	var message database.Message
	messageId := c.PostForm("messageId")
	messageIdInt, err := strconv.Atoi(messageId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid messageId"})
		return
	}
	if !contains(node.Range, messageIdInt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message ID does not exist in node range"})
		return
	}
	if err := db.Where("message_id = ?", messageIdInt).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is not exist"})
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

// DeleteMessage godoc
// @Summary Delete a message by ID
// @Description Delete a specific message associated with a node using its ID
// @Tags messages
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param messageId formData int true "Message ID"
// @Success 200 {object} map[string]interface{} "node object after deletion"
// @Failure 400 {object} map[string]interface{} "Invalid nodeId"
// @Failure 400 {object} map[string]interface{} "Node does not exist"
// @Failure 400 {object} map[string]interface{} "Invalid messageId"
// @Failure 400 {object} map[string]interface{} "Message does not exist"
// @Failure 500 {object} map[string]interface{} "Failed to delete message"
// @Router /messages/delete [post]
func DeleteMessageHandler(c *gin.Context, db *gorm.DB) {
	var node database.Node
	nodeId := c.PostForm("nodeId")
	nodeIdInt, err := strconv.Atoi(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodeId"})
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
	if err := db.Where("message_id = ?", messageIdInt).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is not exist"})
		return
	}
	if err := db.Delete(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete message"})
		return
	}
	node.Range = removeValue(node.Range, messageIdInt)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}

// FlexMessage
