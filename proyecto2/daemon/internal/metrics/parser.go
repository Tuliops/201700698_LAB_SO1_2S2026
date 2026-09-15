package metrics

import (
	"encoding/json"
	"fmt"
	"os"
)

const ProcPath = "/proc/continfo_pr2_so1_201700698"

func ReadProcMetrics() (*SystemMetrics, error) {
	data, err := os.ReadFile(ProcPath)
	if err != nil {
		return nil, fmt.Errorf("error leyendo %s: %w", ProcPath, err)
	}

	var metrics SystemMetrics
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, fmt.Errorf("error des-serializando JSON de proc: %w", err)
	}

	return &metrics, nil
}
