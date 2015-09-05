package pab

// Handler represent a message handler
type Handler struct {
	Method  MessageType
	Pattern string
	Run     func(*Response)
}

// Listen for listesting a message in any chat room which matches a pattern
func Listen(pattern string, f func(*Response)) Handler {
	return Handler{
		Method:  PublicMsg,
		Pattern: pattern,
		Run:     f,
	}
}

// Direct for listesting for a direct message to the bot which matches a pattern
func Direct(pattern string, f func(*Response)) Handler {
	return Handler{
		Method:  DirectMsg,
		Pattern: pattern,
		Run:     f,
	}
}

// Private for listening a private message to the bot
func Private(pattern string, f func(*Response)) Handler {
	return Handler{
		Method:  PrivateMsg,
		Pattern: pattern,
		Run:     f,
	}
}
