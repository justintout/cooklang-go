package main

import (
	"encoding/json"
	"fmt"

	"github.com/justintout/cooklang-go"
)

func main() {
	f := "./example-condensed.cook"
	r := cooklang.MustParseFile(f)
	fmt.Printf("%+v\n", r)
	for _, s := range r.Steps {
		fmt.Printf("\t%+v\n", s)
	}
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
