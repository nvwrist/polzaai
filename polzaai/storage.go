package polzaai

import (
	"context"
	"path/filepath"
	"polzasdk/polzaai/models"
	"strings"
)

type StorageService struct {
	client *Client
}

func (c *Client) Storage() *StorageService {
	return &StorageService{client: c}
}

// UploadFile загружает файл в хранилище Polza.ai через JSON.
func (s *StorageService) UploadFile(ctx context.Context, req models.UploadFileRequest) (*models.UploadFileResponse, error) {
	// Очищаем base64 от префикса, если он есть
	base64Data := req.Data
	if idx := strings.Index(base64Data, ";base64,"); idx != -1 {
		base64Data = base64Data[idx+8:] // пропускаем ";base64,"
	}

	// Формируем JSON-тело
	jsonBody := map[string]interface{}{
		"base64":        base64Data,
		"storagePolicy": req.Policy,
	}

	// Определяем MIME-тип по расширению файла
	if req.Filename != "" {
		ext := strings.ToLower(filepath.Ext(req.Filename))
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
		jsonBody["mimeType"] = mimeType
	}

	// Отправляем POST с JSON
	var result models.UploadFileResponse
	err := s.client.doRequest(ctx, "POST", "/storage/upload", jsonBody, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

// GetFileInfo возвращает метаданные файла.
func (s *StorageService) GetFileInfo(ctx context.Context, fileID string) (*models.FileInfo, error) {
	var result models.FileInfo
	err := s.client.doRequest(ctx, "GET", "/storage/files/"+fileID, nil, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}

func (s *StorageService) FileList(ctx context.Context) ([]models.FileInfo, error) {
	var resp models.FileListResponse
	err := s.client.doRequest(ctx, "GET", "/storage/files", nil, &resp)
	if err != nil {
		return nil, err
	}
	return resp.Data, nil
}

// DeleteFile удаляет файл.
func (s *StorageService) DeleteFile(ctx context.Context, fileID string) error {
	return s.client.doRequest(ctx, "DELETE", "/storage/"+fileID, nil, nil)
}
