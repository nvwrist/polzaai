package models

// models/models.go – добавим структуры для ответа /v1/models
type ModelsResponse struct {
	Data []ModelInfo `json:"data"`
}

type ModelInfo struct {
	ID               string                   `json:"id"`
	Name             string                   `json:"name"`
	Type             string                   `json:"type"` // chat, image, video, tts, stt, music
	ShortDescription string                   `json:"short_description"`
	Created          int64                    `json:"created"`
	Architecture     Architecture             `json:"architecture"`
	TopProvider      ProviderDetails          `json:"top_provider"`
	Endpoints        []string                 `json:"endpoints"`
	Parameters       map[string]ParameterInfo `json:"parameters,omitempty"` // специфичные для модели параметры
}

type Architecture struct {
	Modality         string   `json:"modality"`
	InputModalities  []string `json:"input_modalities"`
	OutputModalities []string `json:"output_modalities"`
	Tokenizer        string   `json:"tokenizer"`
	InstructType     *string  `json:"instruct_type"`
}

type ProviderDetails struct {
	IsModerated         bool                     `json:"is_moderated"`
	ContextLength       *int                     `json:"context_length"`
	MaxCompletionTokens *int                     `json:"max_completion_tokens"`
	Pricing             Pricing                  `json:"pricing"`
	SupportedParameters []string                 `json:"supported_parameters"`
	DefaultParameters   map[string]interface{}   `json:"default_parameters"`
	Parameters          map[string]ParameterInfo `json:"parameters,omitempty"`
}

type ParameterInfo struct {
	Required    bool        `json:"required,omitempty"`
	Description string      `json:"description,omitempty"`
	Min         *float64    `json:"min,omitempty"`
	Max         *float64    `json:"max,omitempty"`
	MaxLength   *int        `json:"max_length,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Values      []string    `json:"values,omitempty"`
}

type Pricing struct {
	PromptPerMillion         string      `json:"prompt_per_million,omitempty"`
	CompletionPerMillion     string      `json:"completion_per_million,omitempty"`
	InputCacheReadPerMillion string      `json:"input_cache_read_per_million,omitempty"`
	Tiers                    []PriceTier `json:"tiers,omitempty"`
	PerRequest               string      `json:"per_request,omitempty"`
	TTSPerMillionCharacters  string      `json:"tts_per_million_characters,omitempty"`
	STTPerMinute             string      `json:"stt_per_minute,omitempty"`
	Currency                 string      `json:"currency"`
}

type PriceTier struct {
	Conditions []string `json:"conditions"`
	CostRub    string   `json:"cost_rub"`
}
