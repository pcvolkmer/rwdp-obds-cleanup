package main

import (
	"fmt"
	"github.com/alecthomas/kong"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"os"
	"strings"
)

var (
	cli     *CLI
	context *kong.Context
)

type CLI struct {
	BootstrapServers string `env:"BOOTSTRAP_SERVERS" help:"Kafka bootstrap servers" required:""`
	InputTopics      string `env:"INPUT_TOPICS" help:"Kafka input topic(s)" required:""`
	OutputTopic      string `env:"OUTPUT_TOPIC" help:"Kafka output topic" required:""`
}

func initCLI() {
	cli = &CLI{}
	context = kong.Parse(cli,
		kong.Name("rwdp-obds-cleanup"),
		kong.Description("A simple pipeline step to clean up incoming kafka messages"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
		}),
	)
}

func main() {

	initCLI()

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cli.BootstrapServers,
		"group.id":          "rwdp-obds-cleanup",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		log.Fatalf("Kafka consumer available: %v", err)
	}

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cli.BootstrapServers,
	})
	if err != nil {
		log.Fatalf("Kafka Producer not available: %v", err)
	}

	inputTopics := strings.Split(cli.InputTopics, ",")
	outputTopic := cli.OutputTopic

	err = consumer.SubscribeTopics(inputTopics, nil)

	var run = true
	for run == true {
		ev := consumer.Poll(100)
		switch e := ev.(type) {
		case *kafka.Message:
			// Prepare message
			e.TopicPartition.Topic = &outputTopic
			e.TopicPartition.Partition = kafka.PartitionAny

			// Manipulate message value
			// ... TBD

			// Send message
			_ = producer.Produce(e, nil)
		case kafka.Error:
			_, _ = fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
		}
	}

}
