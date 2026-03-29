// client.go
package polzaai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/nvwrist/polzaai/polzaai/models"

	"io"
	"net/http"
	"time"
)

// Client представляет клиент Polza.ai API.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// NewClient создаёт новый клиент.
// apiKey: ваш API-ключ, полученный в консоли Polza.ai.
func NewClient(apiKey string) *Client {
	return &Client{
		baseURL: "https://polza.ai/api/v1",
		apiKey:  apiKey,
		httpClient: &http.Client{
			Timeout: 60 * time.Second, // можно переопределить через опции
		},
	}
}

// doRequest выполняет HTTP-запрос с JSON-телом и возвращает ответ.
// ctx - контекст для отмены/таймаута.
// method, path, body - параметры запроса.
// result - указатель на структуру, в которую будет распаршен JSON ответа.
func (c *Client) doRequest(ctx context.Context, method, path string, body interface{}, result interface{}) error {
	url := c.baseURL + path

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return fmt.Errorf("ошибка сериализации запроса: %w", err)
		}
		reqBody = bytes.NewBuffer(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, method, url, reqBody)
	if err != nil {
		return fmt.Errorf("ошибка создания запроса: %w", err)
	}

	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("ошибка выполнения запроса: %w", err)
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("ошибка чтения ответа: %w", err)
	}

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var errResp models.ErrorResponse
		if err := json.Unmarshal(respBody, &errResp); err != nil {
			return fmt.Errorf("API вернул статус %d, тело: %s", resp.StatusCode, string(respBody))
		}
		return fmt.Errorf("API ошибка: %s - %s", errResp.Error.Code, errResp.Error.Message)
	}

	if result != nil {
		if err := json.Unmarshal(respBody, result); err != nil {
			return fmt.Errorf("ошибка парсинга ответа: %w", err)
		}
	}
	return nil
}

// models.go – в client.go добавим
func (c *Client) ListModels(ctx context.Context) (*models.ModelsResponse, error) {
	var result models.ModelsResponse
	err := c.doRequest(ctx, "GET", "/models", nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (c *Client) GetModelInfo(ctx context.Context, modelID string) (*models.ModelInfo, error) {
	resp, err := c.ListModels(ctx)
	if err != nil {
		return nil, err
	}
	for _, m := range resp.Data {
		if m.ID == modelID {
			return &m, nil
		}
	}
	return nil, fmt.Errorf("модель %s не найдена", modelID)
}

// validateMediaRequest – проверяет, что все обязательные параметры указаны и значения допустимы
func (c *Client) validateMediaRequest(ctx context.Context, modelID string, input map[string]interface{}) error {
	modelInfo, err := c.GetModelInfo(ctx, modelID)
	if err != nil {
		return err
	}
	// Если модель не использует media API (например, chat), пропускаем
	if !contains(modelInfo.Endpoints, "/api/v1/media") {
		return nil
	}
	// Получаем спецификацию параметров
	params := modelInfo.Parameters
	if params == nil && modelInfo.TopProvider.Parameters != nil {
		params = modelInfo.TopProvider.Parameters
	}
	for name, spec := range params {
		val, exists := input[name]
		if spec.Required && !exists {
			return fmt.Errorf("обязательный параметр %s отсутствует", name)
		}
		if exists {
			// Проверка типа и допустимых значений
			switch v := val.(type) {
			case string:
				if len(spec.Values) > 0 {
					found := false
					for _, allowed := range spec.Values {
						if v == allowed {
							found = true
							break
						}
					}
					if !found {
						return fmt.Errorf("параметр %s имеет недопустимое значение %s, допустимые: %v", name, v, spec.Values)
					}
				}
				if spec.MaxLength != nil && len(v) > *spec.MaxLength {
					return fmt.Errorf("параметр %s превышает максимальную длину %d", name, *spec.MaxLength)
				}
			case float64:
				if spec.Min != nil && v < *spec.Min {
					return fmt.Errorf("параметр %s меньше минимального %v", name, *spec.Min)
				}
				if spec.Max != nil && v > *spec.Max {
					return fmt.Errorf("параметр %s больше максимального %v", name, *spec.Max)
				}
			}
		}
	}
	return nil
}

// contains проверяет, содержится ли строка в слайсе.
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
