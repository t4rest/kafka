package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	//"github.com/pkg/profile"
	"github.com/sirupsen/logrus"
)

func main() {
	//defer profile.Start(profile.MemProfile, profile.ProfilePath(".")).Stop()

	logrus.Info("main sub")
	go func() {
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			_, _ = fmt.Fprintf(w, "Congratulations! Your Go application has been successfully deployed on Kubernetes.")
		})

		logrus.Info("ListenAndServe:3001")
		_ = http.ListenAndServe(":3001", nil)
	}()

	go func() {
		kfkHostPort := "127.0.0.1:9092"
		if os.Getenv("KFK_HOST_PORT") != "" {
			kfkHostPort = os.Getenv("KFK_HOST_PORT")
		}
		logrus.Info("kfkHostPort: ", kfkHostPort)
		sub, err := newKafkaConsumer("saramaGroupTest", []string{kfkHostPort}, []string{"partitions-test-one"}, &Handler{})
		if err != nil {
			log.Fatal(err)
		}
		defer sub.Stop()
		sub.Start()
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	<-c
}
