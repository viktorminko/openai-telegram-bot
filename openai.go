package main

import (
	"context"

	gpt3 "github.com/sashabaranov/go-gpt3"
)

func openAICompletion(ctx context.Context, apiKey string, model string, prompt string) (string, error) {
	// Set up the GPT-3 client
	client := gpt3.NewClient(apiKey)

	// Set up the request to the OpenAI API
	req := gpt3.CompletionRequest{
		Model:            model,
		Prompt:           prompt,
		MaxTokens:        500,
		Temperature:      0,
		TopP:             1,
		FrequencyPenalty: 0,
		PresencePenalty:  0,

		Echo: true,
	}

	// Send the request to the OpenAI API
	resp, err := client.CreateCompletion(ctx, req)
	if err != nil {
		return "", err
	}

	// Return the generated text
	return resp.Choices[0].Text, nil
}

func generateImage(ctx context.Context, apiKey string, prompt string) (string, error) {

	// Set up the GPT-3 client
	client := gpt3.NewClient(apiKey)

	// Set up the request to the OpenAI API
	req := gpt3.ImageCreateRequest{
		Prompt: prompt,
		Size:   gpt3.ImageSizeMedium,
	}

	// Send the request to the OpenAI API
	resp, err := client.CreateImageCreate(ctx, req)
	if err != nil {
		return "", err
	}

	return resp.URLs[0], nil
}
