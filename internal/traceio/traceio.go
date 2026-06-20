package traceio

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/ValentinTorassa/autoconfine/internal/models"
)

// ReadSyscalls lee una traza JSONL y devuelve el conjunto de syscalls únicas.
func ReadSyscalls(path string) (map[string]int, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("abrir traza: %w", err)
	}
	defer file.Close()
	return scanSyscalls(file)
}

// scanSyscalls cuenta syscalls desde cualquier reader.
func scanSyscalls(r io.Reader) (map[string]int, error) {
	counts := make(map[string]int)
	sc := bufio.NewScanner(r)
	var raw map[string]interface{}
	for sc.Scan() {
		if err := json.Unmarshal(sc.Bytes(), &raw); err != nil {
			continue
		}
		if name, ok := raw["syscall"].(string); ok {
			counts[name]++
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("leer traza: %w", err)
	}
	return counts, nil
}

// ReadEvents lee todos los eventos de una traza.
func ReadEvents(path string) ([]models.SyscallEvent, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("abrir traza: %w", err)
	}
	defer file.Close()

	var events []models.SyscallEvent
	sc := bufio.NewScanner(file)
	for sc.Scan() {
		var evt models.SyscallEvent
		if err := json.Unmarshal(sc.Bytes(), &evt); err != nil {
			continue
		}
		events = append(events, evt)
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("leer eventos: %w", err)
	}
	return events, nil
}

// WriteEvents escribe eventos como JSONL.
func WriteEvents(path string, events []models.SyscallEvent) error {
	file, err := os.Create(path)
	if err != nil {
		return fmt.Errorf("crear traza: %w", err)
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	for _, evt := range events {
		if err := enc.Encode(evt); err != nil {
			return fmt.Errorf("escribir evento: %w", err)
		}
	}
	return nil
}

// UniqueSyscalls devuelve la lista ordenada de syscalls distintas.
func UniqueSyscalls(counts map[string]int) []string {
	seen := make([]string, 0, len(counts))
	for name := range counts {
		seen = append(seen, name)
	}
	// Inserción simple; suficiente para perfiles pequeños.
	for i := 1; i < len(seen); i++ {
		j := i
		for j > 0 && seen[j-1] > seen[j] {
			seen[j-1], seen[j] = seen[j], seen[j-1]
			j--
		}
	}
	return seen
}
