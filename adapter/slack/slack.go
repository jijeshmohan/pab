package slack

import (
	"fmt"
	"os"
	"regexp"
	"time"

	"github.com/jijeshmohan/pab"
	"github.com/nlopes/slack"
)

var directRegex = regexp.MustCompile("^<@(\\w+)>: ")

type slackAdapter struct {
	bot        *pab.Bot
	stop       chan struct{}
	chSender   chan slack.OutgoingMessage
	chReceiver chan slack.SlackEvent
	api        *slack.Slack
	wsAPI      *slack.SlackWS
}

func (s *slackAdapter) Name() string {
	return "slack"
}

func init() {
	pab.RegisterAdapter("slack", func(b *pab.Bot) (pab.Adapter, error) {
		return &slackAdapter{
			bot:        b,
			stop:       make(chan struct{}),
			api:        slack.New(os.Getenv("SLACK_TOKEN")),
			chSender:   make(chan slack.OutgoingMessage),
			chReceiver: make(chan slack.SlackEvent),
		}, nil
	})
}
func (s *slackAdapter) Send(res *pab.Response, msg string) {
	s.chSender <- *s.wsAPI.NewOutgoingMessage(msg, res.Message().ChannelID)
}

func (s *slackAdapter) Receive(msg *pab.Message) {
	s.bot.Receive(msg)
}

func (s *slackAdapter) getSelfID() string {
	return s.api.GetInfo().User.Id
}

func (s *slackAdapter) Run() {
	wsAPI, err := s.api.StartRTM("", "")
	if err != nil {
		fmt.Errorf("%s\n", err.Error())
		panic(err)
	}
	go func(wsAPI *slack.SlackWS, chReceiver chan slack.SlackEvent) {
		defer func() {
			recover()
		}()
		wsAPI.HandleIncomingEvents(chReceiver)
	}(wsAPI, s.chReceiver)
	go wsAPI.Keepalive(20 * time.Second)
	go func(wsAPI *slack.SlackWS, chSender chan slack.OutgoingMessage) {
		for {
			select {
			case msg := <-s.chSender:
				wsAPI.SendMessage(&msg)
			}
		}
	}(wsAPI, s.chSender)
	s.wsAPI = wsAPI
	for {
		select {
		case _ = <-s.stop:
			return
		case msg := <-s.chReceiver:
			switch msg.Data.(type) {
			case *slack.MessageEvent:
				a := msg.Data.(*slack.MessageEvent)
				s.processMessage(a)
			case *slack.PresenceChangeEvent:
				a := msg.Data.(*slack.PresenceChangeEvent)
				fmt.Printf("Presence Change: %v\n", a)
			}
		}
	}
}

func (s *slackAdapter) isPrivateMessage(msg *slack.MessageEvent) bool {
	return msg.ChannelId[0] == 'D'
}

func (s *slackAdapter) isDirectMessage(msg *slack.MessageEvent) bool {
	matches := directRegex.FindStringSubmatch(msg.Text)
	if len(matches) < 2 {
		return false
	}

	if matches[1] == s.getSelfID() {
		return true
	}
	return false
}

func (s *slackAdapter) getDirectMessage(msg *slack.MessageEvent) string {
	return directRegex.ReplaceAllString(msg.Text, "")
}

func (s *slackAdapter) processMessage(msg *slack.MessageEvent) {
	// ignore self message
	if msg.UserId == s.getSelfID() {
		return
	}

	messageType := pab.PublicMsg
	text := msg.Text
	if s.isDirectMessage(msg) {
		messageType = pab.DirectMsg
		text = s.getDirectMessage(msg)
		fmt.Println(text)
	}
	if s.isPrivateMessage(msg) {
		messageType = pab.PrivateMsg
	}
	go s.Receive(&pab.Message{
		Text:      text,
		ChannelID: msg.ChannelId,
		Type:      messageType,
	})
}

func (s *slackAdapter) Stop() {
	s.stop <- struct{}{}
}
