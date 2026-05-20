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
	fmt.Println("Starting Peril server...")

	conString := "amqp://guest:guest@localhost:5672/"

	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, os.Interrupt)

	connection, err := amqp.Dial(conString)
	if err != nil {
		log.Fatalf("could not connect to RabbitMQ: %v", err)
	}
	defer connection.Close()

	channel, err := connection.Channel()
	if err != nil {
		log.Fatalf("error opening the channel: %v", err)
	}
	defer channel.Close()

	gamelogic.PrintServerHelp()

	for {
		words := gamelogic.GetInput()

		if len(words) == 0 {
			continue
		}

		switch words[0] {
		case "pause":
			log.Println("Sending pause message...")

			state := routing.PlayingState{
				IsPaused: true,
			}

			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, state)
			if err != nil {
				log.Fatalf("error publishing json: %v", err)
			}
		case "resume":
			log.Println("Sending resume message...")

			state := routing.PlayingState{
				IsPaused: false,
			}

			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, state)
			if err != nil {
				log.Fatalf("error publishing json: %v", err)
			}
		case "quit":
			log.Println("Exiting game...")
			return
		default:
			log.Println("Unrecognized command.")
		}
	}
}
