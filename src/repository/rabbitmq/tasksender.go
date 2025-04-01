package rabbitmq

import (
	"encoding/json"
	"fmt"
	"nullableocean-postupashki/src/domain"
	"nullableocean-postupashki/src/repository"

	"github.com/streadway/amqp"
)

type RabbitMQTaskSender struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
}

func NewRabbitMQTaskSender(amqpUrl string, queueName string) (repository.TaskSender, error) {
	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq connect error: %s", err)
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, fmt.Errorf("rabbitmq channel error: %s", err)
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, fmt.Errorf("rabbitmq queue declare error: %s", err)
	}

	return &RabbitMQTaskSender{
		conn:      conn,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (rs *RabbitMQTaskSender) Send(task *domain.Task) error {
	taskBody, err := rs.getSendBody(task)
	if err != nil {
		return fmt.Errorf("task send body error: %s", err)
	}

	err = rs.channel.Publish("", rs.queueName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        taskBody,
	})

	if err != nil {
		return fmt.Errorf("rabbitmq publish error: %s", err)
	}

	return nil
}

func (rs *RabbitMQTaskSender) getSendBody(task *domain.Task) ([]byte, error) {
	return json.Marshal(task)
}
