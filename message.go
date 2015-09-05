package pab

type MessageType int

const (
	PublicMsg MessageType = iota
	DirectMsg
	PrivateMsg
)

// Message repesent incoming message from chat service.
type Message struct {
	Text      string
	ChannelID string
	Type      MessageType
}
