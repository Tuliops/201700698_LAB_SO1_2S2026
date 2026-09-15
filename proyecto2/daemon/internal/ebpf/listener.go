package ebpf

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"log"

	"daemon-so1/internal/db"

	"github.com/cilium/ebpf/link"
	"github.com/cilium/ebpf/ringbuf"
	"github.com/cilium/ebpf/rlimit"
)

type KillEvent struct {
	PidSender uint32
	PidTarget uint32
	Sig       int32
}

func StartEBPFAudit(ctx context.Context, valkeyClient *db.ValkeyClient) error {
	if err := rlimit.RemoveMemlock(); err != nil {
		return fmt.Errorf("error removiendo memlock: %w", err)
	}

	objs := bpfObjects{}
	if err := loadBpfObjects(&objs, nil); err != nil {
		return fmt.Errorf("error cargando objetos BPF: %w", err)
	}
	defer objs.Close()

	// Adjuntar el tracepoint al kernel
	tp, err := link.Tracepoint("syscalls", "sys_enter_kill", objs.TraceSysKill, nil)
	if err != nil {
		return fmt.Errorf("error adjuntando tracepoint sys_enter_kill: %w", err)
	}
	defer tp.Close()

	reader, err := ringbuf.NewReader(objs.Events)
	if err != nil {
		return fmt.Errorf("error abriendo ringbuf reader: %w", err)
	}
	defer reader.Close()

	log.Println("[eBPF OK] Tracepoint sys_enter_kill activo. Escuchando eventos...")

	// Goroutine de cierre cuando se cancele el contexto
	go func() {
		<-ctx.Done()
		reader.Close()
	}()

	for {
		record, err := reader.Read()
		if err != nil {
			if errors.Is(err, ringbuf.ErrClosed) {
				log.Println("[eBPF] Ring Buffer cerrado correctamente.")
				return nil
			}
			continue
		}

		var event KillEvent
		if err := binary.Read(bytes.NewReader(record.RawSample), binary.LittleEndian, &event); err != nil {
			log.Printf("[eBPF ERROR] Error leyendo muestra: %v", err)
			continue
		}

		log.Printf("[eBPF EVENTO CAPTURADO] Signal %d enviada por PID %d -> Target PID %d",
			event.Sig, event.PidSender, event.PidTarget)

		if valkeyClient != nil {
			if err := valkeyClient.RecordEBPFKillEvent(ctx, event.PidTarget, event.Sig); err != nil {
				log.Printf("[VALKEY ERROR] Fallo al guardar en Valkey: %v", err)
			}
		}
	}
}
