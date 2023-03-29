package bot

import "github.com/viktorminko/openai-telegram-bot/ai"

type MessageQueue struct {
	Messages []ai.ChatMessage
	MaxSize  int64
	size     int64
}

func (q *MessageQueue) Push(m ai.ChatMessage) {
	messageSize := int64(len(m.Text))
	for len(q.Messages) > 0 && q.size+messageSize > q.MaxSize {
		q.size -= int64(len(q.Messages[0].Text))
		q.Messages = q.Messages[1:]
	}
	q.Messages = append(q.Messages, m)
	q.size += messageSize
}

func (q *MessageQueue) Pop() ai.ChatMessage {
	if len(q.Messages) == 0 {
		return ai.ChatMessage{}
	}
	m := q.Messages[0]
	q.Messages = q.Messages[1:]
	q.size -= int64(len(m.Text))
	return m
}

func (q *MessageQueue) Size() int64 {
	return q.size
}
