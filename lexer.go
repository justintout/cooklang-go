// the lexer is based on rob pike's "lexical scanning in go" talk
// from 2011: https://talks.golang.org/2011/lex.slide#1
// i think he's said at this point its out dated and would do things
// differently. need to research this.
package cooklang

import (
	"strings"
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

const (
	//
	itemCategory itemType = iota
)

type item struct {
	typ itemType
	val string
}

func (i item) String() string {
	switch i.typ {
	case itemEOF:
		return "EOF"
	case itemError:
		return i.val
	}
	return i.val
}

type lexer struct {
	name       string
	input      string
	stepsInput []string
	start      int
	pos        int
	width      int
	items      chan item
	steps      chan Step
}

func lex(name, input string) (*lexer, chan item) {
	l := &lexer{
		name:       name,
		input:      input,
		stepsInput: strings.Split(input, "\n"),
		items:      make(chan item),
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

func (l *lexer) ignore() {
	l.start = l.pos
}

func (l *lexer) backup() {
	l.pos -= l.width
}

func (l *lexer) rewind() {
	l.pos = l.start
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

func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
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
