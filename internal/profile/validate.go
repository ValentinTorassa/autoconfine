package profile

import (
	"encoding/json"
	"fmt"
	"os"
)

// SeccompProfile es la estructura mínima validada.
type SeccompProfile struct {
	DefaultAction string    `json:"defaultAction"`
	Architectures []string  `json:"architectures"`
	Syscalls      []Syscall `json:"syscalls"`
}

// Syscall representa una regla permitida.
type Syscall struct {
	Names  []string `json:"names"`
	Action string   `json:"action"`
}

// ValidateProfile valida que un archivo sea JSON de seccomp OCI parseable.
func ValidateProfile(path string) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("leer perfil: %w", err)
	}
	var profile SeccompProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return fmt.Errorf("parsear JSON: %w", err)
	}
	if profile.DefaultAction == "" {
		return fmt.Errorf("defaultAction es obligatorio")
	}
	if len(profile.Syscalls) == 0 {
		return fmt.Errorf("el perfil no contiene syscalls")
	}
	for i, rule := range profile.Syscalls {
		if rule.Action == "" {
			return fmt.Errorf("syscalls[%d]: action es obligatoria", i)
		}
		if len(rule.Names) == 0 {
			return fmt.Errorf("syscalls[%d]: names no puede estar vacío", i)
		}
	}
	return nil
}

// AllowedSyscalls extrae la lista de syscalls permitidas de un perfil.
func AllowedSyscalls(path string) ([]string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var profile SeccompProfile
	if err := json.Unmarshal(data, &profile); err != nil {
		return nil, err
	}
	var allowed []string
	for _, rule := range profile.Syscalls {
		if rule.Action == "SCMP_ACT_ALLOW" {
			allowed = append(allowed, rule.Names...)
		}
	}
	return allowed, nil
}
