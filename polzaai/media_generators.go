package polzaai

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/nvwrist/polzaai/polzaai/models"
)

// GenerateImage создаёт изображение по тексту, сохраняет его и возвращает JSON с результатом.
func GenerateImage(ctx context.Context, client *Client, model, prompt, aspectRatio, quality string, userid string, imageResolution string) ([]byte, error) {
	input := models.MediaInput{
		Prompt:      &prompt,
		AspectRatio: &aspectRatio,
		Quality:     &quality,
		MaxImages:   intPtr(1),
	}

	// Присваиваем только если не пусто, иначе оставляем nil
	if imageResolution != "" {
		input.ImageResolution = &imageResolution
	}

	req := models.MediaRequest{
		Model: model,
		Input: input,
		User:  &userid,
	}
	return callMediaAndSave(ctx, client, req, "generated_image.png")
}

// EditImage редактирует изображение по тексту и исходному URL, возвращает JSON.
// EditImage редактирует изображение по тексту и исходному URL, возвращает JSON.
func EditImage(ctx context.Context, client *Client, modelname, prompt, imageURL string, userid string, imageResolution string, aspectRatio string, quality string) ([]byte, error) {
	input := models.MediaInput{
		Prompt: &prompt,
		Images: []models.MediaFile{
			{Type: "url", Data: imageURL},
		},
	}
	if imageResolution != "" {
		input.ImageResolution = &imageResolution
	}
	if aspectRatio != "" {
		input.AspectRatio = &aspectRatio
	}
	if quality != "" {
		input.Quality = &quality
	}
	req := models.MediaRequest{
		Model: modelname,
		Input: input,
		User:  &userid,
	}
	return callMediaAndSave(ctx, client, req, "edited_image.png")
}

// GenerateVideo генерирует видео по тексту, возвращает JSON.
func GenerateVideo(ctx context.Context, client *Client, modelname string, prompt string, durationSeconds string, resolution string, fps int, userid string) ([]byte, error) {
	async := true
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Prompt:          &prompt,
			DurationSeconds: &durationSeconds,
			Resolution:      &resolution,
			FPS:             &fps,
		},
		User:  &userid,
		Async: &async,
	}
	return callMediaAsyncAndSave(ctx, client, req, "generated_video.mp4")
}

// AnimateImage создаёт анимацию из изображения, возвращает JSON.
func AnimateImage(ctx context.Context, client *Client, modelname string, prompt, imageURL string, durationSeconds string, fps int, resolution string, userid string) ([]byte, error) {
	async := true
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Prompt: &prompt,
			Images: []models.MediaFile{
				{Type: "url", Data: imageURL},
			},
			DurationSeconds: &durationSeconds,
			FPS:             &fps,
			Resolution:      &resolution,
		},
		User:  &userid,
		Async: &async,
	}
	return callMediaAsyncAndSave(ctx, client, req, "animated_video.mp4")
}

// EditVideo редактирует видео по тексту, возвращает JSON.
func EditVideo(ctx context.Context, client *Client, modelname string, prompt, videoURL string, strength float64, durationSeconds string, userid string) ([]byte, error) {
	async := true
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Prompt: &prompt,
			Videos: []models.MediaFile{
				{Type: "url", Data: videoURL},
			},
			Strength:        &strength,
			DurationSeconds: &durationSeconds,
		},
		User:  &userid,
		Async: &async,
	}
	return callMediaAsyncAndSave(ctx, client, req, "edited_video.mp4")
}

// ExtendVideo продлевает видео, возвращает JSON.
func ExtendVideo(ctx context.Context, client *Client, modelname string, prompt, videoURL string, durationSeconds string, userid string) ([]byte, error) {
	async := true
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Prompt: &prompt,
			Videos: []models.MediaFile{
				{Type: "url", Data: videoURL},
			},
			DurationSeconds: &durationSeconds,
		},
		User:  &userid,
		Async: &async,
	}
	return callMediaAsyncAndSave(ctx, client, req, "extended_video.mp4")
}

// GenerateAudio синтезирует речь (TTS), возвращает JSON.
func GenerateAudio(ctx context.Context, client *Client, modelname, text, voice string, speed float64, format string, userid string) ([]byte, error) {
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Prompt:       &text,
			Voice:        &voice,
			Speed:        &speed,
			OutputFormat: &format,
		},
		User: &userid,
	}
	return callMediaAndSave(ctx, client, req, "generated_audio.mp3")
}

// TranscribeAudio распознаёт речь (STT), возвращает JSON.
func TranscribeAudio(ctx context.Context, client *Client, modelname string, audioURL, language string, userid string) ([]byte, error) {
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Audio: &models.MediaFile{
				Type: "url",
				Data: audioURL,
			},
			Language: &language,
		},
		User: &userid,
	}
	resp, err := client.Media().CreateMedia(ctx, req)
	if err != nil {
		return nil, err
	}
	return json.MarshalIndent(resp, "", "  ")
}

