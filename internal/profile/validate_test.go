package profile

import (
	"encoding/json"
	"os"
	"testing"
)

func TestValidateProfileOK(t *testing.T) {
	prof := SeccompProfile{
		DefaultAction: "SCMP_ACT_ERRNO",
		Architectures: []string{"SCMP_ARCH_X86_64"},
		Syscalls: []Syscall{{
			Names:  []string{"openat", "close"},
			Action: "SCMP_ACT_ALLOW",
		}},
	}
	data, _ := json.Marshal(prof)
	path := t.TempDir() + "/profile.json"
	if err := writeTestFile(path, data); err != nil {
		t.Fatal(err)
	}
	if err := ValidateProfile(path); err != nil {
		t.Fatalf("valid profile rejected: %v", err)
	}
}

func TestValidateProfileMissingAction(t *testing.T) {
	prof := SeccompProfile{DefaultAction: "SCMP_ACT_ERRNO"}
	data, _ := json.Marshal(prof)
	path := t.TempDir() + "/profile.json"
	if err := writeTestFile(path, data); err != nil {
		t.Fatal(err)
	}
	if err := ValidateProfile(path); err == nil {
		t.Fatal("expected error for empty syscalls")
	}
}

func TestAllowedSyscalls(t *testing.T) {
	prof := SeccompProfile{
		DefaultAction: "SCMP_ACT_ERRNO",
		Syscalls: []Syscall{
			{Names: []string{"openat"}, Action: "SCMP_ACT_ALLOW"},
			{Names: []string{"execve"}, Action: "SCMP_ACT_KILL"},
		},
	}
	data, _ := json.Marshal(prof)
	path := t.TempDir() + "/profile.json"
	if err := writeTestFile(path, data); err != nil {
		t.Fatal(err)
	}
	got, err := AllowedSyscalls(path)
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != 1 || got[0] != "openat" {
		t.Fatalf("got %v, want [openat]", got)
	}
}

func writeTestFile(path string, data []byte) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	_, err = f.Write(data)
	f.Close()
	return err
}
