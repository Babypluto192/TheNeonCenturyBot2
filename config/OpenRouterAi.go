package config

import (
	"context"
	"errors"
	"log"
	"regexp"
	"strings"
	"sync"
	"time"

	"github.com/revrost/go-openrouter"
)

var (
	keyMu     sync.Mutex
	activeIdx int
	thinkRe   = regexp.MustCompile(`(?s)<think>.*?</think>`)
)

func GenerateResponse(guildId, question string, clients []*openrouter.Client) string {
	prompt := TNCBotPrompt

	//if IsDevServer(guildId) {
	//	prompt = TNCBotPromptDevServer
	//}

	doc, err := LoadDocument()
	if err != nil {
		return "Something went wrong, Sorry!"
	}

	effort := "low"
	req := openrouter.ChatCompletionRequest{
		Temperature: 0.2,
		Reasoning: &openrouter.ChatCompletionReasoning{
			Effort: &effort,
		},
		Messages: []openrouter.ChatCompletionMessage{
			openrouter.SystemMessage(prompt + "\n\n=== DOCUMENT === \n\n" + doc),
			openrouter.UserMessage(question),
		},
	}

	keyMu.Lock()
	startIdx := activeIdx
	keyMu.Unlock()

	models := []string{firstModel, secondModel}

	for _, model := range models {
		req.Model = model
		for i := range clients {
			idx := (startIdx + i) % len(clients)
			client := clients[idx]

			ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)

			resp, err := client.CreateChatCompletion(ctx, req)
			cancel()

			if err == nil && len(resp.Choices) > 0 {
				answer := thinkRe.ReplaceAllString(resp.Choices[0].Message.Content.Text, "")
				answer = strings.TrimSpace(answer)
				if answer != "" {
					log.Printf("[OK] model=%s key#%d", model, idx+1)
					keyMu.Lock()
					activeIdx = idx
					keyMu.Unlock()
					return answer
				}
			}

			log.Println("OpenRouter key error:", idx+1, err)
			if errors.Is(err, context.DeadlineExceeded) {
				break
			}

		}
	}

	return "Something went wrong, Sorry, contact devs"
}
