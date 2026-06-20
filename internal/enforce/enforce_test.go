package enforce

import (
	"os"
	"path/filepath"
	"testing"
)

func TestLoadProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "profile.json")
	content := `{"defaultAction":"SCMP_ACT_ERRNO","architectures":["SCMP_ARCH_X86_64"]}`
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		t.Fatal(err)
	}

	profile, err := LoadProfile(path)
	if err != nil {
		t.Fatal(err)
	}
	if profile["defaultAction"] != "SCMP_ACT_ERRNO" {
		t.Fatalf("defaultAction = %v", profile["defaultAction"])
	}
}

func TestRunnerMissingProfile(t *testing.T) {
	r := NewRunner(Config{ProfilePath: "/no/existe.json"})
	err := r.Run([]string{"podman", "run", "nginx"})
	if err == nil {
		t.Fatal("se esperaba error por perfil inexistente")
	}
}
