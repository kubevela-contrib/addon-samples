package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	var username, password, host, port string
	flag.StringVar(&username, "user", "admin", "a string")
	flag.StringVar(&password, "password", "admin", "a string")
	flag.StringVar(&host, "host", "localhost", "a string")
	flag.StringVar(&port, "port", "5672", "a string")

	flag.Parse()

	conn, err := amqp.Dial("amqp://" + username + ":" + password + "@" + host + ":" + port + "/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")
	for {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		var body string

		fmt.Println("Enter the message:")
		body, err = bufio.NewReader(os.Stdin).ReadString('\n')

		failOnError(err, "Failed to read message")
		body = strings.TrimSuffix(body, "\n")

		err = ch.PublishWithContext(ctx,
			"",     // exchange
			q.Name, // routing key
			false,  // mandatory
			false,  // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		failOnError(err, "Failed to publish a message")
		log.Printf(" [x] Sent %s\n", body)
		cancel()
	}
}
