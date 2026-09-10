package metrics

import (
	"encoding/json"
	"fmt"
	"os"
)

const ProcFilePath = "/proc/continfo_pr2_so1_201700698"

// ReadProcMetrics abre el archivo /proc, lee su contenido completo
// y deserializa el arreglo JSON en una rebanada (slice) de ProcessMetric.
func ReadProcMetrics() ([]ProcessMetric, error) {
	// Comprobar la existencia del archivo antes de abrir
	if _, err := os.Stat(ProcFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("el archivo de kernel %s no existe. ¿Esta cargado el modulo?", ProcFilePath)
	}

	// Leer el archivo /proc de forma eficiente
	data, err := os.ReadFile(ProcFilePath)
	if err != nil {
		return nil, fmt.Errorf("error al leer %s: %w", ProcFilePath, err)
	}

	if len(data) == 0 {
		return nil, fmt.Errorf("el archivo %s esta vacio", ProcFilePath)
	}

	// Deserialización de la estructura JSON
	var processList []ProcessMetric
	err = json.Unmarshal(data, &processList)
	if err != nil {
		return nil, fmt.Errorf("error al parsear el JSON de /proc: %w", err)
	}

	return processList, nil
}