// GenerateMusic создаёт музыку через Suno, возвращает JSON.
func GenerateMusic(ctx context.Context, client *Client, modelname string, prompt string, customMode, instrumental bool, version string, userid string) ([]byte, error) {
	req := models.MediaRequest{
		Model: modelname,
		Input: models.MediaInput{
			Prompt:       &prompt,
			CustomMode:   &customMode,
			Instrumental: &instrumental,
			Version:      &version,
		},
		User: &userid,
	}
	return callMediaAndSave(ctx, client, req, "generated_music.mp3")
}

// ChatWithAudio отправляет запрос к gpt-audio-mini, возвращает JSON.
func ChatWithAudio(ctx context.Context, client *Client, userMessage, voice, format string, userid string) ([]byte, error) {
	req := models.ChatCompletionRequest{
		Model: "openai/gpt-audio-mini",
		Messages: []models.Message{
			{Role: "user", Content: userMessage},
		},
		Modalities: []string{"text", "audio"},
		Audio: &models.AudioOutputConfig{
			Voice:  voice,
			Format: format,
		},
		User: &userid,
	}
	resp, err := client.Chat().CreateCompletion(ctx, req)
	if err != nil {
		return nil, err
	}

	// Сохраняем аудио в хранилище, если есть
	if len(resp.Choices) > 0 && resp.Choices[0].Message.Audio != nil {
		audioData := resp.Choices[0].Message.Audio.Data
		if audioData != "" {
			savedURL, id, err := saveToStorage(ctx, client, audioData, "generated_audio.wav")
			if err != nil {
				fmt.Printf("Предупреждение: не удалось сохранить аудио в хранилище: %v\n", err)
			} else {
				resp.Choices[0].Message.Audio.Data = savedURL
				// Создаём отдельную структуру для ответа с ID
				responseWithID := map[string]interface{}{
					"completion": resp,
					"stored_id":  id,
				}
				return json.MarshalIndent(responseWithID, "", "  ")
			}
		}
	}
	return json.MarshalIndent(resp, "", "  ")
}

// Внутренние вспомогательные функции
func callMediaAndSave(ctx context.Context, client *Client, req models.MediaRequest, filename string) ([]byte, error) {
	resp, err := client.Media().CreateMedia(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка запроса: %w", err)
	}

	// Извлекаем все исходные URL из ответа
	urls := extractAllMediaURLs(resp)
	if len(urls) > 0 {
		saved := make([]map[string]interface{}, 0, len(urls))
		for i, url := range urls {
			savedURL, id, err := saveToStorage(ctx, client, url, filename)
			if err != nil {
				fmt.Printf("Предупреждение: не удалось сохранить файл %d в хранилище: %v\n", i+1, err)
				continue
			}
			saved = append(saved, map[string]interface{}{
				"original_url": url,
				"stored_url":   savedURL,
				"stored_id":    id,
			})
		}
		if len(saved) > 0 {
			resp.Data = saved
		}
	}
	return json.MarshalIndent(resp, "", "  ")
}

func callMediaAsyncAndSave(ctx context.Context, client *Client, req models.MediaRequest, filename string) ([]byte, error) {
	if req.Async == nil {
		asyncTrue := true
		req.Async = &asyncTrue
	}
	createResp, err := client.Media().CreateMedia(ctx, req)
	if err != nil {
		return nil, fmt.Errorf("ошибка создания задачи: %w", err)
	}
	finalResp, err := waitForMediaCompletion(ctx, client, createResp.ID)
	if err != nil {
		return nil, err
	}
	// Извлекаем все URL
	urls := extractAllMediaURLs(finalResp)
	if len(urls) > 0 {
		saved := []map[string]interface{}{}
		for i, url := range urls {
			savedURL, id, err := saveToStorage(ctx, client, url, filename)
			if err != nil {
				fmt.Printf("Предупреждение: не удалось сохранить файл %d в хранилище: %v\n", i+1, err)
				continue
			}
			saved = append(saved, map[string]interface{}{
				"original_url": url,
				"stored_url":   savedURL,
				"stored_id":    id,
			})
		}
		if len(saved) > 0 {
			finalResp.Data = saved
		}
	}
	return json.MarshalIndent(finalResp, "", "  ")
}

// extractAllMediaURLs извлекает все URL из ответа MediaResponse
func extractAllMediaURLs(resp *models.MediaResponse) []string {
	var urls []string
	if resp.Data == nil {
		return urls
	}

	// Если Data — массив
	if dataSlice, ok := resp.Data.([]interface{}); ok {
		for _, item := range dataSlice {
			if m, ok := item.(map[string]interface{}); ok {
				if u, ok := m["url"].(string); ok && u != "" {
					urls = append(urls, u)
				}
			}
		}
		return urls
	}

	// Если Data — объект
	if dataMap, ok := resp.Data.(map[string]interface{}); ok {
		if u, ok := dataMap["url"].(string); ok && u != "" {
			urls = append(urls, u)
		}
		return urls
	}

	// Если Data — строка (URL)
	if urlStr, ok := resp.Data.(string); ok && urlStr != "" {
		urls = append(urls, urlStr)
	}

	return urls
}

// Вспомогательные функции для создания указателей (скрыты внутри)
func stringPtr(s string) *string    { return &s }
func intPtr(i int) *int             { return &i }
func float64Ptr(f float64) *float64 { return &f }
func boolPtr(b bool) *bool          { return &b }
