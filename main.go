package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/confluentinc/confluent-kafka-go/kafka"
)

var (
	servers = flag.String("servers", "localhost", "bootstrap.servers")
	topic   = flag.String("topic", "example", "")
)

func main() {
	flag.Parse()
	ch := make(chan string, 10)
	go publish(ch)
	r := bufio.NewScanner(os.Stdin)
	for r.Scan() {
		ch <- r.Text()
	}
	if err := r.Err(); err != nil && err != io.EOF {
		panic(err)
	}
}

func publish(ch <-chan string) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": *servers})
	if err != nil {
		panic(err)
	}
	defer p.Close()
	go printDeliveryStatus(p)
	for text := range ch {
		p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{
				Topic:     topic,
				Partition: kafka.PartitionAny,
			},
			Value: []byte(text),
		}, nil)
	}
}

func printDeliveryStatus(p *kafka.Producer) {
	for e := range p.Events() {
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				fmt.Printf("Delivery failed: %v\n", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v\n", ev.TopicPartition)
			}
		}
	}
}
