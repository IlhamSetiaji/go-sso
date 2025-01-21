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

type ISendFindUserProfileMessageRequest struct {
	UserID string `json:"user_id"`
}

type ISendFindUserProfileMessageResponse struct {
	UserProfile map[string]interface{} `json:"user_profile"`
}

type ISendFindUserProfileMessage interface {
	Execute(req *ISendFindUserProfileMessageRequest) (*ISendFindUserProfileMessageResponse, error)
}

type SendFindUserProfileMessage struct {
	Log *logrus.Logger
}

func NewSendFindUserProfileMessage(log *logrus.Logger) ISendFindUserProfileMessage {
	return &SendFindUserProfileMessage{
		Log: log,
	}
}

func (m *SendFindUserProfileMessage) Execute(req *ISendFindUserProfileMessageRequest) (*ISendFindUserProfileMessageResponse, error) {
	payload := map[string]interface{}{
		"user_id": req.UserID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_user_profile_by_user_id",
		MessageData: payload,
		ReplyTo:     "julong_sso",
	}

	log.Printf("INFO: document message: %v", docMsg)

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

	return &ISendFindUserProfileMessageResponse{
		UserProfile: resp.MessageData,
	}, nil
}

func SendFindUserProfileMessageFactory(log *logrus.Logger) ISendFindUserProfileMessage {
	return NewSendFindUserProfileMessage(log)
}
