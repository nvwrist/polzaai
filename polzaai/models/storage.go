// models/storage.go
package models

import "time"

type UploadFileRequest struct {
	Filename string `json:"filename"`
	Data     string `json:"data"`             // base64 или data URI
	Policy   string `json:"policy,omitempty"` // "TEMP" (7 дней) или "PERMANENT"
}

type UploadFileResponse struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	ExpiresAt *int64 `json:"expires_at,omitempty"`
}

type FileListResponse struct {
	Data []FileInfo `json:"data"`
}
type FileInfo struct {
	Id             string    `json:"id"`
	FileType       string    `json:"fileType"`
	MimeType       string    `json:"mimeType"`
	Source         string    `json:"source"`
	StoragePolicy  string    `json:"storagePolicy"`
	Url            string    `json:"url"`
	Size           int       `json:"size"`
	CreatedAt      time.Time `json:"createdAt"`
	UpdatedAt      time.Time `json:"updatedAt"`
	ExternalUserId string    `json:"externalUserId"`
	ExpiresAt      time.Time `json:"expiresAt"`
}
