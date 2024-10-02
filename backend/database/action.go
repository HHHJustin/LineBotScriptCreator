package database

/* Node */
type NodeBaseRequest struct {
	CurrentNodeID int `json:"currentNodeID,omitempty"`
}

type NodeCreateRequest struct {
	NodeBaseRequest
	NewNodeType string `json:"newNodeType" binding:"required"`
}

type NodeDeleteRequest struct {
	NodeBaseRequest
}

type NodeUpdateTitleRequest struct {
	NodeBaseRequest
	NewTitle string `json:"newTitle,omitempty"`
}

type NodeUpdateLocationRequest struct {
	NodeBaseRequest
	LocX float64 `json:"locX"`
	LocY float64 `json:"locY"`
}

type FirstStepRequest struct {
	FirstStepType string `json:"firstStepType,omitempty"`
}

type EditPageRequest struct {
	NodeBaseRequest
}

type LinkCreateRequest struct {
	FromNodeID int `json:"fromNodeID,omitempty"`
	ToNodeID   int `json:"toNodeID,omitempty"`
}

type LinkDeleteRequest struct {
	FromNodeID int `json:"fromNodeID,omitempty"`
	ToNodeID   int `json:"toNodeID,omitempty"`
}

/* Message */
type MessageBaseRequest struct {
	MessageID int `json:"messageID,omitempty"`
}

type MessageCreateRequest struct {
	NodeBaseRequest
	MessageBaseRequest
	MessageType    string `json:"messageType" binding:"required"`
	MessageContent string `json:"messageContent" binding:"required"`
}

type MessageDeleteRequest struct {
	NodeBaseRequest
	MessageBaseRequest
}

type MessageUpdateRequest struct {
	MessageBaseRequest
	MessageContent string `json:"messageContent"`
}

/* Keyword Decision */
type KeywordDecisionBaseRequest struct {
	KeywordDecisionID int `json:"keywordDecisionID,omitempty"`
}

type KeywordDecisionDeleteRequest struct {
	NodeBaseRequest
	KeywordDecisionBaseRequest
}

type KeywordDecisionUpdateRequest struct {
	KeywordDecisionBaseRequest
	Keyword string `json:"Keyword"`
}
