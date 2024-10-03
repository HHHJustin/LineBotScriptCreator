package api

import (
	"LineBotCreator/database"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateBranchHandler(c *gin.Context, db *gorm.DB) {
	var req database.NodeCreateRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var currentNode database.Node
	if err := db.Where("id = ?", req.CurrentNodeID).First(&currentNode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	node, exists := c.Get("node")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node data not found"})
		return
	}
	newNode := node.(database.Node)
	newNode.Title = "New"
	newNode.Type = req.NewNodeType
	newNode.PreviousNode = append(newNode.PreviousNode, req.CurrentNodeID)
	newNode.NextNode = 0
	newNode.LocX = currentNode.LocX + 100
	newNode.LocY = len(currentNode.Range) * 100
	if err := db.Create(&newNode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	if err := branchRangeAppend(&currentNode, &newNode, db); err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, gin.H{"message": "Branch created successfully"})
}

func branchRangeAppend(currentNode *database.Node, newNode *database.Node, db *gorm.DB) error {
	switch currentNode.Type {
	case "KeywordDecision":
		newKWDecision := database.KeywordDecision{
			NodeID:       currentNode.ID,
			NextNode:     newNode.ID,
			NextNodeType: newNode.Type,
		}
		if err := db.Create(&newKWDecision).Error; err != nil {
			return fmt.Errorf("failed to create KeywordDecision: %w", err)
		}
		currentNode.Range = append(currentNode.Range, newKWDecision.KWDecisionID)
		if err := db.Save(&currentNode).Error; err != nil {
			return fmt.Errorf("Failed to update node")
		}
	default:
		return fmt.Errorf("unsupported node type: %s", currentNode.Type)
	}
	return nil
}
