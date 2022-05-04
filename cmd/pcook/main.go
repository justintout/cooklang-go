package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/justintout/cooklang-go"
)

func main() {
	if len(os.Args[1:]) == 0 {
		fmt.Printf("usage: %s file\n", os.Args[0])
		os.Exit(1)
	}
	f := os.Args[1]
	r := cooklang.MustParseFile(f)
	j, err := json.Marshal(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(j))
}
