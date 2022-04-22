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
	// start of comment, -- to end of line
	itemComment
	// plain text
	itemText
	// ingredient, starting with @ and ending in {}
	// with optional quantity and unit
	itemIngredient
	itemLeftIngredient
	//
	itemCookware
	//
	itemTimer
	//
	itemQuantity
)

const (
	//
	itemCategory itemType = iota
)

type item struct {
	tpy itemType
	val string
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

func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
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
	for strings.IndexRune(valid, l.next()) < 0 {
	}
	l.backup()
}

func (l *lexer) emit(t itemType) {
	l.items <- item{t, l.input[l.start:l.pos]}
	l.start = l.pos
}

func lexText(l *lexer) stateFn {
	for {
		if strings.HasPrefix(l.input[l.pos:], leftIngredient) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			return lexIngredient
		}
		if l.next() == eof {
			break
		}
	}
	if l.pos > l.start {
		l.emit(itemText)
	}
	l.emit(itemEOF)
	return nil
}

func lexIngredient(l *lexer) stateFn {
	l.accept("@")
	l.acceptUntil("{ ")
	if l.accept(" ") {

	}
	return nil // HERE IS WHERE YOU STOPPED
}

/*
func lexNumber(l *lexer) stateFn {
	digits := "0123456789"
	l.acceptRun(digits)
	if l.accept(".") {
		l.acceptRun(digits)
	}
	if isAlphaNumeric(l.peek()) {
		l.next()
		return l.errorf("bad number: %q", l.input[l.start:l.pos])
	}
	l.emit(itemNumber)
	return lexInsideQuantity
}
*/
