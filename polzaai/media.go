// media.go
package polzaai

import (
	"context"
	"polzasdk/polzaai/models"
)

// MediaService для работы с /v1/media.
type MediaService struct {
	client *Client
}

func (c *Client) Media() *MediaService {
	return &MediaService{client: c}
}

// CreateMedia отправляет запрос на генерацию медиа.
// Если async = true, возвращает MediaResponse со статусом pending, результат нужно получать через GetMediaStatus.
func (s *MediaService) CreateMedia(ctx context.Context, req models.MediaRequest) (*models.MediaResponse, error) {
	var result models.MediaResponse
	err := s.client.doRequest(ctx, "POST", "/media", req, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetMediaStatus получает текущий статус асинхронной генерации по ID.
func (s *MediaService) GetMediaStatus(ctx context.Context, mediaID string) (*models.MediaResponse, error) {
	var result models.MediaResponse
	err := s.client.doRequest(ctx, "GET", "/media/"+mediaID, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
