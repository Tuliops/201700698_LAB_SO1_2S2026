package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

	"daemon-so1/internal/db"
	"daemon-so1/internal/ebpf"
	"daemon-so1/internal/kernel"
	"daemon-so1/internal/runner"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("   DAEMON DE MONITOREO SO1 - (GO & VALKEY)")
	fmt.Println("==================================================")

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// 1. Verificación de módulo Kernel
	if err := kernel.CheckAndLoadModule(); err != nil {
		log.Printf("[WARN] Error con el modulo del kernel: %v", err)
	}

	// 2. Conexión con Valkey
	valkeyHost := os.Getenv("VALKEY_HOST")
	if valkeyHost == "" {
		valkeyHost = "localhost"
	}
	valkeyPort := os.Getenv("VALKEY_PORT")
	if valkeyPort == "" {
		valkeyPort = "6379"
	}

	valkeyClient, err := db.NewValkeyClient(fmt.Sprintf("%s:%s", valkeyHost, valkeyPort))
	if err != nil {
		log.Printf("[WARN] No se pudo conectar con Valkey: %v", err)
	} else {
		fmt.Println("[+] Conexion con Valkey establecida correctamente.")
	}

	// 3. Iniciar eBPF Audit
	if err := ebpf.StartEBPFAudit(ctx, valkeyClient); err != nil {
		log.Printf("[eBPF WARN] No se pudo iniciar el listener: %v", err)
	}

	// 4. Iniciar Runner del Daemon en Goroutine
	go runner.StartDaemon(ctx, valkeyClient)

	// Graceful Shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	<-sigChan
	fmt.Println("\n[+] Finalizando servicio Daemon...")
	cancel()
	if valkeyClient != nil {
		valkeyClient.Close()
	}
}
