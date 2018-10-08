package main

import (
	"fmt"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"regexp"
)

type ListeningChat struct {
	Goshujinsama ubm_api.CID
	Regexp       *regexp.Regexp
}

var listeningChatMap = make(map[string]ListeningChat)

func listen(cid ubm_api.CID, goshujinsama ubm_api.CID, pattern string) {
	if pattern == "" {
		pattern = `[\s\S]*`
	}
	reg := regexp.MustCompile(pattern)
	listeningChatMap[cid.String()+goshujinsama.String()] = ListeningChat{
		Goshujinsama: goshujinsama,
		Regexp:       reg,
	}
	sendText(
		goshujinsama,
		"要开始了哦！",
	)
}

func stopListen(cid ubm_api.CID, goshujinsama ubm_api.CID) {
	delete(listeningChatMap, cid.String()+goshujinsama.String())
	sendText(
		goshujinsama,
		"索然无味！",
	)
}

func onListen(message *ubm_api.Message) {
	if message.Type != "rich_text" || message.RichText == nil {
		return
	}
	lc, ok := listeningChatMap[message.Chat.CID.String()]
	if !ok {
		return
	}
	text := ""
	for _, elem := range *message.RichText {
		if elem.Type == "text" {
			text += elem.Text
		}
	}
	if lc.Regexp.MatchString(text) {
		msg := append(ubm_api.RichText{
			{
				Type: "text",
				Text: fmt.Sprintf("Chat: %v, From: %v \n", message.Chat, message.From),
			},
		}, *message.RichText...)
		sendMessage(
			lc.Goshujinsama,
			ubm_api.Message{
				Type:     "rich_text",
				RichText: &msg,
			},
		)
	}
}
