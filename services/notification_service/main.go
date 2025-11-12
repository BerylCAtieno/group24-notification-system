package main

import (
	"log"

	emailservice "github.com/uloamaka/notification-service/email_service/sender"
	"github.com/uloamaka/notification-service/email_worker/consumer"
	"github.com/uloamaka/notification-service/email_worker/processor"
)

func main() {
	smtp := emailservice.NewSMTPProvider()
	p := processor.NewEmailProcessor(smtp)
	c := consumer.NewKafkaEmailConsumer(p)

	if err := c.Start("email.jobs"); err != nil {
		log.Fatalf("failed to start consumer: %v", err)
	}
}
