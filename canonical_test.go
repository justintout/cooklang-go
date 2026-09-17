package cooklang_test

import (
	"os"
	"reflect"
	"testing"

	"github.com/justintout/cooklang-go"
	"go.yaml.in/yaml/v3"
)

type canonicalTest struct {
	Source string
	Result struct {
		Metadata map[string]string
		Steps    [][]cooklang.DirectionItem
	}
}

type canonicalTests struct {
	Version string
	Tests   map[string]canonicalTest
}

// TestCanonical runs the official spec tests. canonical.yaml is the current
// suite; the other two are what the official tests held while versions 5 and
// 6 were current, kept so the older rules stay honest. Both of those read
// under the version 5 rules, which is the claim their file names make.
func TestCanonical(t *testing.T) {
	suites := []struct {
		file string
		spec cooklang.Spec
	}{
		{"canonical.yaml", cooklang.SpecV7},
		{"canonical-v6.yaml", cooklang.SpecV5},
		{"canonical-v5.yaml", cooklang.SpecV5},
	}
	for _, suite := range suites {
		t.Run(suite.file, func(t *testing.T) {
			b, err := os.ReadFile(suite.file)
			if err != nil {
				t.Fatalf("failed to read %s: %v", suite.file, err)
			}
			var tests canonicalTests
			if err := yaml.Unmarshal(b, &tests); err != nil {
				t.Fatalf("failed to unmarshal %s: %v", suite.file, err)
			}
			t.Logf("canonical tests version %s\n", tests.Version)
			for name, test := range tests.Tests {
				t.Run(name, func(t *testing.T) {
					r := cooklang.MustParseSpec(test.Source, suite.spec)

					if !equalMetadata(r.Metadata, test.Result.Metadata) {
						t.Errorf("wrong metadata, got: %v, want: %v", r.Metadata, test.Result.Metadata)
					}

					got := make([][]cooklang.DirectionItem, 0, len(r.Steps))
					for _, s := range r.Steps {
						step := make([]cooklang.DirectionItem, 0, len(s.DirectionItems))
						for _, d := range s.DirectionItems {
							step = append(step, d.DirectionItem())
						}
						got = append(got, step)
					}
					if len(test.Result.Steps) == 0 {
						test.Result.Steps = [][]cooklang.DirectionItem{}
					}
					if !reflect.DeepEqual(got, test.Result.Steps) {
						t.Errorf("wrong steps for source %q\n got: %s\nwant: %s", test.Source, dumpSteps(got), dumpSteps(test.Result.Steps))
					}
				})
			}
		})
	}
}

// equalMetadata treats a nil map and an empty one as the same thing.
func equalMetadata(got cooklang.Metadata, want map[string]string) bool {
	if len(got) != len(want) {
		return false
	}
	for k, v := range want {
		if got[k] != v {
			return false
		}
	}
	return true
}

func dumpSteps(steps [][]cooklang.DirectionItem) string {
	b, err := yaml.Marshal(steps)
	if err != nil {
		return err.Error()
	}
	return "\n" + string(b)
}
