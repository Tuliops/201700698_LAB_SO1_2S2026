package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/streadway/amqp"
)

type WarReport struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

func updateMinMax(ctx context.Context, rdb *redis.Client, maxKey, minKey string, val int32) {
	val64 := int64(val)
	if current, err := rdb.Get(ctx, maxKey).Int64(); err != nil || val64 > current {
		rdb.Set(ctx, maxKey, val64, 0)
	}
	if current, err := rdb.Get(ctx, minKey).Int64(); err != nil || val64 < current {
		rdb.Set(ctx, minKey, val64, 0)
	}
}

func main() {
	log.Println("[Consumer] Iniciando servicio...")

	rdb := redis.NewClient(&redis.Options{
		Addr: "my-valkey-valkey-service:6379",
	})
	if err := rdb.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("[Consumer] Error conectando a Valkey: %v", err)
	}
	log.Println("[Consumer] Conectado a Valkey exitosamente.")

	var conn *amqp.Connection
	var err error
	rabbitURL := "amqp://guest:guest@rabbitmq-service:5672/"
	for i := 0; i < 10; i++ {
		conn, err = amqp.Dial(rabbitURL)
		if err == nil {
			break
		}
		log.Printf("[Consumer] Esperando RabbitMQ (intento %d)...", i+1)
		time.Sleep(3 * time.Second)
	}
	if err != nil {
		log.Fatalf("[Consumer] Error conectando a RabbitMQ: %v", err)
	}
	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("[Consumer] Error abriendo canal: %v", err)
	}
	defer ch.Close()

	_, err = ch.QueueDeclare("war_reports", true, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Consumer] Error declarando cola: %v", err)
	}

	msgs, err := ch.Consume("war_reports", "", false, false, false, false, nil)
	if err != nil {
		log.Fatalf("[Consumer] Error al consumir: %v", err)
	}

	log.Println("[Consumer] Esperando mensajes de RabbitMQ...")

	go func() {
		sigchan := make(chan os.Signal, 1)
		signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)
		<-sigchan
		log.Println("[Consumer] Cerrando...")
		ch.Close()
		conn.Close()
		os.Exit(0)
	}()

	for d := range msgs {
		var report WarReport
		if err := json.Unmarshal(d.Body, &report); err != nil {
			log.Printf("[Consumer] Error parseando JSON: %v", err)
			d.Nack(false, false)
			continue
		}

		log.Printf("[Consumer] Mensaje recibido: %+v\n", report)

		ctx := context.Background()

		err := rdb.LPush(ctx, "war_reports_list", d.Body).Err()
		if err != nil {
			log.Printf("[Consumer] Error guardando en Valkey: %v", err)
			d.Nack(false, true)
			continue
		}

		rdb.Incr(ctx, "count:"+report.Country)

		updateMinMax(ctx, rdb, "max_planes", "min_planes", report.WarplanesInAir)
		updateMinMax(ctx, rdb, "max_ships", "min_ships", report.WarshipsInWater)

		rdb.HIncrBy(ctx, "total_planes", report.Country, int64(report.WarplanesInAir))
		rdb.HIncrBy(ctx, "total_ships", report.Country, int64(report.WarshipsInWater))

		freqPlanesKey := fmt.Sprintf("freq_planes:%d", report.WarplanesInAir)
		newCountP, _ := rdb.Incr(ctx, freqPlanesKey).Result()
		currentMaxCountP, _ := rdb.Get(ctx, "max_count_planes").Int64()
		if newCountP > currentMaxCountP {
			rdb.Set(ctx, "max_count_planes", newCountP, 0)
			rdb.Set(ctx, "moda_planes", report.WarplanesInAir, 0)
		}

		freqShipsKey := fmt.Sprintf("freq_ships:%d", report.WarshipsInWater)
		newCountS, _ := rdb.Incr(ctx, freqShipsKey).Result()
		currentMaxCountS, _ := rdb.Get(ctx, "max_count_ships").Int64()
		if newCountS > currentMaxCountS {
			rdb.Set(ctx, "max_count_ships", newCountS, 0)
			rdb.Set(ctx, "moda_ships", report.WarshipsInWater, 0)
		}

		// --- LÓGICA ESPECÍFICA PARA GTM (Llaves simples para Grafana) ---
		if report.Country == "gtm" {
			//  Opción 2: Guardar los valores actuales en llaves simples (para GET)
			rdb.Set(ctx, "gtm_total_planes", report.WarplanesInAir, 0)
			rdb.Set(ctx, "gtm_total_ships", report.WarshipsInWater, 0)

			// Mantener el Sorted Set (serie temporal )
			t, err := time.Parse(time.RFC3339, report.Timestamp)
			if err == nil {
				score := float64(t.UnixNano()) / 1e6
				member := fmt.Sprintf("%d:%d:%s", report.WarplanesInAir, report.WarshipsInWater, report.Timestamp)
				rdb.ZAdd(ctx, "timeline:gtm", redis.Z{
					Score:  score,
					Member: member,
				})
			}
		}

		d.Ack(false)
	}
}
