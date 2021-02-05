package main

import (
	"log"

	"github.com/dominikbraun/buneary"
)

func main() {
	if err := buneary.RootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
