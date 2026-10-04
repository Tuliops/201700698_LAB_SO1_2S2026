package db

import (
	"context"
	"encoding/json"

	"daemon-so1/internal/metrics"

	"github.com/valkey-io/valkey-go"
)

type ValkeyClient struct {
	Client valkey.Client
}

func NewValkeyClient(addr string) (*ValkeyClient, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
	})
	if err != nil {
		return nil, err
	}
	return &ValkeyClient{Client: client}, nil
}

func (v *ValkeyClient) Close() {
	if v.Client != nil {
		v.Client.Close()
	}
}

func (v *ValkeyClient) SaveMetricsSnapshot(ctx context.Context, m *metrics.SystemMetrics) error {
	data, err := json.Marshal(m)
	if err != nil {
		return err
	}

	cmd := v.Client.B().Set().Key("system_metrics").Value(string(data)).Build()
	return v.Client.Do(ctx, cmd).Error()
}

// RecordEBPFKillEvent registra la auditoría de un evento SIGKILL interceptado por eBPF
func (v *ValkeyClient) RecordEBPFKillEvent(ctx context.Context, pid uint32, sig int32) error {
	// Incrementa atómicamente el contador global de eliminaciones para Grafana
	cmd := v.Client.B().Incr().Key("ebpf_kill_total_count").Build()
	return v.Client.Do(ctx, cmd).Error()
}
