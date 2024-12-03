package messaging

import (
	"app/go-sso/internal/repository"
	"encoding/json"
	"log"

	"github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

type MessageConsumer struct {
	JobRepository repository.IJobRepository
}

type IMessageConsumer interface {
	ListenForJobRequest()
}

func NewMessageConsumer(jobRepository repository.IJobRepository) IMessageConsumer {
	return &MessageConsumer{
		JobRepository: jobRepository,
	}
}

func (mc *MessageConsumer) ListenForJobRequest() {
	rabbitmq.ConsumeMessages("job.validation.request", func(d amqp091.Delivery) {
		log.Printf("Received message: %s", d.Body)

		var request map[string]interface{}
		err := json.Unmarshal(d.Body, &request)
		if err != nil {
			log.Printf("Failed to unmarshal request: %v", err)
			return
		}

		jobID, ok := request["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'user_id'")
			return
		}

		exists, err := mc.JobRepository.FindById(uuid.MustParse(jobID))

		if err != nil {
			log.Printf("Failed to validate job: %v", err)
			return
		}

		response := map[string]interface{}{
			"job_id": jobID,
			"exists": exists,
		}
		responseBody, _ := json.Marshal(response)
		rabbitmq.PublishMessage("", "user.validation.response", responseBody)
	})
}
