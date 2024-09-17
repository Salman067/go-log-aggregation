package domain

import "log-aggregation/producer/models"

type ProductionRepoInterface interface {
	DeclareQueues() error
}

type ProducerServiceInterface interface {
	DeclareQueues(logs []*models.LogMessage) error
	CollectLogs() []*models.LogMessage
	PublishLogsToRabbitMQ(logs []*models.LogMessage)
}
