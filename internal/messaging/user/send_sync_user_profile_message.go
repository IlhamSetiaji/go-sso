package messaging

import (
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/http/response"
	"app/go-sso/utils"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type ISendSyncUserProfileMessageRequest struct {
	UserID        string `json:"user_id"`
	Name          string `json:"name"`
	Age           int    `json:"age"`
	BirthDate     string `json:"birth_date"`
	BirthPlace    string `json:"birth_place"`
	MaritalStatus string `json:"marital_status"`
	PhoneNumber   string `json:"phone_number"`
}

type ISendSyncUserProfileMessageResponse struct {
	Message string `json:"message"`
}

type ISendSyncUserProfileMessage interface {
	Execute(req *ISendSyncUserProfileMessageRequest) (*ISendSyncUserProfileMessageResponse, error)
}

type SendSyncUserProfileMessage struct {
	Log *logrus.Logger
}

func NewSendSyncUserProfileMessage(log *logrus.Logger) ISendSyncUserProfileMessage {
	return &SendSyncUserProfileMessage{
		Log: log,
	}
}

func (m *SendSyncUserProfileMessage) Execute(req *ISendSyncUserProfileMessageRequest) (*ISendSyncUserProfileMessageResponse, error) {
	payload := map[string]interface{}{
		"user_id":        req.UserID,
		"name":           req.Name,
		"age":            req.Age,
		"birth_date":     req.BirthDate,
		"birth_place":    req.BirthPlace,
		"marital_status": req.MaritalStatus,
		"phone_number":   req.PhoneNumber,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "sync_user_profile",
		MessageData: payload,
		ReplyTo:     "julong_sso",
	}

	m.Log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_recruitment",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	log.Printf("INFO: message sent to queue: %v", msg)

	// wait for reply
	resp, err := utils.WaitForReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendSyncUserProfileMessage] " + errMsg)
	}

	log.Printf("INFO: response: %v", resp)

	return &ISendSyncUserProfileMessageResponse{
		Message: "Successfully synced user profile",
	}, nil
}

func SendSyncUserProfileMessageFactory(log *logrus.Logger) ISendSyncUserProfileMessage {
	return NewSendSyncUserProfileMessage(log)
}
