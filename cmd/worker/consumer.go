package main

import (
	"app/go-sso/internal/config"
	messaging "app/go-sso/internal/messaging/job"
	"encoding/json"

	"github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
)

func main() {
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)
	err := rabbitmq.InitializeConnection(viperConfig.GetString("rabbitmq.url"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.CloseConnection()

	rabbitmq.ConsumeMessages("job.validation.request", func(d amqp091.Delivery) {
		log.Printf("Received message: %s", d.Body)

		var request map[string]interface{}
		err := json.Unmarshal(d.Body, &request)
		if err != nil {
			log.Printf("Failed to unmarshal request: %v", err)
			return
		}

		switch request["message_type"] {
		case "check_job_exists":
			jobID, ok := request["job_id"].(string)
			if !ok {
				log.Printf("Invalid request format: missing 'user_id'")
				return
			}

			if err != nil {
				log.Printf("Failed to validate job: %v", err)
				return
			}
			messageFactory := messaging.CheckJobExistMessageFactory(log)
			message, err := messageFactory.Execute(messaging.ICheckJobExistMessageRequest{
				JobID: uuid.MustParse(jobID),
			})

			if err != nil {
				log.Printf("Failed to execute message: %v", err)
				return
			}

			response := map[string]interface{}{
				"job_id": jobID,
				"exists": message.Exists,
			}
			responseBody, _ := json.Marshal(response)
			rabbitmq.PublishMessage("", "user.validation.response", responseBody)
		default:
			log.Printf("Unknown message type: %s", request["message_type"])
		}
	})
}
