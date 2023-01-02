package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/go-telegram-bot-api/telegram-bot-api"
)

func readEnvVar(name string) string {
	val, ok := os.LookupEnv(name)
	if !ok {
		fmt.Printf("Error: %s environment variable not found\n", name)
		os.Exit(1)
	}
	return val
}

func main() {
	// Read the environment variables
	botToken := readEnvVar("TELEGRAM_BOT_TOKEN")
	apiKey := readEnvVar("OPENAI_API_KEY")
	model := readEnvVar("MODEL")

	bot, err := tgbotapi.NewBotAPI(botToken)
	if err != nil {
		panic(err)
	}

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		// Read the message sent by the user
		msg := update.Message.Text

		command := update.Message.Command()

		if command == "image" {
			sendMessage(bot, update.Message.Chat.ID, "generating image...")
			imageURL, err := generateImage(context.Background(), apiKey, msg)
			if err != nil {
				fmt.Println(fmt.Errorf("error generating image: %v", err))
				sendMessage(bot, update.Message.Chat.ID, "error")
				continue
			}

			if err := sendImage(bot, update.Message.Chat.ID, imageURL); err != nil {
				log.Println(fmt.Errorf("error sending image: %v", err))
				sendMessage(bot, update.Message.Chat.ID, "error")
				continue
			}
			continue
		}

		// Call the OpenAI API
		text, err := openAICompletion(context.Background(), apiKey, model, msg)
		if err != nil {
			sendMessage(bot, update.Message.Chat.ID, "Error calling OpenAI API")
			continue
		}

		// Send the response to the user
		sendMessage(bot, update.Message.Chat.ID, text)
	}
}
