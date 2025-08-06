package summary

import (
	"context"
	"errors"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/sashabaranov/go-openai"
)

type OpenAISummarizer struct {
	client  *openai.Client
	prompt  string
	model   string
	enabled bool
	mu      sync.Mutex
}

func NewOpenAISummarizer(apiKey, model, prompt string) *OpenAISummarizer {
	s := &OpenAISummarizer{
		client:  openai.NewClient(apiKey),
		prompt:  prompt,
		model:   model,
		enabled: true,
	}

	log.Printf("openai summarizer is enabled: %v", apiKey != "")

	if apiKey != "" {
		s.enabled = true
	}

	return s
}

func (s *OpenAISummarizer) Summarizer(text string) (string, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if !s.enabled {
		log.Println("openai summarizer is disabled, returning original text")
		return SmartTrim(text, 100), nil
	}

	request := openai.ChatCompletionRequest{
		Model: s.model,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: s.prompt,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: text,
			},
		},
		MaxTokens:   400,
		Temperature: 1,
		TopP:        1,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Minute)
	defer cancel()

	resp, err := s.client.CreateChatCompletion(ctx, request)
	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("no choices in openai response")
	}

	rawSummary := strings.TrimSpace(resp.Choices[0].Message.Content)
	if strings.HasSuffix(rawSummary, ".") {
		return rawSummary, nil
	}

	// cut all after the last ".":
	sentences := strings.Split(rawSummary, ".")

	return strings.Join(sentences[:len(sentences)-1], ".") + ".", nil
}

func SmartTrim(text string, maxWords int) string {
	sentences := strings.FieldsFunc(text, func(r rune) bool {
		return r == '.' || r == '!' || r == '?'
	})

	var result []string
	wordCount := 0

	for _, sentence := range sentences {
		trimmed := strings.TrimSpace(sentence)
		words := strings.Fields(trimmed)
		if wordCount+len(words) > maxWords {
			break
		}
		result = append(result, trimmed)
		wordCount += len(words)
	}

	return strings.Join(result, ". ") + "."
}
