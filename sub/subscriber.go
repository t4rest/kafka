package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kafkaConsumer struct {
	saramaConsumer sarama.ConsumerGroup
	handler        *Handler
	closeChan      chan struct{}
	consume        bool
	topics         []string
}

func (k *kafkaConsumer) Start() {
	k.consume = true
	k.closeChan = make(chan struct{})
	k.handler.ready = make(chan bool)
	var err error

	go func() {
		for err = range k.saramaConsumer.Errors() {
			logrus.Errorf("consume error: %s", err)
		}
	}()

	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}
	wg.Add(1)

	go func() {
		defer wg.Done()
		for k.consume {
			err = k.saramaConsumer.Consume(ctx, k.topics, k.handler)
			if err != nil {
				logrus.Fatalf("Error from consumer: %v", err)
			}

			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				logrus.Errorf("ctx.Err()", err)

				return
			}

			k.handler.ready = make(chan bool)
		}
	}()

	<-k.handler.ready // Await till the consumer has been set up

	select {
	case <-ctx.Done():
		logrus.Info("terminating: context cancelled")
	case <-k.closeChan:
		logrus.Info("terminating: via signal")
	}

	cancel()
	wg.Wait()

	if k.consume {
		k.Start()
	}
}

func (k *kafkaConsumer) Stop() {
	k.consume = false
	close(k.closeChan)
	_ = k.saramaConsumer.Close()
}

func newKafkaConsumer(name string, addr, topics []string, handler *Handler) (*kafkaConsumer, error) {
	//
	saramaConfig := sarama.NewConfig()
	saramaConfig.Version = sarama.V2_0_0_0
	saramaConfig.Consumer.Return.Errors = true
	saramaConfig.Consumer.Offsets.Initial = sarama.OffsetNewest

	consumer, err := sarama.NewConsumerGroup(addr, name, saramaConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create consumer group %s", err)
	}
	return &kafkaConsumer{
		saramaConsumer: consumer,
		handler:        handler,
		topics:         topics,
	}, nil
}
