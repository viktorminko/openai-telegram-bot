package messenger

const (
	CommandImage       = "image"
	CommandExportAudio = "export"
	CommandReset       = "reset"
)

type Message struct {
	Audio            []byte
	ChatID           int64
	FromUserID       int
	RequestMessageID int
	Command          string
	Text             string
}
