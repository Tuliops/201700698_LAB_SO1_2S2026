package db

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"daemon-so1/internal/metrics"

	"github.com/valkey-io/valkey-go"
)

type ValkeyClient struct {
	client valkey.Client
}

func NewValkeyClient(addr string) (*ValkeyClient, error) {
	client, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{addr},
	})
	if err != nil {
		return nil, fmt.Errorf("error al conectar con Valkey: %w", err)
	}
	return &ValkeyClient{client: client}, nil
}

func (v *ValkeyClient) SaveMetrics(ctx context.Context, processes []metrics.ProcessMetric) error {
	data, err := json.Marshal(processes)
	if err != nil {
		return fmt.Errorf("error al serializar metricas para valkey: %w", err)
	}

	key := fmt.Sprintf("metrics:%d", time.Now().Unix())

	// CORRECCION: Se pasa time.Hour en lugar del entero 3600
	err = v.client.Do(ctx, v.client.B().Set().Key(key).Value(string(data)).Ex(time.Hour).Build()).Error()
	if err != nil {
		return fmt.Errorf("error al escribir en Valkey: %w", err)
	}

	fmt.Printf("[VALKEY] Metricas guardadas exitosamente bajo la clave: %s\n", key)
	return nil
}

func (v *ValkeyClient) Close() {
	if v.client != nil {
		v.client.Close()
	}
}
