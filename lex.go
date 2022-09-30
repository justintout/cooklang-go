package cooklang

import (
	"strings"
)

const special = "@#~{}\n"

func lexText(l *lexer) stateFn {
	for {
		// if l.pos == len(l.input) {
		// 	break
		// }
		if strings.HasPrefix(l.input[l.pos:], leftMetadata) {
			if l.pos > l.start {
				l.emit(itemText)
			}
			l.start = l.pos
			return lexMetadata
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
		if strings.HasPrefix(l.input[l.pos:], leftLineComment) {
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
	l.accept(leftLineComment)
	l.acceptUntil("\n")
	l.emit(itemComment)
	return lexText
}

func lexBlockComment(l *lexer) stateFn {
	l.accept(leftBlockComment)
	l.acceptUntil(rightBlockComment)
	l.accept(rightBlockComment)
	l.emit(itemComment)
	return lexText
}

func lexMetadata(l *lexer) stateFn {
	l.accept(leftMetadata)
	l.acceptUntil("\n")
	l.emit(itemMetadata)
	l.ignore()
	// l.lineStart = l.pos
	l.start = l.pos
	return lexText
}

func lexIngredient(l *lexer) stateFn {
	l.accept(leftIngredient)
	return lexQuantifiedItem(l, itemIngredient)
}

func lexCookware(l *lexer) stateFn {
	l.accept(leftCookware)
	return lexQuantifiedItem(l, itemCookware)
}

func lexTimer(l *lexer) stateFn {
	l.accept(leftTimer)
	return lexQuantifiedItem(l, itemTimer)
}

func lexQuantifiedItem(l *lexer, typ itemType) stateFn {
	l.acceptUntil(" " + leftQuantity + "\n")
	if l.accept(leftQuantity) {
		l.acceptUntil(rightQuantity)
		l.accept(rightQuantity)
		l.emit(typ)
		return lexText
	}
	if l.peek() == '\n' {
		// single word default amount ingredient
		l.emit(typ)
		return lexText
	}
	if l.accept(" ") && l.peekSpecial() != '{' {
		// single word default amount ingredient
		l.emit(typ)
		return lexText
	}
	l.acceptUntil("}")
	l.accept("}")
	l.emit(typ)
	return lexText

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
