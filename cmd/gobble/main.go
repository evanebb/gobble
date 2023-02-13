package main

import (
	"github.com/evanebb/gobble/server"
	"log"
)

func main() {
	s, err := server.NewServer()
	if err != nil {
		log.Fatal(err)
	}

	s.Run()
}
