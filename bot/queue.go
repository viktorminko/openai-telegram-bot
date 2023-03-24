package bot

import "github.com/viktorminko/openai-telegram-bot/ai"

type MessageQueue struct {
	Messages []ai.ChatMessage
	MaxSize  int
}

func (q *MessageQueue) Push(m ai.ChatMessage) {
	if len(q.Messages) == q.MaxSize {
		q.Messages = q.Messages[1:]
	}
	q.Messages = append(q.Messages, m)
}

func (q *MessageQueue) Pop() ai.ChatMessage {
	if len(q.Messages) == 0 {
		return ai.ChatMessage{}
	}
	m := q.Messages[0]
	q.Messages = q.Messages[1:]
	return m
}
