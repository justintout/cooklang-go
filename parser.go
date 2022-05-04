package cooklang

import (
	"fmt"
	"os"
	"unicode/utf8"
)

func MustParseFile(path string) Recipe {
	r, err := ParseFile(path)
	if err != nil {
		panic(err)
	}
	return r
}

func ParseFile(path string) (Recipe, error) {
	c, err := os.ReadFile(path)
	if err != nil {
		return Recipe{}, fmt.Errorf("failed to parse %q: %v", path, err)
	}
	return parse(string(c))
}

func MustParse(input string) Recipe {
	r, err := Parse(input)
	if err != nil {
		panic(err)
	}
	return r
}

func Parse(input string) (Recipe, error) {
	return parse(input)
}

// TODO: scanners aren't gonna be good because we can't backtrack
// gotta do this better...
func parse(input string) (Recipe, error) {
	_, items := lex("recipe", input)
	recipe := NewRecipe("recipe")
	step := &Step{}
	for item := range items {
		switch item.typ {
		case itemMetadata:
			recipe.Metadata.Add(item.val)
		case itemComment:
		case itemText:
			step.AddText(NewText(item.val))
		case itemIngredient:
			step.AddIngredient(NewIngredient(item.val))
		case itemCookware:
			step.AddCookware(NewCookware(item.val))
		case itemTimer:
			step.AddTimer(NewTimer(item.val))
		case itemStep:
			recipe.AddStep(step)
			step = &Step{}
		}
	}

	// HERE IS WHERE YOU STOPPED
	// find way to print step out in order (raw step? idk)
	// or iterate through each thing and find the right pos
	// probably need to track that somewhere...
	return recipe, nil

	/*
		lines := bufio.NewScanner(r)
		recipe := new(Recipe)
		recipe.Metadata = make(map[string]string)
		for lines.Scan() {
			line := strings.TrimSpace(lines.Text())
			if strings.HasPrefix(line, ">>") {
				line := strings.Replace(strings.TrimPrefix(line, ">>"), " ", "", -1)
				s := strings.SplitN(line, ":", 2)
				if len(s) == 1 {
					s[1] = ""
				}
				recipe.Metadata[s[0]] = s[1]
				if strings.EqualFold(s[0], "servings") {
					s := strings.Split(s[1], "|")
					for i, ss := range s {
						srv := Serving{
							Idx: i,
							S:   ss,
						}
						var err error
						srv.N, err = strconv.Atoi(ss)
						if err == nil {
							recipe.Servings = append(recipe.Servings, srv)
							continue
						}
						// TODO: handle? ignore? invalidate whole servings?
						fmt.Printf("error parsing serving: %q\n", ss)
					}
				}
				continue
			}
			// TODO: don't discard line comments
			if strings.HasPrefix(line, "--") {
				continue
			}
			// step := Step{
			// 	raw: line,
			// }
			// rel := relationships{
			// 	step: &step,
			// 	recipe: recipe,
			// }
		}
		return *recipe, nil
	*/
}

// isSpace patched to ignore '\n'
// https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/bufio/scan.go;l=370
func isSpace(r rune, ignoreNewline bool) bool {
	if r == '\n' && ignoreNewline {
		return false
	}
	if r <= '\u00FF' {
		// Obvious ASCII ones: \t through \r plus space. Plus two Latin-1 oddballs.
		switch r {
		case ' ', '\t', '\n', '\v', '\f', '\r':
			return true
		case '\u0085', '\u00A0':
			return true
		}
		return false
	}
	// High-valued ones.
	if '\u2000' <= r && r <= '\u200a' {
		return true
	}
	switch r {
	case '\u1680', '\u2028', '\u2029', '\u202f', '\u205f', '\u3000':
		return true
	}
	return false
}

// ScanWords patched to scan '\n' as an empty token instead
// https://cs.opensource.google/go/go/+/refs/tags/go1.18.1:src/bufio/scan.go;l=396
func scanWordsAndNewline(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !isSpace(r, true) {
			break
		}
	}
	// Scan until space, or newline, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if isSpace(r, false) {
			if r == '\n' {
				return i + width, append(data[start:i], '\n'), nil
			}
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		return len(data), data[start:], nil
	}
	// Request more data.
	return start, nil, nil
}
