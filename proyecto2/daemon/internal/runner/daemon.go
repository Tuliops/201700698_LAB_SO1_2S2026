package runner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"daemon-so1/internal/db"
	"daemon-so1/internal/metrics"
)

const MemoryThresholdMB = 15.0 // Umbral de consumo en MB

// Lista blanca de procesos del sistema e infraestructura que NUNCA deben ser eliminados
var processWhitelist = map[string]bool{
	"main":            true,
	"go":              true,
	"gopls":           true,
	"dockerd":         true,
	"containerd":      true,
	"valkey-server":   true,
	"systemd":         true,
	"init":            true,
	"bash":            true,
	"MainThread":      true,
	"wslservice":      true,
	"weston":          true, // Interfaz grafica WSLg
	"Xwayland":        true, // Servidor grafico
	"networkd-dispat": true,
	"snapd":           true,
	"unattended-upgr": true,
}

func RunDaemon(ctx context.Context, interval time.Duration, valkeyClient *db.ValkeyClient) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	fmt.Printf("[+] Daemon iniciado. Polling cada %v. Umbral RAM: %.1f MB\n", interval, MemoryThresholdMB)

	collectAndEnforce(ctx, valkeyClient)

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[+] Deteniendo bucle del Daemon...")
			return
		case t := <-ticker.C:
			fmt.Printf("\n--- [Ciclo de Recoleccion y Evaluacion: %s] ---\n", t.Format("15:04:05"))
			collectAndEnforce(ctx, valkeyClient)
		}
	}
}

func collectAndEnforce(ctx context.Context, valkeyClient *db.ValkeyClient) {
	processList, err := metrics.ReadProcMetrics()
	if err != nil {
		fmt.Printf("[ERROR] %v\n", err)
		return
	}

	fmt.Printf("[INFO] Procesos recolectados: %d\n", len(processList))

	// 1. Guardar en Valkey
	if valkeyClient != nil {
		if err := valkeyClient.SaveMetrics(ctx, processList); err != nil {
			fmt.Printf("[ERROR VALKEY] %v\n", err)
		}
	}

	// 2. Evaluar reglas autónomas con Whitelist
	for _, p := range processList {
		rssMB := p.GetRSSMB()

		// Verificar si el proceso está en la lista blanca
		if processWhitelist[p.Name] {
			continue // Ignorar procesos vitales del sistema
		}

		if rssMB > MemoryThresholdMB {
			fmt.Printf("[AUTONOMO] PROCESO EN ALERTA: PID %d (%s) consume %.2f MB (Excede %.1f MB)\n",
				p.PID, p.Name, rssMB, MemoryThresholdMB)

			// Ejecutar eliminación segura
			terminateTarget(p.PID, p.Name)
		}
	}
}

func terminateTarget(pid int, name string) {
	// Si el proceso es de un contenedor generado por nuestro script de prueba
	if strings.Contains(name, "stress") || strings.Contains(name, "alpine") || strings.Contains(name, "dd") {
		fmt.Printf("[ACTION] Eliminando contenedor/proceso de prueba PID: %d (%s)...\n", pid, name)

		cmd := exec.Command("kill", "-9", fmt.Sprintf("%d", pid))

		// USAR '_' descarta la variable de salida y evita el error de compilacion
		_, err := cmd.CombinedOutput()
		if err != nil {
			// Si falla por PID de namespace, se intenta limpiar via Docker
			fmt.Printf("[WARN] Fallo kill directo (%v). Limpiando contenedores huerfanos con Docker...\n", err)
			cleanupDockerContainers()
			return
		}
		fmt.Printf("[SUCCESS] Proceso PID %d terminado exitosamente.\n", pid)
	} else {
		fmt.Printf("[SKIP] El proceso '%s' (PID %d) excede el umbral pero no es marcado como contenedor de prueba.\n", name, pid)
	}
}

// Limpia contenedores con el prefijo stress_test_
func cleanupDockerContainers() {
	cmd := exec.Command("bash", "-c", "docker rm -f $(docker ps -q --filter name=stress_test_) 2>/dev/null")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("[DOCKER ACTION] No se encontraron contenedores para eliminar o ejecucion limpia.\n")
		return
	}
	fmt.Printf("[DOCKER ACTION] Limpieza automatica de contenedores ejecutada. Salida: %s\n", strings.TrimSpace(string(output)))
}
