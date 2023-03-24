package bot

import "github.com/viktorminko/openai-telegram-bot/ai"

type ChatContext struct {
	queue *MessageQueue
}

func NewChatContext(queue *MessageQueue) *ChatContext {
	return &ChatContext{
		queue: queue,
	}
}

func (c *ChatContext) AddMessage(msg ai.ChatMessage) {
	c.queue.Push(msg)
}

func (c *ChatContext) GetMessages() []ai.ChatMessage {
	return c.queue.Messages
}
