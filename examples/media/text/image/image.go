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

	// Генерация изображения
	data, err := polzaai.GenerateImage(
		ctx,
		client,
		"google/gemini-2.5-flash-image",
		"Космический корабль в стиле киберпанк", // промпт
		"16:9", // соотношение сторон
		"high", // качество
	)
	if err != nil {
		log.Fatal(err)
	}
	println(string(data))

}
