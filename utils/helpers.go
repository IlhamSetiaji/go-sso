package utils

import (
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/http/response"
	mqResponse "app/go-sso/internal/http/response"
	"errors"
	"log"
	"time"

	"github.com/gin-gonic/gin"
)

type TemplateHelper struct {
	Ctx *gin.Context
}

func NewTemplateHelper(c *gin.Context) *TemplateHelper {
	return &TemplateHelper{
		Ctx: c,
	}
}

func (h *TemplateHelper) HasPermission(requiredPermission string) bool {
	permissions, err := GetUserPermissions(h.Ctx)
	if err != nil || permissions == nil {
		return false
	}

	for _, permission := range permissions {
		if permission.Name == requiredPermission {
			return true
		}
	}
	return false
}

func (h *TemplateHelper) HasRole(requiredRole string) bool {
	roles, err := GetUserRoles(h.Ctx)
	if err != nil || roles == nil {
		return false
	}

	for _, role := range roles {
		if role.Name == requiredRole {
			return true
		}
	}
	return false
}

var Rchans = make(map[string](chan response.RabbitMQResponse))
var Pchan = make(chan RabbitMsg, 10)

type RabbitMsg struct {
	QueueName string                  `json:"queueName"`
	Message   request.RabbitMQRequest `json:"message"`
}

func WaitReply(uid string, rchan chan mqResponse.RabbitMQResponse, ctx *gin.Context) {
	for {
		select {
		case docReply := <-rchan:
			// responses received
			log.Printf("INFO: received reply: %v uid: %s", docReply, uid)

			// send response back to client
			SuccessResponse(ctx, 200, "Success", docReply)
			// remove channel from rchans
			delete(Rchans, uid)
			return
		case <-time.After(10 * time.Second):
			// timeout
			log.Printf("ERROR: request timeout uid: %s", uid)

			// send response back to client
			ErrorResponse(ctx, 408, "Request Timeout", "Request Timeout")
			// remove channel from rchans
			delete(Rchans, uid)
			return
		}
	}
}

func WaitForReply(id string, rchan chan response.RabbitMQResponse) (response.RabbitMQResponse, error) {
	for {
		select {
		case docReply := <-rchan:
			// responses received
			log.Printf("INFO: received reply: %v uid: %s", docReply, id)

			delete(Rchans, id)
			return docReply, nil
		case <-time.After(10 * time.Second):
			// timeout
			log.Printf("ERROR: request timeout uid: %s", id)

			// remove channel from rchans
			delete(Rchans, id)
			return response.RabbitMQResponse{}, errors.New("request timeout")
		}
	}
}
