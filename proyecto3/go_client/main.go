package main

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"

	// Importa el paquete generado
	pb "go_client/proto"
)

// estructrua del servidor grpc , se genera en protoc

type Server struct {
	pb.UnimplementedWarReportServiceServer
}

//Funcion sendreport , definida en .protoc

func (s *Server) SendReport(ctx context.Context, req *pb.WarReportRequest) (*pb.WarReportResponse, error) {

	//Print reporte
	fmt.Printf("[gRPC Server] ¡Reporte recibido!\n")
	fmt.Printf("  -> País: %v\n", req.GetCountry())
	fmt.Printf("  -> Aviones en aire: %d\n", req.GetWarplanesInAir())
	fmt.Printf("  -> Barcos en agua: %d\n", req.GetWarshipsInWater())
	fmt.Printf("  -> Timestamp: %s\n", req.GetTimestamp())

	// Respuesta OK
	return &pb.WarReportResponse{
		Status: "Success - recibido en go grpc Server",
	}, nil

}

func main() {

	port := ":50051"
	listener, err := net.Listen("tcp", port)

	if err != nil {
		log.Fatalf("No se pudo inicial el listener tcp en el puerto %s: %v", port, err)
	}

	// Instancia del servidor grpc

	grpcServer := grpc.NewServer()

	// Registro de nuestro Servidor

	pb.RegisterWarReportServiceServer(grpcServer, &Server{})
	fmt.Printf("Servidor gRPC en Go en el Puerto %v ... \n ", port)

	// Arraque del servidor

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Falla al ejecutar el servidro gRPC:  %v", err)
	}

}
