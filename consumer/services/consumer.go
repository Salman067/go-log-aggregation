package services

import (
	"fmt"
	"log"
	"log-aggregation/consumer/domain"
	filewrite "log-aggregation/consumer/file_write"
	"log-aggregation/consumer/types"
	"runtime"
	"time"

	"github.com/streadway/amqp"
)

var queueNames = map[string]bool{}

type ConsumerService struct {
	channelRabbitMQ *amqp.Channel
}

func ConsumerServiceInstance(channelRabbitMQ *amqp.Channel) domain.ConsumerServiceInterface {
	return &ConsumerService{
		channelRabbitMQ: channelRabbitMQ,
	}
}

func (cs *ConsumerService) ConsumeFromQueue(logLevel string) {
	queueName, ok := queueNames[logLevel]
	if !ok && !queueName {
		queueNames[logLevel] = true
		types.LogFilePaths[logLevel] = fmt.Sprintf("/home/salman/Desktop/vivasoft/go-log-aggregator/consumer/logs/%s.log", logLevel)
	}

	err := cs.channelRabbitMQ.Qos(50000, 0, false)
	if err != nil {
		log.Fatalf("Failed to set QoS: %s", err)
	}

	messages, err := cs.channelRabbitMQ.Consume(
		logLevel, // queue name
		"",       // consumer
		false,    // auto-ack
		false,    // exclusive
		false,    // no local
		false,    // no wait
		nil,      // arguments
	)
	if err != nil {
		log.Fatalf("Failed to consume from the queue: %v", err)
	}

	bulkSize := 50000
	logs := []amqp.Delivery{}
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	ackChannel := make(chan amqp.Delivery, bulkSize)
	go func() {
		for message := range ackChannel {
			if err := message.Ack(false); err != nil {
				log.Printf("Failed to acknowledge message: %v", err)
			}
		}
	}()

	go func() {
		for {
			select {
			case message := <-messages:
				logs = append(logs, message)
				if len(logs) >= bulkSize {
					processAndWriteLogs(logs, logLevel, ackChannel)
					logs = []amqp.Delivery{}
				}

			case <-ticker.C:
				if len(logs) > 0 {
					processAndWriteLogs(logs, logLevel, ackChannel)
					logs = []amqp.Delivery{}
				}
			}
		}
	}()

	select {}
}

func processAndWriteLogs(logs []amqp.Delivery, logLevel string, ackChannel chan amqp.Delivery) {
	processLogs(logs, logLevel)

	err := filewrite.ProcessAndRetryWriteLogs(logLevel, logs, 0, 5)
	if err != nil {
		log.Printf("Failed to write logs to log file: %v", err)
		return
	}

	for _, message := range logs {
		ackChannel <- message
	}

	runtime.GC()
}
func processLogs(logs []amqp.Delivery, logLevel string) {
	for _, message := range logs {
		logMessage := string(message.Body)
		log.Printf(" > Received message from %s queue: %s\n", logLevel, logMessage)
	}
}
