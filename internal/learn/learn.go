package learn

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/ValentinTorassa/autoconfine/internal/bpf"
	"github.com/ValentinTorassa/autoconfine/internal/models"
)

// Config define parámetros de la fase de aprendizaje.
type Config struct {
	Image    string
	Duration time.Duration
	Output   string
}

// Tracer encapsula la observación de syscalls.
type Tracer struct {
	cfg Config
	bpf bpf.Probe
}

// NewTracer crea un nuevo tracer.
func NewTracer(cfg Config) *Tracer {
	return &Tracer{
		cfg: cfg,
		bpf: bpf.NewNoopProbe(), // cambiar por NewEBPFProbe en entornos con kernel/BTF.
	}
}

// Run ejecuta el aprendizaje y persiste la traza.
func (t *Tracer) Run() error {
	if t.cfg.Image == "" {
		return fmt.Errorf("se requiere --image")
	}

	raised := make(chan error, 1)
	events := make(chan models.SyscallEvent, 1024)

	go func() {
		raised <- t.bpf.Attach(t.cfg.Image, events)
	}()

	go func() {
		time.Sleep(t.cfg.Duration)
		t.bpf.Detach()
	}()

	file, err := os.Create(t.cfg.Output)
	if err != nil {
		return fmt.Errorf("crear traza: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	for evt := range events {
		if err := enc.Encode(evt); err != nil {
			return fmt.Errorf("escribir evento: %w", err)
		}
	}

	if err := <-raised; err != nil {
		return err
	}
	return nil
}
