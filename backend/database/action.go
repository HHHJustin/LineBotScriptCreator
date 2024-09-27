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

type NodeUpdateTitleRequest struct {
	CurrentNodeID int    `json:"currentNodeID,omitempty"`
	NewTitle      string `json:"newTitle,omitempty"`
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

type MessageCreateRequest struct {
	CurrentNodeID  int    `json:"currentNodeID,omitempty"`
	MessageType    string `json:"messageType" binding:"required"`
	MessageContent string `json:"messageContent" binding:"required"`
}

type MessageDeleteRequest struct {
	CurrentNodeID int `json:"currentNodeID,omitempty"`
	MessageID     int `json:"messageID,omitempty"`
}

type MessageUpdateRequest struct {
	MessageID      int    `json:"messageID"`
	MessageContent string `json:"messageContent"`
}
