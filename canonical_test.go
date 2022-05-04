package cooklang_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/justintout/cooklang-go"
	"gopkg.in/yaml.v3"
)

// https://raw.githubusercontent.com/cooklang/spec/main/tests/canonical.yaml
// version: 5

type CanonicalTests struct {
	Version string
	Tests   map[string]struct {
		Source string
		Result struct {
			Steps    [][]cooklang.CanonicalDirectionItem
			Metadata map[string]interface{}
		}
	}
}

func TestCanonical(t *testing.T) {
	t.Skip("this doesn't work yet.")
	cy, err := os.ReadFile("canonical.yaml")
	if err != nil {
		t.Fatalf("failed to read canonical.yaml: %v", err)
	}
	var ct CanonicalTests
	if err := yaml.Unmarshal([]byte(cy), &ct); err != nil {
		t.Fatalf("failed to unmarshal canonical.yaml: %v", err)
	}

	for i, tt := range ct.Tests {
		fmt.Println(tt.Source)
		r, err := cooklang.Parse(tt.Source)
		if err != nil {
			t.Errorf("error parsing canonical test %q: %v", i, err)
		}
		y, err := yaml.Marshal(&r)
		fmt.Printf("- %s\n", tt.Source)
		w, _ := yaml.Marshal(tt.Result)
		fmt.Printf("%s\n", string(y))
		fmt.Printf("\twant: %s\n", string(w))
	}
}
