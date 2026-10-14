package main

import (
	"context"
	"fmt"
	"log"
	"net"

	// importar para rabbitmq
	amqp "github.com/rabbitmq/amqp091-go"
	"google.golang.org/grpc"

	// Importa el paquete generado
	pb "go_client/proto"
)

// estructrua del servidor grpc , se genera en protoc
// manejo de RabbitMQ
type Server struct {
	pb.UnimplementedWarReportServiceServer
	rabbitChannel *amqp.Channel
	queueName     string
}

//Funcion sendreport , definida en .protoc

func (s *Server) SendReport(ctx context.Context, req *pb.WarReportRequest) (*pb.WarReportResponse, error) {

	//Print reporte
	fmt.Printf("[gRPC Server] ¡Reporte recibido!\n")
	fmt.Printf("  -> País: %v\n", req.GetCountry())
	fmt.Printf("  -> Aviones en aire: %d\n", req.GetWarplanesInAir())
	fmt.Printf("  -> Barcos en agua: %d\n", req.GetWarshipsInWater())
	fmt.Printf("  -> Timestamp: %s\n", req.GetTimestamp())

	//Mensaje para enviar a la Cola  RabbitMQ
	messageBody := fmt.Sprintf(
		`{"country": %d, "warplanes_in_air": %d, "warships_in_water": %d, "timestamp": "%s"}`,
		req.Country, req.WarplanesInAir, req.WarshipsInWater, req.Timestamp,
	)

	// RabbitMQ

	// Publicamos el mensaje en la cola de RabbitMQ (
	err := s.rabbitChannel.PublishWithContext(
		ctx,
		"",          // exchange por defecto
		s.queueName, // nombre de la cola
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        []byte(messageBody),
		},
	)

	if err != nil {
		log.Printf(" Error al publicar mensaje en RabbitMQ: %v", err)
		return &pb.WarReportResponse{Status: "Error - No se pudo encolar"}, err
	}

	log.Println(" [RabbitMQ] Reporte encolado exitosamente en la cola:", s.queueName)

	return &pb.WarReportResponse{
		Status: "Success - Recibido en Go gRPC y encolado en RabbitMQ",
	}, nil
}

func main() {

	// Conexion con RabbitMQ
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		log.Fatalf("No se pudo conectar con RabbitMQ: %v", err)
	}

	defer conn.Close()

	ch, err := conn.Channel()
	if err != nil {
		log.Fatalf("No se pudo abrir el canal en RabbitMQ: %v", err)
	}
	defer ch.Close()

	//Se declara o crea la cola

	queueName := "war_reports_queue"

	_, err = ch.QueueDeclare(
		queueName, // nombre
		true,      // durable (sobrevive a reinicios del broker)
		false,     // auto-delete
		false,     // exclusive
		false,     // no-wait
		nil,       // argumentos
	)

	if err != nil {
		log.Fatalf("No se pudo declarar la cola en RabbitMQ %v", err)

	}

	log.Println("Conectado a RabbitMQ con exito. Cola lista : ", queueName)

	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("[Error] al escuchar el puerto : %v", err)

	}
	grpcServer := grpc.NewServer()

	//Registro del SErvicio por el canal y la cola
	pb.RegisterWarReportServiceServer(grpcServer, &Server{
		rabbitChannel: ch,
		queueName:     queueName,
	})

	log.Println("[ok] Servidor gRPC en go Corriendo y escuchando en el puerto : 50051 ")

	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("[Error] en el servidor gRPC: %v", err)
	}

}
