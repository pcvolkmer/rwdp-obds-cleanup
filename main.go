/*
 * This file is part of rwdp-obds-cleanup
 *
 * Copyright (c) 2024  Paul-Christian Volkmer
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as published
 * by the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with this program.  If not, see <https://www.gnu.org/licenses/>.
 */

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
			if recordValue, err := ParseRecordValue(e.Value); err == nil {
				recordValue.RemoveLeadingZerosFromPatientIds()
				fixedRecordValue, err := recordValue.ToJson()
				if err == nil {
					e.Value = fixedRecordValue
				}
			}

			// Send message
			_ = producer.Produce(e, nil)
		case kafka.Error:
			_, _ = fmt.Fprintf(os.Stderr, "%% Error: %v\n", e)
			run = false
		default:
		}
	}

}
