// Command cookupgrade rewrites a recipe written for spec version 5 so that
// version 7 reads it the same way.
//
// Two things changed. Metadata moved from ">> key: value" lines into a YAML
// front matter block, and a step became a paragraph instead of a line. A
// version 5 recipe made every line its own step, and a blank line did
// nothing, so the upgrade separates every line with a blank line. The result
// has the same steps and the same metadata, and is longer to read.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

const (
	fence        = "---"
	metadataIn   = ">>"
	lineComment  = "--"
	blockComment = "[-"
	blockEnd     = "-]"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s recipe.cook\n\nThe upgraded recipe is written to stdout.\n", os.Args[0])
	}
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}
	b, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fail(err)
	}
	out, err := upgrade(string(b))
	if err != nil {
		fail(fmt.Errorf("%s: %v", flag.Arg(0), err))
	}
	fmt.Print(out)
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	os.Exit(1)
}

// upgrade rewrites a version 5 recipe for version 7.
func upgrade(input string) (string, error) {
	if strings.HasPrefix(input, fence) {
		return "", errors.New("already opens with YAML front matter, so it is not a version 5 recipe")
	}

	lines := strings.Split(strings.TrimSuffix(input, "\n"), "\n")

	var front []string
	var body []string
	inComment := false
	for _, line := range lines {
		startsInComment := inComment
		inComment = scanComments(line, inComment)
		if entry, ok := metadataLine(line); ok && !startsInComment {
			front = append(front, entry)
			continue
		}
		body = append(body, line)
	}

	out := []string{}
	if len(front) > 0 {
		out = append(out, fence)
		out = append(out, front...)
		out = append(out, fence, "")
	}
	out = append(out, separated(body)...)

	return strings.Join(out, "\n") + "\n", nil
}

// metadataLine returns the "key: value" text of a version 5 metadata line.
func metadataLine(line string) (string, bool) {
	if !strings.HasPrefix(line, metadataIn) {
		return "", false
	}
	return strings.TrimSpace(strings.TrimPrefix(line, metadataIn)), true
}

// scanComments reports whether a line ends inside a block comment, given
// whether it began inside one. A block comment runs from "[-" to the first
// "-]" and outlives its line, so a ">>" between the two is comment text and
// not metadata. A line comment ends at the newline, so it only hides a "[-"
// that shares its line.
func scanComments(line string, inComment bool) bool {
	for i := 0; i < len(line); {
		if inComment {
			j := strings.Index(line[i:], blockEnd)
			if j < 0 {
				return true
			}
			inComment = false
			i += j + len(blockEnd)
			continue
		}
		blockAt := strings.Index(line[i:], blockComment)
		lineAt := strings.Index(line[i:], lineComment)
		if blockAt < 0 || (lineAt >= 0 && lineAt < blockAt) {
			return false
		}
		inComment = true
		i += blockAt + len(blockComment)
	}
	return inComment
}

// separated puts a blank line between each pair of lines and drops the blank
// lines of the original, where they carried no meaning.
func separated(lines []string) []string {
	out := []string{}
	for _, line := range lines {
		if strings.TrimSpace(line) == "" {
			continue
		}
		if len(out) > 0 {
			out = append(out, "")
		}
		out = append(out, line)
	}
	return out
}
