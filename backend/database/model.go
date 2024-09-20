package database

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"
)

type IntArray []int

func (ia IntArray) Value() (driver.Value, error) {
	return json.Marshal(ia)
}

func (ia *IntArray) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}

	return json.Unmarshal(bytes, ia)
}

type Node struct {
	ID           int           `gorm:"primaryKey;autoIncrement"`
	Title        string        `gorm:"size:255;not null"`
	Type         string        `gorm:"size:255;not null"`
	Range        IntArray      `gorm:"type:jsonb;not null"`
	PreviousNode int           `gorm:"index"`
	NextNode     int           `gorm:"index"`
	Loc          string        `gorm:"size:50;default:'0 0'" json:"loc"`
	Messages     []Message     `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	QuickReplies []QuickReply  `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	KeyDecisions []KeyDecision `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	TagDecisions []TagDecision `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Randoms      []Random      `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Link struct {
	From int `json:"from" gorm:"index"`
	To   int `json:"to" gorm:"index"`
}

type Graph struct {
	Nodes []GraphNode `json:"nodes"`
	Links []Link      `json:"links"`
}

type GraphNode struct {
	Key   int    `json:"key"`
	Text  string `json:"text"`
	Color string `json:"color"`
	Loc   string `json:"loc"`
}

type Message struct {
	MessageID int    `gorm:"primaryKey;autoIncrement"`
	Type      string `gorm:"size:255;not null"`
	Content   string `gorm:"type:text;not null"`
	NodeID    int    `gorm:"not null;index"`
	Node      Node   `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type QuickReply struct {
	QuickReplyID int    `gorm:"primaryKey;autoIncrement"`
	ButtonName   string `gorm:"size:255;not null"`
	Reply        string `gorm:"type:text;not null"`
	NodeID       int    `gorm:"not null;index"`
	Node         Node   `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type KeyDecision struct {
	KWDecisionID int    `gorm:"primaryKey;autoIncrement"`
	Keyword      string `gorm:"size:255;not null"`
	NextNode     int    `gorm:"not null;index"`
	NodeID       int    `gorm:"not null;index"`
	Node         Node   `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type TagDecision struct {
	TagDecisionID int      `gorm:"primaryKey;autoIncrement"`
	Tags          IntArray `gorm:"type:json;not null;default:NULL;"`
	NextNode      int      `gorm:"not null;index"`
	NodeID        int      `gorm:"not null;index"`
	Node          Node     `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Random struct {
	RandomID  int    `gorm:"primaryKey;autoIncrement"`
	Weight    int    `gorm:"not null"`
	Condition string `gorm:"size:255;not null"`
	NextNode  int    `gorm:"not null;index"`
	NodeID    int    `gorm:"not null;index"`
	Node      Node   `gorm:"foreignKey:NodeID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}

type Tag struct {
	TagID int    `gorm:"primaryKey;autoIncrement"`
	Name  string `gorm:"size:255;not null"`
}

type UserSession struct {
	Index     int    `gorm:"primaryKey;autoIncrement"`
	UserID    string `gorm:"size:255;not null"`
	CurrentID int    `gorm:"not null"`
	Time      time.Time
}

type FirstStep struct {
	FirstStepID int    `gorm:"primaryKey;autoIncrement"`
	Type        string `gorm:"size:255;not null"`
	NextNode    int    `gorm:"not null"`
}
