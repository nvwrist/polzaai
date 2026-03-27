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

	// Генерация видео
	data, err := polzasdk2.GenerateVideo(
		ctx,
		client,
		"openai/sora-2",
		"Кот играет с лазерной указкой",
		"5",    // длительность
		"720p", // разрешение
		24,     // fps
	)
	if err != nil {
		log.Fatal(err)
	}
	println(string(data))
}
