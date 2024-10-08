package api

import (
	"LineBotCreator/database"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/line/line-bot-sdk-go/v8/linebot"
	"gorm.io/gorm"
)

func CallbackHandler(c *gin.Context, bot *linebot.Client, db *gorm.DB) {
	events, err := bot.ParseRequest(c.Request)
	if err != nil {
		if err == linebot.ErrInvalidSignature {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid signature"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse request"})
		}
		return
	}
	for _, event := range events {
		switch event.Type {
		case linebot.EventTypeFollow:
			nextNode, err := addFriendHandler(event, db)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			gotoNextNode(nextNode, db, bot, event)
		case linebot.EventTypeJoin:
		case linebot.EventTypeMessage:
			switch message := event.Message.(type) {
			case *linebot.TextMessage:
				nextNode, err := checkMessageCondition(event, db, message)
				if err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
					return
				}
				gotoNextNode(nextNode, db, bot, event)

			default:
				log.Printf("不支援的訊息類型: %T", message)
			}
		default:
			log.Printf("不支援的事件類型: %s", event.Type)
		}
	}
	c.Status(http.StatusOK)
}

func addFriendHandler(event *linebot.Event, db *gorm.DB) (int, error) {
	var node database.Node
	if err := db.Where("type = ? AND title = ?", "FirstStep", "AddFriend").First(&node).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch information: %v", err)
	}

	user := database.UserSession{
		UserID:    event.Source.UserID,
		CurrentID: node.NextNode,
		Time:      time.Now(),
	}
	if err := db.Create(&user).Error; err != nil {
		return 0, fmt.Errorf("failed to create user: %v", err)
	}

	return node.NextNode, nil
}

func gotoNextNode(nextNode int, db *gorm.DB, bot *linebot.Client, event *linebot.Event) error {
	var user database.UserSession
	userID := event.Source.UserID
	replyToken := event.ReplyToken
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		return fmt.Errorf("failed to fetch user session: %v", err)
	}
	for nextNode != 0 {
		var next database.Node
		if err := db.Where("id = ?", nextNode).First(&next).Error; err != nil {
			return fmt.Errorf("failed to fetch next node: %v", err)
		}

		user.CurrentID = nextNode
		switch next.Type {
		case "Message":
			var messages []database.Message
			if err := db.Where("node_id = ?", next.ID).Find(&messages).Error; err != nil {
				return fmt.Errorf("failed to fetch messages for node: %v", err)
			}
			var replyMessages []linebot.SendingMessage
			for _, message := range messages {
				replyMessages = append(replyMessages, linebot.NewTextMessage(message.Content))
			}
			if _, err := bot.ReplyMessage(replyToken, replyMessages...).Do(); err != nil {
				return fmt.Errorf("failed to reply message to user: %v", err)
			}
			nextNode = next.NextNode

		case "QuickReply":
			user.CurrentID = next.NextNode
			var quickReplies []database.QuickReply
			if err := db.Where("node_id = ?", next.ID).Find(&quickReplies).Error; err != nil {
				return fmt.Errorf("failed to fetch quick replies for node: %v", err)
			}

			var quickReplyItems []*linebot.QuickReplyButton
			for _, reply := range quickReplies {
				quickReplyItems = append(quickReplyItems, linebot.NewQuickReplyButton(
					"",
					linebot.NewMessageAction(reply.ButtonName, reply.Reply),
				))
			}

			quickReply := linebot.NewQuickReplyItems(quickReplyItems...)
			message := linebot.NewTextMessage(next.Title).WithQuickReplies(quickReply)
			if _, err := bot.ReplyMessage(replyToken, message).Do(); err != nil {
				return fmt.Errorf("failed to send quick reply to user: %v", err)
			}
			nextNode = 0

		case "KeywordDecision":
			nextNode = 0

		default:
			return fmt.Errorf("unsupported node type: %s", next.Type)
		}
	}
	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user session: %v", err)
	}
	return nil
}

func checkMessageCondition(event *linebot.Event, db *gorm.DB, message linebot.Message) (int, error) {
	var user database.UserSession
	if err := db.Where("user_id = ?", event.Source.UserID).First(&user).Error; err != nil {
		var node database.Node
		if err := db.Where("type = ? AND title = ?", "FirstStep", "AddFriend").First(&node).Error; err != nil {
			return 0, fmt.Errorf("failed to fetch information: %v", err)
		}
		user := database.UserSession{
			UserID:    event.Source.UserID,
			CurrentID: node.NextNode,
			Time:      time.Now(),
		}
		if err := db.Create(&user).Error; err != nil {
			return 0, fmt.Errorf("failed to create user: %v", err)
		}
	}

	var node database.Node
	if err := db.Where("id = ?", user.CurrentID).First(&node).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch information: %v", err)
	}

	switch node.Type {
	case "Message":
		return node.NextNode, nil

	case "KeywordDecision":
		var keywordDecisions []database.KeywordDecision
		if err := db.Where("node_id = ?", node.ID).Find(&keywordDecisions).Error; err != nil {
			return node.ID, fmt.Errorf("not set any keyword decision: %v", err)
		}
		textMessage, ok := message.(*linebot.TextMessage)
		if !ok {
			return node.ID, fmt.Errorf("message type is not text")
		}
		for _, decision := range keywordDecisions {
			if textMessage.Text == decision.Keyword {
				return decision.NextNode, nil
			}
		}
		return node.ID, nil

	default:
		return 0, fmt.Errorf("unsupported node type: %s", node.Type)
	}
}
