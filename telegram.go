package main

import (
	"io/ioutil"
	"log"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// sendMessage sends a message to a chat
func sendMessage(bot *tgbotapi.BotAPI, chatID int64, message string) {
	if _, err := bot.Send(tgbotapi.NewMessage(chatID, message)); err != nil {
		log.Printf("error sending message: %v", err)
	}
}

func sendImage(bot *tgbotapi.BotAPI, chatID int64, imageURL string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	inputFile := tgbotapi.FileBytes{
		Name:  res.Header.Get("Content-Disposition"),
		Bytes: data,
	}

	_, err = bot.Send(tgbotapi.NewPhotoUpload(int64(chatID), inputFile))
	if err != nil {
		return err
	}

	return nil
}
