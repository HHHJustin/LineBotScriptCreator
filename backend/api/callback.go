package api

import (
	"LineBotCreator/database"
	"fmt"
	"time"

	"github.com/line/line-bot-sdk-go/v8/linebot"
	"gorm.io/gorm"
)

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

func gotoNextNode(nextNode int, db *gorm.DB, bot *linebot.Client, userID string) error {
	var user database.UserSession
	if err := db.Where("user_id = ?", userID).First(&user).Error; err != nil {
		fmt.Printf("Error fetching database.UserSession node: %v\n", err)
		return fmt.Errorf("failed to fetch user session: %v", err)
	}

	var next database.Node
	if err := db.Where("id = ?", nextNode).First(&next).Error; err != nil {
		fmt.Printf("Error fetching database.Node node: %v\n", err)
		return fmt.Errorf("failed to fetch next node: %v", err)
	}
	user.CurrentID = nextNode

	switch next.Type {
	case "Message":
		var messages []database.Message
		if err := db.Where("node_id = ?", next.ID).Find(&messages).Error; err != nil {
			return fmt.Errorf("failed to fetch messages for node: %v", err)
		}

		for _, message := range messages {
			if _, err := bot.PushMessage(userID, linebot.NewTextMessage(message.Content)).Do(); err != nil {
				return fmt.Errorf("failed to send message to user: %v", err)
			}
		}
		if next.NextNode != 0 {
			return gotoNextNode(next.NextNode, db, bot, userID)
		}

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
		if _, err := bot.PushMessage(userID, message).Do(); err != nil {
			return fmt.Errorf("failed to send quick reply to user: %v", err)
		}

	case "KeywordDecision":
		return nil

	default:
		return fmt.Errorf("unsupported node type: %s", next.Type)
	}

	if err := db.Save(&user).Error; err != nil {
		return fmt.Errorf("failed to update user session: %v", err)
	}
	return nil
}

func checkMessageCondition(event *linebot.Event, db *gorm.DB, message linebot.Message) (int, error) {
	var user database.UserSession
	if err := db.Where("user_id = ?", event.Source.UserID).First(&user).Error; err != nil {
		return 0, fmt.Errorf("failed to fetch information: %v", err)
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
		return 0, fmt.Errorf("no matching keyword found")

	default:
		return 0, fmt.Errorf("unsupported node type: %s", node.Type)
	}
}
