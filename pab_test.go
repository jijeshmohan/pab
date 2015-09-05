package pab

import "testing"

type mockAdapter struct {
	bot  *Bot
	stop chan struct{}
}

func (s *mockAdapter) Name() string {
	return "mock"
}

func init() {
	RegisterAdapter("mock", func(b *Bot) (Adapter, error) {
		return &mockAdapter{b, make(chan struct{})}, nil
	})
}

func (s *mockAdapter) Send(res *Response, msg string) {
}

func (s *mockAdapter) Receive(msg *Message) {
	s.bot.Receive(msg)
}

func (s *mockAdapter) Run() {
	<-s.stop
}

func (s *mockAdapter) newMessage(line string) *Message {
	return &Message{
		Text: line,
		Type: PrivateMsg,
	}
}

func (s *mockAdapter) Stop() {
	s.stop <- struct{}{}
}

func TestBot(t *testing.T) {
	conf := NewConfig()
	conf.Adapter = "mock"
	bot, err := NewBot(conf)
	if err != nil {
		t.Fatal(err)
	}

	var count int
	bot.AddHandlers(
		Listen(".*", func(r *Response) {
			count = count + 1
		}),
	)
	go bot.Run()
	defer bot.Stop()

}
