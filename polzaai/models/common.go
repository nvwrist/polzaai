// models/common.go
package models

// ErrorResponse структура ошибки API.
type ErrorResponse struct {
	Error struct {
		Code    string `json:"code"`
		Message string `json:"message"`
	} `json:"error"`
}

// Message представляет сообщение в диалоге.
type Message struct {
	Role    string      `json:"role"`    // "system", "user", "assistant", "tool", "developer"
	Content interface{} `json:"content"` // может быть string или []ContentPart для мультимодальности
	Name    *string     `json:"name,omitempty"`
}

// ContentPart используется для мультимодальных сообщений.
type ContentPart struct {
	Type     string     `json:"type"` // "text", "image_url", "input_audio", "video_url", "file"
	Text     *string    `json:"text,omitempty"`
	ImageURL *ImageURL  `json:"image_url,omitempty"`
	Audio    *AudioData `json:"input_audio,omitempty"`
	VideoURL *VideoURL  `json:"video_url,omitempty"`
	File     *FileData  `json:"file,omitempty"`
}

type ImageURL struct {
	URL    string `json:"url"`
	Detail string `json:"detail,omitempty"`
}

type AudioData struct {
	Data   string `json:"data"`   // base64
	Format string `json:"format"` // "mp3", "wav", "pcm16"
}

type VideoURL struct {
	URL string `json:"url"`
}

type FileData struct {
	Filename string `json:"filename"`
	FileData string `json:"file_data"` // data URI или base64
}

// ProviderConfig настройки роутинга между провайдерами.
type ProviderConfig struct {
	AllowFallbacks *bool           `json:"allow_fallbacks,omitempty"`
	Order          []string        `json:"order,omitempty"`
	Only           []string        `json:"only,omitempty"`
	Ignore         []string        `json:"ignore,omitempty"`
	Sort           *string         `json:"sort,omitempty"` // "price", "latency"
	MaxPrice       *MaxPriceConfig `json:"max_price,omitempty"`
}

type MaxPriceConfig struct {
	Prompt     *float64 `json:"prompt,omitempty"`
	Completion *float64 `json:"completion,omitempty"`
	Image      *float64 `json:"image,omitempty"`
	Audio      *float64 `json:"audio,omitempty"`
	Request    *float64 `json:"request,omitempty"`
}

// Usage информация об использованных токенах/единицах.
type Usage struct {
	PromptTokens            int                      `json:"prompt_tokens"`
	CompletionTokens        int                      `json:"completion_tokens"`
	TotalTokens             int                      `json:"total_tokens"`
	CostRub                 *float64                 `json:"cost_rub,omitempty"`
	Cost                    *float64                 `json:"cost,omitempty"`
	CompletionTokensDetails *CompletionTokensDetails `json:"completion_tokens_details,omitempty"`
	PromptTokensDetails     *PromptTokensDetails     `json:"prompt_tokens_details,omitempty"`
	ServerToolUse           *ServerToolUse           `json:"server_tool_use,omitempty"`
}

type CompletionTokensDetails struct {
	ReasoningTokens          *int `json:"reasoning_tokens,omitempty"`
	AudioTokens              *int `json:"audio_tokens,omitempty"`
	ImageTokens              *int `json:"image_tokens,omitempty"`
	AcceptedPredictionTokens *int `json:"accepted_prediction_tokens,omitempty"`
	RejectedPredictionTokens *int `json:"rejected_prediction_tokens,omitempty"`
}

type PromptTokensDetails struct {
	CachedTokens *int `json:"cached_tokens,omitempty"`
	AudioTokens  *int `json:"audio_tokens,omitempty"`
	VideoTokens  *int `json:"video_tokens,omitempty"`
}

type ServerToolUse struct {
	WebSearchRequests *int `json:"web_search_requests,omitempty"`
}
