package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/justintout/cooklang-go"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "usage: %s [-spec version] recipe.cook\n\nThe parsed recipe is written to stdout as JSON.\nversion is 5, 6, 7, auto or latest.\n", os.Args[0])
	}
	version := flag.String("spec", "latest", "spec version to read the recipe as")
	flag.Parse()
	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	b, err := os.ReadFile(flag.Arg(0))
	if err != nil {
		fail(err)
	}
	src := string(b)
	spec, err := resolve(*version, src)
	if err != nil {
		fail(err)
	}

	r := cooklang.MustParseSpec(src, spec)
	j, err := json.Marshal(r)
	if err != nil {
		fail(err)
	}
	fmt.Println(string(j))
}

func fail(err error) {
	fmt.Fprintf(os.Stderr, "%s: %v\n", os.Args[0], err)
	os.Exit(1)
}

// resolve reports the spec version to read src as. Version 6 resolves to the
// version 5 rules: the official test files for the two hold the same
// expectations, so this parser has no way to tell a recipe of one from the
// other.
func resolve(version, src string) (cooklang.Spec, error) {
	switch strings.ToLower(version) {
	case "5", "6":
		return cooklang.SpecV5, nil
	case "7", "latest":
		return cooklang.SpecV7, nil
	case "auto":
		return cooklang.Detect(src), nil
	}
	return 0, fmt.Errorf("unknown spec version %q", version)
}
