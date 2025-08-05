package ai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"redhat-bot/config"
)

type DeepSeekClient struct {
	apiKey     string
	refererURL string
	siteTitle  string
	client     *http.Client
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model    string    `json:"model"`
	Messages []Message `json:"messages"`
}

type ChatResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func NewDeepSeekClient() *DeepSeekClient {
	return &DeepSeekClient{
		apiKey:     config.AppConfig.DeepSeekToken,
		refererURL: "<YOUR_SITE_URL>",  // تنظیم از کانفیگ
		siteTitle:  "<YOUR_SITE_NAME>", // تنظیم از کانفیگ
		client:     &http.Client{},
	}
}

func (d *DeepSeekClient) AskQuestion(question string) (string, error) {
	messages := []Message{
		{
			Role:    "user",
			Content: question,
		},
	}
	return d.makeRequest(messages)
}

func (d *DeepSeekClient) makeRequest(messages []Message) (string, error) {
	requestBody := ChatRequest{
		Model:    "deepseek/deepseek-r1-0528:free",
		Messages: messages,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		return "", fmt.Errorf("خطا در تبدیل درخواست: %v", err)
	}

	req, err := http.NewRequest(
		"POST",
		"https://openrouter.ai/api/v1/chat/completions",
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return "", fmt.Errorf("خطا در ایجاد درخواست: %v", err)
	}

	// تنظیم هدرهای مورد نیاز
	req.Header.Set("Authorization", "Bearer "+d.apiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("HTTP-Referer", d.refererURL)
	req.Header.Set("X-Title", d.siteTitle)

	resp, err := d.client.Do(req)
	if err != nil {
		return "", fmt.Errorf("درخواست ناموفق بود: %v", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("خطا در خواندن پاسخ: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("خطای API: %s", body)
	}

	var chatResp ChatResponse
	if err := json.Unmarshal(body, &chatResp); err != nil {
		return "", fmt.Errorf("خطا در تجزیه پاسخ: %v", err)
	}

	if len(chatResp.Choices) == 0 {
		return "", fmt.Errorf("پاسخی تولید نشد")
	}

	return chatResp.Choices[0].Message.Content, nil
}
