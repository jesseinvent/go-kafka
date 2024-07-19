package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/IBM/sarama"
)

func main() {
	topic := "comments"

	consumer, err := connectConsumer([]string{"localhost:29092"})

	if err != nil {
		panic(err)
	}

	consumerPartition, err := consumer.ConsumePartition(topic, 0, sarama.OffsetOldest)

	if err != nil {
		panic(err)
	}

	fmt.Println("consumer started")

	sigChan := make(chan os.Signal, 1)

	signal.Notify(sigChan, syscall.SIGTTIN, syscall.SIGTERM)

	msgCount := 0

	done := make(chan struct{})

	go func() {
		for {
			select {
			case err := <-consumerPartition.Errors():
				fmt.Println(err)

			case msg := <-consumerPartition.Messages():
				msgCount++
				fmt.Printf("Received message count: %d: | Topic (%s) | Message(%s)\n", msgCount, string(msg.Topic), string(msg.Value))

			case <-sigChan:
				fmt.Println("Interruption detected")
				done <- struct{}{}
			}
		}
	}()

	<-done

	fmt.Println("Processed", msgCount, "messages")

	if err := consumer.Close(); err != nil {
		panic(err)
	}

}

// Connects to Kafka consumer
func connectConsumer(brokersUrl []string) (sarama.Consumer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	conn, err := sarama.NewConsumer(brokersUrl, config)

	if err != nil {
		return nil, err
	}

	return conn, nil
}
