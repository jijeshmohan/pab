package slack

import (
	"fmt"
	"os"
	"regexp"

	"github.com/jijeshmohan/pab"
	"github.com/nlopes/slack"
)

var directRegex = regexp.MustCompile("^<@(\\w+)>: ")

type slackAdapter struct {
	bot      *pab.Bot
	stop     chan struct{}
	chSender chan slack.OutgoingMessage
	api      *slack.Client
	wsAPI    *slack.RTM
}

func (s *slackAdapter) Name() string {
	return "slack"
}

func init() {
	pab.RegisterAdapter("slack", func(b *pab.Bot) (pab.Adapter, error) {
		return &slackAdapter{
			bot:      b,
			stop:     make(chan struct{}),
			api:      slack.New(os.Getenv("SLACK_TOKEN")),
			chSender: make(chan slack.OutgoingMessage),
		}, nil
	})
}
func (s *slackAdapter) Send(res *pab.Response, msg string) {
	params := slack.PostMessageParameters{}
	s.api.PostMessage(res.Message().ChannelID, msg, params)
}

func (s *slackAdapter) Receive(msg *pab.Message) {
	s.bot.Receive(msg)
}

func (s *slackAdapter) getSelfID() string {
	return s.wsAPI.GetInfo().User.ID
}

func (s *slackAdapter) Run() {
	wsAPI := s.api.NewRTM()
	go wsAPI.ManageConnection()
	s.wsAPI = wsAPI
	for {
		select {
		case _ = <-s.stop:
			return
		case msg := <-s.wsAPI.IncomingEvents:
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
	return msg.Channel[0] == 'D' //..ChannelId[0] == 'D'
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
	if msg.User == s.getSelfID() {
		return
	}

	messageType := pab.PublicMsg
	text := msg.Text
	if s.isDirectMessage(msg) {
		messageType = pab.DirectMsg
		text = s.getDirectMessage(msg)
	}
	if s.isPrivateMessage(msg) {
		messageType = pab.PrivateMsg
	}
	go s.Receive(&pab.Message{
		Text:      text,
		ChannelID: msg.Channel,
		Type:      messageType,
	})
}

func (s *slackAdapter) Stop() {
	s.stop <- struct{}{}
}
