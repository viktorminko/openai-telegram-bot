package ai_test

import (
	"context"
	"strings"
	"testing"

	"github.com/viktorminko/openai-telegram-bot/ai"
)

const apiKey = "<API_KEY>"

func TestChatCompletion(t *testing.T) {
	cfg := ai.Config{
		ChatModel: "text-davinci-002",
		ImageSize: "256x256",
	}
	client := ai.NewClient(apiKey, cfg)

	messages := []ai.ChatMessage{
		{Role: ai.RoleUser, Text: "tell me a joke"},
	}

	response, err := client.ChatCompletion(context.Background(), messages)
	if err != nil {
		t.Fatalf("Error in ChatCompletion: %v", err)
	}

	if strings.TrimSpace(response) == "" {
		t.Error("Expected a non-empty response, got empty response")
	}
}

func TestGetImageURL(t *testing.T) {
	cfg := ai.Config{
		ChatModel: "text-davinci-002",
		ImageSize: "256x256",
	}
	client := ai.NewClient(apiKey, cfg)

	prompt := "A picture of a beautiful sunset over the ocean"

	url, err := client.GetImageURL(context.Background(), prompt)
	if err != nil {
		t.Fatalf("Error in GetImageURL: %v", err)
	}

	if strings.TrimSpace(url) == "" {
		t.Error("Expected a non-empty URL, got empty URL")
	}
}
