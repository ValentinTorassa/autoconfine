package bpf

import (
	"testing"
	"time"

	"github.com/ValentinTorassa/autoconfine/internal/models"
)

func TestNoopProbeEmitsEvents(t *testing.T) {
	probe := NewNoopProbe()
	events := make(chan models.SyscallEvent, 100)

	done := make(chan error, 1)
	go func() {
		done <- probe.Attach("nginx", events)
	}()

	time.Sleep(350 * time.Millisecond)
	probe.Detach()

	if err := <-done; err != nil {
		t.Fatalf("attach inesperado: %v", err)
	}

	count := 0
	for range events {
		count++
		if count > 10 {
			break
		}
	}
	if count == 0 {
		t.Fatal("no se emitieron eventos")
	}
}
