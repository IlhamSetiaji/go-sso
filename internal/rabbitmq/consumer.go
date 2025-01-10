package rabbitmq

import (
	"app/go-sso/internal/http/request"
	"app/go-sso/internal/http/response"
	empMessaging "app/go-sso/internal/messaging/employee"
	messaging "app/go-sso/internal/messaging/job"
	orgMessaging "app/go-sso/internal/messaging/organization"
	userMessaging "app/go-sso/internal/messaging/user"
	"app/go-sso/internal/service"
	"app/go-sso/utils"
	"encoding/json"
	"errors"
	"os"

	"github.com/google/uuid"
	"github.com/rabbitmq/amqp091-go"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func InitConsumer(viper *viper.Viper, log *logrus.Logger) {
	// conn
	log.Printf("CHECKING: url on rabbitmq: ", viper.GetString("rabbitmq.url"))
	conn, err := amqp091.Dial(viper.GetString("rabbitmq.url"))
	if err != nil {
		log.Printf("ERROR: fail init consumer: %s", err.Error())
		log.Printf("ERROR: fail init consumer conn: ", viper.GetString("rabbitmq.url"))
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
			docRply := &response.RabbitMQResponse{}
			err = json.Unmarshal(msg.Body, docMsg)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received docMsg: %v", docMsg)

			err = json.Unmarshal(msg.Body, docRply)
			if err != nil {
				log.Printf("ERROR: fail unmarshl: %s", msg.Body)
				continue
			}
			log.Printf("INFO: received docRply: %v", docRply)

			// ack for message
			err = msg.Ack(true)
			if err != nil {
				log.Printf("ERROR: fail to ack: %s", err.Error())
			}

			if rchan, ok := utils.Rchans[docRply.ID]; ok {
				rchan <- *docRply
			}

			// handle docMsg
			handleMsg(docMsg, log, viper)
		}
	}
}

