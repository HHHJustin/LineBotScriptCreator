package api

import (
	"LineBotCreator/database"
	"errors"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func FirstStepPageHandler(c *gin.Context, db *gorm.DB) {
	var firstSteps []database.Node
	if err := db.Where("type = ?", "FirstStep").Find(&firstSteps).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch information"})
		return
	}
	var firstStepWithIndex []struct {
		Index     int
		FirstStep database.FirstStep
	}
	for i, fs := range firstSteps {
		firstStep := database.FirstStep{
			Type:     fs.Title,
			NextNode: fs.NextNode,
		}
		firstStepWithIndex = append(firstStepWithIndex, struct {
			Index     int
			FirstStep database.FirstStep
		}{
			Index:     i + 1,
			FirstStep: firstStep,
		})
	}
	c.HTML(http.StatusOK, "first_step.html", gin.H{
		"FirstSteps": firstStepWithIndex,
	})
}

func FirstStepHandler(c *gin.Context, db *gorm.DB) {
	var req database.FirstStepRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var existingNode database.Node
	if err := db.Where("title = ?", req.FirstStepType).First(&existingNode).Error; err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Node with this title already exists"})
		return
	} else if !errors.Is(err, gorm.ErrRecordNotFound) {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check existing node"})
		return
	}
	node, exists := c.Get("node")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node data not found"})
		return
	}
	var count int64
	if err := db.Model(&database.Node{}).Where("type = ?", "FirstStep").Count(&count).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to count nodes"})
		return
	}
	newNode := node.(database.Node)
	newNode.Title = req.FirstStepType
	newNode.Type = "FirstStep"
	newNode.PreviousNode = 0
	newNode.NextNode = 0
	newNode.LocX = 0
	newNode.LocY = (int(count)) * 100
	if err := db.Create(&newNode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Node": newNode})
}

func DeleteFirstStepHandler(c *gin.Context, db *gorm.DB) {
	var req database.FirstStepRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var node database.Node
	if err := db.Where("type = ? AND title = ?", "FirstStep", req.FirstStepType).First(&node).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "Node not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch node"})
		}
		return
	}

	if err := db.Delete(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node deleted successfully"})

}
