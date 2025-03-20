package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type GrokRequest struct {
	Messages   []Message `json:"messages"`
	Model      string    `json:"model"`
	Stream     bool      `json:"stream"`
	Temperature float32  `json:"temperature"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type GrokResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func CallGrokAPI(userInput string) (string, error) {
	const apiKey = "xai-eYQveJd3YQJOZkRbJFc9kNGhSVL4Fwob4T7QuERqg8XPMydAWysjLMRiKZKLgdBL1immW5yvmUmPrZKW"
	const url = "https://api.x.ai/v1/chat/completions"

	// 构建请求数据
	reqBody := GrokRequest{
		Messages: []Message{
			{
				Role:    "system",
				Content: "You are a test assistant.",
			},
			{
				Role:    "user",
				Content: userInput,
			},
		},
		Model:      "grok-2-latest",
		Stream:     false,
		Temperature: 0,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return "", err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// 调试：打印状态码和响应内容
	fmt.Printf("状态码: %d\n", resp.StatusCode)
	fmt.Printf("响应内容: %s\n", string(body))

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(body))
	}

	var grokResp GrokResponse
	err = json.Unmarshal(body, &grokResp)
	if err != nil {
		return "", fmt.Errorf("JSON解析失败: %v, 响应内容: %s", err, string(body))
	}

	if len(grokResp.Choices) > 0 {
		return grokResp.Choices[0].Message.Content, nil
	}
	return "Grok未返回有效回答", nil
}