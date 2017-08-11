package main

import (
	"context"
	"flag"
	"log"
	"strings"

	elastic "gopkg.in/olivere/elastic.v5"

	cluster "github.com/bsm/sarama-cluster"
)

var (
	kafkaAddr = flag.String("kafka", "127.0.0.1:9092", "kafka address list")
	esAddr    = flag.String("es", "http://127.0.0.1:9200", "es address list")
	topic     = flag.String("topic", "falcon", "kafka topic")
)

func main() {
	flag.Parse()
	kafkaList := strings.Split(*kafkaAddr, ",")
	esList := strings.Split(*esAddr, ",")
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	consumer, err := cluster.NewConsumer(kafkaList, "falcon-saver", []string{*topic}, config)
	if err != nil {
		log.Fatal(err)
	}

	client, err := elastic.NewClient(elastic.SetURL(esList...))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case msg := <-consumer.Messages():
			_, err := client.Index().
				Index("falcon").
				Type("falcon").
				BodyString(string(msg.Value)).
				Do(context.TODO())
			if err != nil {
				log.Print(err)
			}
			log.Printf("%s", msg.Value)
		case err := <-consumer.Errors():
			log.Print(err)
		}
	}
}
