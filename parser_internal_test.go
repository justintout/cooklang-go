package cooklang

import "testing"

func TestIsSpaceButNotNewline(t *testing.T) {
	notSpaces := []rune{'\n'}
	for _, r := range notSpaces {
		if isSpace(r, true) {
			t.Errorf("%q is not a space", r)
		}
	}
}
