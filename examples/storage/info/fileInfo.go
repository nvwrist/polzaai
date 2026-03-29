package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/nvwrist/polzaai/polzaai"
)

func main() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Ошибка загрузки .env файла:", err)
	}
	apiKey := os.Getenv("POLZA_API")
	if apiKey == "" {
		log.Fatal("POLZA_API не задан в .env")
	}
	client := polzaai.NewClient(apiKey)
	ctx := context.Background()
	files, err := client.Storage().FileList(ctx)
	if err != nil {
		log.Fatal(err)
	}
	printJSON(files)
}

// printJSON выводит объект в формате JSON
func printJSON(v interface{}) {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		fmt.Printf("Ошибка форматирования JSON: %v\n", err)
		return
	}
	fmt.Println(string(data))
}
