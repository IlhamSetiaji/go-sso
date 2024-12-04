package main

import (
	"app/go-sso/internal/config"
	messaging "app/go-sso/internal/messaging/job"
	"encoding/json"
	"os"
	"os/signal"
	"syscall"

	"github.com/IlhamSetiaji/go-rabbitmq-utils/rabbitmq"
	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
)

func main() {
	// Initialize configurations and logger
	viperConfig := config.NewViper()
	log := config.NewLogrus(viperConfig)

	// Initialize RabbitMQ connection
	err := rabbitmq.InitializeConnection(viperConfig.GetString("rabbitmq.url"))
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ: %v", err)
	}
	defer rabbitmq.CloseConnection()

	// Declare the queue before consuming messages
	channel := rabbitmq.GetChannel()
	defer channel.Close()

	queueName := "julong_queue_request"
	_, err = channel.QueueDeclare(
		queueName, // name
		true,      // durable
		false,     // delete when unused
		false,     // exclusive
		false,     // no-wait
		nil,       // arguments
	)
	if err != nil {
		log.Fatalf("Failed to declare queue: %v", err)
	}

	// Start consuming messages
	err = rabbitmq.ConsumeMessages(queueName, func(d amqp091.Delivery) {
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
				log.Printf("Invalid request format: missing 'job_id'")
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
			err = rabbitmq.PublishMessage("", "julong_queue_response", responseBody)
			if err != nil {
				log.Printf("Failed to publish response: %v", err)
			}
		default:
			log.Printf("Unknown message type: %s", request["message_type"])
		}
	})

	if err != nil {
		log.Fatalf("Failed to start consuming messages: %v", err)
	}

	// Wait for shutdown signal to gracefully exit
	waitForShutdown(log)
}

func waitForShutdown(log *logrus.Logger) {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	log.Println("Waiting for shutdown signal...")
	<-quit
	log.Println("Shutting down gracefully...")
}
