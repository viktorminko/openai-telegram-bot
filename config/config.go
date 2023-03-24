package config

import (
	"fmt"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken        string        `mapstructure:"TELEGRAM_BOT_TOKEN" validate:"required"`
	OpenAIApiKey            string        `mapstructure:"OPENAI_API_KEY" validate:"required"`
	OpenAIChatModel         string        `mapstructure:"OPENAI_CHAT_MODEL" validate:"required"`
	OpenAITranscriptModel   string        `mapstructure:"OPENAI_TRANSCRIPT_MODEL" validate:"required"`
	ImageSize               string        `mapstructure:"IMAGE_SIZE" validate:"required"`
	ContextSizeBytes        int           `mapstructure:"CONTEXT_SIZE_BYTES" validate:"min=1,max=8000"`
	MaxVoiceMessageDuration time.Duration `mapstructure:"MAX_VOICE_MESSAGE_DURATION" validate:"required,max=60s"`
}

func LoadConfig() (*Config, error) {
	var cfg Config

	viper.AutomaticEnv()

	envVars := []string{
		"TELEGRAM_BOT_TOKEN",
		"OPENAI_API_KEY",
		"OPENAI_CHAT_MODEL",
		"OPENAI_TRANSCRIPT_MODEL",
		"IMAGE_SIZE",
		"CONTEXT_SIZE_BYTES",
		"MAX_VOICE_MESSAGE_DURATION",
	}

	for _, envVar := range envVars {
		if err := viper.BindEnv(envVar); err != nil {
			return nil, fmt.Errorf("bind environment variable %s: %v", envVar, err)
		}
	}

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unmarshal config: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("validate config: %v", err)
	}

	return &cfg, nil
}
