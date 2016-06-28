package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"

	"github.com/kyleconroy/pb/lint"
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	for _, file := range flag.Args() {
		// TODO: Support diretories
		blob, err := ioutil.ReadFile(file)
		if err != nil {
			log.Fatal(err)
		}
		path := filepath.Base(file)
		problems, err := lint.Lint(path, blob)
		if err != nil {
			log.Fatal(err)
		}

		for _, p := range problems {
			fmt.Printf("%s:%d %s\n", path, 0, p.Text)
		}
	}
}
