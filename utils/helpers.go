package utils

import (
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/http/response"
	mqResponse "app/go-sso/internal/http/response"
	"crypto/rand"
	"errors"
	"log"
	"math/big"
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

func (h *TemplateHelper) IsAuthenticated() bool {
	return h.Ctx.GetBool("isAuthenticated")
}

func (h *TemplateHelper) NotInArrays(value string, list []string) bool {
	for _, item := range list {
		if value == item {
			return false
		}
	}
	return true
}

func (h *TemplateHelper) CreateSlice(values ...string) []string {
	return values
}

func (h *TemplateHelper) DateFormatter(date time.Time) string {
	return date.Format("2006-01-02")
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

func GenerateRandomIntToken(digits int) (int64, error) {
	max := new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil).Sub(new(big.Int).Exp(big.NewInt(10), big.NewInt(int64(digits)), nil), big.NewInt(1))
	n, err := rand.Int(rand.Reader, max)
	if err != nil {
		return 0, err
	}
	return n.Int64(), nil
}

func GenerateRandomStringToken(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		n, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		b[i] = charset[n.Int64()]
	}
	return string(b)
}
