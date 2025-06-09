package types

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
	Name    string `json:"name,omitempty"`
}

type FunctionCall struct {
	Name string `json:"name"`
}

type ToolChoice struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type ToolFunctionParameters map[string]any

type ToolFunction struct {
	Name        string                 `json:"name"`
	Description string                 `json:"description,omitempty"`
	Parameters  ToolFunctionParameters `json:"parameters"`
}

type Tool struct {
	Type     string       `json:"type"`
	Function ToolFunction `json:"function"`
}

type ApproximateLocation struct {
	City     string `json:"city,omitempty"`
	Country  string `json:"country,omitempty"`
	Region   string `json:"region,omitempty"`
	TimeZone string `json:"timezone,omitempty"`
}

type WebSearchUserLocation struct {
	Type        string              `json:"type"`
	Approximate ApproximateLocation `json:"approximate"`
}

type WebSearchOptions struct {
	SearchContextSize string                 `json:"search_context_size,omitempty"`
	UserLocation      *WebSearchUserLocation `json:"user_location,omitempty"`
}

type ChatCompletionRequest struct {
	Messages            []ChatMessage     `json:"messages"`
	Model               string            `json:"model"`
	FrequencyPenalty    *float64          `json:"frequency_penalty,omitempty"`
	FunctionCall        *FunctionCall     `json:"function_call,omitempty"`
	LogProbs            bool              `json:"logprobs,omitempty"`
	MaxCompletionTokens *int64            `json:"max_completion_tokens,omitempty"`
	N                   *int              `json:"n,omitempty"`
	ParallelToolCalls   bool              `json:"parallel_tool_calls,omitempty"`
	PresencePenalty     *float64          `json:"presence_penalty,omitempty"`
	ReasoningEffort     string            `json:"reasoning_effort,omitempty"`
	ResponseFormat      string            `json:"response_format,omitempty"`
	Stream              *bool             `json:"stream,omitempty"`
	Temperature         *float64          `json:"temperature,omitempty"`
	Tools               []Tool            `json:"tools,omitempty"`
	TopLogProbs         int               `json:"top_logprobs,omitempty"`
	TopP                *float64          `json:"top_p,omitempty"`
	User                string            `json:"user,omitempty"`
	WebSearchOptions    *WebSearchOptions `json:"web_search_options,omitempty"`
}

type TokenLogProb struct {
	Token   string  `json:"token"`
	LogProb float64 `json:"logprob"`
	Bytes   []byte  `json:"bytes,omitempty"`
}

type ChatCompletionLogProb struct {
	TokenLogProb
	TopLogProbs []TokenLogProb `json:"top_logprobs,omitempty"`
}

type ChatCompletionLogProbs struct {
	Content []ChatCompletionLogProb `json:"content,omitempty"`
	Refusal []ChatCompletionLogProb `json:"refusal,omitempty"`
}

type ChatCompletionsMessageFunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments"`
}

type ChatCompletionsMessageToolCall struct {
	ID       string                             `json:"id"`
	Type     string                             `json:"type"`
	Function ChatCompletionsMessageFunctionCall `json:"function"`
}

type URLCitation struct {
	StartIndex int    `json:"start_index"`
	EndIndex   int    `json:"end_index"`
	URL        string `json:"url"`
	Title      string `json:"title"`
}

type ChatCompletionResponseMessageAnnotation struct {
	Type        string      `json:"type"`
	URLCitation URLCitation `json:"url_citation"`
}

type ChatCompletionResponseMessage struct {
	Content string `json:"content"`
	Refusal bool   `json:"refusal,omitempty"`
}

type ChatCompletionChoice struct {
	FinishReason string                  `json:"finish_reason"`
	Index        int                     `json:"index"`
	LogProbs     *ChatCompletionLogProbs `json:"logprobs,omitempty"`
	Message      ChatMessage             `json:"message"`
}

type ChatCompletionResponse struct {
	ID      string                 `json:"id"`
	Choices []ChatCompletionChoice `json:"choices"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Object  string                 `json:"object"`
}
