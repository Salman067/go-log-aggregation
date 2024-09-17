package containers

import (
	"fmt"
	"log"
	"log-aggregation/producer/config"
	"log-aggregation/producer/connection"
	"log-aggregation/producer/models"
	"log-aggregation/producer/services"

	"github.com/labstack/echo/v4"
)

func Serve(e *echo.Echo) {
	config.SetConfig()
	connectRabbitMQ, err := connection.ConnectToRabbitMQ()
	if err != nil {
		log.Fatalf("Failed to connect to RabbitMQ after multiple retries: %v", err)
	}

	// Open a channel to RabbitMQ
	channelRabbitMQ, err := connectRabbitMQ.Channel()
	if err != nil {
		log.Fatalf("Failed to open a channel: %v", err)
	}

	producer := services.ProducerServiceInstance(channelRabbitMQ)

	var logData []*models.LogMessage

	for i := 0; i < 100000; i++ {
		logs := producer.CollectLogs()
		logData = append(logData, logs...)
	}

	go producer.PublishLogsToRabbitMQ(logData)

	log.Fatal(e.Start(fmt.Sprintf(":%s", config.LocalConfig.SERVER_PORT)))

}
