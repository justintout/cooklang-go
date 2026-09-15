package main

import (
	"testing"

	"github.com/justintout/cooklang-go"
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
