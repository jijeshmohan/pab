package pab

type Response struct {
	bot   *Bot
	msg   *Message
	Match []string
}

func (r *Response) Send(str string) {
	r.bot.Adapter.Send(r, str)
}

func (r *Response) Message() *Message {
	return r.msg
}
