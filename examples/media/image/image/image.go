package main

import (
	"context"
	"log"
	"os"
	polzasdk2 "polzasdk/polzaai"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal(err)
	}
	apiKey := os.Getenv("POLZA_API")
	client := polzasdk2.NewClient(apiKey)
	ctx := context.Background()

	data, err := polzasdk2.EditImage(
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
