package main

import (
	"testing"

	"github.com/justintout/cooklang-go"
)

func TestResolve(t *testing.T) {
	tests := []struct {
		version string
		src     string
		want    cooklang.Spec
	}{
		{"5", "Add @salt.", cooklang.SpecV5},
		{"6", "Add @salt.", cooklang.SpecV5},
		{"7", "Add @salt.", cooklang.SpecV7},
		{"latest", "Add @salt.", cooklang.SpecV7},
		{"AUTO", ">> sourced: babooshka\nAdd @salt.", cooklang.SpecV5},
		{"auto", "---\nsourced: babooshka\n---\n\nAdd @salt.", cooklang.SpecV7},
	}
	for _, tt := range tests {
		got, err := resolve(tt.version, tt.src)
		if err != nil {
			t.Errorf("resolve(%q) returned an error: %v", tt.version, err)
			continue
		}
		if got != tt.want {
			t.Errorf("resolve(%q) = %d, want %d", tt.version, got, tt.want)
		}
	}
}

func TestResolveRejectsUnknownVersion(t *testing.T) {
	if _, err := resolve("8", "Add @salt."); err == nil {
		t.Error("expected an error for an unknown spec version")
	}
}
