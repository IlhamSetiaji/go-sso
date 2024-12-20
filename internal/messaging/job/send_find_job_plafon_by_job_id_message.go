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

type IFindJobPlafonByJobIDMessageRequest struct {
	JobID string `json:"job_id"`
}

type IFindJobPlafonByJobIDMessageResponse struct {
	ID     string `json:"id"`
	Plafon int    `json:"plafon"`
}

type IFindJobPlafonByJobIDMessage interface {
	Execute(req *IFindJobPlafonByJobIDMessageRequest) (*IFindJobPlafonByJobIDMessageResponse, error)
}

type FindJobPlafonByJobIDMessage struct {
	Log *logrus.Logger
}

func NewFindJobPlafonByJobIDMessage(log *logrus.Logger) IFindJobPlafonByJobIDMessage {
	return &FindJobPlafonByJobIDMessage{
		Log: log,
	}
}

func (m *FindJobPlafonByJobIDMessage) Execute(req *IFindJobPlafonByJobIDMessageRequest) (*IFindJobPlafonByJobIDMessageResponse, error) {
	payload := map[string]interface{}{
		"job_id": req.JobID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "find_job_plafon_by_job_id",
		MessageData: payload,
		ReplyTo:     "julong_sso",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_manpower",
		Message:   *docMsg,
	}
	utils.Pchan <- msg

	log.Printf("Message published")

	// wait for reply
	resp, err := utils.WaitForReply(docMsg.ID, rchan)
	if err != nil {
		return nil, err
	}

	log.Printf("INFO: response: %v", resp)

	if errMsg, ok := resp.MessageData["error"].(string); ok && errMsg != "" {
		return nil, errors.New("[SendFindJobLevelByJobIDMessage] " + errMsg)
	}

	return &IFindJobPlafonByJobIDMessageResponse{
		ID:     resp.MessageData["id"].(string),
		Plafon: int(resp.MessageData["plafon"].(float64)),
	}, nil
}

func FindJobPlafonByJobIDMessageFactory(log *logrus.Logger) IFindJobPlafonByJobIDMessage {
	return NewFindJobPlafonByJobIDMessage(log)
}
