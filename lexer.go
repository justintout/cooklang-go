// the lexer is based on rob pike's "lexical scanning in go" talk
// from 2011: https://talks.golang.org/2011/lex.slide#1
// i think he's said at this point its out dated and would do things
// differently. need to research this.
package cooklang

import (
	"strings"
	"unicode"
	"unicode/utf8"
)

type stateFn func(*lexer) stateFn

type itemType int

const (
	itemError itemType = iota
	itemEOF
	itemText
	itemComment
	itemMetadata
	itemStep
	itemIngredient
	//
	itemCookware
	//
	itemTimer
	//
)

type item struct {
	typ itemType
	val string
}

type lexer struct {
	input     string
	spec      Spec
	start     int
	lineStart int
	pos       int
	width     int
	items     chan item
}

func lex(input string, spec Spec) (*lexer, chan item) {
	l := &lexer{
		input: input,
		spec:  spec,
		items: make(chan item),
	}
	go l.run()
	return l, l.items
}

func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.items)
}

func (l *lexer) next() (r rune) {
	if l.pos >= len(l.input) {
		l.width = 0
		return eof
	}
	r, l.width = utf8.DecodeRuneInString(l.input[l.pos:])
	l.pos += l.width
	return r
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) peekSpecial() rune {
	p := l.pos
	for {
		r := l.next()
		if r == eof {
			l.pos = p
			return eof
		}
		if strings.ContainsRune(special, r) {
			l.pos = p
			return r
		}
	}
}

func (l *lexer) accept(valid string) bool {
	if strings.ContainsRune(valid, l.next()) {
		return true
	}
	l.backup()
	return false
}

// acceptString consumes a multi-character prefix, unlike accept which matches
// a single rune against a set of them.
func (l *lexer) acceptString(prefix string) bool {
	if strings.HasPrefix(l.input[l.pos:], prefix) {
		l.pos += len(prefix)
		return true
	}
	return false
}

// atFence reports whether s opens with a line that is exactly "---", the
// delimiter of the YAML front matter block.
func atFence(s string) bool {
	rest, ok := strings.CutPrefix(s, metadataFence)
	return ok && (rest == "" || strings.HasPrefix(rest, "\n"))
}

// atMetadataFence reports whether the lexer sits on a front matter fence.
func (l *lexer) atMetadataFence() bool {
	return l.pos == l.lineStart && atFence(l.input[l.pos:])
}

// atLineComment reports whether the lexer is on a line comment. A run of
// three or more dashes is not one, so that the metadata fence stays text
// wherever it appears away from the start of a file.
func (l *lexer) atLineComment() bool {
	if !strings.HasPrefix(l.input[l.pos:], leftLineComment) {
		return false
	}
	rest := l.input[l.pos+len(leftLineComment):]
	if strings.HasPrefix(rest, "-") {
		return false
	}
	// the scan reaches the second dash of a longer run, so check behind too
	return l.pos == 0 || l.input[l.pos-1] != '-'
}

// nameless reports whether a marker is followed by nothing that could be a
// name, in which case the marker is plain text.
func (l *lexer) nameless() bool {
	r := l.peek()
	return r == eof || unicode.IsSpace(r)
}

// cutAtWordEnd backs up to the first punctuation or space in the word
// scanned since floor, which is where a name that is not delimited by braces
// ends. It reports whether any name is left, which it is not when the word
// opens with punctuation.
func (l *lexer) cutAtWordEnd(floor int) bool {
	for i, r := range l.input[floor:l.pos] {
		if unicode.IsPunct(r) || unicode.IsSpace(r) {
			l.pos = floor + i
			return i > 0
		}
	}
	return l.pos > floor
}

func (l *lexer) acceptUntil(valid string) {
	for !strings.ContainsRune(valid, l.next()) && l.peek() != eof {
	}
	l.backup()
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}
