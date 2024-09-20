package models

import (
	"DualModelIterativeReasoning/message"
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
	//ModelName: "Qwen/Qwen2.5-14B-Instruct-GPTQ-Int4",
	//ModelName: "neuralmagic/Meta-Llama-3.1-8B-Instruct-quantized.w4a16",
	//ModelName: "AMead10/Mistral-Small-Instruct-2409-awq",
	ModelName: "casperhansen/mistral-nemo-instruct-2407-awq",
	Url:       "http://gpu.lan:8007/v1/chat/completions",
	ApiKey:    "token-deaf",
}

var SLM2 = &Model{
	//ModelName: "neuralmagic/Qwen2-7B-Instruct-quantized.w8a16",
	//Url:       "http://gpu.lan:8006/v1/chat/completions",
	//ModelName: "neuralmagic/Phi-3-medium-128k-instruct-quantized.w4a16",
	//ModelName: "neuralmagic/gemma-2-9b-it-quantized.w4a16",
	//ModelName: "shuyuej/Mistral-Nemo-Instruct-2407-GPTQ",
	ModelName: "Qwen/Qwen2.5-14B-Instruct-AWQ",
	//ModelName: "Qwen/Qwen2.5-32B-Instruct-AWQ",
	Url:    "http://gpu.lan:8003/v1/chat/completions",
	ApiKey: "token-deaf",
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

func (m *Model) AskLLM(temperature float64, stream bool, msg ...*message.Message) (r *message.Message, err error) {
	messages := make([]*message.Message, 0, len(msg))
	for _, _msg := range msg {
		if _msg != nil {
			messages = append(messages, _msg)
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

func handleStreamResponse(resp *http.Response) (*message.Message, error) {
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
	return message.Assistant(fullContent.String()), nil
}

func handleNonStreamResponse(resp *http.Response) (*message.Message, error) {
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
		return message.Assistant(response.Choices[0].Message.Content), nil
	}

	return nil, nil
}
