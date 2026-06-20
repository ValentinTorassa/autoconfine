package enforce

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/ValentinTorassa/autoconfine/internal/drift"
)

// Config define parámetros de ejecución protegida.
type Config struct {
	ProfilePath   string
	Audit         bool
	DriftReporter drift.Reporter
}

// Runner ejecuta Podman con un perfil seccomp.
type Runner struct {
	cfg Config
}

// NewRunner crea un runner.
func NewRunner(cfg Config) *Runner {
	return &Runner{cfg: cfg}
}

// Run ejecuta el contenedor con el perfil aplicado.
func (r *Runner) Run(args []string) error {
	if len(args) == 0 {
		return fmt.Errorf("se requiere comando de Podman")
	}

	info, err := os.Stat(r.cfg.ProfilePath)
	if err != nil {
		return fmt.Errorf("perfil no encontrado: %w", err)
	}
	if info.IsDir() {
		return fmt.Errorf("la ruta del perfil es un directorio")
	}

	podmanArgs := append([]string{"run", "--security-opt", fmt.Sprintf("seccomp=%s", r.cfg.ProfilePath)}, args[1:]...)
	cmd := exec.Command(args[0], podmanArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if r.cfg.Audit {
		if err := r.reportAuditHeader(); err != nil {
			return err
		}
	}

	return cmd.Run()
}

func (r *Runner) reportAuditHeader() error {
	return r.cfg.DriftReporter.Report(drift.Event{Syscall: "audit_mode_enabled"})
}

// LoadProfile lee un perfil seccomp desde disco.
func LoadProfile(path string) (map[string]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var profile map[string]interface{}
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	return profile, nil
}
