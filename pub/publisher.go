package main

import (
	"encoding/json"
	"os"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kafkaPub struct {
	producer sarama.SyncProducer
}

type testSarama struct {
	Error string `json:"error"`
	Value int    `json:"value"`
}

// New creates new kafka connection
func New() (*kafkaPub, error) {

	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // Wait for all in-sync replicas to ack the message
	config.Producer.Retry.Max = 10                   // Retry up to 10 times to produce the message
	config.Producer.Return.Successes = true

	kfkHostPort := "127.0.0.1:9092"
	if os.Getenv("KFK_HOST_PORT") != "" {
		kfkHostPort = os.Getenv("KFK_HOST_PORT")
	}
	logrus.Info("kfkHostPort: ", kfkHostPort)

	producer, err := sarama.NewSyncProducer([]string{kfkHostPort}, config)
	if err != nil {
		return nil, err
	}

	return &kafkaPub{producer: producer}, nil
}

// Publish publish event
func (kfk *kafkaPub) Publish(tx testSarama) error {

	data, _ := json.Marshal(tx)

	_, _, err := kfk.producer.SendMessage(&sarama.ProducerMessage{
		Topic: "partitions-test-one",
		Value: sarama.ByteEncoder(data),
		Key:   sarama.StringEncoder(tx.Value),
	})

	return err
}

// Close close connection
func (kfk *kafkaPub) Close() {
	if err := kfk.producer.Close(); err != nil {
		logrus.Fatalf("failed to shut down data collector cleanly. Error: %s", err)
	}
}
