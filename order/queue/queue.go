package queue

import (
	"fmt"
	"os"

	"github.com/streadway/amqp"
)


func Connect() *amqp.Channel {
	dns := os.Getenv("RABBITMQ_CONNECT_URL")

	conn, err := amqp.Dial(dns)
	if err != nil {
		panic(err.Error())
	}

	channel, err := conn.Channel()
	if err != nil {
		panic(err.Error())
	}

	return channel
}


func Publisher(payload []byte, exchange string, key string, ch *amqp.Channel) {
	err := ch.Publish(
		exchange,
		key,
		false,
		false,
		amqp.Publishing {
			ContentType: "application/json",
			Body: []byte(payload),
		})

	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Mensagem publicada: " + string(payload))
}


func Consumer(queue string, ch *amqp.Channel, in chan []byte) {

	q, err := ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		panic(err.Error())
	}

	msgs, err := ch.Consume(
		q.Name,
		"checkout",
		true,
		false,
		false,
		false,
		nil)

	if err != nil {
		panic(err.Error())
	}

	go func ()  {
		for m := range msgs {
			in <- []byte(m.Body)
		}
		close(in)
	}()
}
