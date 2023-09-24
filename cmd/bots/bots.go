package main

import (
	"github.com/disgoorg/log"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Debug(err)
	}

	RunTyrant()
}
