package containers

import (
	"log"
	"log-aggregation/consumer/config"
	"log-aggregation/consumer/connection"
	"log-aggregation/consumer/services"
	"os"
)

func Serve() {
	config.SetConfig()
	connectRabbitMQ, err := connection.ConnectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after multiple retries: %v", err)
	}

	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	consumer := services.ConsumerServiceInstance(channelRabbitMQ)

	consumerName := os.Getenv("CONSUMER_NAME")
	if consumerName == "" {
		log.Fatalf("Please set LOG_LEVEL environment variable")
		return
	}
	go consumer.ConsumeFromQueue(consumerName)

	log.Println("Successfully connected to RabbitMQ")
	log.Println("Waiting for messages...")

	forever := make(chan bool)
	<-forever

}
