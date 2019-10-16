package main

import (
	"encoding/json"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

// Handler .
type Handler struct {
	total int
	m     sync.Mutex
	ready chan bool
}

type testSarama struct {
	Error string `json:"error"`
	Value int    `json:"value"`
}

func (h *Handler) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(h.ready)
	return nil
}
func (h *Handler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

func (h *Handler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {

	for {
		select {
		case msg, ok := <-claim.Messages():
			if !ok {
				return nil
			}

			logrus.Println("msg.Offset - ", msg.Offset)

			func() {
				tx := testSarama{}
				_ = json.Unmarshal(msg.Value, &tx)

				err := h.processTx(tx)
				if err == nil {
					session.MarkMessage(msg, "")
				}
			}()

		case <-session.Context().Done():
			return nil
		}
	}
}

func (h *Handler) processTx(tx testSarama) error {
	h.m.Lock()
	defer h.m.Unlock()
	h.total++

	logrus.Println("processTx ", tx.Value, " error ", tx.Error, " total ", h.total)
	//time.Sleep(1 * time.Second)

	if tx.Error != "" {
		return fmt.Errorf(tx.Error)
	}

	return nil
}
