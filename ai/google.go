package ai

import (
	"bytes"
	"context"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	texttospeech "cloud.google.com/go/texttospeech/apiv1"
	"cloud.google.com/go/texttospeech/apiv1/texttospeechpb"
	language "cloud.google.com/go/translate"
	"google.golang.org/api/option"
)

type GoogleTextToSpeechClient struct {
	textToSpeechClient   *texttospeech.Client
	detectLanguageClient *language.Client
	voiceSelections      voiceSelector
	audioConfig          *texttospeechpb.AudioConfig
}

type voiceSelector interface {
	GetVoice(language string, id int) (string, error)
}

type Option func(*GoogleTextToSpeechClient) error

func WithVoiceSelections(voiceSelector voiceSelector) Option {
	return func(c *GoogleTextToSpeechClient) error {
		c.voiceSelections = voiceSelector

		return nil
	}
}

func WithAudioEncoding(encoding texttospeechpb.AudioEncoding) Option {
	return func(c *GoogleTextToSpeechClient) error {
		c.audioConfig = &texttospeechpb.AudioConfig{
			AudioEncoding: encoding,
		}

		return nil
	}
}

func NewGoogleTextToSpeechClient(ctx context.Context, credentialsJsonBase64 string, opts ...Option) (*GoogleTextToSpeechClient, error) {
	credentialsJson, err := base64.StdEncoding.DecodeString(credentialsJsonBase64)
	if err != nil {
		return nil, fmt.Errorf("base64 decode credentials: %v", err)
	}

	textToSPeechClient, err := texttospeech.NewClient(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, fmt.Errorf("texttospeech.NewClient: %v", err)
	}

	detectLanguageClient, err := language.NewClient(ctx, option.WithCredentialsJSON(credentialsJson))
	if err != nil {
		return nil, fmt.Errorf("language.NewClient: %v", err)
	}

	c := &GoogleTextToSpeechClient{
		textToSpeechClient:   textToSPeechClient,
		detectLanguageClient: detectLanguageClient,
	}

	for _, opt := range opts {
		if err := opt(c); err != nil {
			return nil, fmt.Errorf("option: %v", err)
		}
	}

	return c, nil
}

func (c *GoogleTextToSpeechClient) ConvertTextToSpeech(ctx context.Context, messages []ChatMessage) (io.Reader, error) {
	var audioBuffer bytes.Buffer

	audioConf := c.audioConfig
	if c.audioConfig == nil {
		audioConf = &texttospeechpb.AudioConfig{
			AudioEncoding: texttospeechpb.AudioEncoding_MP3,
		}
	}

	// Process each message separately.
	for _, message := range messages {
		synthesisInput := &texttospeechpb.SynthesisInput{
			InputSource: &texttospeechpb.SynthesisInput_Text{
				Text: message.Text,
			},
		}

		// Determine the lang of the text.
		lang, err := c.detectLanguage(ctx, message.Text)
		if err != nil {
			return nil, fmt.Errorf("detect lang: %v", err)
		}

		// Choose a voice for the text based on its lang and role.
		voice, err := c.selectVoice(lang, message.FromUserID)
		if err != nil {
			return nil, fmt.Errorf("select voice: %v", err)
		}

		// Perform the text-to-speech synthesis request.
		resp, err := c.textToSpeechClient.SynthesizeSpeech(ctx, &texttospeechpb.SynthesizeSpeechRequest{
			Input:       synthesisInput,
			AudioConfig: audioConf,
			Voice:       voice,
		})
		if err != nil {
			return nil, fmt.Errorf("SynthesizeSpeech: %+v, %+v, %+v, %w", synthesisInput, audioConf, voice, err)
		}

		// Concatenate the audio content to the buffer.
		if _, err := audioBuffer.Write(resp.AudioContent); err != nil {
			return nil, fmt.Errorf("write to audio buffer: %v", err)
		}
	}

	// Return the audio content as a bytes.Buffer.
	return &audioBuffer, nil
}

func (c *GoogleTextToSpeechClient) detectLanguage(ctx context.Context, inputText string) (string, error) {
	resp, err := c.detectLanguageClient.DetectLanguage(ctx, []string{inputText})
	if err != nil {
		return "", fmt.Errorf("detect language: %v", err)
	}

	if len(resp) == 0 || len(resp[0]) == 0 {
		return "", errors.New("detectLanguage return value empty")
	}

	return resp[0][0].Language.String(), nil
}

func (c *GoogleTextToSpeechClient) Close() error {
	if err := c.textToSpeechClient.Close(); err != nil {
		return fmt.Errorf("textToSpeechClient.Close: %v", err)
	}

	if err := c.detectLanguageClient.Close(); err != nil {
		return fmt.Errorf("detectLanguageClient.Close: %v", err)
	}

	return nil
}

func (c *GoogleTextToSpeechClient) selectVoice(language string, id int) (*texttospeechpb.VoiceSelectionParams, error) {
	voice, err := c.voiceSelections.GetVoice(language, id)
	if err != nil {
		return nil, fmt.Errorf("get voice for language %q and id %q: %v", language, id, err)
	}

	// If no matching voice is found, return a default voice for the language.
	return &texttospeechpb.VoiceSelectionParams{
		LanguageCode: language,
		Name:         voice,
	}, nil
}
