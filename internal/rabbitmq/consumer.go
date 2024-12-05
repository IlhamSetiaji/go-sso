package rabbitmq

import (
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/http/response"
	messaging "app/go-sso/internal/messaging/job"
	"encoding/json"
	"os"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConsumer(viper *viper.Viper, log *logrus.Logger) {
	// conn
	conn, err := amqp091.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		log.Printf("ERROR: fail init consumer: %s", err.Error())
		os.Exit(1)
	}

	log.Printf("INFO: done init consumer conn")

	// create channel
	amqpChannel, err := conn.Channel()
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// create queue
	queue, err := amqpChannel.QueueDeclare(
		"julong_sso", // channelname
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		log.Printf("ERROR: fail create queue: %s", err.Error())
		os.Exit(1)
	}

	// channel
	msgChannel, err := amqpChannel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		log.Printf("ERROR: fail create channel: %s", err.Error())
		os.Exit(1)
	}

	// consume
	for {
		select {
		case msg := <-msgChannel:
			// unmarshal
			docMsg := &request.RabbitMQRequest{}
			err = json.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received msg: %v", docMsg)

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: fail to ack: %s", err.Error())
			}

			// handle docMsg
			handleMsg(docMsg, log)
		}
	}
}

func handleMsg(docMsg *request.RabbitMQRequest, log *logrus.Logger) {
	// switch case
	var msgData map[string]interface{}
	switch docMsg.MessageType {
	case "check_job_exist":
		jobID, ok := docMsg.MessageData["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_id'")
		}

		messageFactory := messaging.CheckJobExistMessageFactory(log)
		message, err := messageFactory.Execute(messaging.ICheckJobExistMessageRequest{
			JobID: uuid.MustParse(jobID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
		}

		msgData = map[string]interface{}{
			"job_id": jobID,
			"exists": message.Exists,
		}
	default:
		log.Printf("Unknown message type: %s", docMsg.MessageType)
	}
	// reply
	reply := response.RabbitMQResponse{
		ID:          docMsg.ID,
		MessageType: docMsg.MessageType,
		MessageData: msgData,
	}
	msg := RabbitMsg{
		QueueName: docMsg.ReplyTo,
		Reply:     reply,
	}
	rchan <- msg
}
