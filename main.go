package main

import (
	"context"
	"log"

	"github.com/viktorminko/openai-telegram-bot/ai"
	"github.com/viktorminko/openai-telegram-bot/bot"
	"github.com/viktorminko/openai-telegram-bot/config"
	"github.com/viktorminko/openai-telegram-bot/messenger"

	tgapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("error loading config: %v", err)
	}

	tgBot, err := tgapi.NewBotAPI(cfg.TelegramBotToken)
	if err != nil {
		log.Fatalf("new telegram bot api: %s", err)
	}

	errch, err := bot.NewBot(
		messenger.NewTelegram(
			tgBot,
		),
		ai.NewClient(cfg.OpenAIApiKey,
			ai.Config{
				ChatModel:       cfg.OpenAIChatModel,
				TranscriptModel: cfg.OpenAITranscriptModel,
				ImageSize:       cfg.ImageSize,
			},
		),
		int64(cfg.ContextSizeBytes),
		cfg.MaxVoiceMessageDuration,
	).Run(context.Background())
	if err != nil {
		log.Fatalf("new bot: %s", err)
	}

	for err := range errch {
		log.Printf("error: %s", err)
	}

}
