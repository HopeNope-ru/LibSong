package main

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
)

func init() {
	if err := godotenv.Load("app.env"); err != nil {
		log.Fatal().Msg("environment not present")
	}
}

func main() {
	GetApp()
}
