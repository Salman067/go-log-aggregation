package services

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"log-aggregation/producer/domain"
	"log-aggregation/producer/models"
	"os"

	"github.com/streadway/amqp"
)

var queueNames = map[string]bool{}

type ProducerService struct {
	channelRabbitMQ *amqp.Channel
}

func ProducerServiceInstance(channelRabbitMQ *amqp.Channel) domain.ProducerServiceInterface {
	return &ProducerService{
		channelRabbitMQ: channelRabbitMQ,
	}
}

func (ps *ProducerService) DeclareQueues(logs []*models.LogMessage) error {
	for _, logData := range logs {
		queue, ok := queueNames[logData.AppName]

		if !ok && !queue {
			_, err := ps.channelRabbitMQ.QueueDeclare(
				logData.AppName, // queue name
				true,            // durable
				false,           // auto-delete
				false,           // exclusive
				false,           // no-wait
				nil,             // arguments
			)
			if err != nil {
				return err
			}
			queueNames[logData.AppName] = true
		}

	}
	return nil
}

func (ps *ProducerService) CollectLogs() []*models.LogMessage {
	logFilePath := "/home/salman/Desktop/vivasoft/go-log-aggregator/data.json"
	//logFilePath := "/app/data.json"
	file, err := os.Open(logFilePath)
	if err != nil {
		log.Fatalf("Failed to open log file: %v", err)
	}
	defer file.Close()

	fileContent, err := ioutil.ReadAll(file)
	if err != nil {
		log.Fatalf("Failed to read log file: %v", err)
	}

	var allLogs []*models.LogMessage

	if err := json.Unmarshal(fileContent, &allLogs); err != nil {
		log.Fatalf("Failed to unmarshal log messages: %v", err)
	}

	return allLogs

}

func (ps *ProducerService) PublishLogsToRabbitMQ(logs []*models.LogMessage) {
	err := ps.DeclareQueues(logs)
	if err != nil {
		log.Fatalf("Failed to declare queues: %v", err)
		return
	}
	for _, logMessage := range logs {
		_, ok := queueNames[logMessage.AppName]
		if !ok {
			log.Panic("this app name is not exist")
			return
		}

		formattedMessage := fmt.Sprintf("%s [%s] [%s : %s]: %s", logMessage.Timestamp, logMessage.LogLevel, logMessage.AppName, logMessage.ServiceName, logMessage.Message)

		message := amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "application/json",
			Body:         []byte(formattedMessage),
		}

		err = ps.channelRabbitMQ.Publish(
			"",                 // exchange
			logMessage.AppName, // routing key (queue name)
			false,              // mandatory
			false,              // immediate
			message,
		)
		if err != nil {
			log.Printf("Failed to publish log message to %s: %v", logMessage.AppName, err)
		} else {
			log.Printf("Log sent to %s: %s", logMessage.AppName, formattedMessage)
		}
	}
}
