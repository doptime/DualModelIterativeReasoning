package models

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"strings"
)

type Model struct {
	ModelName string
	Url       string
	ApiKey    string
}

var SLM1 = &Model{
	ModelName: "neuralmagic/Meta-Llama-3.1-8B-Instruct-quantized.w4a16",
	Url:       "http://gpu.lan:8007/v1/chat/completions",
	ApiKey:    "token-deaf",
}

var SLM2 = &Model{
	//ModelName: "neuralmagic/Qwen2-7B-Instruct-quantized.w8a16",
	//Url:       "http://gpu.lan:8006/v1/chat/completions",
	//ModelName: "neuralmagic/Phi-3-medium-128k-instruct-quantized.w4a16",
	//ModelName: "neuralmagic/gemma-2-9b-it-quantized.w4a16",
	ModelName: "shuyuej/Mistral-Nemo-Instruct-2407-GPTQ",
	Url:       "http://gpu.lan:8003/v1/chat/completions",
	ApiKey:    "token-deaf",
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func (msg *Message) String() string {
	if msg == nil {
		return ""
	}
	return msg.Role + ": " + msg.Content
}
func MsgOfUser(msg string) *Message {
	return &Message{Role: "user", Content: msg}
}
func MsgOfAssistant(msg string) *Message {
	return &Message{Role: "assistant", Content: msg}
}

type ChatGPTResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Created int    `json:"created"`
	Model   string `json:"model"`
	Choices []struct {
		Index   int `json:"index"`
		Message struct {
			Role      string        `json:"role"`
			Content   string        `json:"content"`
			ToolCalls []interface{} `json:"tool_calls"`
		} `json:"message"`
		Delta struct {
			Content string `json:"content"`
		} `json:"delta"`
		Logprobs     interface{} `json:"logprobs"`
		FinishReason string      `json:"finish_reason"`
		StopReason   interface{} `json:"stop_reason"`
	} `json:"choices"`
	Usage struct {
		PromptTokens     int `json:"prompt_tokens"`
		CompletionTokens int `json:"completion_tokens"`
		TotalTokens      int `json:"total_tokens"`
	} `json:"usage"`
	PromptLogprobs interface{} `json:"prompt_logprobs"`
}

func (m *Model) AskLLM(temperature float64, stream bool, message ...*Message) (r *Message, err error) {
	messages := make([]Message, 0, len(message))
	for _, msg := range message {
		if msg != nil {
			messages = append(messages, *msg)
		}
	}
	// Prepare the payload
	payload := map[string]interface{}{
		"model":       m.ModelName,
		"messages":    messages,
		"temperature": temperature,
		"stream":      stream,
	}

	// Convert payload to JSON
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	// Create a new request
	req, err := http.NewRequest("POST", m.Url, bytes.NewBuffer(payloadBytes))
	if err != nil {
		return nil, err
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+m.ApiKey)

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Read the response
	if stream {
		return handleStreamResponse(resp)
	} else {
		return handleNonStreamResponse(resp)
	}
}

func handleStreamResponse(resp *http.Response) (*Message, error) {
	var fullContent strings.Builder
	decoder := json.NewDecoder(resp.Body)
	for {
		var chunk ChatGPTResponse
		err := decoder.Decode(&chunk)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		if len(chunk.Choices) > 0 {
			fullContent.WriteString(chunk.Choices[0].Delta.Content)
		}
	}
	return MsgOfAssistant(fullContent.String()), nil
}

func handleNonStreamResponse(resp *http.Response) (*Message, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var response ChatGPTResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}

	if len(response.Choices) > 0 {
		return MsgOfAssistant(response.Choices[0].Message.Content), nil
	}

	return nil, nil
}
