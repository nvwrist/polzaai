// models/media.go
package models

// MediaRequest параметры для /v1/media.
type MediaRequest struct {
	Model    string          `json:"model"`
	Input    MediaInput      `json:"input"`
	Async    *bool           `json:"async,omitempty"`
	User     *string         `json:"user,omitempty"`
	Provider *ProviderConfig `json:"provider,omitempty"`
}

// MediaInput содержит все возможные поля для генерации изображений, видео, аудио.
// Поля, не поддерживаемые конкретной моделью, будут проигнорированы с предупреждением.
type MediaInput struct {
	Prompt                    *string     `json:"prompt,omitempty"`
	AspectRatio               *string     `json:"aspect_ratio,omitempty"` // "16:9", "1:1", "9:16"
	Images                    []MediaFile `json:"images,omitempty"`       // для img2img, edit
	Videos                    []MediaFile `json:"videos,omitempty"`       // для vid2vid
	Audio                     *MediaFile  `json:"audio,omitempty"`        // для TTS или STT
	CallBackUrl               *string     `json:"callBackUrl,omitempty"`  // webhook после завершения
	Seed                      *int64      `json:"seed,omitempty"`
	Watermark                 *string     `json:"watermark,omitempty"`
	ImageResolution           *string     `json:"image_resolution,omitempty"` // "1K", "2K", "4K"
	Quality                   *string     `json:"quality,omitempty"`          // "standard", "high"
	OutputFormat              *string     `json:"output_format,omitempty"`    // "png", "jpg", "mp4", "mp3"
	MaxImages                 *int        `json:"max_images,omitempty"`       // до 10
	IsEnhance                 *bool       `json:"isEnhance,omitempty"`
	GuidanceScale             *float64    `json:"guidance_scale,omitempty"` // 1-20
	Strength                  *float64    `json:"strength,omitempty"`       // 0-1
	EnableSafetyChecker       *bool       `json:"enable_safety_checker,omitempty"`
	UpscaleFactor             *string     `json:"upscale_factor,omitempty"` // "2", "4"
	FontInputs                []FontInput `json:"font_inputs,omitempty"`
	SuperResolutionReferences []string    `json:"super_resolution_references,omitempty"`
	// Поля для аудио TTS
	Voice    *string  `json:"voice,omitempty"`    // "alloy", "echo", ...
	Speed    *float64 `json:"speed,omitempty"`    // 0.25-4.0
	Language *string  `json:"language,omitempty"` // код языка
	// Поля для видео
	DurationSeconds *string `json:"duration,omitempty"`
	FPS             *int    `json:"fps,omitempty"`
	Resolution      *string `json:"resolution,omitempty"` // "720p", "1080p"
	//для SUNO
	CustomMode   *bool   `json:"customMode,omitempty"`
	Instrumental *bool   `json:"instrumental,omitempty"`
	Version      *string `json:"version,omitempty"`
}

type MediaFile struct {
	Type string `json:"type"` // "url" или "base64"
	Data string `json:"data"` // URL или base64-строка (с data URI или без)
}

type FontInput struct {
	FontURL string `json:"font_url"`
	Text    string `json:"text"`
}

// MediaResponse ответ на создание медиа-генерации.
type MediaResponse struct {
	ID               string      `json:"id"`
	Object           string      `json:"object"`
	Status           string      `json:"status"` // pending, processing, completed, failed, cancelled
	Created          int64       `json:"created"`
	Model            string      `json:"model"`
	CompletedAt      *int64      `json:"completed_at,omitempty"`
	Data             interface{} `json:"data,omitempty"` // может быть URL или base64
	Usage            *MediaUsage `json:"usage,omitempty"`
	Error            *ErrorInfo  `json:"error,omitempty"`
	Content          *string     `json:"content,omitempty"` // текстовый ответ модели
	ReasoningSummary *string     `json:"reasoning_summary,omitempty"`
	Warnings         []string    `json:"warnings,omitempty"`
}

type MediaUsage struct {
	InputUnits      *int     `json:"input_units,omitempty"`
	OutputUnits     *int     `json:"output_units,omitempty"`
	DurationSeconds *float64 `json:"duration_seconds,omitempty"`
	InputTokens     *int     `json:"input_tokens,omitempty"`
	OutputTokens    *int     `json:"output_tokens,omitempty"`
	TotalTokens     *int     `json:"total_tokens,omitempty"`
	CostRub         *float64 `json:"cost_rub,omitempty"`
	Cost            *float64 `json:"cost,omitempty"`
}

type ErrorInfo struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
