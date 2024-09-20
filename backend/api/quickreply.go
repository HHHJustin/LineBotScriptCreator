package api

import (
	"LineBotCreator/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateQuickReplyHandler godoc
// @Summary Create a new quick reply for a specific node
// @Description Create a new quick reply associated with a given node
// @Tags quickreplies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param buttonName formData string true "Button Name"
// @Param reply formData string true "Reply Content"
// @Success 200 {object} map[string]interface{} "Quick reply created successfully"
// @Failure 400 {object} map[string]interface{} "Invalid nodeId"
// @Failure 400 {object} map[string]interface{} "Node does not exist"
// @Failure 400 {object} map[string]interface{} "Node is wrong type"
// @Failure 500 {object} map[string]interface{} "Create QuickReply fail"
// @Failure 500 {object} map[string]interface{} "Failed to update node"
// @Router /quickreplies/create [post]
func CreateQuickReplyHandler(c *gin.Context, db *gorm.DB) {
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

	if node.Type != "quickreply" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is wrong type"})
		return
	}
	buttomName := c.PostForm("buttonName")
	reply := c.PostForm("reply")
	newQuickReply := database.QuickReply{
		ButtonName: buttomName,
		Reply:      reply,
		NodeID:     nodeIdInt,
	}
	if err := db.Create(&newQuickReply).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Create QuickReply fail"})
		return
	}
	node.Range = append(node.Range, newQuickReply.QuickReplyID)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"QuickReply": newQuickReply})
}

func ReadQuickReplyHandler(c *gin.Context, db *gorm.DB) {
	c.HTML(http.StatusOK, ".html", nil)
}

// UpdateQuickReplyHandler godoc
// @Summary Update a quick reply by ID
// @Description Update the button name and reply content of a specific quick reply associated with a node
// @Tags quickreplies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param quickReplyId formData int true "Quick Reply ID"
// @Param newButtonName formData string true "New Button Name"
// @Param newReply formData string true "New Reply Content"
// @Success 200 {object} map[string]interface{} "Quick reply updated successfully"
// @Failure 400 {object} map[string]interface{} "Invalid nodeId"
// @Failure 400 {object} map[string]interface{} "Node does not exist"
// @Failure 400 {object} map[string]interface{} "Node is wrong type"
// @Failure 400 {object} map[string]interface{} "Invalid quickReplyId"
// @Failure 400 {object} map[string]interface{} "QuickReply ID does not exist in node range"
// @Failure 400 {object} map[string]interface{} "QuickReply does not exist"
// @Failure 500 {object} map[string]interface{} "Failed to update quick reply"
// @Router /quickreplies/update [post]
func UpdateQuickReplyHandler(c *gin.Context, db *gorm.DB) {
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
	if node.Type != "quickreply" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is wrong type"})
		return
	}
	var quickReply database.QuickReply
	quickReplyId := c.PostForm("quickReplyId")
	quickReplyIdInt, err := strconv.Atoi(quickReplyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quickReplyId"})
		return
	}
	if !contains(node.Range, quickReplyIdInt) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QuickReply ID does not exist in node range"})
		return
	}
	if err := db.Where("quick_reply_id = ?", quickReplyIdInt).First(&quickReply).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QuickReply is not exist"})
		return
	}
	newButtonName := c.PostForm("newButtonName")
	newReply := c.PostForm("newReply")
	quickReply.ButtonName = newButtonName
	quickReply.Reply = newReply
	if err := db.Save(&quickReply).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update message"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"quickReply": quickReply})
}

// DeleteQuickReplyHandler godoc
// @Summary Delete a quick reply by ID
// @Description Delete a specific quick reply associated with a node using its ID
// @Tags quickreplies
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param quickReplyId formData int true "Quick Reply ID"
// @Success 200 {object} map[string]interface{} "Node updated successfully after deleting quick reply"
// @Failure 400 {object} map[string]interface{} "Invalid nodeId"
// @Failure 400 {object} map[string]interface{} "Node does not exist"
// @Failure 400 {object} map[string]interface{} "Invalid quickReplyId"
// @Failure 400 {object} map[string]interface{} "QuickReply does not exist"
// @Failure 500 {object} map[string]interface{} "Failed to delete QuickReply"
// @Failure 500 {object} map[string]interface{} "Failed to update node"
// @Router /quickreplies/delete [post]
func DeleteQuickReplyHandler(c *gin.Context, db *gorm.DB) {
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
	var quickReply database.QuickReply
	quickReplyId := c.PostForm("quickReplyId")
	quickReplyIdInt, err := strconv.Atoi(quickReplyId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid quickReplyId"})
		return
	}
	if err := db.Where("quick_reply_id = ?", quickReplyIdInt).First(&quickReply).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "QuickReply is not exist"})
		return
	}
	if err := db.Delete(&quickReply).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete QuickReply"})
		return
	}
	node.Range = removeValue(node.Range, quickReplyIdInt)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}
