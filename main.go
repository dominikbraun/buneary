package main

import (
	"log"
)

func main() {
	if err := rootCommand().Execute(); err != nil {
		log.Fatal(err)
	}
}
