package service

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/nlsra/message-system/pkg/model"
	"github.com/nlsra/message-system/pkg/repository/message"
	"log"
	"net/http"
	"os"
	"sync"
	"time"
)

type MessageService struct {
	MessageRepository *message.Repository
	wg                sync.WaitGroup
	stopSendingChan   chan struct{}
}

func NewMessageService(p *message.Repository) MessageService {
	return MessageService{
		MessageRepository: p,
		stopSendingChan:   make(chan struct{}),
	}
}

func (msg *MessageService) FindAllSentMessages() ([]model.Message, error) {
	return msg.MessageRepository.FindAllSentMessages()
}

func (msg *MessageService) Start() bool {
	select {
	case <-msg.stopSendingChan:
		msg.stopSendingChan = make(chan struct{})
	default:
		return false
	}

	msg.StartSending()
	return true
}

func (msg *MessageService) Stop() {
	close(msg.stopSendingChan)
	msg.wg.Wait()
}

func (msg *MessageService) Migrate() error {
	return msg.MessageRepository.Migrate()
}

func (msg *MessageService) StartSending() {
	msg.wg.Add(1)
	go func() {
		defer msg.wg.Done()
		for {
			select {
			case <-msg.stopSendingChan:
				return
			default:
				msg.SendMessagesToWebhook(2)
				time.Sleep(2 * time.Minute)
			}
		}
	}()
}

func (msg *MessageService) SendMessagesToWebhook(count int) {
	messages, err := msg.getPendingMessages(count)
	if err != nil {
		log.Printf("Error getting pending messages: %v", err)
		return
	}

	webhookURL := os.Getenv("EXTERNAL_API_URL")

	for _, message := range messages {
		if err := msg.sendMessage(&message, webhookURL); err != nil {
			log.Printf("Error processing message: %v", err)
		}
	}
}

func (msg *MessageService) getPendingMessages(count int) ([]model.Message, error) {
	var messages []model.Message

	if count <= 0 {
		return messages, msg.MessageRepository.DB.Model(&model.Message{}).Where("status = ?", "pending").Find(&messages).Error
	}

	return messages, msg.MessageRepository.DB.Model(&model.Message{}).Where("status = ?", "pending").Limit(count).Find(&messages).Error
}

func (msg *MessageService) sendMessage(message *model.Message, webhookURL string) error {
	payload, err := json.Marshal(map[string]string{
		"to":      message.To,
		"content": message.Content,
	})
	if err != nil {
		return err
	}

	resp, err := http.Post(webhookURL, "application/json", bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	return msg.handleResponse(message, resp)
}

func (msg *MessageService) handleResponse(message *model.Message, resp *http.Response) error {
	var body map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&body); err != nil {
		return err
	}

	responseMessage, ok := body["message"].(string)
	if !ok {
		return errors.New("response message not found")
	}

	if resp.StatusCode == http.StatusAccepted && responseMessage == "Accepted" {
		if err := msg.updateMessageStatus(message, body); err != nil {
			return err
		}
	} else {
		if err := msg.updateMessageStatus(message, nil); err != nil {
			return err
		}
	}

	log.Printf("Webhook Response Message: %s", responseMessage)
	return nil
}

func (msg *MessageService) updateMessageStatus(message *model.Message, body map[string]interface{}) error {
	var status string
	if body != nil {
		status = "success"
		if messageId, ok := body["messageId"].(string); ok {
			sendTime := time.Now().Format(time.RFC3339)
			redisKey := "messages:" + messageId
			if err := msg.MessageRepository.RedisClient.Set(redisKey, sendTime, 0).Err(); err != nil {
				log.Printf("Redis Error: %v", err)
			}
		}
	} else {
		status = "failure"
	}

	return msg.MessageRepository.UpdateMessageStatus(message, status)
}
