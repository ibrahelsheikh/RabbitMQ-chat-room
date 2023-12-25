/*
USAGE:
1) you can push to ONLY ONE TOPIC, you may modify this code to emit to multiple topics, loop on publish
2) THIS CAN'T HAVE * OR #

go run emit_logs_topic.go anonymous.info # this sends to anonymous.info
go run emit_logs_topic.go kernel.info  # this sends to kernel.info
*/
package main

import (
	"fmt"
	"github.com/streadway/amqp"
	"load_balancer/utils"
	"log"
	"os"
	"time"
)

func main() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	utils.FailOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	utils.FailOnError(err, "Failed to open a channel")
	defer ch.Close()

	err = ch.ExchangeDeclare(
		"logs_topic", // name
		"topic",      // type
		true,         // durable
		false,        // auto-deleted
		false,        // internal
		false,        // no-wait
		nil,          // arguments
	)
	utils.FailOnError(err, "Failed to declare an exchange")

	topic := severityFrom(os.Args) // a command line arg like web.info (this is a topic)
	// this is the topic I'm interested in
	for i := 0; ; {

		body := fmt.Sprintf("Message %d sent form topic %s", i, topic)

		err = ch.Publish(
			"logs_topic", // exchange
			topic,        // routing key (THIS CAN'T HAVE * OR #)
			false,        // mandatory
			false,        // immediate
			amqp.Publishing{
				ContentType: "text/plain",
				Body:        []byte(body),
			})
		utils.FailOnError(err, "Failed to publish a message")

		log.Printf(" [x] Sent %s", body)
		time.Sleep(400 * time.Millisecond)
		i += 1
	}
}

func severityFrom(args []string) string {
	// EXAMPLE: go run emit_log_topic.go anonymous.info, you can push to ONLY ONE TOPIC
	// (THIS CAN'T HAVE * OR #)
	// the default queue is anonymous.info
	// you may try something else
	var s string
	if (len(args) < 2) || os.Args[1] == "" {
		s = "anonymous.info"
	} else {
		s = os.Args[1]
	}
	return s
}
