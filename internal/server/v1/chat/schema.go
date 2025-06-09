package chat

type ChatMessage struct {
	Role    string `json:"role" validate:"required,oneof=user assistant system"`
	Content string `json:"content" validate:"required"`
	Name    string `json:"name,omitempty"`
}

type FunctionCall struct {
	Name string `json:"name" validate:"required"`
}

type ToolChoice struct {
	Name string `json:"name" validate:"required"`
	Type string `json:"type" validate:"required,oneof=function"`
}

type ToolFunctionParameters map[string]any

type ToolFunction struct {
	Name        string                 `json:"name" validate:"required"`
	Description string                 `json:"description,omitempty"`
	Parameters  ToolFunctionParameters `json:"parameters" validate:"required,dive" swagger:"object"`
}

type Tool struct {
	Type     string       `json:"type" validate:"required,oneof=function"`
	Function ToolFunction `json:"function" validate:"required,dive"`
}

type ApproximateLocation struct {
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty" validate:"omitempty,max=2,min=2"`
	Region   string `json:"region,omitempty"`
	TimeZone string `json:"timezone,omitempty" validate:"omitempty,max=64"`
}

type WebSearchUserLocation struct {
	Type        string              `json:"type" validate:"required,oneof=approximate"`
	Approximate ApproximateLocation `json:"approximate" validate:"required,dive"`
}

type WebSearchOptions struct {
	SearchContextSize string                 `json:"search_context_size,omitempty" validate:"omitempty,oneof=small medium large"`
	UserLocation      *WebSearchUserLocation `json:"user_location,omitempty" validate:"omitempty,dive"`
}

type ChatCompletionRequest struct {
	Messages            []ChatMessage     `json:"messages" validate:"required,dive"`
	Model               string            `json:"model" validate:"required"`
	FrequencyPenalty    float64           `json:"frequency_penalty,omitempty" validate:"omitempty,min=-2.0,max=2.0,default=0.0"`
	FunctionCall        *FunctionCall     `json:"function_call,omitempty" validate:"omitempty,dive"`
	LogProbs            bool              `json:"logprobs,omitempty" validate:"omitempty"`
	MaxCompletionTokens int               `json:"max_completion_tokens,omitempty" validate:"omitempty,min=1,max=4096"`
	N                   int               `json:"n,omitempty" validate:"omitempty,min=1,max=10,default=1"`
	ParallelToolCalls   bool              `json:"parallel_tool_calls,omitempty" validate:"omitempty"`
	PresencePenalty     float64           `json:"presence_penalty,omitempty" validate:"omitempty,min=-2.0,max=2.0,default=0.0"`
	ReasoningEffort     string            `json:"reasoning_effort,omitempty" validate:"omitempty,oneof=low medium high"`
	ResponseFormat      string            `json:"response_format,omitempty" validate:"omitempty,oneof=text json"`
	Stream              bool              `json:"stream,omitempty" validate:"omitempty, default=false"`
	Temperature         float64           `json:"temperature,omitempty" validate:"omitempty,min=0.0,max=2.0,default=1.0"`
	Tools               []Tool            `json:"tools,omitempty" validate:"omitempty,dive"`
	TopLogProbs         int               `json:"top_logprobs,omitempty" validate:"omitempty,min=0,max=20,"`
	TopP                float64           `json:"top_p,omitempty" validate:"omitempty,min=0.0,max=1.0,default=1.0"`
	User                string            `json:"user,omitempty" validate:"omitempty"`
	WebSearchOptions    *WebSearchOptions `json:"web_search_options,omitempty" validate:"omitempty,dive"`
}

type TokenLogProb struct {
	Token   string  `json:"token" validate:"required"`
	LogProb float64 `json:"logprob" validate:"required"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

type ChatCompletionLogProb struct {
	TokenLogProb
	TopLogProbs []TokenLogProb `json:"top_logprobs,omitempty" validate:"omitempty,dive"`
}

type ChatCompletionLogProbs struct {
	Content []ChatCompletionLogProb `json:"content,omitempty" validate:"omitempty,dive"`
	Refusal []ChatCompletionLogProb `json:"refusal,omitempty" validate:"omitempty,dive"`
}

type ChatCompletionsMessageFunctionCall struct {
	Name      string `json:"name" validate:"required"`
	Arguments string `json:"arguments" validate:"required"`
}

type ChatCompletionsMessageToolCall struct {
	ID       string                             `json:"id" validate:"required"`
	Type     string                             `json:"type" validate:"required,oneof=function"`
	Function ChatCompletionsMessageFunctionCall `json:"function" validate:"required,dive"`
}

type URLCitation struct {
	StartIndex int    `json:"start_index" validate:"required,min=0"`
	EndIndex   int    `json:"end_index" validate:"required,min=0"`
	URL        string `json:"url" validate:"required,url"`
	Title      string `json:"title" validate:"required"`
}

type ChatCompletionResponseMessageAnnotation struct {
	Type        string      `json:"type" validate:"required,oneof=url_citation"`
	URLCitation URLCitation `json:"url_citation" validate:"required,dive"`
}

type ChatCompletionResponseMessage struct {
	Content string `json:"content" validate:"required"`
	Refusal bool   `json:"refusal,omitempty" validate:"omitempty"`
}

type ChatCompletionChoice struct {
	FinishReason string                  `json:"finish_reason" validate:"required,oneof=stop length tool_calls content_filter function_call"`
	Index        int                     `json:"index" validate:"required,min=0"`
	LogProbs     *ChatCompletionLogProbs `json:"logprobs,omitempty" validate:"omitempty,dive"`
	Message      ChatMessage             `json:"message" validate:"required,dive"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id" validate:"required"`
	Choices []ChatCompletionChoice `json:"choices" validate:"required,dive"`
	Created int64                  `json:"created" validate:"required"`
	Model   string                 `json:"model" validate:"required"`
	Object  string                 `json:"object" validate:"required,oneof=chat.completion"`
}
