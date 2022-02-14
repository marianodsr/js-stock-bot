package main

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/streadway/amqp"
)

const PORT = 9000

const QUEUE_NAME = "stock_prices"

func main() {
	r := chi.NewRouter()

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	conn, err := amqp.Dial("amqp://guest:guest@rabbitMQ-container:5672/")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	chann, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	defer chann.Close()

	_, err = chann.QueueDeclare(
		QUEUE_NAME,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		panic(err)
	}

	server := NewServer(r, chann)

	server.ServeHTTP()

	go http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", PORT), r)
	fmt.Println("listening on port ", PORT)

	for {
	}

}
