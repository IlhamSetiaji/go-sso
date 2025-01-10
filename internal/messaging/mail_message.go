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

type IMailMessage interface {
	SendMail(req *request.MailRequest) (string, error)
}

type MailMessage struct {
	Log *logrus.Logger
}

func NewMailMessage(log *logrus.Logger) IMailMessage {
	return &MailMessage{
		Log: log,
	}
}

func MailMessageFactory(log *logrus.Logger) IMailMessage {
	return NewMailMessage(log)
}

func (m *MailMessage) SendMail(req *request.MailRequest) (string, error) {
	payload := map[string]interface{}{
		"to":      req.To,
		"subject": req.Subject,
		"body":    req.Body,
		"from":    req.From,
		"email":   req.Email,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "send_mail",
		MessageData: payload,
		ReplyTo:     "gift-redeem-be",
	}

	log.Printf("INFO: document message: %v", docMsg)

	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "gift-redeem-be",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	// wait for reply
	resp, err := utils.WaitForReply(docMsg.ID, rchan)
	if err != nil {
		return "", err
	}

	log.Printf("INFO: response from mail message: %v", resp)

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return "", errors.New("[SendFindOrganizationByIDMessage] " + errMsg)
	}

	return "Success	", nil
}
