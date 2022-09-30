package cooklang_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/justintout/cooklang-go"
	"gopkg.in/yaml.v3"
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

func TestCanonical(t *testing.T) {
	b, err := os.ReadFile("canonical.yaml")
	if err != nil {
		panic("failed to read canonical.yaml: " + err.Error())
	}
	var tests canonicalTests
	err = yaml.Unmarshal(b, &tests)
	if err != nil {
		panic("failed to unmarshal tests: " + err.Error())
	}
	t.Logf("canonical tests version %s\n", tests.Version)
	for name, test := range tests.Tests {
		t.Run(name, func(t *testing.T) {

			r := cooklang.MustParse(test.Source)

			_, err := yaml.Marshal(r)
			if err != nil {
				t.Errorf("error marshaling recipe: %v", err)
			}

			for tk, tv := range test.Result.Metadata {
				found := false
				for k, v := range *r.Metadata {
					if k == tk {
						found = true
						if v != tv {
							t.Errorf("%q value incorrect, got: %v:%v, want: %v:%v", tk, k, v, tk, tv)
						}
					}
				}
				if !found {
					t.Errorf("did not find metadata key %q, got: %v, want: %v", tk, r.Metadata, test.Result.Metadata)
				}
			}

			if len(test.Result.Steps) != len(r.Steps) {
				fmt.Printf("--- %s\n", name)
				fmt.Printf("from: %s\n", test.Source)
				fmt.Printf("got: %+v\n", r.Steps)
				fmt.Printf("want: %+v\n", test.Result.Steps)
				fmt.Printf("---\n\n")
				t.Errorf("wrong number of steps, got: %d, want: %d", len(r.Steps), len(test.Result.Steps))
			}

		})
	}
}
