package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"daemon-so1/internal/db"
	"daemon-so1/internal/ebpf"
	"daemon-so1/internal/metrics"
)

const (
	// Ruta del módulo del Kernel por carnet
	ProcFilePath = "/proc/continfo_pr2_so1_201700698"

	// Intervalo de ejecución del bucle
	IntervaloLoop = 30 * time.Second

	// Umbral de memoria RSS en KB (15 MB * 1024 = 15360 KB)
	MemoryThresholdKB = 15 * 1024
)

func main() {
	log.Println("==================================================")
	log.Println("   DAEMON DE MONITOREO SO1 - (GO & VALKEY)")
	log.Println("==================================================")

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 1. Conexión a Valkey usando dirección concatenada ("host:port")
	valkeyHost := getEnv("VALKEY_HOST", "127.0.0.1")
	valkeyPort := getEnv("VALKEY_PORT", "6379")
	valkeyAddr := fmt.Sprintf("%s:%s", valkeyHost, valkeyPort)

	valkeyClient, err := db.NewValkeyClient(valkeyAddr)
	if err != nil {
		log.Fatalf("[FATAL] No se pudo conectar a Valkey: %v", err)
	}
	defer valkeyClient.Close()
	log.Println("[+] Conexión con Valkey establecida correctamente.")

	// 2. Iniciar el auditor eBPF en segundo plano
	go func() {
		if err := ebpf.StartEBPFAudit(ctx, valkeyClient); err != nil {
			log.Printf("[ERROR eBPF] El auditor falló: %v", err)
		}
	}()

	// 3. Ticker del bucle principal
	ticker := time.NewTicker(IntervaloLoop)
	defer ticker.Stop()

	// Primera ejecución inmediata
	ejecutarIteracion(ctx, valkeyClient)

	for {
		select {
		case <-ctx.Done():
			log.Println("[DAEMON] Apagado controlado recibido. Cerrando...")
			return
		case <-ticker.C:
			ejecutarIteracion(ctx, valkeyClient)
		}
	}
}

func ejecutarIteracion(ctx context.Context, valkeyClient *db.ValkeyClient) {
	log.Printf("[LOOP] --- Inicio de iteración (%s) ---", time.Now().Format("15:04:05"))

	// A. Leer /proc
	rawContent, err := os.ReadFile(ProcFilePath)
	if err != nil {
		log.Printf("[ERROR PROC] No se pudo leer %s: %v", ProcFilePath, err)
		return
	}

	// B. Deserializar el contenido JSON a SystemMetrics
	metrics, err := metrics.ParseProcContent(rawContent)
	if err != nil {
		log.Printf("[ERROR PROC] Error al deserializar métricas: %v", err)
		return
	}

	// C. Análisis de procesos y remediación
	for _, process := range metrics.Processes {
		// Omitir procesos protegidos
		if process.Name == "grafana" || process.Name == "valkey_metrics" {
			continue
		}

		// Evaluar consumo de RAM RSS mayor a 15 MB
		if process.RSSKB > MemoryThresholdKB {
			log.Printf("[REMEDIACION] Proceso/Contenedor '%s' (PID %d) excede el límite con %.2f MB. Eliminando...",
				process.Name, process.PID, process.GetRSSMB())

			// Enviar SIGKILL para activar el Tracepoint de eBPF
			if err := syscall.Kill(process.PID, syscall.SIGKILL); err != nil {
				log.Printf("[ERROR KILL] Fallo al enviar SIGKILL a PID %d: %v", process.PID, err)
			}
		}
	}

	// D. Almacenar snapshot actual en Valkey
	if err := valkeyClient.SaveMetricsSnapshot(ctx, metrics); err != nil {
		log.Printf("[ERROR VALKEY] Error guardando snapshot: %v", err)
	} else {
		log.Println("[VALKEY OK] Métricas actualizadas en Valkey.")
	}
}

func getEnv(key, fallback string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return fallback
}
