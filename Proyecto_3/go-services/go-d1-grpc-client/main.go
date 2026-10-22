package main

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"time"

	pb "client/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ReportPayload struct {
	Country         string `json:"country"`
	WarplanesInAir  int32  `json:"warplanes_in_air"`
	WarshipsInWater int32  `json:"warships_in_water"`
	Timestamp       string `json:"timestamp"`
}

func stringToCountry(c string) pb.Countries {
	switch c {
	case "usa":
		return pb.Countries_usa
	case "rus":
		return pb.Countries_rus
	case "chn":
		return pb.Countries_chn
	case "esp":
		return pb.Countries_esp
	case "gtm":
		return pb.Countries_gtm
	default:
		return pb.Countries_countries_unknown
	}
}

func main() {
	conn, _ := grpc.NewClient("go-d2-service:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	defer conn.Close()
	client := pb.NewWarReportServiceClient(conn)

	http.HandleFunc("/procesar", func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		var payload ReportPayload
		json.Unmarshal(body, &payload)

		req := &pb.WarReportRequest{
			Country:         stringToCountry(payload.Country),
			WarplanesInAir:  payload.WarplanesInAir,
			WarshipsInWater: payload.WarshipsInWater,
			Timestamp:       payload.Timestamp,
		}

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		res, err := client.SendReport(ctx, req)
		if err != nil {
			log.Printf("[D1-gRPC] Error: %v\n", err)
			http.Error(w, "Error gRPC", 500)
			return
		}
		log.Printf("[D1-gRPC] Go D2 respondió: %s\n", res.Status)
		w.WriteHeader(200)
	})

	log.Println("[D1-gRPC] Escuchando en :8083")
	http.ListenAndServe(":8083", nil)
}
