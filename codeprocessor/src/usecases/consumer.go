package usecases

type Consumer interface {
	Consume(amqpUrl string, queueName string) error
	Stop()
}
