package shell

import (
	"bufio"
	"fmt"
	"os"

	"github.com/jijeshmohan/pab"
)

type shellAdapter struct {
	bot  *pab.Bot
	stop chan struct{}
}

func (s *shellAdapter) Name() string {
	return "shell"
}

func init() {
	pab.RegisterAdapter("shell", func(b *pab.Bot) (pab.Adapter, error) {
		return &shellAdapter{b, make(chan struct{})}, nil
	})
}

func (s *shellAdapter) Send(res *pab.Response, msg string) {
	fmt.Printf("%s: %s\n", s.bot.Name, msg)
}

func (s *shellAdapter) Receive(msg *pab.Message) {
	s.bot.Receive(msg)
}

func (s *shellAdapter) Run() {
	fmt.Println("shell running")
	reader := bufio.NewReader(os.Stdin)
	prompt()
	go func() {
		for {
			if line, _, err := reader.ReadLine(); err == nil {
				msg := s.newMessage(string(line))
				s.Receive(msg)
				prompt()
			}
		}
	}()
	<-s.stop
}

func (s *shellAdapter) newMessage(line string) *pab.Message {
	return &pab.Message{
		Text: line,
		Type: pab.PrivateMsg, // in shell every message is a private msg
	}
}

func prompt() {
	fmt.Print("> ")
}

func (s *shellAdapter) Stop() {
	s.stop <- struct{}{}
}
