package cooklang

import (
	"reflect"
	"testing"
)

func TestMetadata(t *testing.T) {
	t.Run("Metadata#Add", func(t *testing.T) {
		m := Metadata{}
		m.Add("one: thing")
		if v, ok := m["one"]; !ok || v != "thing" {
			t.Errorf(`expected key/value "one: thing", got: %+v`, m)
		}
		m.Add("   another:      different thing     ")

		one, ok := m["one"]
		if !ok {
			t.Errorf(`key/value "one: thing" removed by Add call: %+v`, m)
		}
		if one != "thing" {
			t.Errorf(`key/value "one: thing" changed by Add call to different key: %+v`, m)
		}

		if v, ok := m["another"]; !ok || v != "different thing" {
			t.Errorf(`expected new key/value "another: different thing", got: %+v`, m)
		}

		m.Add("one: actually this thing")
		if v, ok := m["one"]; !ok || v != "actually this thing" {
			t.Errorf(`expected updated ey/value "one: actually this thing", got: %+v`, m)
		}
	})

	t.Run("YAML block list tags and nested source map", func(t *testing.T) {
		src := `---
title: Classic Pasta
tags:
  - pasta
  - quick
  - italian
source:
  name: Grandma's Cookbook
  author: Nonna
  url: https://example.com/pasta
---

Boil @pasta{200%g}.
`
		r := MustParse(src)
		wantTags := []string{"pasta", "quick", "italian"}
		if !reflect.DeepEqual(r.Metadata.Tags(), wantTags) {
			t.Errorf("Tags() = %v, want %v", r.Metadata.Tags(), wantTags)
		}
		if r.Metadata["tags"] != "pasta, quick, italian" {
			t.Errorf("Metadata[\"tags\"] = %q, want %q", r.Metadata["tags"], "pasta, quick, italian")
		}
		if r.Metadata["source.name"] != "Grandma's Cookbook" {
			t.Errorf("source.name = %q, want %q", r.Metadata["source.name"], "Grandma's Cookbook")
		}
		if r.Metadata["source.author"] != "Nonna" {
			t.Errorf("source.author = %q, want %q", r.Metadata["source.author"], "Nonna")
		}
		if r.Metadata["source.url"] != "https://example.com/pasta" {
			t.Errorf("source.url = %q, want %q", r.Metadata["source.url"], "https://example.com/pasta")
		}
		wantSource := Source{
			Name:   "Grandma's Cookbook",
			Author: "Nonna",
			URL:    "https://example.com/pasta",
		}
		if got := r.Metadata.Source(); got != wantSource {
			t.Errorf("Source() = %+v, want %+v", got, wantSource)
		}
	})

	t.Run("YAML flow sequence tags", func(t *testing.T) {
		src := `---
tags: [2022, baking]
---

Bake @bread.
`
		r := MustParse(src)
		wantTags := []string{"2022", "baking"}
		if !reflect.DeepEqual(r.Metadata.Tags(), wantTags) {
			t.Errorf("Tags() = %v, want %v", r.Metadata.Tags(), wantTags)
		}
		if r.Metadata["tags"] != "2022, baking" {
			t.Errorf("Metadata[\"tags\"] = %q, want %q", r.Metadata["tags"], "2022, baking")
		}
	})
}
