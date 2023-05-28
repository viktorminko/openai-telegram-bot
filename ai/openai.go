package ai

import (
	"context"
	"fmt"
	"io"
	"os"

	"github.com/sashabaranov/go-openai"
)

type Config struct {
	ChatModel       string
	TranscriptModel string
	ImageSize       string
}

type Role string

const (
	RoleUser      Role = "user"
	RoleAssistant Role = "assistant"
)

type ChatMessage struct {
	Role       Role
	FromUserID int
	Text       string
}

type OpenAI struct {
	cfg    Config
	client *openai.Client
}

func NewClient(apiKey string, cfg Config) *OpenAI {
	return &OpenAI{
		client: openai.NewClient(apiKey),
		cfg:    cfg,
	}
}

func (c *OpenAI) ChatCompletion(ctx context.Context, messages []ChatMessage) (string, error) {
	msgs := make([]openai.ChatCompletionMessage, len(messages))
	for i, m := range messages {
		msgs[i] = openai.ChatCompletionMessage{
			Role:    string(m.Role),
			Content: m.Text,
		}
	}

	resp, err := c.client.CreateChatCompletion(
		ctx,
		openai.ChatCompletionRequest{
			Model:    c.cfg.ChatModel,
			Messages: msgs,
		},
	)

	if err != nil {
		return "", fmt.Errorf("openAI chat completion: %v", err)
	}

	return resp.Choices[0].Message.Content, nil
}

func (c *OpenAI) Transcript(ctx context.Context, audio io.Reader) (string, error) {
	f, err := os.CreateTemp("", "audio_*.mp3")
	if err != nil {
		return "", fmt.Errorf("openAI create temp file: %v", err)
	}
	defer os.Remove(f.Name())

	_, err = io.Copy(f, audio)
	if err != nil {
		return "", fmt.Errorf("openAI read audio: %v", err)
	}

	if err := f.Close(); err != nil {
		return "", fmt.Errorf("failed to close temp file: %v", err)
	}

	resp, err := c.client.CreateTranscription(
		ctx,
		openai.AudioRequest{
			Model:    c.cfg.TranscriptModel,
			FilePath: f.Name(),
		},
	)

	if err != nil {
		return "", fmt.Errorf("openAI create transcription: %v", err)
	}

	return resp.Text, nil
}

func (c *OpenAI) GetImageURL(ctx context.Context, prompt string) (string, error) {

	resp, err := c.client.CreateImage(ctx, openai.ImageRequest{
		Prompt: prompt,
		Size:   c.cfg.ImageSize,
	})

	if err != nil {
		return "", fmt.Errorf("openAI image: %v", err)
	}

	return resp.Data[0].URL, nil
}
