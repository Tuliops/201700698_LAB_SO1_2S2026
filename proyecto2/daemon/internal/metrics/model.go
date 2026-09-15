package metrics

type SystemMetrics struct {
	TotalRAMMB uint64          `json:"total_ram_mb"`
	FreeRAMMB  uint64          `json:"free_ram_mb"`
	UsedRAMMB  uint64          `json:"used_ram_mb"`
	Processes  []ProcessMetric `json:"processes"`
}

type ProcessMetric struct {
	PID        int    `json:"pid"`
	Name       string `json:"name"`
	CPUJiffies uint64 `json:"cpu_jiffies"`
	VSZKB      uint64 `json:"vsz_kb"`
	RSSKB      uint64 `json:"rss_kb"`
}

func (p *ProcessMetric) GetRSSMB() float64 {
	return float64(p.RSSKB) / 1024.0
}
