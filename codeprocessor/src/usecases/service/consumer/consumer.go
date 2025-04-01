package consumer

import (
	"codeproccesor/src/domain"
	"codeproccesor/src/usecases"
	"encoding/json"
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

type TaskConsumer struct {
	pool      usecases.ProcessPool
	stopCh    chan struct{}
	isStopped bool

	logger *log.Logger
}

func NewTaskConsumer(pool usecases.ProcessPool, logger *log.Logger) usecases.Consumer {
	return &TaskConsumer{
		pool:      pool,
		stopCh:    make(chan struct{}),
		isStopped: false,
	}
}

func (c *TaskConsumer) Consume(amqpUrl string, queueName string) error {
	if c.isStopped {
		return fmt.Errorf("consumer stopped")
	}

	conn, err := amqp.Dial(amqpUrl)
	if err != nil {
		return fmt.Errorf("amqp connection error: %s", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		return fmt.Errorf("amqp channel error: %s", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare(queueName, false, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("amqp queue declare error: %s", err)
	}

	msgCh, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		return fmt.Errorf("amqp consume channel error: %s", err)
	}

LISTENER:
	for {
		select {
		case dmsg := <-msgCh:
			c.processDeliveryMessage(dmsg)
		case <-c.stopCh:
			break LISTENER
		}
	}

	return nil
}

func (c *TaskConsumer) processDeliveryMessage(dsmg amqp.Delivery) {
	task := domain.Task{}

	err := json.Unmarshal(dsmg.Body, &task)
	if err != nil {
		c.logger.Printf("[CONSUME ERROR] task body json error: %s\n", err)
	}

	c.logger.Printf("start process task %s\n", task.Uuid)

	err = c.pool.Push(task)
	if err != nil {
		c.logger.Printf("[POOL PUSH ERROR] error: %s\n", err)
	}
}

func (c *TaskConsumer) Stop() {
	if !c.isStopped {
		c.isStopped = true
		close(c.stopCh)
	}
}
