package cooklang

import (
	"strings"
)

const special = "@#~{}\n"

func lexText(l *lexer) stateFn {
	for {
		// Front matter only counts at the very start of the file; a "---"
		// line anywhere else is ordinary text.
		if l.spec == SpecV7 && l.pos == 0 && l.atMetadataFence() {
			l.pos += len(metadataFence)
			l.accept("\n")
			l.start = l.pos
			// The fence counts as a line of its own, or lexFrontMatter cannot
			// recognize a closing fence that follows it immediately.
			l.lineStart = l.pos
			return lexFrontMatter
		}
		// Version 5 metadata is a ">> key: value" line, wherever it appears.
		// The marker is dropped here so that the token holds the same
		// "key: value" text that front matter holds under version 7.
		if l.spec == SpecV5 && l.pos == l.lineStart && strings.HasPrefix(l.input[l.pos:], leftMetadata) {
			l.pos += len(leftMetadata)
			l.start = l.pos
			l.acceptUntil("\n")
			if l.pos > l.start {
				l.emit(itemMetadata)
			}
			continue
		}
		if l.pos == l.lineStart && strings.HasPrefix(l.input[l.pos:], leftSection) {
			l.start = l.pos
			return lexSection
		}
		if l.pos == l.lineStart && strings.HasPrefix(l.input[l.pos:], leftNote) && !strings.HasPrefix(l.input[l.pos:], leftMetadata) {
			l.start = l.pos
			return lexNote
		}
		if strings.HasPrefix(l.input[l.pos:], lineBreakMarker+"\n") || strings.HasPrefix(l.input[l.pos:], lineBreakMarker+"\r\n") {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.acceptString(lineBreakMarker)
			l.accept("\r")
			l.accept("\n")
			l.emit(itemLineBreak)
			l.lineStart = l.pos
			continue
		}
		if strings.HasPrefix(l.input[l.pos:], leftIngredient) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.start = l.pos
			return lexIngredient
		}
		if strings.HasPrefix(l.input[l.pos:], leftCookware) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.start = l.pos
			return lexCookware
		}
		if strings.HasPrefix(l.input[l.pos:], leftTimer) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.start = l.pos
			return lexTimer
		}
		if l.atLineComment() {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.start = l.pos
			return lexLineComment
		}
		if strings.HasPrefix(l.input[l.pos:], leftBlockComment) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.start = l.pos
			return lexBlockComment
		}
		if strings.HasPrefix(l.input[l.pos:], "\n") {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.accept("\n")
			l.emit(itemStep)
			l.lineStart = l.pos
			// The line is finished. Without this the loop falls through to
			// l.next(), which consumes the first character of the next line
			// and stops the prefix checks below from ever seeing a line
			// start.
			continue
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

func lexLineComment(l *lexer) stateFn {
	l.acceptString(leftLineComment)
	l.acceptUntil("\n")
	l.emit(itemComment)
	return lexText
}

// lexBlockComment reads a "[- ... -]" comment. The text of one may hold
// dashes and brackets, so the scan looks for the two characters "-]" rather
// than for the first of either; a comment left open runs to the end of the
// input.
func lexBlockComment(l *lexer) stateFn {
	l.acceptString(leftBlockComment)
	if i := strings.Index(l.input[l.pos:], rightBlockComment); i >= 0 {
		l.pos += i + len(rightBlockComment)
	} else {
		l.pos = len(l.input)
	}
	l.emit(itemComment)
	return lexText
}

// lexFrontMatter reads the YAML front matter block that a recipe may open
// with: a "---" line, then YAML metadata content, closed by another "---"
// line. The front matter body is emitted as an itemMetadata token.
func lexFrontMatter(l *lexer) stateFn {
	for {
		if l.atMetadataFence() {
			if l.pos > l.start {
				l.emit(itemMetadata)
			}
			l.pos += len(metadataFence)
			l.start = l.pos
			l.lineStart = l.pos
			return lexText
		}
		if strings.HasPrefix(l.input[l.pos:], "\n") {
			l.accept("\n")
			l.lineStart = l.pos
			continue
		}
		if l.next() == eof {
			if l.pos > l.start {
				l.emit(itemMetadata)
			}
			l.emit(itemEOF)
			return nil
		}
	}
}

func lexSection(l *lexer) stateFn {
	for l.peek() != '\n' && l.peek() != eof {
		if l.atLineComment() || strings.HasPrefix(l.input[l.pos:], leftBlockComment) {
			break
		}
		l.next()
	}
	val := strings.TrimSpace(l.input[l.start:l.pos])
	val = strings.Trim(val, leftSection)
	val = strings.TrimSpace(val)
	l.items <- item{itemSection, val}
	l.start = l.pos
	return lexText
}

func lexNote(l *lexer) stateFn {
	l.acceptString(leftNote)
	for l.accept(" \t") {
	}
	l.start = l.pos
	for l.peek() != '\n' && l.peek() != eof {
		if strings.HasPrefix(l.input[l.pos:], lineBreakMarker+"\n") || strings.HasPrefix(l.input[l.pos:], lineBreakMarker+"\r\n") {
			break
		}
		if l.atLineComment() || strings.HasPrefix(l.input[l.pos:], leftBlockComment) {
			break
		}
		l.next()
	}
	l.emit(itemNote)
	return lexText
}

func lexIngredient(l *lexer) stateFn {
	return lexMarker(l, leftIngredient, itemIngredient)
}

func lexCookware(l *lexer) stateFn {
	return lexMarker(l, leftCookware, itemCookware)
}

func lexTimer(l *lexer) stateFn {
	return lexMarker(l, leftTimer, itemTimer)
}

// lexMarker reads a "@name", "#name" or "~name" item. A marker with no name
// after it is ordinary text.
func lexMarker(l *lexer, marker string, typ itemType) stateFn {
	l.acceptString(marker)
	if typ == itemIngredient || typ == itemCookware {
		l.accept(optionalMarker)
	}
	if l.nameless() {
		l.emit(itemText)
		return lexText
	}
	return lexQuantifiedItem(l, typ)
}

func lexQuantifiedItem(l *lexer, typ itemType) stateFn {
	nameStart := l.pos
	l.acceptUntil(" " + leftQuantity + "\n")
	if l.accept(leftQuantity) {
		l.acceptUntil(rightQuantity)
		l.accept(rightQuantity)
		if typ == itemIngredient {
			l.acceptPreparation()
		}
		l.emit(typ)
		return lexText
	}
	if l.peek() == '\n' || l.peek() == eof {
		// single word default amount ingredient
		if !l.cutAtWordEnd(nameStart) {
			l.emit(itemText)
			return lexText
		}
		if typ == itemIngredient {
			l.acceptPreparation()
		}
		l.emit(typ)
		return lexText
	}
	if p := l.pos; l.accept(" ") {
		if l.peekSpecial() != '{' {
			// single word default amount ingredient: the space is not part
			// of the name, it belongs to the text that follows
			l.pos = p
			if !l.cutAtWordEnd(nameStart) {
				l.emit(itemText)
				return lexText
			}
			if typ == itemIngredient {
				l.acceptPreparation()
			}
			l.emit(typ)
			return lexText
		}
	}
	l.acceptUntil("}")
	l.accept("}")
	if typ == itemIngredient {
		l.acceptPreparation()
	}
	l.emit(typ)
	return lexText
}
