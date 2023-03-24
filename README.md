# README.md

## OpenAI Telegram Bot

This Telegram bot is designed to interact with users in a conversational manner using the OpenAI API. The bot can process text and audio messages, and it can also respond with images when given specific commands.

### Key Features

1. Engage in text-based conversations with users using the OpenAI API.
2. Respond to image commands with relevant images.
3. Convert audio messages to text, process the text, and generate appropriate responses.

### Configuration

Before running the application, ensure that you have set up the necessary configurations. The main fields that need to be configured include:

1. `TELEGRAM_BOT_TOKEN`: Your Telegram bot token.
2. `OPENAI_API_KEY`: Your OpenAI API key.
3. `OPENAI_CHAT_MODEL`: The OpenAI chat model you want to use (e.g., "davinci-codex").
4. `OPENAI_TRANSCRIPT_MODEL`: The OpenAI transcription model you want to use (e.g., "facebook/wav2vec2-large-960h").
5. `IMAGE_SIZE`: The size of the images to be returned by the bot.
6. `CONTEXT_SIZE`: The maximum context size for conversations.
7. `MAX_VOICE_MESSAGE_DURATION`: The maximum duration for voice messages.

### How to Run the App

1. Pull the latest Docker image from the GitHub Container Registry using the following command:

```docker pull ghcr.io/viktorminko/openai-telegram-bot:latest```

2. Run the Docker container using the following command:
```
docker run -d
-e TELEGRAM_BOT_TOKEN=<your_telegram_bot_token>
-e OPENAI_API_KEY=<your_openai_api_key>
-e OPENAI_CHAT_MODEL=<openai_chat_model>
-e OPENAI_TRANSCRIPT_MODEL=<openai_transcript_model>
-e IMAGE_SIZE=<image_size>
-e CONTEXT_SIZE=<context_size>
-e MAX_VOICE_MESSAGE_DURATION=<max_voice_message_duration>
--name openai-telegram-bot
ghcr.io/viktorminko/openai-telegram-bot:latest
```


Replace `<your_telegram_bot_token>` with your Telegram bot token, `<your_openai_api_key>` with your OpenAI API key, `<openai_chat_model>` with the desired OpenAI chat model, `<openai_transcript_model>` with the desired OpenAI transcription model, `<image_size>` with the desired image size, `<context_size>` with the maximum context size for conversations, and `<max_voice_message_duration>` with the maximum duration for voice messages.
