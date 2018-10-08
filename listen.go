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

var listeningChatMap = make(map[string]map[string]ListeningChat)

func listen(cid ubm_api.CID, from *ubm_api.User, pattern string) {
	if ok, _ := isMyMaster(from.UID); !ok {
		sendText(cid, "只有主人才可以跟咱做那种事情哦～！")
		return
	}
	goshujinsama  := from.PrivateChat
	if pattern == "" {
		pattern = `[\s\S]*`
	}
	reg := regexp.MustCompile(pattern)
	if lc, ok := listeningChatMap[cid.String()]; !ok || lc == nil {
		listeningChatMap[cid.String()] = make(map[string]ListeningChat)
	}
	listeningChatMap[cid.String()][goshujinsama.String()] = ListeningChat{
		Goshujinsama: goshujinsama,
		Regexp:       reg,
	}
	sendText(
		goshujinsama,
		"要开始了哦！",
	)
}

func stopListen(cid ubm_api.CID, goshujinsama ubm_api.CID) {
	if lc, ok := listeningChatMap[cid.String()]; !ok {
		return
	} else if lc == nil {
		delete(listeningChatMap, cid.String())
	} else {
		delete(lc, goshujinsama.String())
		if len(lc) == 0 {
			delete(listeningChatMap, cid.String())
		}
	}
	sendText(
		goshujinsama,
		"索然无味！",
	)
}

func onListen(message *ubm_api.Message) {
	if message.Type != "rich_text" || message.RichText == nil {
		return
	}
	lc2, ok := listeningChatMap[message.Chat.CID.String()]
	if !ok || lc2 == nil || len(lc2) == 0 {
		return
	}
	text := ""
	for _, elem := range *message.RichText {
		if elem.Type == "text" {
			text += elem.Text
		}
	}
	msg := append(ubm_api.RichText{
		{
			Type: "text",
			Text: fmt.Sprintf("Chat: %v, From: %v \n", message.Chat, message.From),
		},
	}, *message.RichText...)
	for _, lc := range lc2 {
		if lc.Regexp.MatchString(text) {
			sendMessage(
				lc.Goshujinsama,
				ubm_api.Message{
					Type:     "rich_text",
					RichText: &msg,
				},
			)
		}
	}
}
