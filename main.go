package main

import (
	"github.com/honeybadger-io/honeybadger-go"
	"os"
)

func main() {
	defer honeybadger.Monitor()

	s := Server{
		Host: os.Getenv("HOST"),
		Port: os.Getenv("PORT"),
	}
	s.Start()
}
