package main

import (
	"os"
	"reflect"
	"testing"

	"github.com/justintout/cooklang-go"
	"go.yaml.in/yaml/v3"
)

func TestUpgrade(t *testing.T) {
	tests := []struct {
		name string
		in   string
		want string
	}{
		{
			"metadata moves into front matter",
			">> Author: Jamie Oliver\nAdd @salt.\n",
			"---\nAuthor: Jamie Oliver\n---\n\nAdd @salt.\n",
		},
		{
			"each line becomes its own step",
			"Chop the @onion{}.\nAdd the @tomato{}.\n",
			"Chop the @onion{}.\n\nAdd the @tomato{}.\n",
		},
		{
			"blank lines carry no meaning so they are rebuilt",
			"one\n\n\ntwo\n",
			"one\n\ntwo\n",
		},
		{
			"a comment line stays where it is",
			"one\n-- note\ntwo\n",
			"one\n\n-- note\n\ntwo\n",
		},
		{
			"metadata can appear anywhere in the file",
			"one\n>> k: v\ntwo\n",
			"---\nk: v\n---\n\none\n\ntwo\n",
		},
		{
			"a block comment spans lines, so it hides the metadata inside it",
			"[- notes\n>> k: v\n-]\nAdd @salt.\n",
			"[- notes\n\n>> k: v\n\n-]\n\nAdd @salt.\n",
		},
		{
			"a block comment runs to the first -], so an inner dash hides nothing",
			"[- 9-inch pan\n>> k: v\n-]\nAdd @salt.\n",
			"[- 9-inch pan\n\n>> k: v\n\n-]\n\nAdd @salt.\n",
		},
		{
			"a block comment that closed hides nothing",
			"[- note -]\n>> k: v\nAdd @salt.\n",
			"---\nk: v\n---\n\n[- note -]\n\nAdd @salt.\n",
		},
		{
			"a bracket inside a line comment opens no block comment",
			"-- see [- here\n>> k: v\n",
			"---\nk: v\n---\n\n-- see [- here\n",
		},
	}
	for _, tt := range tests {
		got, err := upgrade(tt.in)
		if err != nil {
			t.Errorf("%s: unexpected error: %v", tt.name, err)
			continue
		}
		if got != tt.want {
			t.Errorf("%s:\n got: %q\nwant: %q", tt.name, got, tt.want)
		}
	}
}

func TestUpgradeRefusesFrontMatter(t *testing.T) {
	if _, err := upgrade("---\nk: v\n---\n\nAdd @salt.\n"); err == nil {
		t.Error("expected an error for a recipe that already has front matter")
	}
}

// TestUpgradeKeepsTheRecipe is the point of the rewrite: version 7 reads the
// upgraded recipe the way version 5 read the original.
func TestUpgradeKeepsTheRecipe(t *testing.T) {
	in := ">> Author: Jamie Oliver\n" +
		"Crack the @eggs{3} into a blender.\n" +
		"Pour into a #bowl and stand for ~{15%minutes}.\n" +
		"Serve straightaway.\n"
	out, err := upgrade(in)
	if err != nil {
		t.Fatal(err)
	}
	r := cooklang.MustParse(out)
	if len(r.Steps) != 3 {
		t.Errorf("got %d steps, want 3", len(r.Steps))
	}
	if r.Metadata["Author"] != "Jamie Oliver" {
		t.Errorf("metadata lost, got: %v", r.Metadata)
	}
	if len(r.Ingredients["eggs"]) != 1 || len(r.Cookware["bowl"]) != 1 {
		t.Errorf("direction items lost, got ingredients: %v cookware: %v", r.Ingredients, r.Cookware)
	}
}

// TestUpgradeReadsTheSameV5AndV7 runs every case of the official version 5
// tests through the rewrite and checks that version 7 reads the result the way
// version 5 read the source. That is the contract of both the rewrite and the
// parser's version 5 mode, so they are held to it together.
func TestUpgradeReadsTheSameV5AndV7(t *testing.T) {
	b, err := os.ReadFile("../../canonical-v5.yaml")
	if err != nil {
		t.Fatal(err)
	}
	var tests struct {
		Tests map[string]struct {
			Source string
		}
	}
	if err := yaml.Unmarshal(b, &tests); err != nil {
		t.Fatal(err)
	}
	if len(tests.Tests) == 0 {
		t.Fatal("no tests found in canonical-v5.yaml")
	}
	for name, tt := range tests.Tests {
		t.Run(name, func(t *testing.T) {
			out, err := upgrade(tt.Source)
			if err != nil {
				t.Fatal(err)
			}
			before := cooklang.MustParseSpec(tt.Source, cooklang.SpecV5)
			after := cooklang.MustParseSpec(out, cooklang.SpecV7)
			if !reflect.DeepEqual(map[string]string(before.Metadata), map[string]string(after.Metadata)) {
				t.Errorf("metadata changed by the rewrite:\nbefore: %v\nafter: %v", before.Metadata, after.Metadata)
			}
			if b, a := items(before), items(after); !reflect.DeepEqual(b, a) {
				t.Errorf("steps changed by the rewrite of %q\nbefore: %v\nafter: %v", tt.Source, b, a)
			}
		})
	}
}

// items flattens the steps of a recipe for comparison.
func items(r cooklang.Recipe) [][]cooklang.DirectionItem {
	out := make([][]cooklang.DirectionItem, 0, len(r.Steps))
	for _, s := range r.Steps {
		step := make([]cooklang.DirectionItem, 0, len(s.DirectionItems))
		for _, d := range s.DirectionItems {
			step = append(step, d.DirectionItem())
		}
		out = append(out, step)
	}
	return out
}
