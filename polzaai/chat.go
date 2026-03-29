// chat.go
package polzaai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"strings"

	"github.com/nvwrist/polzaai/polzaai/models"
)

// Chat возвращает вспомогательный объект для работы с чатом.
func (c *Client) Chat() *ChatService {
	return &ChatService{client: c}
}

type ChatService struct {
	client *Client
}

// CreateCompletion отправляет запрос на генерацию текста (синхронно).
func (s *ChatService) CreateCompletion(ctx context.Context, req models.ChatCompletionRequest) (*models.ChatCompletionResponse, error) {
	var result models.ChatCompletionResponse
	err := s.client.doRequest(ctx, "POST", "/chat/completions", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// CreateCompletionStream запускает потоковую генерацию и возвращает канал с чанками.
// При ошибке в канал будет отправлена ошибка, затем канал закроется.
func (s *ChatService) CreateCompletionStream(ctx context.Context, req models.ChatCompletionRequest) (<-chan models.ChatCompletionChunk, <-chan error) {
	chunkChan := make(chan models.ChatCompletionChunk)
	errChan := make(chan error, 1)

	// Принудительно включаем стриминг
	streamTrue := true
	req.Stream = &streamTrue

	url := s.client.baseURL + "/chat/completions"
	jsonData, err := json.Marshal(req)
	if err != nil {
		errChan <- fmt.Errorf("сериализация запроса: %w", err)
		close(chunkChan)
		close(errChan)
		return chunkChan, errChan
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		errChan <- fmt.Errorf("создание запроса: %w", err)
		close(chunkChan)
		close(errChan)
		return chunkChan, errChan
	}
	httpReq.Header.Set("Authorization", "Bearer "+s.client.apiKey)
	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Accept", "text/event-stream")

	resp, err := s.client.httpClient.Do(httpReq)
	if err != nil {
		errChan <- fmt.Errorf("выполнение запроса: %w", err)
		close(chunkChan)
		close(errChan)
		return chunkChan, errChan
	}

	go func() {
		defer resp.Body.Close()
		defer close(chunkChan)
		defer close(errChan)

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			line := scanner.Text()
			if !strings.HasPrefix(line, "data: ") {
				continue
			}
			data := strings.TrimPrefix(line, "data: ")
			if data == "[DONE]" {
				return
			}
			var chunk models.ChatCompletionChunk
			if err := json.Unmarshal([]byte(data), &chunk); err != nil {
				errChan <- fmt.Errorf("ошибка парсинга чанка: %w", err)
				return
			}
			chunkChan <- chunk
		}
		if err := scanner.Err(); err != nil {
			errChan <- fmt.Errorf("ошибка чтения стрима: %w", err)
		}
	}()

	return chunkChan, errChan
}
