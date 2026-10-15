package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

type WarReport struct {
	Country         int    `json:"country"`
	WarplanesInAir  int    `json:"warplanes_in_air"`
	WarshipsInWater int    `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

// Mapa opcional para traducir el número del enum a texto legible en los logs
var countryMap = map[int]string{
	1: "USA",
	2: "RUS",
	3: "CHN",
	4: "ESP",
	5: "GTM",
}

var ctx = context.Background()

func main() {

	// Configuracion y conexion en Valkey / Redis

	valkeyURL := os.Getenv("VALKEY_URL")

	if valkeyURL == "" {
		valkeyURL = "localhost:6379"
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     valkeyURL,
		Password: "",
		DB:       0,
	})

	//pruebas de conexion de valkey
	_, err := rdb.Ping(ctx).Result()

	if err != nil {
		log.Fatalf("[ERROR] Conexion con Valkey: %v", err)
	}

	fmt.Println("[Valkey] Conexion Exitosa con la base de datos ")

	// Leer la URL de RabbitMQ desde variable de entorno o usar el valor por defecto local
	rabbitMQURL := os.Getenv("RABBITMQ_URL")
	if rabbitMQURL == "" {
		rabbitMQURL = "amqp://guest:guest@localhost:5672/"
	}

	// Conectar a RabbintMQ -- tiempo por si levanta antes que el broker

	var conn *amqp.Connection

	// bucle para intentos de conectarse
	for i := 0; i < 5; i++ {

		conn, err = amqp.Dial(rabbitMQURL)

		if err != nil {
			break
		}
		log.Printf("[RabbitMQ]  %d: Esperando a que RabbitMQ esté disponible...", i)
		time.Sleep(3 * time.Second)

	}
	//manejo de errores
	if err != nil {
		log.Fatalf("[ERROR] No se pudo conectar a RabbitMQ: %v", err)
	}
	defer conn.Close()

	// canal de RabbitMQ
	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[Error] al abrir canal en RabbitMQ: %v", err)
	}
	defer ch.Close()

	//Declaracion de  la cola
	queueName := "war_reports_queue"
	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Error] al declarar la cola: %v", err)
	}

	//Registro del Consumidor
	msgs, err := ch.Consume(queueName, "", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Error] al registrar el consumidor: %v", err)
	}

	fmt.Printf("[Go Consumer] Escuchando mensajes en la cola: %s...\n", queueName)

	forever := make(chan bool)

	go func() {
		for d := range msgs {
			var report WarReport
			// deserializacion del Mensaje
			err := json.Unmarshal(d.Body, &report)
			if err != nil {
				log.Printf(" [Error] al deserializar mensaje JSON: %v", err)
				continue
			}
			countryName := countryMap[report.Country]
			if countryName == "" {
				countryName = "UNKNOWN"
			}

			//valkey

			// Guardamos el reporte serializado en una lista o como clave individual
			reportJSON, _ := json.Marshal(report)
			reportKey := fmt.Sprintf("report:%d:%d", report.Country, time.Now().UnixNano())

			// Guardar el reporte en Valkey con un TTL opcional para evitar saturación de memoria
			err = rdb.Set(ctx, reportKey, reportJSON, 24*time.Hour).Err()
			if err != nil {
				log.Printf(" [Error] al guardar en Valkey: %v", err)
			}

			// gregarlo a una lista general o acumuladores para los topes
			rdb.LPush(ctx, "reports_list", reportJSON)

			// Mostrar el reporte consumido con éxito
			fmt.Println("--------------------------------------------------")
			fmt.Printf("[Consumer] ¡Reporte extraído de la cola!\n")
			fmt.Printf("   País: %s \n", countryName)
			fmt.Printf("   Aviones en aire: %d\n", report.WarplanesInAir)
			fmt.Printf("   Barcos en agua: %d\n", report.WarshipsInWater)
			fmt.Printf("   Timestamp: %s\n", report.Timestamp)
			fmt.Println("--------------------------------------------------")

		}
	}()

	<-forever

}
