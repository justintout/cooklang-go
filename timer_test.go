package cooklang

import (
	"io"
	"os"
	"testing"

	"go.yaml.in/yaml/v3"
)

func TestNewTimer(t *testing.T) {
	tests := []struct {
		source string
		name   string
		s      string
		units  string
	}{
		{"~{10%minutes}", "", "10", "minutes"},
		{"~potato{42%minutes}", "potato", "42", "minutes"},
		{"~rest", "rest", "", ""},
		// no "%" divider: the duration is kept as written
		{"~steep{5}", "steep", "5", ""},
	}
	for _, tt := range tests {
		timer := NewTimer(tt.source)
		if timer.Name != tt.name {
			t.Errorf("wrong name for %q: got: %q, want: %q", tt.source, timer.Name, tt.name)
		}
		if timer.S != tt.s {
			t.Errorf("wrong string quantity for %q: got: %q, want: %q", tt.source, timer.S, tt.s)
		}
		if timer.Units != tt.units {
			t.Errorf("wrong units for %q: got: %q, want: %q", tt.source, timer.Units, tt.units)
		}
	}
}

// TestLibraryIsQuiet guards the contract that parsing and marshaling never
// write to stdout.
func TestLibraryIsQuiet(t *testing.T) {
	out := captureStdout(t, func() {
		NewTimer("~steep{5}")
		r := MustParse("Add @salt.\n")
		if _, err := yaml.Marshal(&r); err != nil {
			t.Errorf("failed to marshal: %v", err)
		}
	})
	if out != "" {
		t.Errorf("library wrote to stdout: %q", out)
	}
}

func captureStdout(t *testing.T, f func()) string {
	t.Helper()
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = orig
	out, _ := io.ReadAll(r)
	return string(out)
}
