package messenger

import (
	"context"
	"fmt"
	"io"
	"net/http"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type Telegram struct {
	bot *tgbotapi.BotAPI
}

func NewTelegram(bot *tgbotapi.BotAPI) *Telegram {
	return &Telegram{
		bot: bot,
	}
}

func (b *Telegram) RecvMessages(ctx context.Context) (<-chan Message, error) {
	upd, err := b.bot.GetUpdatesChan(tgbotapi.NewUpdate(0))
	if err != nil {
		return nil, fmt.Errorf("get updates channel: %w", err)
	}

	res := make(chan Message)
	go func() {
		<-ctx.Done()
		close(res)
	}()

	go func() {
		for msg := range upd {
			select {
			case <-ctx.Done():
				return
			default:

			}

			if msg.Message == nil {
				continue
			}

			var audio []byte
			if msg.Message.Voice != nil {
				audio, err = b.downloadVoiceMessage(msg.Message.Voice.FileID)
				if err != nil {
					fmt.Printf("failed to download voice message: %v /n", err)
					// TODO handle error
					continue
				}
			}

			res <- Message{
				ChatID:           msg.Message.Chat.ID,
				FromUserID:       msg.Message.From.ID,
				RequestMessageID: msg.Message.MessageID,
				Command:          msg.Message.Command(),
				Text:             msg.Message.Text,
				Audio:            audio,
			}
		}
	}()

	return res, nil
}

// SendMessage sends a message to a chat
func (b *Telegram) SendMessage(ctx context.Context, chatID int64, replyToMessageID int, message string) (*Message, error) {
	msg := tgbotapi.NewMessage(chatID, message)
	msg.ReplyToMessageID = replyToMessageID

	replyMsg, err := b.bot.Send(msg)
	if err != nil {
		return nil, fmt.Errorf("send message: %w", err)
	}

	return &Message{
		ChatID:           chatID,
		FromUserID:       replyMsg.From.ID,
		RequestMessageID: replyToMessageID,
		Text:             msg.Text,
	}, nil
}

func (b *Telegram) SendImage(ctx context.Context, chatID int64, imageURL string) error {
	client := &http.Client{}
	req, err := http.NewRequest("GET", imageURL, nil)
	if err != nil {
		return err
	}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	data, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	inputFile := tgbotapi.FileBytes{
		Name:  res.Header.Get("Content-Disposition"),
		Bytes: data,
	}

	_, err = b.bot.Send(tgbotapi.NewPhotoUpload(int64(chatID), inputFile))
	if err != nil {
		return err
	}

	return nil
}

func (b *Telegram) SendVoiceMessage(ctx context.Context, chatID int64, audio io.Reader) error {
	// read audio bytes from the io.Reader
	audioBytes, err := io.ReadAll(audio)
	if err != nil {
		return fmt.Errorf("read audio bytes: %w", err)
	}

	// create the message with the audio bytes
	msg := tgbotapi.NewVoiceUpload(chatID, tgbotapi.FileBytes{Name: "audio.ogg", Bytes: audioBytes})

	// send the message
	if _, err := b.bot.Send(msg); err != nil {
		return fmt.Errorf("send voice message: %w", err)
	}

	return nil
}

func (b *Telegram) downloadVoiceMessage(fileID string) ([]byte, error) {
	fileConfig := tgbotapi.FileConfig{FileID: fileID}
	file, err := b.bot.GetFile(fileConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to get file: %w", err)
	}

	resp, err := http.Get(file.Link(b.bot.Token))
	if err != nil {
		return nil, fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	bts, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read file: %w", err)
	}

	return bts, nil
}
