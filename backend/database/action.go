package database

type MessageWithIndex struct {
	Message
	Index int
}

type NodeCreateRequest struct {
	CurrentNodeID int    `json:"currentNodeID,omitempty"`
	NewNodeType   string `json:"newNodeType" binding:"required"`
}

type NodeDeleteRequest struct {
	CurrentNodeID int `json:"currentNodeID,omitempty"`
}

type EditPageRequest struct {
	CurrentNodeID int `json:"currentNodeID,omitempty"`
}

type LinkCreateRequest struct {
	FromNodeID int `json:"fromNodeID,omitempty"`
	ToNodeID   int `json:"toNodeID,omitempty"`
}

type LinkDeleteRequest struct {
	FromNodeID int `json:"fromNodeID,omitempty"`
	ToNodeID   int `json:"toNodeID,omitempty"`
}
