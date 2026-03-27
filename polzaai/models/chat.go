// models/chat.go
package models

// ChatCompletionRequest параметры запроса к /chat/completions.
type ChatCompletionRequest struct {
	Model               string             `json:"model"`
	Messages            []Message          `json:"messages,omitempty"`
	Prompt              *string            `json:"prompt,omitempty"`
	MaxTokens           *int               `json:"max_tokens,omitempty"`
	MaxCompletionTokens *int               `json:"max_completion_tokens,omitempty"`
	Temperature         *float64           `json:"temperature,omitempty"`
	TopP                *float64           `json:"top_p,omitempty"`
	TopK                *int               `json:"top_k,omitempty"`
	FrequencyPenalty    *float64           `json:"frequency_penalty,omitempty"`
	PresencePenalty     *float64           `json:"presence_penalty,omitempty"`
	Stop                []string           `json:"stop,omitempty"`
	Seed                *int               `json:"seed,omitempty"`
	N                   *int               `json:"n,omitempty"`
	Stream              *bool              `json:"stream,omitempty"`
	Logprobs            *bool              `json:"logprobs,omitempty"`
	TopLogprobs         *int               `json:"top_logprobs,omitempty"`
	LogitBias           map[string]int     `json:"logit_bias,omitempty"`
	ParallelToolCalls   *bool              `json:"parallel_tool_calls,omitempty"`
	ResponseFormat      *ResponseFormat    `json:"response_format,omitempty"`
	Provider            *ProviderConfig    `json:"provider,omitempty"`
	Tools               []Tool             `json:"tools,omitempty"`
	ToolChoice          interface{}        `json:"tool_choice,omitempty"` // string или ToolChoiceObject
	Reasoning           *ReasoningConfig   `json:"reasoning,omitempty"`
	Plugins             interface{}        `json:"plugins,omitempty"`
	WebSearchOptions    *WebSearchOptions  `json:"web_search_options,omitempty"`
	User                *string            `json:"user,omitempty"`
	ImageConfig         *ImageConfig       `json:"image_config,omitempty"`
	Modalities          []string           `json:"modalities,omitempty"` // "text", "image", "audio"
	Audio               *AudioOutputConfig `json:"audio,omitempty"`
}

type ResponseFormat struct {
	Type       string         `json:"type"` // "text", "json_object", "json_schema", "grammar"
	JSONSchema *JSONSchema    `json:"json_schema,omitempty"`
	Grammar    *GrammarConfig `json:"grammar,omitempty"`
}

type JSONSchema struct {
	Name   string                 `json:"name"`
	Schema map[string]interface{} `json:"schema"`
	Strict *bool                  `json:"strict,omitempty"`
}

type GrammarConfig struct {
	Grammar string `json:"grammar"` // GBNF
}

type Tool struct {
	Type     string       `json:"type"` // "function"
	Function FunctionTool `json:"function"`
}

type FunctionTool struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  map[string]interface{} `json:"parameters"`
	Strict      *bool                  `json:"strict,omitempty"`
}

type ToolChoiceObject struct {
	Type     string             `json:"type"`
	Function ToolChoiceFunction `json:"function"`
}

type ToolChoiceFunction struct {
	Name string `json:"name"`
}

type ReasoningConfig struct {
	Effort    *string `json:"effort,omitempty"` // "xhigh", "high", "medium", "low", "minimal", "none"
	MaxTokens *int    `json:"max_tokens,omitempty"`
	Summary   *string `json:"summary,omitempty"` // "auto", "concise", "detailed"
	Enabled   *bool   `json:"enabled,omitempty"`
	Exclude   *bool   `json:"exclude,omitempty"`
}

type WebSearchOptions struct {
	SearchContextSize *string `json:"search_context_size,omitempty"` // "small", "medium", "large"
}

type ImageConfig struct {
	Quality *string `json:"quality,omitempty"` // "standard", "high"
	Size    *int    `json:"size,omitempty"`
}

type AudioOutputConfig struct {
	Voice  string `json:"voice"`            // "alloy", "echo", "fable", "onyx", "nova", "shimmer"
	Format string `json:"format,omitempty"` // "pcm16", "mp3", "opus"
}

// ChatCompletionResponse ответ от API.
type ChatCompletionResponse struct {
	ID                string       `json:"id"`
	Object            string       `json:"object"`
	Created           int64        `json:"created"`
	Model             string       `json:"model"`
	Provider          *string      `json:"provider,omitempty"`
	Choices           []ChatChoice `json:"choices"`
	SystemFingerprint *string      `json:"system_fingerprint,omitempty"`
	Usage             *Usage       `json:"usage,omitempty"`
}

type ChatChoice struct {
	Index        int             `json:"index"`
	Message      ResponseMessage `json:"message"`
	FinishReason string          `json:"finish_reason"`
	Logprobs     *Logprobs       `json:"logprobs,omitempty"`
}

type ResponseMessage struct {
	Role      string       `json:"role"`
	Content   *string      `json:"content,omitempty"`
	Name      *string      `json:"name,omitempty"`
	ToolCalls []ToolCall   `json:"tool_calls,omitempty"`
	Refusal   *string      `json:"refusal,omitempty"`
	Reasoning *string      `json:"reasoning,omitempty"`
	Audio     *AudioOutput `json:"audio,omitempty"`
}

type ToolCall struct {
	ID       string       `json:"id"`
	Type     string       `json:"type"` // "function"
	Function FunctionCall `json:"function"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"` // JSON строка
}

type AudioOutput struct {
	ID         string `json:"id"`
	Data       string `json:"data"` // base64
	Transcript string `json:"transcript"`
	ExpiresAt  int64  `json:"expires_at"`
}

type Logprobs struct {
	Content []LogprobItem `json:"content,omitempty"`
	Refusal []LogprobItem `json:"refusal,omitempty"`
}

type LogprobItem struct {
	Token       string           `json:"token"`
	Logprob     float64          `json:"logprob"`
	Bytes       []int            `json:"bytes,omitempty"`
	TopLogprobs []TopLogprobItem `json:"top_logprobs,omitempty"`
}

type TopLogprobItem struct {
	Token   string  `json:"token"`
	Logprob float64 `json:"logprob"`
	Bytes   []int   `json:"bytes,omitempty"`
}

// Chunk используется для стриминга.
type ChatCompletionChunk struct {
	ID      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []ChoiceDelta `json:"choices"`
	Usage   *Usage        `json:"usage,omitempty"`
}

type ChoiceDelta struct {
	Index        int          `json:"index"`
	Delta        MessageDelta `json:"delta"`
	FinishReason *string      `json:"finish_reason,omitempty"`
}

type MessageDelta struct {
	Role      *string         `json:"role,omitempty"`
	Content   *string         `json:"content,omitempty"`
	ToolCalls []ToolCallDelta `json:"tool_calls,omitempty"`
}

type ToolCallDelta struct {
	Index    *int               `json:"index,omitempty"`
	ID       *string            `json:"id,omitempty"`
	Type     *string            `json:"type,omitempty"`
	Function *FunctionCallDelta `json:"function,omitempty"`
}

type FunctionCallDelta struct {
	Name      *string `json:"name,omitempty"`
	Arguments *string `json:"arguments,omitempty"`
}
