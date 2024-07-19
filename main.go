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
	"github.com/alecthomas/kong"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
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
	// Features
	EnableRemovePatientIdZeros bool `env:"ENABLE_REMOVE_PATIENT_ID_ZERO" help:"Feature: Remove leading patient id zeros" default:"true"`
	EnableDropNonObds2         bool `env:"ENABLE_DROP_NON_OBDS_2" help:"Feature: Drop records not containing oBDS 2.x message"`
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

	log.Printf("Trying to connect to bootstrap servers: %s\n", cli.BootstrapServers)

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": cli.BootstrapServers,
		"group.id":          "rwdp-obds-cleanup",
		"auto.offset.reset": "smallest",
	})
	if err != nil {
		log.Fatalf("Kafka consumer not available: %v", err)
	}
	log.Println("Kafka consumer is ready")

	producer, err := kafka.NewProducer(&kafka.ConfigMap{
		"bootstrap.servers": cli.BootstrapServers,
	})
	if err != nil {
		log.Fatalf("Kafka producer not available: %v", err)
	}
	log.Println("Kafka producer is ready")

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

				if cli.EnableDropNonObds2 && !recordValue.IsObdsVersion2x() {
					log.Printf("Dropping non oBDS 2.x record - Key: %v\n", string(e.Key))
					continue
				}

				if cli.EnableRemovePatientIdZeros {
					recordValue.RemoveLeadingZerosFromPatientIds()
				}

				fixedRecordValue, err := recordValue.ToJson()
				if err == nil {
					e.Value = fixedRecordValue
				}
			}

			// Send message
			if err = producer.Produce(e, nil); err == nil {
				log.Printf("Record processed - Key: %v\n", string(e.Key))
			} else {
				log.Fatalf("Kafka producer failed: %v", err)
			}
		case kafka.Error:
			log.Printf("Error: %v\n", e)
			run = false
		default:
			// Nop
		}
	}

}
