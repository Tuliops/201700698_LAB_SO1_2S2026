package runner

import (
	"context"
	"fmt"
	"time"

	"daemon-so1/internal/db"
)

// StartDaemon inicia el ciclo de monitoreo periodico del sistema y contenedores
func StartDaemon(ctx context.Context, valkeyClient *db.ValkeyClient) {
	ticker := time.NewTicker(20 * time.Second)
	defer ticker.Stop()

	fmt.Println("[+] Daemon iniciado. Polling cada 20s. Umbral RAM: 15.0 MB")

	for {
		select {
		case <-ctx.Done():
			fmt.Println("[+] Deteniendo ciclo del Daemon...")
			return
		case <-ticker.C:
			// Logica de recoleccion desde /proc/continfo_pr2_so1_201700698 y evaluacion de umbrales
			fmt.Println("[INFO] Ejecutando ciclo de recoleccion de metricas...")
		}
	}
}