func handleMsg(docMsg *request.RabbitMQRequest, log *logrus.Logger, viper *viper.Viper) {
	// switch case
	var msgData map[string]interface{}

	switch docMsg.MessageType {
	case "reply":
		log.Printf("INFO: received reply message")
		return
	case "check_job_exist":
		jobID, ok := docMsg.MessageData["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := messaging.CheckJobExistMessageFactory(log)
		message, err := messageFactory.Execute(messaging.ICheckJobExistMessageRequest{
			JobID: uuid.MustParse(jobID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"job_id": jobID,
			"exists": message.Exists,
		}
	case "find_organization_by_id":
		organizationID, ok := docMsg.MessageData["organization_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'organization_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'organization_id'").Error(),
			}
			break
		}

		messageFactory := orgMessaging.FindOrganizationByIDMessageFactory(log)
		message, err := messageFactory.Execute(orgMessaging.IFindOrganizationByIDMessageRequest{
			OrganizationID: uuid.MustParse(organizationID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"organization_id":       organizationID,
			"name":                  message.Name,
			"organization_category": message.OrganizationCategory,
		}
	case "find_job_by_id":
		jobID, ok := docMsg.MessageData["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := messaging.FindJobByIDMessageFactory(log)
		message, err := messageFactory.Execute(messaging.IFindJobByIDMessageRequest{
			JobID: uuid.MustParse(jobID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"job_id": jobID,
			"name":   message.Name,
		}
	case "find_user_by_id":
		userID, ok := docMsg.MessageData["user_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'user_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := userMessaging.FindUserByIDMessageFactory(log)
		message, err := messageFactory.Execute(userMessaging.IFindUserByIDMessageRequest{
			UserID: uuid.MustParse(userID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"user_id": userID,
			"name":    message.Name,
		}
	case "find_organization_location_by_id":
		organizationLocationID, ok := docMsg.MessageData["organization_location_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'organization_location_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := orgMessaging.FindOrganizationLocationByIDMessageFactory(log)
		message, err := messageFactory.Execute(orgMessaging.IFindOrganizationLocationByIDMessageRequest{
			OrganizationLocationID: uuid.MustParse(organizationLocationID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"organization_location_id": organizationLocationID,
			"name":                     message.Name,
		}
	case "find_job_level_by_id":
		jobLevelID, ok := docMsg.MessageData["job_level_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_level_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := messaging.FindJobLevelByIDMessageFactory(log)
		message, err := messageFactory.Execute(messaging.IFindJobLevelByIDMessageRequest{
			JobLevelID: uuid.MustParse(jobLevelID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"job_level_id": jobLevelID,
			"name":         message.Name,
			"level":        message.Level,
		}
	case "check_job_by_job_level":
		jobID, ok := docMsg.MessageData["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		jobLevelID, ok := docMsg.MessageData["job_level_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_level_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_level_id'").Error(),
			}
			break
		}

		messageFactory := messaging.CheckJobByJobLevelMessageFactory(log)
		message, err := messageFactory.Execute(messaging.ICheckJobByJobLevelMessageRequest{
			JobID:      uuid.MustParse(jobID),
			JobLevelID: uuid.MustParse(jobLevelID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"job_id":       jobID,
			"job_level_id": jobLevelID,
			"exists":       message.Exists,
		}
	case "find_organization_structure_by_id":
		organizationStructureID, ok := docMsg.MessageData["organization_structure_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'organization_structure_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'organization_structure_id'").Error(),
			}
			break
		}

		messageFactory := orgMessaging.FindOrganizationStructureByIDMessageFactory(log)
		message, err := messageFactory.Execute(orgMessaging.IFindOrganizationStructureByIDMessageRequest{
			OrganizationStructureID: uuid.MustParse(organizationStructureID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"organization_structure_id": organizationStructureID,
			"name":                      message.Name,
		}
	case "get_user_me":
		userID, ok := docMsg.MessageData["user_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'user_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'user_id'").Error(),
			}
			break
		}

		messageFactory := userMessaging.GetUserMeMessageFactory(log)
		message, err := messageFactory.Execute(userMessaging.IGetUserMeMessageRequest{
			UserID: uuid.MustParse(userID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"user": message.User,
		}
	case "find_job_data_by_id":
		jobID, ok := docMsg.MessageData["job_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'job_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'job_id'").Error(),
			}
			break
		}

		messageFactory := messaging.FindJobDataByIdMessageFactory(log)
		message, err := messageFactory.Execute(messaging.IFindJobDataByIdMessageRequest{
			JobID: uuid.MustParse(jobID),
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"job_id": jobID,
			"job":    message.Job,
		}
	case "find_employee_by_id":
		employeeID, ok := docMsg.MessageData["employee_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'employee_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'employee_id'").Error(),
			}
			break
		}

		messageFactory := empMessaging.FindEmployeeByIDMessageFactory(log)
		message, err := messageFactory.Execute(empMessaging.IFindEmployeeByIDMessageRequest{
			EmployeeID: employeeID,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"employee_id": employeeID,
			"employee":    message.Employee,
		}
	case "find_organization_locations_paginated":
		page, ok := docMsg.MessageData["page"].(float64)
		if !ok {
			log.Printf("Invalid request format: missing 'page'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'page'").Error(),
			}
			break
		}

		pageSize, ok := docMsg.MessageData["page_size"].(float64)
		if !ok {
			log.Printf("Invalid request format: missing 'page_size'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'page_size'").Error(),
			}
			break
		}

		search, ok := docMsg.MessageData["search"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'search'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'search'").Error(),
			}
			break
		}

		organizationID, ok := docMsg.MessageData["organization_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'organization_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'organization_id'").Error(),
			}
			break
		}

		includedIDsInterface, ok := docMsg.MessageData["included_ids"].([]interface{})
		if !ok {
			log.Printf("Invalid request format: missing 'included_ids'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'included_ids'").Error(),
			}
			break
		}

		isNull, ok := docMsg.MessageData["is_null"].(bool)
		if !ok {
			log.Printf("Invalid request format: missing 'is_null'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'is_null'").Error(),
			}
			break
		}

		includedIDs := make([]string, len(includedIDsInterface))
		for i, v := range includedIDsInterface {
			str, ok := v.(string)
			if !ok {
				log.Printf("Invalid request format: missing 'included_ids'")
				msgData = map[string]interface{}{
					"error": errors.New("missing 'included_ids'").Error(),
				}
				break
			}
			includedIDs[i] = str
		}

		// includedIds, ok := docMsg.MessageData["included_ids"].([]string)
		// if !ok {
		// 	log.Printf("Invalid request format: missing 'included_ids'")
		// 	msgData = map[string]interface{}{
		// 		"error": errors.New("missing 'included_ids'").Error(),
		// 	}
		// 	break
		// }

		messageFactory := orgMessaging.FindAllOrgLocationPaginatedUseCaseFactory(log)
		message, err := messageFactory.Execute(&orgMessaging.IFindAllOrgLocationPaginatedRequest{
			Page:           int(page),
			PageSize:       int(pageSize),
			Search:         search,
			IncludedIDs:    includedIDs,
			IsNull:         isNull,
			OrganizationID: organizationID,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"organization_locations": message.OrganizationLocations,
			"total":                  message.Total,
		}
	case "get_all_job_data":
		messageFactory := messaging.GetAllJobMessageFactory(log)
		message, err := messageFactory.Execute()

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"jobs": message.Jobs,
		}
	case "find_all_organization":
		includedIDsInterface, ok := docMsg.MessageData["included_ids"].([]interface{})
		if !ok {
			log.Printf("Invalid request format: missing 'included_ids'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'included_ids'").Error(),
			}
			break
		}
		includedIDs := make([]string, len(includedIDsInterface))
		for i, v := range includedIDsInterface {
			str, ok := v.(string)
			if !ok {
				log.Printf("Invalid request format: missing 'included_ids'")
				msgData = map[string]interface{}{
					"error": errors.New("missing 'included_ids'").Error(),
				}
				break
			}
			includedIDs[i] = str
		}

		messageFactory := orgMessaging.FindAllOrgMessageFactory(log)
		message, err := messageFactory.Execute(&orgMessaging.IFindAllOrgMessageRequest{
			IncludedIDs: includedIDs,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"organizations": message.Organizations,
		}
	case "find_all_jobs_by_ids":
		includedIDsInterface, ok := docMsg.MessageData["included_ids"].([]interface{})
		if !ok {
			log.Printf("Invalid request format: missing 'included_ids'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'included_ids'").Error(),
			}
			break
		}

		includedIDs := make([]string, len(includedIDsInterface))
		for i, v := range includedIDsInterface {
			str, ok := v.(string)
			if !ok {
				log.Printf("Invalid request format: missing 'included_ids'")
				msgData = map[string]interface{}{
					"error": errors.New("missing 'included_ids'").Error(),
				}
				break
			}
			includedIDs[i] = str
		}

		messageFactory := messaging.FindAllJobsMessageFactory(log)
		message, err := messageFactory.Execute(&messaging.IFindAllJobsMessageRequest{
			IncludedIDs: includedIDs,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"jobs": message.Jobs,
		}
	case "find_all_organization_locations_by_ids":
		includedIDsInterface, ok := docMsg.MessageData["included_ids"].([]interface{})
		if !ok {
			log.Printf("Invalid request format: missing 'included_ids'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'included_ids'").Error(),
			}
			break
		}

		includedIDs := make([]string, len(includedIDsInterface))
		for i, v := range includedIDsInterface {
			str, ok := v.(string)
			if !ok {
				log.Printf("Invalid request format: missing 'included_ids'")
				msgData = map[string]interface{}{
					"error": errors.New("missing 'included_ids'").Error(),
				}
				break
			}
			includedIDs[i] = str
		}

		messageFactory := orgMessaging.FindAllOrgLocationsByIDsMessageFactory(log)
		message, err := messageFactory.Execute(&orgMessaging.IFindAllOrgLocationsByIDsMessageRequest{
			IncludedIDs: includedIDs,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"organization_locations": message.OrgLocations,
		}
	case "find_all_org_structure_children_ids":
		parentID, ok := docMsg.MessageData["parent_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'parent_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'parent_id'").Error(),
			}
			break
		}

		messageFactory := orgMessaging.FindAllOrgStructureChildrenIDsMessageFactory(log)
		message, err := messageFactory.Execute(&orgMessaging.IFindAllOrgStructureChildrenIDsMessageRequest{
			ParentID: parentID,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"children_ids": message.ChildrenIDs,
		}
	case "find_all_jobs_by_organization_id":
		organizationID, ok := docMsg.MessageData["organization_id"].(string)
		if !ok {
			log.Printf("Invalid request format: missing 'organization_id'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'organization_id'").Error(),
			}
			break
		}

		messageFactory := messaging.FindAllJobsByOrganizationIDMessageFactory(log)
		message, err := messageFactory.Execute(&messaging.IFindAllJobsByOrganizationIDMessageRequest{
			OrganizationID: organizationID,
		})

		if err != nil {
			log.Printf("Failed to execute message: %v", err)
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		}

		msgData = map[string]interface{}{
			"jobs": message.Jobs,
		}
	case "send_mail":
		to, ok := docMsg.MessageData["to"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'to'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'to'").Error(),
			}
			break
		}
		subject, ok := docMsg.MessageData["subject"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'subject'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'subject'").Error(),
			}
			break
		}
		body, ok := docMsg.MessageData["body"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'body'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'body'").Error(),
			}
			break
		}
		from, ok := docMsg.MessageData["from"].(string)
		if !ok {
			log.Errorf("Invalid request format: missing 'from'")
			msgData = map[string]interface{}{
				"error": errors.New("missing 'from'").Error(),
			}
			break
		}

		mailService := service.MailServiceFactory(log, viper)
		err := mailService.SendMail(service.MailData{
			From:    from,
			To:      []string{to},
			Subject: subject,
			Body:    body,
		})
		if err != nil {
			log.Errorf("ERROR: fail send mail: %s", err.Error())
			msgData = map[string]interface{}{
				"error": err.Error(),
			}
			break
		} else {
			log.Printf("INFO: success send mail")
		}

		msgData = map[string]interface{}{
			"message": "success",
		}
	default:
		log.Printf("Unknown message type, please recheck your type: %s", docMsg.MessageType)

		msgData = map[string]interface{}{
			"error": errors.New("unknown message type").Error(),
		}
	}
	// reply
	reply := response.RabbitMQResponse{
		ID: docMsg.ID,
		// MessageType: docMsg.MessageType,
		MessageType: "reply",
		MessageData: msgData,
	}
	msg := RabbitMsg{
		QueueName: docMsg.ReplyTo,
		Reply:     reply,
	}
	rchan <- msg
}
