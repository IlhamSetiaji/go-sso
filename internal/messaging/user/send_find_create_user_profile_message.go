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

type ISendFindCreateUserProfileMessageRequest struct {
	UserID uuid.UUID `json:"user_id"`
	Name   string    `json:"name"`
}

type ISendFindCreateUserProfileMessageResponse struct {
	Message string `json:"message"`
}

type ISendFindCreateUserProfileMessage interface {
	Execute(req *ISendFindCreateUserProfileMessageRequest) (*ISendFindCreateUserProfileMessageResponse, error)
}

type SendFindCreateUserProfileMessage struct {
	Log *logrus.Logger
}

func NewSendFindCreateUserProfileMessage(log *logrus.Logger) ISendFindCreateUserProfileMessage {
	return &SendFindCreateUserProfileMessage{
		Log: log,
	}
}

func (m *SendFindCreateUserProfileMessage) Execute(req *ISendFindCreateUserProfileMessageRequest) (*ISendFindCreateUserProfileMessageResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.UserID,
		"name":    req.Name,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "create_user_profile",
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

	log.Printf("Message published")

	// wait for reply
	resp, err := utils.WaitForReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindUserProfileMessage] " + errMsg)
	}

	log.Printf("INFO: response: %v", resp)

	return &ISendFindCreateUserProfileMessageResponse{
		Message: "success",
	}, nil
}

func SendFindCreateUserProfileMessageFactory(log *logrus.Logger) ISendFindCreateUserProfileMessage {
	return NewSendFindCreateUserProfileMessage(log)
}
