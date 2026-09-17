package cooklang

import "strings"

// Spec identifies a version of the Cooklang specification.
type Spec int

const (
	// SpecV5 reads metadata from ">> key: value" lines and ends a step at
	// every newline. Version 6 reads the same, so it has no constant of its
	// own.
	SpecV5 Spec = 5
	// SpecV7 reads metadata from YAML front matter and ends a step at a blank
	// line.
	SpecV7 Spec = 7
)

// latest is the version that the entry points without a spec read.
const latest = SpecV7

// known reports whether s names a version this package reads.
func (s Spec) known() bool {
	switch s {
	case SpecV5, SpecV7:
		return true
	}
	return false
}

// Detect reports the spec version a recipe appears to be written for: SpecV7
// for a source that opens with a front matter fence, SpecV5 for one with a
// ">> key: value" line, and SpecV7 otherwise.
//
// Those two are the only signals there are. A version 5 recipe with no
// metadata reads as version 7, which merges its lines into paragraphs, and
// no amount of looking at the text can separate it from version 7 prose.
// Pass a version explicitly where the answer matters.
func Detect(input string) Spec {
	// Ahead of the metadata check: front matter holds lines verbatim, so one
	// of them can read ">> something: value" and look like version 5.
	if atFence(input) {
		return SpecV7
	}
	_, items := lex(input, SpecV5)
	found := false
	for it := range items {
		// The lexer blocks on an unbuffered channel, so the tokens keep
		// coming until the end of the input even once there is an answer.
		if it.typ == itemMetadata && strings.Contains(it.val, ":") {
			found = true
		}
	}
	if found {
		return SpecV5
	}
	return SpecV7
}
