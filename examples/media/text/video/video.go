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

	// Генерация видео
	data, err := polzaai.GenerateVideo(
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
