package main

import (
	"io"
	"log"
	"net/http"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
)

func main() {
	var ch *amqp.Channel
	var err error

	// Bucle de reintento para conectarse a RabbitMQ
	for i := 0; i < 10; i++ {
		conn, e := amqp.Dial("amqp://guest:guest@rabbitmq-service:5672/")
		if e == nil {
			ch, err = conn.Channel()
			if err == nil {
				ch.QueueDeclare("war_reports", true, false, false, false, nil)
				break
			}
		}
		log.Printf("[Sidecar] Esperando RabbitMQ...")
		time.Sleep(3 * time.Second)
	}
	if ch == nil {
		log.Fatal("[Sidecar] No se pudo conectar a RabbitMQ")
	}

	http.HandleFunc("/publicar", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		err := ch.Publish("", "war_reports", false, false, amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
		if err != nil {
			log.Printf("[Sidecar] Error publicando: %v", err)
			http.Error(w, "Error", 500)
			return
		}
		log.Printf("[Sidecar] Publicado en RabbitMQ: %s", string(body))
		w.WriteHeader(200)
	})

	log.Println("[Sidecar] Escuchando en :8082")
	http.ListenAndServe(":8082", nil)
}
