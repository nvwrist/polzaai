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

	// Генерация изображения
	data, err := polzasdk2.GenerateImage(
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
