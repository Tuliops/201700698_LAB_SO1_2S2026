package metrics

import (
	"strconv"
	"strings"
)

// ProcessMetric representa la métrica individual de cada proceso expuesta por /proc
type ProcessMetric struct {
	PID          int    `json:"pid"`
	Name         string `json:"name"`
	CPUJiffies   uint64 `json:"cpu_jiffies"`
	VSZ          string `json:"vsz"`
	RSS          string `json:"rss"`
	IOReadBytes  uint64 `json:"io_read_bytes"`
	IOWriteBytes uint64 `json:"io_write_bytes"`
}

// GetRSSMB convierte la cadena de RSS (ej: "14.19 MB" o "1060 KB") a un valor numérico float64 en Megabytes
func (p *ProcessMetric) GetRSSMB() float64 {
	parts := strings.Fields(p.RSS)
	if len(parts) < 2 {
		return 0
	}

	val, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	// Si la unidad es KB, convertir a MB dividiendo entre 1024
	if strings.ToUpper(parts[1]) == "KB" {
		return val / 1024.0
	}

	return val
}
