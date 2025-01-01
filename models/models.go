package models

// models to share across project
// to avoid circular imports
type ASSET struct {
	Asset     string `json:"asset"`
	Extension string `json:"extension"`
	Blob      []byte `json:"blob"`
}

type DiaryEntry struct {
	Ids     []int   `json:"ids"`
	Content string  `json:"content"`
	Date    string  `json:"date"`
	Asset   []ASSET `json:"asset"`
}

// --LLM model response Struct
type ChatCompletion struct {
	ID                string   `json:"id"`
	Object            string   `json:"object"`
	Created           int64    `json:"created"`
	Model             string   `json:"model"`
	Choices           []Choice `json:"choices"`
	Usage             Usage    `json:"usage"`
	SystemFingerprint string   `json:"system_fingerprint"`
	XGroq             XGroq    `json:"x_groq"`
}

type Choice struct {
	Index        int       `json:"index"`
	Message      Message   `json:"message"`
	Logprobs     *Logprobs `json:"logprobs"`
	FinishReason string    `json:"finish_reason"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Logprobs struct {
	// Define fields if needed, currently it's null in the provided JSON
}

type Usage struct {
	QueueTime        float64 `json:"queue_time"`
	PromptTokens     int     `json:"prompt_tokens"`
	PromptTime       float64 `json:"prompt_time"`
	CompletionTokens int     `json:"completion_tokens"`
	CompletionTime   float64 `json:"completion_time"`
	TotalTokens      int     `json:"total_tokens"`
	TotalTime        float64 `json:"total_time"`
}

type XGroq struct {
	ID string `json:"id"`
}

//--------

// ----- Model for Log json =

type Link struct {
	Href string `json:"href"`
}

// DiaryEntry represents the structure of each diary entry in the JSON.
type LogEntry struct {
	Title       string `json:"title"`
	Dates       string `json:"dates"`
	Description string `json:"description"`
	Image       string `json:"image"`
	Tags        string `json:"tags"`
	Links       []Link `json:"links"`
}

//--------
