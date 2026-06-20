package profile

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/ValentinTorassa/autoconfine/internal/traceio"
)

// Stats resume el análisis de reducción de syscalls.
type Stats struct {
	DefaultAllowed   int      `json:"default_allowed"`
	GeneratedAllowed int      `json:"generated_allowed"`
	ReductionPercent float64  `json:"reduction_percent"`
	UniqueSyscalls   []string `json:"unique_syscalls"`
}

// AnalyzeTrace compara syscalls de una traza contra un perfil por defecto.
func AnalyzeTrace(tracePath string, defaultAllowed int) (*Stats, error) {
	counts, err := traceio.ReadSyscalls(tracePath)
	if err != nil {
		return nil, fmt.Errorf("analyze trace: %w", err)
	}
	unique := traceio.UniqueSyscalls(counts)
	reduction := 0.0
	if defaultAllowed > 0 {
		reduction = float64(defaultAllowed-len(unique)) / float64(defaultAllowed) * 100
	}
	return &Stats{
		DefaultAllowed:   defaultAllowed,
		GeneratedAllowed: len(unique),
		ReductionPercent: reduction,
		UniqueSyscalls:   unique,
	}, nil
}

// WriteReport genera un reporte markdown de análisis.
func WriteReport(stats *Stats, tracePath, outPath string) error {
	file, err := os.Create(outPath)
	if err != nil {
		return fmt.Errorf("crear reporte: %w", err)
	}
	defer file.Close()

	fmt.Fprintf(file, "# Reporte de análisis - AutoConfine\n\n")
	fmt.Fprintf(file, "- **Traza analizada:** `%s`\n", tracePath)
	fmt.Fprintf(file, "- **Syscalls permitidas por defecto:** %d\n", stats.DefaultAllowed)
	fmt.Fprintf(file, "- **Syscalls en perfil generado:** %d\n", stats.GeneratedAllowed)
	fmt.Fprintf(file, "- **Reducción:** %.2f%%\n\n", stats.ReductionPercent)
	fmt.Fprintf(file, "## Syscalls permitidas\n\n")
	for _, s := range stats.UniqueSyscalls {
		fmt.Fprintf(file, "- `%s`\n", s)
	}
	return nil
}

// WriteReportJSON escribe el análisis como JSON.
func WriteReportJSON(stats *Stats, outPath string) error {
	data, err := json.MarshalIndent(stats, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(outPath, data, 0644)
}
