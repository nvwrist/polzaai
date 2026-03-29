package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nvwrist/polzaai/polzaai"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	apiKey := os.Getenv("POLZA_API")
	client := polzaai.NewClient(apiKey)
	ctx := context.Background()

	data, err := polzaai.EditImage(
		ctx,
		client,
		"google/gemini-2.5-flash-image",
		"Измени корабль на сосиску",
		"https://example.com/example.png",
	)
	if err != nil {
		log.Fatal(err)
	}
	println(string(data))

}
