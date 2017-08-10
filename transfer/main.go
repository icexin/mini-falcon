package main

import (
	"bufio"
	"context"
	"flag"
	"log"
	"net"

	elastic "gopkg.in/olivere/elastic.v5"

	"github.com/Shopify/sarama"
)

var (
	listenAddr = flag.String("listen", ":7070", "listen address")
	topic      = flag.String("topic", "falcon", "kafka topic")
	kafkaAddrs = flag.String("kafka", "127.0.0.1:9092", "kafka address list")
)

type Server struct {
	msgch chan<- *sarama.ProducerMessage
}

func NewServer(msgch chan<- *sarama.ProducerMessage) *Server {
	return &Server{
		msgch: msgch,
	}
}

func (s *Server) ListenAndServe(addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	for {
		conn, err := l.Accept()
		if err != nil {
			log.Fatal(err)
		}

		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn net.Conn) {
	defer conn.Close()
	r := bufio.NewReader(conn)
	for {
		line, err := r.ReadBytes('\n')
		if err != nil {
			log.Print(err)
			return
		}
		line = line[:len(line)-1]
		log.Print(string(line))
		_, err = client.Index().
			Index("falcon").
			Type("falcon").
			BodyString(string(line)).
			Do(context.TODO())
		if err != nil {
			log.Print(err)
		}
		continue
		msg := &sarama.ProducerMessage{
			Topic: *topic,
			Key:   nil,
			Value: sarama.ByteEncoder(line),
		}
		s.msgch <- msg
	}
}

var (
	client *elastic.Client
)

func main() {
	flag.Parse()
	var err error
	// producer, err := sarama.NewAsyncProducer(strings.Split(*kafkaAddrs, ","), nil)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	client, err = elastic.NewClient(elastic.SetURL("http://127.0.0.1:9200"))
	if err != nil {
		log.Fatal(err)
	}

	// go func() {
	// 	for err := range producer.Errors() {
	// 		log.Print(err)
	// 	}
	// }()

	server := NewServer(nil)
	log.Fatal(server.ListenAndServe(*listenAddr))
}
