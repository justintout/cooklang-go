package cooklang

import (
	"reflect"
	"testing"
)

// TestMetadataLines covers what the version 5 suites have nothing to say
// about: the marker counts only at a line start, an empty one is no entry,
// and a marker that is comment text or mid-line keeps its text.
func TestMetadataLines(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		metadata map[string]string
		steps    []string
	}{
		{
			"an entry holds the text after the marker",
			">> sourced: babooshka\nAdd @salt.",
			map[string]string{"sourced": "babooshka"},
			[]string{"Add ."},
		},
		{
			"the last line needs no newline",
			">> sourced: babooshka",
			map[string]string{"sourced": "babooshka"},
			nil,
		},
		{
			"the marker counts only at a line start",
			"hello >> sourced: babooshka",
			map[string]string{},
			[]string{"hello >> sourced: babooshka"},
		},
		{
			"a bare marker is no entry",
			">>\nAdd @salt.",
			map[string]string{},
			[]string{"Add ."},
		},
		{
			"a marker inside a comment is comment text",
			"[- >> sourced: babooshka -]\nAdd @salt.",
			map[string]string{},
			[]string{"Add ."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MustParseSpec(tt.source, SpecV5)
			if !reflect.DeepEqual(map[string]string(r.Metadata), tt.metadata) {
				t.Errorf("wrong metadata for %q: got: %v, want: %v", tt.source, r.Metadata, tt.metadata)
			}
			var steps []string
			for _, s := range r.Steps {
				text := ""
				for _, t := range s.Text {
					text += t.Value
				}
				steps = append(steps, text)
			}
			if !reflect.DeepEqual(steps, tt.steps) {
				t.Errorf("wrong steps for %q: got: %q, want: %q", tt.source, steps, tt.steps)
			}
		})
	}
}

// TestEmptyFrontMatter covers a front matter block holding nothing, where the
// closing fence follows the opening one.
func TestEmptyFrontMatter(t *testing.T) {
	r := MustParse("---\n---\nAdd @salt.")
	if len(r.Metadata) != 0 {
		t.Errorf("empty front matter produced metadata: %v", r.Metadata)
	}
	if len(r.Steps) != 1 {
		t.Errorf("got %d steps, want 1", len(r.Steps))
	}
}
