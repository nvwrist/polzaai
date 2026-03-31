package polzaai

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"path/filepath"

	"strings"
	"time"

	"github.com/nvwrist/polzaai/polzaai/models"
)

// saveToStorage скачивает файл (если URL) или принимает base64, загружает в хранилище.
// Возвращает (storedURL, storedID, error)
func saveToStorage(ctx context.Context, client *Client, urlOrBase64, filename string) (string, string, error) {
	var fileData string
	var actualFilename string

	if strings.HasPrefix(urlOrBase64, "http://") || strings.HasPrefix(urlOrBase64, "https://") {
		resp, err := http.Get(urlOrBase64)
		if err != nil {
			return "", "", fmt.Errorf("ошибка скачивания файла: %w", err)
		}
		defer resp.Body.Close()

		fileBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			return "", "", fmt.Errorf("ошибка чтения файла: %w", err)
		}

		base64Data := base64.StdEncoding.EncodeToString(fileBytes)

		if filename == "" {
			parts := strings.Split(urlOrBase64, "/")
			actualFilename = parts[len(parts)-1]
			if actualFilename == "" {
				actualFilename = fmt.Sprintf("file_%d", time.Now().Unix())
			}
		} else {
			actualFilename = filename
		}

		ext := strings.ToLower(filepath.Ext(actualFilename))
		mimeType := "application/octet-stream"
		switch ext {
		case ".mp3":
			mimeType = "audio/mpeg"
		case ".wav":
			mimeType = "audio/wav"
		case ".mp4":
			mimeType = "video/mp4"
		case ".png":
			mimeType = "image/png"
		case ".jpg", ".jpeg":
			mimeType = "image/jpeg"
		}
		fileData = "data:" + mimeType + ";base64," + base64Data
	} else {
		fileData = urlOrBase64
		actualFilename = filename
		if actualFilename == "" {
			actualFilename = fmt.Sprintf("file_%d", time.Now().Unix())
		}
	}

	uploadResp, err := client.Storage().UploadFile(ctx, models.UploadFileRequest{
		Filename: actualFilename,
		Data:     fileData,
		Policy:   "PERMANENT",
	})
	if err != nil {
		return "", "", fmt.Errorf("ошибка загрузки в хранилище: %w", err)
	}
	return uploadResp.URL, uploadResp.ID, nil
}

// getMediaURL извлекает URL из поля Data ответа MediaResponse.
func getMediaURL(resp *models.MediaResponse) string {
	if resp.Data == nil {
		return ""
	}
	switch v := resp.Data.(type) {
	case string:
		return v
	case map[string]interface{}:
		if url, ok := v["url"].(string); ok {
			return url
		}
	case []interface{}:
		if len(v) > 0 {
			if first, ok := v[0].(map[string]interface{}); ok {
				if url, ok := first["url"].(string); ok {
					return url
				}
			}
		}
	}
	return ""
}

// waitForMediaCompletion ожидает завершения асинхронной задачи.
func waitForMediaCompletion(ctx context.Context, client *Client, mediaID string) (*models.MediaResponse, error) {
	// Увеличиваем интервал до 10 секунд. Для видео это оптимально.
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	// 90 попыток по 10 секунд = 900 секунд (15 минут).
	// Это как раз совпадает с таймаутом в твоем UniversalHandler.
	maxAttempts := 90

	for attempt := 1; attempt <= maxAttempts; attempt++ {
		select {
		case <-ctx.Done():
			// Если время в context (15 мин) выйдет раньше, Go сам прервет цикл здесь
			return nil, ctx.Err()
		case <-ticker.C:
			statusResp, err := client.Media().GetMediaStatus(ctx, mediaID)
			if err != nil {
				// Можно добавить лог, чтобы видеть ошибки сети во время ожидания
				fmt.Printf("[SDK DEBUG] Ошибка опроса статуса (попытка %d): %v\n", attempt, err)
				continue // Пропускаем ошибку и пробуем в следующий раз
			}

			fmt.Printf("[SDK DEBUG] Статус задачи %s: %s (попытка %d/%d)\n",
				mediaID, statusResp.Status, attempt, maxAttempts)

			switch statusResp.Status {
			case "completed":
				return statusResp, nil
			case "failed", "error": // Добавляем обработку обоих вариантов ошибок
				errMsg := "неизвестная ошибка"
				if statusResp.Error != nil {
					errMsg = statusResp.Error.Message
				}
				return nil, fmt.Errorf("задача завершилась с ошибкой: %s", errMsg)
			case "cancelled":
				return nil, fmt.Errorf("задача была отменена")
			}
		}
	}
	return nil, fmt.Errorf("превышено время ожидания завершения (ждали 15 минут)")
}
