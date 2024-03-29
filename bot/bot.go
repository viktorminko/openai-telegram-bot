package bot

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"time"

	"github.com/viktorminko/openai-telegram-bot/ai"
	"github.com/viktorminko/openai-telegram-bot/messenger"
)

type bot interface {
	RecvMessages(ctx context.Context) (<-chan messenger.Message, error)
	SendMessage(ctx context.Context, chatID int64, replyToMessageID int, message string) (*messenger.Message, error)
	SendImage(ctx context.Context, chatID int64, imageURL string) error
	SendVoiceMessage(ctx context.Context, chatID int64, audio io.Reader) error
}

type aiClient interface {
	GetImageURL(ctx context.Context, text string) (string, error)
	ChatCompletion(ctx context.Context, msgs []ai.ChatMessage) (string, error)
	Transcript(ctx context.Context, audio io.Reader) (string, error)
}

type textToSpeechClient interface {
	ConvertTextToSpeech(ctx context.Context, messages []ai.ChatMessage) (io.Reader, error)
}

type Bot struct {
	bot                     bot
	client                  aiClient
	ttsClient               textToSpeechClient
	chatContexts            map[int64]*ChatContext
	maxContextSize          int64
	maxVoiceMessageDuration time.Duration
}

func NewBot(bot bot, client aiClient, ttsClient textToSpeechClient, maxContextSize int64, maxVoiceMessageDuration time.Duration) *Bot {
	return &Bot{
		bot:                     bot,
		client:                  client,
		ttsClient:               ttsClient,
		chatContexts:            make(map[int64]*ChatContext),
		maxContextSize:          maxContextSize,
		maxVoiceMessageDuration: maxVoiceMessageDuration,
	}
}

func (b *Bot) Run(ctx context.Context) (<-chan error, error) {
	updates, err := b.bot.RecvMessages(ctx)
	if err != nil {
		return nil, fmt.Errorf("get updates channel: %w", err)
	}

	errch := make(chan error)
	go func() {
		<-ctx.Done()
		close(errch)
	}()

	go func() {
		for msg := range updates {
			if err := b.processMessage(ctx, msg); err != nil {
				errch <- fmt.Errorf("process message: %w", err)
			}
		}
	}()

	return errch, nil
}

func (b *Bot) processMessage(ctx context.Context, msg messenger.Message) error {
	switch msg.Command {
	case messenger.CommandImage:
		return b.handleImageCommand(ctx, msg)
	case messenger.CommandExportAudio:
		return b.handleExportAudioCommand(ctx, msg.ChatID)
	case messenger.CommandReset:
		return b.handleResetCommand(ctx, msg.ChatID)
	default:
		return b.handleUserMessage(ctx, msg)
	}
}

func (b *Bot) handleImageCommand(ctx context.Context, msg messenger.Message) error {
	url, err := b.client.GetImageURL(ctx, msg.Text)
	if err != nil {
		return fmt.Errorf("get image url: %w", err)
	}

	if err := b.bot.SendImage(ctx, msg.ChatID, url); err != nil {
		return fmt.Errorf("send image: %w", err)
	}

	return nil
}

func (b *Bot) handleResetCommand(ctx context.Context, chatID int64) error {
	delete(b.chatContexts, chatID)

	_, err := b.bot.SendMessage(ctx, chatID, 0, "Chat context has been reset")
	if err != nil {
		return fmt.Errorf("send message: %w", err)
	}

	return nil
}

func (b *Bot) handleExportAudioCommand(ctx context.Context, chatID int64) error {
	chatCtx, ok := b.chatContexts[chatID]
	if !ok {
		return fmt.Errorf("no chat context found for chatID: %d", chatID)
	}

	// Get all messages in the conversation.
	messages := chatCtx.GetMessages()

	// Convert the messages to speech using the TextToSpeechClient.
	audio, err := b.ttsClient.ConvertTextToSpeech(ctx, messages)
	if err != nil {
		return fmt.Errorf("text-to-speech conversion failed: %w", err)
	}

	// send audio to user
	if err := b.bot.SendVoiceMessage(ctx, chatID, audio); err != nil {
		return fmt.Errorf("send audio: %w", err)
	}

	return nil
}

func (b *Bot) handleAudioMessage(ctx context.Context, msg messenger.Message) error {
	mp3, err := convertOggToMp3(msg.Audio, b.maxVoiceMessageDuration)
	if err != nil {
		return fmt.Errorf("convert ogg to mp3: %v", err)
	}

	res, err := b.client.Transcript(ctx, bytes.NewReader(mp3))
	if err != nil {
		return fmt.Errorf("transcript: %v", err)
	}

	if _, err := b.bot.SendMessage(
		ctx,
		msg.ChatID,
		msg.RequestMessageID,
		fmt.Sprintf("transcript: %s", res),
	); err != nil {
		return fmt.Errorf("send image: %w", err)
	}

	if err := b.handleTextMessage(ctx, messenger.Message{
		FromUserID: msg.FromUserID,
		ChatID:     msg.ChatID,
		Text:       res,
	}); err != nil {
		return fmt.Errorf("handle text message: %v", err)
	}

	return nil
}

func (b *Bot) handleTextMessage(ctx context.Context, msg messenger.Message) error {
	chatID := msg.ChatID

	chatCtx, ok := b.chatContexts[chatID]
	if !ok {
		b.chatContexts[chatID] = NewChatContext(
			&MessageQueue{
				MaxSize: b.maxContextSize,
			},
		)
		chatCtx = b.chatContexts[chatID]
	}

	chatCtx.AddMessage(ai.ChatMessage{
		Role:       ai.RoleUser,
		FromUserID: msg.FromUserID,
		Text:       msg.Text,
	})

	res, err := b.client.ChatCompletion(ctx, chatCtx.GetMessages())
	if err != nil {
		return fmt.Errorf("chat completion: %v", err)
	}

	replyMessage, err := b.bot.SendMessage(ctx, msg.ChatID, msg.RequestMessageID, res)
	if err != nil {
		return fmt.Errorf("send image: %w", err)
	}

	chatCtx.AddMessage(ai.ChatMessage{
		Role:       ai.RoleAssistant,
		FromUserID: replyMessage.FromUserID,
		Text:       res,
	})

	return nil
}

func (b *Bot) handleUserMessage(ctx context.Context, msg messenger.Message) error {
	if len(msg.Audio) != 0 {
		return b.handleAudioMessage(ctx, msg)
	}

	return b.handleTextMessage(ctx, msg)
}
