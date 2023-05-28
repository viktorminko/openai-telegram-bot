package config

import (
	"fmt"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramBotToken        string                    `mapstructure:"TELEGRAM_BOT_TOKEN" validate:"required"`
	OpenAIApiKey            string                    `mapstructure:"OPENAI_API_KEY" validate:"required"`
	OpenAIChatModel         string                    `mapstructure:"OPENAI_CHAT_MODEL" validate:"required"`
	OpenAITranscriptModel   string                    `mapstructure:"OPENAI_TRANSCRIPT_MODEL" validate:"required"`
	ImageSize               string                    `mapstructure:"IMAGE_SIZE" validate:"required"`
	ContextSizeBytes        int                       `mapstructure:"CONTEXT_SIZE_BYTES" validate:"min=1,max=8000"`
	MaxVoiceMessageDuration time.Duration             `mapstructure:"MAX_VOICE_MESSAGE_DURATION" validate:"required,max=60s"`
	GoogleTextToSpeech      string                    `mapstructure:"GOOGLE_TEXT_TO_SPEECH" validate:"required"`
	GoogleVoicesForUsers    map[string]map[int]string `mapstructure:"GOOGLE_VOICES" validate:"required"`
	GoogleDefaultVoice      string                    `mapstructure:"GOOGLE_DEFAULT_VOICE" validate:"required"`
}

func LoadConfig(configPath string) (*Config, error) {
	var cfg Config

	viper.SetConfigFile(configPath)
	viper.SetConfigType("yaml")

	// Read from .yml file
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("unable to read config from file: %v", err)
	}

	// Read from environment variables
	viper.AutomaticEnv()
	viper.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, fmt.Errorf("unable to unmarshal config: %v", err)
	}

	validate := validator.New()
	if err := validate.Struct(&cfg); err != nil {
		return nil, fmt.Errorf("unable to validate config: %v", err)
	}

	return &cfg, nil
}
