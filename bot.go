package pab

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"
)

// Bot represent a chatter bot.
type Bot struct {
	Name     string
	Adapter  Adapter
	handlers []Handler
	osChan   chan os.Signal
}

// Run is bots run loop
func (b *Bot) Run() error {
	go b.Adapter.Run()

	signal.Notify(b.osChan, syscall.SIGINT, syscall.SIGTERM)

	<-b.osChan
	b.Stop()
	return nil
}

// Stop to stop the bot
func (b *Bot) Stop() {
	b.Adapter.Stop()
}

func (b *Bot) Receive(msg *Message) {
	res := &Response{
		bot: b,
		msg: msg,
	}
	for _, s := range b.handlers {
		if msg.Type < s.Method {
			continue
		}
		re := regexp.MustCompile(s.Pattern)
		matches := re.FindStringSubmatch(msg.Text)
		if len(matches) != 0 {
			res.Match = matches[1:]
			s.Run(res)
		}
	}
}

// AddHandlers to add different scripts to bot
func (b *Bot) AddHandlers(s ...Handler) {
	b.handlers = append(b.handlers, s...)
}

// NewBot create a new bot with provided configuration
func NewBot(conf *Config) (*Bot, error) {
	var err error
	b := &Bot{
		Name:     conf.Name,
		osChan:   make(chan os.Signal, 1),
		handlers: make([]Handler, 0),
	}
	b.Adapter, err = newAdapter(conf.Adapter, b)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	return b, nil
}
