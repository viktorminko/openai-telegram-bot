package messenger

const CommandImage = "image"

type Message struct {
	Audio            []byte
	ChatID           int64
	RequestMessageID int
	Command          string
	Text             string
}
