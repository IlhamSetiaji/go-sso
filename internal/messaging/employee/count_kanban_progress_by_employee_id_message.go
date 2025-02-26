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

type ICountKanbanProgressByEmployeeIDMessageRequest struct {
	EmployeeID string `json:"employee_id"`
}

type ICountKanbanProgressByEmployeeIDMessageResponse struct {
	EmployeeID string `json:"employee_id"`
	TotalTask  int    `json:"total_task"`
	ToDo       int    `json:"to_do"`
	InProgress int    `json:"in_progress"`
	NeedReview int    `json:"need_review"`
	Completed  int    `json:"completed"`
}

type ICountKanbanProgressByEmployeeIDMessage interface {
	Execute(req *ICountKanbanProgressByEmployeeIDMessageRequest) (*ICountKanbanProgressByEmployeeIDMessageResponse, error)
}

type CountKanbanProgressByEmployeeIDMessage struct {
	Log *logrus.Logger
}

func NewCountKanbanProgressByEmployeeIDMessage(log *logrus.Logger) ICountKanbanProgressByEmployeeIDMessage {
	return &CountKanbanProgressByEmployeeIDMessage{
		Log: log,
	}
}

func (m *CountKanbanProgressByEmployeeIDMessage) Execute(req *ICountKanbanProgressByEmployeeIDMessageRequest) (*ICountKanbanProgressByEmployeeIDMessageResponse, error) {
	payload := map[string]interface{}{
		"employee_id": req.EmployeeID,
	}

	docMsg := &request.RabbitMQRequest{
		ID:          uuid.New().String(),
		MessageType: "count_kanban_progress_by_employee_id",
		MessageData: payload,
		ReplyTo:     "julong_sso",
	}

	log.Printf("INFO: document message: %v", docMsg)

	// create channel and add to rchans with uid
	rchan := make(chan response.RabbitMQResponse)
	utils.Rchans[docMsg.ID] = rchan

	// publish rabbit message
	msg := utils.RabbitMsg{
		QueueName: "julong_onboarding",
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
		return nil, errors.New("[CountKanbanProgressByEmployeeIDMessage] " + errMsg)
	}

	log.Printf("INFO: response: %v", resp)

	totalTask, ok := resp.MessageData["total_task"].(float64)
	if !ok {
		return nil, errors.New("total_task is not a number")
	}
	toDo, ok := resp.MessageData["to_do"].(float64)
	if !ok {
		return nil, errors.New("to_do is not a number")
	}
	inProgress, ok := resp.MessageData["in_progress"].(float64)
	if !ok {
		return nil, errors.New("in_progress is not a number")
	}
	needReview, ok := resp.MessageData["need_review"].(float64)
	if !ok {
		return nil, errors.New("need_review is not a number")
	}
	completed, ok := resp.MessageData["completed"].(float64)
	if !ok {
		return nil, errors.New("completed is not a number")
	}

	return &ICountKanbanProgressByEmployeeIDMessageResponse{
		EmployeeID: req.EmployeeID,
		TotalTask:  int(totalTask),
		ToDo:       int(toDo),
		InProgress: int(inProgress),
		NeedReview: int(needReview),
		Completed:  int(completed),
	}, nil
}

func CountKanbanProgressByEmployeeIDMessageFactory(log *logrus.Logger) ICountKanbanProgressByEmployeeIDMessage {
	return NewCountKanbanProgressByEmployeeIDMessage(log)
}
