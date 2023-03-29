package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/viktorminko/openai-telegram-bot/ai"
)

func TestMessageQueue_Push(t *testing.T) {
	queue := &MessageQueue{
		Messages: make([]ai.ChatMessage, 0),
		MaxSize:  30,
	}

	// Test adding messages to the queue
	queue.Push(ai.ChatMessage{Text: "message 1"})
	assert.Equal(t, 1, len(queue.Messages))
	assert.Equal(t, int64(9), queue.Size())

	queue.Push(ai.ChatMessage{Text: "message 2"})
	assert.Equal(t, 2, len(queue.Messages))
	assert.Equal(t, int64(18), queue.Size())

	queue.Push(ai.ChatMessage{Text: "message 3"})
	assert.Equal(t, 3, len(queue.Messages))
	assert.Equal(t, int64(27), queue.Size())

	// Test adding messages that exceed the max size
	queue.Push(ai.ChatMessage{Text: "message 4"})
	queue.Push(ai.ChatMessage{Text: "message 5"})
	queue.Push(ai.ChatMessage{Text: "message 6"})
	queue.Push(ai.ChatMessage{Text: "message 7"})
	assert.Equal(t, 3, len(queue.Messages))
	assert.Equal(t, int64(27), queue.Size())

	// Test adding messages that exceed the max size and max number of messages
	queue = &MessageQueue{
		Messages: make([]ai.ChatMessage, 0),
		MaxSize:  20,
	}
	queue.Push(ai.ChatMessage{Text: "message 1"})
	queue.Push(ai.ChatMessage{Text: "message 2"})
	queue.Push(ai.ChatMessage{Text: "message 3"})
	queue.Push(ai.ChatMessage{Text: "message 4"})
	assert.Equal(t, 2, len(queue.Messages))
	assert.Equal(t, int64(18), queue.Size())
}

func TestMessageQueue_Pop(t *testing.T) {
	queue := &MessageQueue{
		Messages: make([]ai.ChatMessage, 0),
		MaxSize:  30,
	}

	queue.Push(ai.ChatMessage{Text: "message 1"})
	queue.Push(ai.ChatMessage{Text: "message 2"})
	queue.Push(ai.ChatMessage{Text: "message 3"})

	// Test removing messages from the queue
	message := queue.Pop()
	assert.Equal(t, "message 1", message.Text)
	assert.Equal(t, 2, len(queue.Messages))
	assert.Equal(t, int64(18), queue.Size())

	message = queue.Pop()
	assert.Equal(t, "message 2", message.Text)
	assert.Equal(t, 1, len(queue.Messages))
	assert.Equal(t, int64(9), queue.Size())

	message = queue.Pop()
	assert.Equal(t, "message 3", message.Text)
	assert.Equal(t, 0, len(queue.Messages))
	assert.Equal(t, int64(0), queue.Size())

	// Test removing messages from an empty queue
	message = queue.Pop()
	assert.Equal(t, "", message.Text)
	assert.Equal(t, 0, len(queue.Messages))
	assert.Equal(t, int64(0), queue.Size())
}
