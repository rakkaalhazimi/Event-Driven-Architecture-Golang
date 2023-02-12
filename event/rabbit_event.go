package event

import (
	"context"
	"fmt"

	amqp "github.com/rabbitmq/amqp091-go"

	utils "event-driven/utils"
)

type RabbitEvent struct {
	Conn *amqp.Connection
}

func (e *RabbitEvent) ConsumeMessage(topic string) (interface{}, error) {

	ch, err := e.Conn.Channel()
	utils.FailsOnError(err, "Cannot create RabbitMQ channel")

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailsOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		true,   // auto-ack
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	utils.FailsOnError(err, "Failed to declare a consumer")

	return msgs, err
}

func (e *RabbitEvent) PublishMessage(message interface{}, topic string) error {

	ch, err := e.Conn.Channel()
	utils.FailsOnError(err, "Cannot create RabbitMQ channel")

	q, err := ch.QueueDeclare(
		topic, // name
		false, // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	utils.FailsOnError(err, "Failed to declare a queue")

	err = ch.PublishWithContext(
		context.TODO(), // context
		"",             // exchange
		q.Name,         // routing key,
		false,          // mandatory
		false,          // immediate
		amqp.Publishing{
			ContentType: "text/plain",
			Body:        message.([]byte),
		},
	)
	utils.FailsOnError(err, "Failed to publish message")

	return err
}

func NewRabbitEvent(uri string) *RabbitEvent {
	conn, err := amqp.Dial(uri)
	utils.FailsOnError(err, fmt.Sprintf(`Cannot connect to RabbitMQ with address: '%v'`, uri))

	return &RabbitEvent{
		Conn: conn,
	}
}
