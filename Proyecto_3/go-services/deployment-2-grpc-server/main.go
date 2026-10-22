package main

import (
	"bytes"
	"context"
	"encoding/json"
	"log"
	"net"
	"net/http"

	pb "d2grpc/proto"

	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedWarReportServiceServer
}

func (s *server) SendReport(ctx context.Context, req *pb.WarReportRequest) (*pb.WarReportResponse, error) {
	log.Printf("[D2 gRPC] Recibido: country=%v, planes=%d\n", req.Country, req.WarplanesInAir)

	msg := map[string]interface{}{
		"country":           req.Country.String(),
		"warplanes_in_air":  req.WarplanesInAir,
		"warships_in_water": req.WarshipsInWater,
		"timestamp":         req.Timestamp,
	}
	body, _ := json.Marshal(msg)

	// Enviar al Sidecar (contenedor 2) por localhost
	resp, err := http.Post("http://localhost:8082/publicar", "application/json", bytes.NewBuffer(body))
	if err != nil {
		log.Printf("[D2 gRPC] Error contactando al sidecar: %v\n", err)
		return &pb.WarReportResponse{Status: "error"}, nil
	}
	defer resp.Body.Close()

	return &pb.WarReportResponse{Status: "OK"}, nil
}

func main() {
	lis, _ := net.Listen("tcp", ":50051")
	grpcServer := grpc.NewServer()
	pb.RegisterWarReportServiceServer(grpcServer, &server{})
	log.Println("[D2 gRPC] Escuchando en :50051")
	grpcServer.Serve(lis)
}
