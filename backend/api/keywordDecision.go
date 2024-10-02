package api

import (
	"LineBotCreator/database"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateKWDecisionHandler(c *gin.Context, db *gorm.DB) {

}

func DeleteKWDecisionHandler(c *gin.Context, db *gorm.DB) {
	var req database.KeywordDecisionDeleteRequest
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
	var keywordDecision database.KeywordDecision
	keywordDecisionId := req.KeywordDecisionID
	if err := db.Where("kw_decision_id = ?", keywordDecisionId).First(&keywordDecision).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword Decision is not exist"})
		return
	}
	if err := db.Delete(&keywordDecision).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete Keyword Decision"})
		return
	}
	node.Range = removeValue(node.Range, keywordDecisionId)
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"node": node})
}

func UpdateKWDecisionHandler(c *gin.Context, db *gorm.DB) {
	var req database.KeywordDecisionUpdateRequest
	if err := c.Bind(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var keywordDecision database.KeywordDecision
	if err := db.Where("kw_decision_id = ?", req.KeywordDecisionID).First(&keywordDecision).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword Decision does not exist"})
		return
	}
	keywordDecision.Keyword = req.Keyword

	if err := db.Save(&keywordDecision).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update keyword"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Keyword updated successfully"})
}
