package profile

import (
	"fmt"

	"github.com/ValentinTorassa/autoconfine/internal/models"
	"github.com/ValentinTorassa/autoconfine/internal/traceio"
)

// MergeTraces combina múltiples trazas JSONL en una sola, deduplicando por evento exacto.
func MergeTraces(paths []string, outPath string) error {
	seen := make(map[string]struct{})
	var merged []models.SyscallEvent

	for _, path := range paths {
		events, err := traceio.ReadEvents(path)
		if err != nil {
			return fmt.Errorf("leer %s: %w", path, err)
		}
		for _, e := range events {
			key := fmt.Sprintf("%d-%s-%s-%s", e.PID, e.Comm, e.Syscall, e.Timestamp.Format("20060102150405.000"))
			if _, ok := seen[key]; ok {
				continue
			}
			seen[key] = struct{}{}
			merged = append(merged, e)
		}
	}

	return traceio.WriteEvents(outPath, merged)
}
