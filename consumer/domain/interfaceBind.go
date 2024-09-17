package domain

type ConsumerServiceInterface interface {
	ConsumeFromQueue(logLevel string)
}
