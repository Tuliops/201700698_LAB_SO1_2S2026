package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"daemon-so1/internal/db"
	"daemon-so1/internal/kernel"
	"daemon-so1/internal/runner"
)

const (
	PollingInterval     = 20 * time.Second
	KernelModuleRelPath = "../kernel_module"
	ValkeyAddr          = "127.0.0.1:6379"
)

func main() {
	fmt.Println("==================================================")
	fmt.Println("    DAEMON DE MONITOREO SO1 - DIA 4 (GO & VALKEY)")
	fmt.Println("==================================================")

	execDir, err := os.Getwd()
	if err != nil {
		fmt.Printf("[FATAL] No se pudo obtener el directorio actual: %v\n", err)
		os.Exit(1)
	}

	moduleAbsolutePath := filepath.Join(execDir, KernelModuleRelPath)

	// 1. Cargar Módulo de Kernel
	if err := kernel.EnsureModuleLoaded(moduleAbsolutePath); err != nil {
		fmt.Printf("[FATAL] Fallo en la inicializacion del kernel: %v\n", err)
		os.Exit(1)
	}

	// 2. Conectar con Valkey
	valkeyClient, err := db.NewValkeyClient(ValkeyAddr)
	if err != nil {
		fmt.Printf("[WARN] No se pudo conectar con Valkey (%v). Continuando sin persistencia...\n", err)
	} else {
		defer valkeyClient.Close()
		fmt.Println("[+] Conexion con Valkey establecida correctamente.")
	}

	// 3. Contexto para Graceful Shutdown
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 4. Iniciar Bucle Autónomo
	go runner.RunDaemon(ctx, PollingInterval, valkeyClient)

	<-ctx.Done()
	fmt.Println("\n[+] Senal de apagado recibida.")

	// 5. Limpieza Final
	fmt.Println("[+] Limpiando recursos...")
	if err := kernel.UnloadModule(moduleAbsolutePath); err != nil {
		fmt.Printf("[WARN] No se pudo remover el modulo: %v\n", err)
	}

	fmt.Println("[+] Daemon finalizado correctamente.")
}
