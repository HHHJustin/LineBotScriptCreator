package api

import (
	"LineBotCreator/database"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func CreateNodeMiddleware(c *gin.Context) {
	newNode := database.Node{
		Title: "New",
		Range: database.IntArray{},
	}
	c.Set("node", newNode)
	c.Next()
}

func CreateNextNodeHandler(c *gin.Context, db *gorm.DB) {
	node, exists := c.Get("node")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node data not found"})
		return
	}
	newNode := node.(database.Node)
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
	newNode.Type = req.NewNodeType
	newNode.PreviousNode = append(newNode.PreviousNode, currentNode.ID)
	newNode.NextNode = currentNode.NextNode
	newNode.LocX = currentNode.LocX + 200
	newNode.LocY = currentNode.LocY
	if err := db.Create(&newNode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	if currentNode.NextNode != 0 {
		var nextNode database.Node
		if err := db.Where("id = ?", currentNode.NextNode).First(&nextNode).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
			return
		}
		nextNode.PreviousNode = removeValue(nextNode.PreviousNode, currentNode.ID)
		nextNode.PreviousNode = append(nextNode.PreviousNode, newNode.ID)
		if err := db.Save(&nextNode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
			return
		}
	}
	currentNode.NextNode = newNode.ID
	if err := db.Save(&currentNode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"Node": newNode})
}

func CreatePreviousNodeHandler(c *gin.Context, db *gorm.DB) {
	node, exists := c.Get("node")
	if !exists {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node data not found"})
		return
	}
	newNode := node.(database.Node)
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
	newNode.Type = req.NewNodeType
	newNode.PreviousNode = currentNode.PreviousNode
	newNode.NextNode = req.CurrentNodeID
	newNode.LocX = currentNode.LocX - 200
	newNode.LocY = currentNode.LocY - 200
	if err := db.Create(&newNode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create node"})
		return
	}
	if currentNode.PreviousNode != nil {
		var previousNode database.Node
		if err := db.Where("id = ?", currentNode.PreviousNode).First(&previousNode).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
			return
		}
		previousNode.NextNode = newNode.ID
		if err := db.Save(&previousNode).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
			return
		}
	}
	currentNode.PreviousNode = append(currentNode.PreviousNode, newNode.ID)
	if err := db.Save(&currentNode).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node"})
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
		Loc := fmt.Sprintf("%d %d", node.LocX, node.LocY)
		graph.Nodes = append(graph.Nodes, database.GraphNode{
			Key:   node.ID,
			Text:  node.Title,
			Color: getColorByType(node.Type),
			Loc:   Loc,
		})
		switch node.Type {
		case "Message", "FirstStep", "QuickReply":
			if node.NextNode != 0 {
				graph.Links = append(graph.Links, database.Link{From: node.ID, To: node.NextNode})
			}
		case "KeywordDecision":
			for _, nextID := range node.Range {
				var keywordDecision database.KeywordDecision
				if err := db.Where("kw_decision_id = ?", nextID).First(&keywordDecision).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						c.JSON(http.StatusConflict, gin.H{"error": fmt.Sprintf("Keyword decision with ID %d does not exist.", nextID)})
						return
					}
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
					return
				}
				graph.Links = append(graph.Links, database.Link{
					From: node.ID,
					To:   keywordDecision.NextNode,
				})
			}

		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Error node type."})
			return
		}

	}
	c.JSON(http.StatusOK, graph)
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
	var req database.NodeUpdateTitleRequest
	var node database.Node
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
		fmt.Println("BindJSON Error:", err) // 打印出錯誤
		return
	}
	nodeId := req.CurrentNodeID
	if err := db.Where("id = ?", nodeId).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	nodeTitle := req.NewTitle
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
	var req database.NodeDeleteRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	nodeIdInt := req.CurrentNodeID
	var node database.Node
	if err := db.Where("id = ?", nodeIdInt).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	if node.Type == "KeywordDecision" {
		for _, rangeID := range node.Range {
			var keywordDecision database.KeywordDecision
			if err := db.Where("kw_decision_id = ?", rangeID).First(&keywordDecision).Error; err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Keyword decision is not exist"})
				return
			}
			if keywordDecision.NextNode != 0 {
				var nextNode database.Node
				if err := db.Where("id = ?", keywordDecision.NextNode).First(&nextNode).Error; err == nil {
					nextNode.PreviousNode = removeValue(nextNode.PreviousNode, node.ID)
					if err := db.Save(&nextNode).Error; err != nil {
						c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update previous node"})
						return
					}
				}
			}

		}
	} else {
		if node.NextNode != 0 {
			var nextNode database.Node
			if err := db.Where("id = ?", node.NextNode).First(&nextNode).Error; err == nil {
				nextNode.PreviousNode = removeValue(nextNode.PreviousNode, node.ID)
				if err := db.Save(&nextNode).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update next node"})
					return
				}
			}
		}
	}
	if node.PreviousNode != nil {
		for _, prevNodeID := range node.PreviousNode {
			var prevNode database.Node
			if err := db.Where("id = ?", prevNodeID).First(&prevNode).Error; err == nil {
				prevNode.NextNode = 0
				if err := db.Save(&prevNode).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update previous node"})
					return
				}
			}
		}
	}

	if err := db.Delete(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to delete node"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Node deleted successfully"})
}

