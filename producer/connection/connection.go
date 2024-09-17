package connection

import (
	"log"
	"log-aggregation/producer/config"

	"time"

	"github.com/streadway/amqp"
)

func ConnectToRabbitMQ() (*amqp.Connection, error) {
	config := config.LocalConfig
	maxRetries := 10
	var conn *amqp.Connection
	var err error

	for i := 1; i <= maxRetries; i++ {
		conn, err = amqp.Dial(config.RABBITMQ_URL)
		if err == nil {
			return conn, nil
		}

		log.Printf("Failed to connect to RabbitMQ, retrying (%d/%d)...", i, maxRetries)
		time.Sleep(2 * time.Second)
	}

	return nil, err
}
