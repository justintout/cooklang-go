package cooklang

import "testing"

func TestDetect(t *testing.T) {
	tests := []struct {
		name   string
		source string
		want   Spec
	}{
		{"front matter", "---\nsourced: babooshka\n---\n\nAdd @salt.", SpecV7},
		{"a metadata line", ">> sourced: babooshka\nAdd @salt.", SpecV5},
		{"neither signal", "Add @salt.", SpecV7},
		{"a marker inside a comment", "[- >> sourced: babooshka -]\nAdd @salt.", SpecV7},
		{"front matter holding a marker", "---\n>> note: keep the pan hot\n---\n\nAdd @salt.", SpecV7},
		{"a dash line that is no fence", "--- see below\nAdd @salt.", SpecV7},
		{"a marker with no colon", ">> see below\nAdd @salt.", SpecV7},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Detect(tt.source); got != tt.want {
				t.Errorf("Detect(%q) = %d, want %d", tt.source, got, tt.want)
			}
		})
	}
}

func TestParseSpecRejectsUnknownVersion(t *testing.T) {
	if _, err := ParseSpec("Add @salt.", Spec(0)); err == nil {
		t.Error("expected an error for an unknown spec version")
	}
}
