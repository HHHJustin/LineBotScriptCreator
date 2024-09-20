package api

import (
	"LineBotCreator/database"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// CreateNode godoc
// @Summary Create a new node
// @Description Create a new node with the specified type and range
// @Tags nodes
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeTitle formData string true "nodeTitle for the Node"
// @Param nodeType formData string true "nodeType for the Node"
// @Param location formData string true "Location for the Node"
// @Success 201 {object} map[string]interface{} "Successfully created Node"
// @Failure 400 {object} map[string]interface{} "Create Node fail"
// @Router /nodes/create [post]
func CreateNodeHandler(c *gin.Context, db *gorm.DB) {
	nodeTitle := c.PostForm("nodeTitle")
	nodeType := c.PostForm("nodeType")
	location := c.PostForm("location")
	newNode := database.Node{
		Title:        nodeTitle,
		Type:         nodeType,
		Range:        database.IntArray{},
		PreviousNode: 0,
		NextNode:     0,
		Loc:          location,
	}

	if err := db.Create(&newNode).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Create Node fail"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Node": newNode})
}

func ReadNodeHandler(c *gin.Context, db *gorm.DB) {
	var nodes []database.Node
	if err := db.Find(&nodes).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	graph := database.Graph{
		Nodes: []database.GraphNode{},
		Links: []database.Link{},
	}
	for _, node := range nodes {
		graph.Nodes = append(graph.Nodes, database.GraphNode{
			Key:   node.ID,
			Text:  node.Title,
			Color: getColorByType(node.Type),
			Loc:   node.Loc,
		})
		if node.NextNode != 0 {
			graph.Links = append(graph.Links, database.Link{From: node.ID, To: node.NextNode})
		}
	}
	c.JSON(http.StatusOK, graph)
}

// UpdateNodePreviousHandler godoc
// @Summary Update the previous node of a specific node
// @Description Update the previous node ID for the specified node
// @Tags nodes
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param nodePrevious formData int true "Previous Node ID"
// @Success 200 {object} map[string]interface{} "Updated node"
// @Failure 400 {object} map[string]interface{} "Invalid nodePreviousId"
// @Failure 400 {object} map[string]interface{} "Node is not exist"
// @Failure 500 {object} map[string]interface{} "Failed to update node"
// @Router /nodes/previous [post]
func UpdateNodePreviousHandler(c *gin.Context, db *gorm.DB) {
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
	nodePreviousId := c.PostForm("nodePrevious")
	nodePreviousIdInt, err := strconv.Atoi(nodePreviousId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodePreviousId"})
		return
	}
	node.PreviousNode = nodePreviousIdInt
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Node": node})
}

// UpdateNodeNextHandler godoc
// @Summary Update the next node of a specific node
// @Description Update the next node ID for the specified node
// @Tags nodes
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param nodeNext formData int true "Next Node ID"
// @Success 200 {object} map[string]interface{} "Updated node"
// @Failure 400 {object} map[string]interface{} "Invalid nodeNextId"
// @Failure 400 {object} map[string]interface{} "Node is not exist"
// @Failure 500 {object} map[string]interface{} "Failed to update node"
// @Router /nodes/next [post]
func UpdateNodeNextHandler(c *gin.Context, db *gorm.DB) {
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
	nodeNextId := c.PostForm("nodeNext")
	nodeNextIdInt, err := strconv.Atoi(nodeNextId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodePreviousId"})
		return
	}
	node.NextNode = nodeNextIdInt
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Node": node})
}

// UpdateNodeTitleHandler godoc
// @Summary Update the title of a specific node
// @Description Update the title for the specified node
// @Tags nodes
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData int true "Node ID"
// @Param nodeTitle formData string true "New Node Title"
// @Success 200 {object} map[string]interface{} "Updated node"
// @Failure 400 {object} map[string]interface{} "Invalid nodeId"
// @Failure 400 {object} map[string]interface{} "Node is not exist"
// @Failure 500 {object} map[string]interface{} "Failed to update node"
// @Router /nodes/title [post]
func UpdateNodeTitleHandler(c *gin.Context, db *gorm.DB) {
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
	nodeTitle := c.PostForm("nodeTitle")
	node.Title = nodeTitle
	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"Node": node})
}

// DeleteNode godoc
// @Summary Delete a new node
// @Description Delete a new node with the specified ID
// @Tags nodes
// @Accept x-www-form-urlencoded
// @Produce json
// @Param nodeId formData string true "nodeID for the Node"
// @Success 201 {object} map[string]interface{} "Successfully delete Node"
// @Failure 400 {object} map[string]interface{} "Delete Node fail"
// @Router /nodes/delete [post]
func DeleteNodeHandler(c *gin.Context, db *gorm.DB) {
	var node database.Node
	nodeId := c.PostForm("nodeId")
	nodeIdInt, err := strconv.Atoi(nodeId)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid nodeId format"})
		return
	}
	if err := db.Where("id = ?", nodeIdInt).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	if node.PreviousNode != 0 {
		var previousNode database.Node
		if err := db.Where("id = ?", node.PreviousNode).First(&previousNode).Error; err == nil {
			previousNode.NextNode = node.NextNode
			if err := db.Save(&previousNode).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update previous node"})
				return
			}
		}
	}
	if node.NextNode != 0 {
		var nextNode database.Node
		if err := db.Where("id = ?", node.NextNode).First(&nextNode).Error; err == nil {
			nextNode.PreviousNode = node.PreviousNode
			if err := db.Save(&nextNode).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update next node"})
				return
			}
		}
	}
	if err := db.Delete(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Node deleted successfully"})
}
