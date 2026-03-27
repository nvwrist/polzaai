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

	// Просто передаём значения
	data, err := polzasdk2.GenerateAudio(
		ctx,
		client,
		"elevenlabs/text-to-speech-turbo-2-5",
		"Привет, мир!", // текст
		"Laura",        // голос
		1.0,            // скорость
		"mp3",          // формат
	)
	if err != nil {
		log.Fatal(err)
	}
	println(string(data))

}
