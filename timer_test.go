package cooklang

import (
	"io"
	"os"
	"testing"
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

// TestNewTimerDegradesQuietly checks that a malformed quantity does not panic
// and, per the library contract, does not print to stdout.
func TestNewTimerDegradesQuietly(t *testing.T) {
	orig := os.Stdout
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("failed to create pipe: %v", err)
	}
	os.Stdout = w

	// no "%unit" divider, so strictParseQuantity returns an error
	timer := NewTimer("~steep{5}")

	w.Close()
	os.Stdout = orig
	out, _ := io.ReadAll(r)

	if len(out) != 0 {
		t.Errorf("NewTimer printed to stdout: %q", out)
	}
	if timer.Name != "steep" {
		t.Errorf("wrong name: got: %q, want: %q", timer.Name, "steep")
	}
	if timer.S != "5" {
		t.Errorf("expected lenient parse to recover quantity, got: %q", timer.S)
	}
}
