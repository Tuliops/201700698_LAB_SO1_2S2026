package db

import (
	"context"
	"fmt"
	"log"

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
		return nil, fmt.Errorf("error al conectar con valkey en %s: %w", addr, err)
	}
	return &ValkeyClient{Client: client}, nil
}

// RecordEBPFKillEvent incrementa el contador global en Valkey
func (v *ValkeyClient) RecordEBPFKillEvent(ctx context.Context, pidTarget uint32, sig int32) error {
	if v == nil || v.Client == nil {
		return fmt.Errorf("cliente valkey no inicializado")
	}

	cmd := v.Client.B().Incr().Key("ebpf_kill_total_count").Build()
	err := v.Client.Do(ctx, cmd).Error()
	if err != nil {
		log.Printf("[VALKEY ERROR] Fallo al ejecutar INCR: %v", err)
		return err
	}

	log.Println("[VALKEY OK] Clave 'ebpf_kill_total_count' incrementada exitosamente.")
	return nil
}

func (v *ValkeyClient) Close() {
	if v != nil && v.Client != nil {
		v.Client.Close()
	}
}
