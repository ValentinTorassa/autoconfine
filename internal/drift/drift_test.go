package drift

import (
	"bytes"
	"encoding/json"
	"testing"
)

func TestProfileSet(t *testing.T) {
	set := ProfileSet([]string{"openat", "close"})
	if !Allowed(set, "openat") {
		t.Error("openat debería estar permitida")
	}
	if Allowed(set, "execve") {
		t.Error("execve no debería estar permitida")
	}
}

func TestJSONReporter(t *testing.T) {
	var buf bytes.Buffer
	rep := NewJSONReporter(&buf)
	if err := rep.Report(Event{Syscall: "openat"}); err != nil {
		t.Fatal(err)
	}

	var evt Event
	if err := json.Unmarshal(buf.Bytes(), &evt); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if evt.Syscall != "openat" {
		t.Fatalf("syscall = %s, want openat", evt.Syscall)
	}
	if evt.Timestamp.IsZero() {
		t.Error("timestamp no asignado")
	}
}
