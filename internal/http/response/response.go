package response

type RabbitMQResponse struct {
	ID          string                 `json:"id"`
	MessageType string                 `json:"message_type"`
	MessageData map[string]interface{} `json:"message_data"`
}