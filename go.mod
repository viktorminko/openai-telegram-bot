module openai-telegram-bot

go 1.18

require (
	github.com/go-telegram-bot-api/telegram-bot-api v4.6.4+incompatible
	github.com/sashabaranov/go-gpt3 v0.0.0-20221216095610-1c20931ead68
)

require github.com/technoweenie/multipartstreamer v1.0.1 // indirect

replace github.com/sashabaranov/go-gpt3 => github.com/viktorminko/go-gpt3 v0.0.1
