package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	fmt.Println("Starting Peril client...")

	conString := "amqp://guest:guest@localhost:5672/"

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	connection, err := amqp.Dial(conString)
	if err != nil {
		log.Fatalf("Could not connect to RabbitMQ: %v", err)
	}
	defer connection.Close()

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		log.Fatalf("could not retrieve username: %v", err)
	}

	queueName := routing.PauseKey + "." + username

	_, _, err = pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, queueName, routing.PauseKey, pubsub.Transient)
	if err != nil {
		log.Fatalf("error creating or binding queue: %v", err)
	}

	fmt.Println("Client connected successfully!")

	<-signalChan

	fmt.Println("\nShutting Client down...")
}
