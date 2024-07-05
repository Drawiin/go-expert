package rabbitmq

import (
	amqp "github.com/rabbitmq/amqp091-go"
)

func OpenChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}

	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}

	return ch
}

// / Simple consumer, in the real wold this would be a lot more complex
// / Hanving error handling, logging, reconnecting, etc.
func Consume(ch *amqp.Channel, out chan<- amqp.Delivery, queue, consumer string) error {
	msgs, err := ch.Consume(queue, consumer, false, false, false, false, nil)
	if err != nil {
		return err
	}
	for msg := range msgs {
		out <- msg
	}
	return nil
}

func Publish(ch *amqp.Channel, body, exchange string ) error {
	return ch.Publish(exchange, "", false, false, amqp.Publishing{
		ContentType: "text/plain",
		Body:        []byte(body),
	})
}
