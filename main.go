package main

import (
	"log"

	"gha-publish-evidence/cmd"
)

func main() {

	if err := cmd.Execute(); err != nil {
		log.Fatal(err)
	}
}
