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

	// Просто передаём значения
	data, err := polzaai.GenerateAudio(
		ctx,
		client,
		"elevenlabs/text-to-speech-turbo-2-5",
		"Привет, мир. Как дела?!", // текст
		"Laura", // голос
		1.0,     // скорость
		"mp3",   // формат
		"111000111",
	)
	if err != nil {
		log.Fatal(err)
	}
	println(string(data))

}
