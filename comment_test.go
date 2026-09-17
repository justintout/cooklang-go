package cooklang

import (
	"reflect"
	"testing"
)

// TestComments checks that no text of a comment reaches a step. A block
// comment runs to the two characters "-]", which its own text may hold only
// singly, and a line comment runs to the newline.
func TestComments(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   []string
	}{
		{
			"a dash inside a block comment does not end it",
			"[- use a non-stick pan -]\nAdd @salt.",
			[]string{"Add ."},
		},
		{
			"a bracket inside a block comment does not end it",
			"[- see [1] -]\nAdd @salt.",
			[]string{"Add ."},
		},
		{
			"a block comment spans lines",
			"[- a note\nover two lines -]\nAdd @salt.",
			[]string{"Add ."},
		},
		{
			"an unclosed block comment runs to the end of the input",
			"Add @salt.\n[- unclosed\nmore text",
			[]string{"Add ."},
		},
		{
			"a bracketed dash inside a line comment opens no block comment",
			"Add @salt.\n-- a bracket [- in a line comment\nPepper.",
			[]string{"Add .", "Pepper."},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := MustParse(tt.source)
			var got []string
			for _, s := range r.Steps {
				text := ""
				for _, t := range s.Text {
					text += t.Value
				}
				got = append(got, text)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("wrong steps for %q\n got: %q\nwant: %q", tt.source, got, tt.want)
			}
		})
	}
}