func GetNodeTypeHandler(c *gin.Context, db *gorm.DB) {
	currentNodeID := c.Query("currentNodeID")
	nodeIdInt, err := strconv.Atoi(currentNodeID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid node ID"})
		return
	}
	var node database.Node
	if err := db.Where("id = ?", nodeIdInt).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"nodeID":   node.ID,
		"nodeType": node.Type,
	})
}

func EditPageHandler(c *gin.Context, db *gorm.DB) {
	nodeID := c.Param("nodeID")
	nodeType := c.Param("nodeType")
	var node database.Node
	if err := db.Where("id = ?", nodeID).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node is not exist"})
		return
	}
	switch nodeType {
	case "Message":
		var messages []database.Message
		if err := db.Where("node_id = ?", nodeID).Find(&messages).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
			return
		}
		var messageWithIndex []struct {
			Index   int
			Message database.Message
		}
		for i, v := range messages {
			messageWithIndex = append(messageWithIndex, struct {
				Index   int
				Message database.Message
			}{
				Index:   i + 1,
				Message: v,
			})
		}
		c.HTML(http.StatusOK, "message.html", gin.H{
			"Node":     node,
			"Messages": messageWithIndex,
		})
	case "QuickReply":
		var quickReplies []database.QuickReply
		if err := db.Where("node_id = ?", nodeID).Find(&quickReplies).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
			return
		}
		var quickReplyWithIndex []struct {
			Index      int
			QuickReply database.QuickReply
		}
		for i, v := range quickReplies {
			quickReplyWithIndex = append(quickReplyWithIndex, struct {
				Index      int
				QuickReply database.QuickReply
			}{
				Index:      i + 1,
				QuickReply: v,
			})
		}
		c.HTML(http.StatusOK, "quickReply.html", gin.H{
			"Node":         node,
			"QuickReplies": quickReplyWithIndex,
		})
	case "KeywordDecision":
		var keywordDecisions []database.KeywordDecision
		if err := db.Where("node_id = ?", nodeID).Find(&keywordDecisions).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch messages"})
			return
		}
		var KWDecionWithIndex []struct {
			Index           int
			KeywordDecision database.KeywordDecision
		}
		for i, v := range keywordDecisions {
			KWDecionWithIndex = append(KWDecionWithIndex, struct {
				Index           int
				KeywordDecision database.KeywordDecision
			}{
				Index:           i + 1,
				KeywordDecision: v,
			})
		}
		c.HTML(http.StatusOK, "keywordDecision.html", gin.H{
			"Node":             node,
			"KeywordDecisions": KWDecionWithIndex,
		})
	case "TagDecision":
		c.HTML(http.StatusOK, "tagDecision.html", gin.H{
			"nodeID": nodeID,
		})
	case "TagOperation":
		c.HTML(http.StatusOK, "tagOperation.html", gin.H{
			"nodeID": nodeID,
		})
	case "Random":
		c.HTML(http.StatusOK, "random.html", gin.H{
			"nodeID": nodeID,
		})
	case "FirstStep":
		c.HTML(http.StatusOK, "firstStep.html", nil)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "Unsupported node type"})
	}
}

func UpdateLocationHandler(c *gin.Context, db *gorm.DB) {
	var req database.NodeUpdateLocationRequest
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input data"})
		return
	}
	var node database.Node
	if err := db.Where("id = ?", req.CurrentNodeID).First(&node).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Node not found"})
		return
	}

	node.LocX = int(req.LocX)
	node.LocY = int(req.LocY)

	if err := db.Save(&node).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update node location"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Node location updated successfully"})
}
