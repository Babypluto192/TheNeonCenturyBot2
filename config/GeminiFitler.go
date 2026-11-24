package config

import (
	"context"
	"strings"
	"time"

	"google.golang.org/genai"
)

func GeminiFitler(client *genai.Client, geminiFilterPrompt, question string) string {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	cfg := &genai.GenerateContentConfig{
		SystemInstruction: &genai.Content{
			Parts: []*genai.Part{
				{Text: geminiFilterPrompt},
			},
		},
	}

	resp, err := client.Models.GenerateContent(ctx, "gemini-3.5-flash-lite", genai.Text(question), cfg)
	if err != nil {
		return ""
	}
	answer := strings.TrimSpace(resp.Text())

	return answer
}
