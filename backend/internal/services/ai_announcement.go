package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const announcementGenTimeout = 15 * time.Second

const announcementSystemPrompt = `你是公告撰寫助手。根據使用者的簡短描述，產生一則正式但友善的公告。

規則：
- 使用繁體中文
- title 簡潔明瞭，不超過 50 字
- content 使用 Markdown 格式，段落分明，語氣友善專業
- 不要使用 emoji
- 不要加入日期時間（管理員會自行設定）`

var announcementResponseSchema = json.RawMessage(`{
  "type": "object",
  "properties": {
    "title":   { "type": "string" },
    "content": { "type": "string" }
  },
  "required": ["title", "content"]
}`)

type AnnouncementGenResult struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type AnnouncementGenerator struct {
	apiKey  string
	model   string
	baseURL string
	client  *http.Client
}

func NewAnnouncementGenerator(apiKey, model, baseURL string) *AnnouncementGenerator {
	if baseURL == "" {
		baseURL = geminiDefaultBaseURL
	}
	if model == "" {
		model = "gemini-2.5-flash"
	}
	return &AnnouncementGenerator{
		apiKey:  apiKey,
		model:   model,
		baseURL: strings.TrimRight(baseURL, "/"),
		client:  &http.Client{Timeout: announcementGenTimeout + 5*time.Second},
	}
}

func (g *AnnouncementGenerator) Generate(ctx context.Context, description string) (*AnnouncementGenResult, error) {
	if g.apiKey == "" {
		return nil, fmt.Errorf("gemini api key not configured")
	}

	callCtx, cancel := context.WithTimeout(ctx, announcementGenTimeout)
	defer cancel()

	reqBody := geminiRequest{
		SystemInstruction: &geminiContent{
			Parts: []geminiPart{{Text: announcementSystemPrompt}},
		},
		Contents: []geminiContent{{
			Role:  "user",
			Parts: []geminiPart{{Text: description}},
		}},
		GenerationConfig: geminiGenerationCfg{
			ResponseMimeType: "application/json",
			ResponseSchema:   announcementResponseSchema,
			Temperature:      0.7,
		},
	}

	payload, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("marshal request: %w", err)
	}

	endpoint := fmt.Sprintf("%s/v1beta/models/%s:generateContent?key=%s", g.baseURL, g.model, g.apiKey)
	req, err := http.NewRequestWithContext(callCtx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("build request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := g.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("gemini call: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("gemini http %d: %s", resp.StatusCode, string(body))
	}

	var gemResp geminiResponse
	if err := json.Unmarshal(body, &gemResp); err != nil {
		return nil, fmt.Errorf("unmarshal response: %w", err)
	}

	if gemResp.Error != nil {
		return nil, fmt.Errorf("gemini error %d: %s", gemResp.Error.Code, gemResp.Error.Message)
	}

	if len(gemResp.Candidates) == 0 || len(gemResp.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("empty gemini response")
	}

	raw := gemResp.Candidates[0].Content.Parts[0].Text
	var result AnnouncementGenResult
	if err := json.Unmarshal([]byte(raw), &result); err != nil {
		return nil, fmt.Errorf("parse generated announcement: %w", err)
	}

	return &result, nil
}
