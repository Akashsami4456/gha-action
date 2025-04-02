package main

import (
	"log"

	"github.com/calculi-corp/gha-publish-evidence-item/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
