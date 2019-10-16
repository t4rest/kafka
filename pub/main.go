package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"

	"github.com/sirupsen/logrus"
)

var producer *kafkaPub

func main() {
	logrus.Info("main pub")
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Congratulations! Your Go application has been successfully deployed on Kubernetes.")
		})

		http.HandleFunc("/pub", func(w http.ResponseWriter, r *http.Request) {
			var bullets = 2500

			for b := 0; b < bullets; b++ {
				err := producer.Publish(testSarama{Value: b, Error: ""})
				if err != nil {
					log.Fatal(err)
				}
				logrus.Println("Publish ok - ", b)
			}

			_, _ = fmt.Fprintf(w, "Congratulations! Your Go application has been successfully published.")
		})

		logrus.Info("ListenAndServe:3000")
		_ = http.ListenAndServe(":3000", nil)
	}()

	var err error
	producer, err = New()
	if err != nil {
		log.Fatal(err)
	}
	defer producer.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
