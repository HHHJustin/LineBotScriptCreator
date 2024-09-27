package api

import (
	"LineBotCreator/database"
	"fmt"
	"net/http"

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
	var req database.MessageCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var node database.Node
	nodeId := req.CurrentNodeID
	if err := db.Where("id = ?", nodeId).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	if node.Type != "Message" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is wrong type"})
		return
	}
	messageType := req.MessageType
	messageContent := req.MessageContent
	newMessage := database.Message{
		Type:    messageType,
		Content: messageContent,
		NodeID:  nodeId,
	}
	if err := db.Create(&newMessage).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Create Message fail"})
		return
	}
	if node.Range == nil {
		node.Range = make([]int, 0)
	}
	node.Range = append(node.Range, newMessage.MessageID)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Message": newMessage})
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
	var req database.MessageUpdateRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var message database.Message
	if err := db.Where("message_id = ?", req.MessageID).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Message does not exist"})
		return
	}

	message.Content = req.MessageContent

	if err := db.Save(&message).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Message updated successfully"})
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
	var req database.MessageDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	fmt.Println(req)
	var node database.Node
	nodeId := req.CurrentNodeID
	if err := db.Where("id = ?", nodeId).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	var message database.Message
	messageId := req.MessageID
	if err := db.Where("message_id = ?", messageId).First(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "message is not exist"})
		return
	}
	if err := db.Delete(&message).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete message"})
		return
	}
	node.Range = removeValue(node.Range, messageId)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}

// FlexMessage
