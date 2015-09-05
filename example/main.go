package main

import (
	"errors"
	"fmt"
	"net"

	"github.com/jijeshmohan/pab"
	_ "github.com/jijeshmohan/pab/adapter/shell"
	_ "github.com/jijeshmohan/pab/adapter/slack"
)

func main() {
	conf := pab.NewConfig()
	conf.Adapter = "shell" // you can use slack, but make sure to set SLACK_TOKEN environment variable while running
	bot, err := pab.NewBot(conf)
	if err != nil {
		fmt.Println("Error while crating bot")
	}
	bot.AddHandlers(
		pab.Listen("deploy to (\\w+)", func(res *pab.Response) {
			res.Send(fmt.Sprintf("deploying in %s", res.Match[0]))
		}),
		pab.Listen("^current time$", func(res *pab.Response) {
			res.Send("Time is just an illusion")
		}),
		pab.Direct("record (\\d+) (?:hrs|hours)", func(res *pab.Response) {
			res.Send("Recorded your timesheet entry")
		}),
		pab.Private("what is my password", func(res *pab.Response) {
			res.Send("secret")
		}),
		pab.Listen("show your ip", func(res *pab.Response) {
			ip, err := externalIP()
			if err != nil {
				res.Send(err.Error())
				return
			}

			res.Send(ip)
		}),
	)

	if err = bot.Run(); err != nil {
		fmt.Println("Unable to run bot")
	}
}

func externalIP() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, iface := range ifaces {
		if iface.Flags&net.FlagUp == 0 {
			continue // interface down
		}
		if iface.Flags&net.FlagLoopback != 0 {
			continue // loopback interface
		}
		addrs, err := iface.Addrs()
		if err != nil {
			return "", err
		}
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if ip == nil || ip.IsLoopback() {
				continue
			}
			ip = ip.To4()
			if ip == nil {
				continue // not an ipv4 address
			}
			return ip.String(), nil
		}
	}
	return "", errors.New("are you connected to the network?")
}
