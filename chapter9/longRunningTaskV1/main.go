package main

import (
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/streadway/amqp"
)

const queueName string = "jobQueue"
const hostString string = "127.0.0.1:8000"

func handleError(err error, message string) {
	if err != nil {
		log.Fatalf("%s: %s", message, err)
	}
}

func getServer(name string) JobServer {
	/*
		Creates a server object and initiates
		the channel and Queue details to publish messages
	*/
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	handleError(err, "Dialing failed to RabbitMQ broker")

	channel, err := conn.Channel()
	handleError(err, "Fztching channel failed")

	jobQueue, err := channel.QueueDeclare(
		name,
		false,
		false,
		false,
		false,
		nil,
	)
	handleError(err, "Fetching channel failed")
	return JobServer{Conn: conn, Channel: channel, Queue: jobQueue}
}

func main() {
	jobServer := getServer(queueName)

	// Start Workers
	go func(conn *amqp.Connection) {
		workerProcess := Workers{
			conn: jobServer.Conn,
		}
		workerProcess.run()
	}(jobServer.Conn)

	router := mux.NewRouter()
	// Attach handlers
	router.HandleFunc("/job/database", jobServer.asyncDBHandler)
	router.HandleFunc("/job/mail", jobServer.asyncMailHandler)
	router.HandleFunc("/job/callback", jobServer.asyncCallbackHandler)

	httpServer := &http.Server{
		Handler:      router,
		Addr:         hostString,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	// Run HTTP server
	log.Fatal(httpServer.ListenAndServe())

	// Cleanup resources
	defer jobServer.Channel.Close()
	defer jobServer.Conn.Close()
}
