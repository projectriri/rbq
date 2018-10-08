package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	"strings"
	"time"
)

var scdoers = make(map[string]time.Time)

func addScdoer(cid ubm_api.CID, uid ubm_api.UID, args []string) {
	plus1s := func() time.Duration {
		duration := 15 * time.Minute
		if len(args) > 0 {
			d, err := time.ParseDuration(args[0])
			if err == nil {
				duration = d
			}
		}
		scdoers[cid.String()] = time.Now().Add(duration)
		return duration
	}
	switch checkScdoPrivilege(cid, uid) {
	case "":
		sendText(cid, "汝不是咱的主人！")
	case "scdoer":
		plus1s()
		sendText(cid, "我是绒布球哦！")
	default:
		d := plus1s()
		go func() {
			for {
				<-time.After(d)
				if time.Now().After(scdoers[cid.String()]) {
					delete(scdoers, cid.String())
					return
				}
			}
		}()
		sendText(cid, "主人允许了哦～")
	}
}

func checkScdoPrivilege(cid ubm_api.CID, uid ubm_api.UID) string {
	if _, ok := scdoers[cid.String()]; ok {
		return "scdoer"
	}
	if ok, _ := isMyMaster(uid); ok {
		return "master"
	}
	return ""
}

func scdo(cmd cmd.Command) {
	if checkScdoPrivilege(cmd.Message.Chat.CID, cmd.Message.From.UID) == "" {
		sendText(cmd.Message.Chat.CID, "只有主人才可以跟咱做那种事情哦～")
		return
	}

	if len(cmd.ArgsTxt) > 0 && cmd.ArgsTxt[0] != "--" {
		cmd.Message.Chat.CID.ChatID = cmd.ArgsTxt[0]
		cmd.ArgsTxt = cmd.ArgsTxt[1:]
	}
	if len(cmd.ArgsTxt) > 0 && cmd.ArgsTxt[0] != "--" {
		cmd.Message.Chat.CID.ChatType = cmd.ArgsTxt[0]
		cmd.ArgsTxt = cmd.ArgsTxt[1:]
	}
	if len(cmd.ArgsTxt) > 0 && cmd.ArgsTxt[0] != "--" {
		cmd.Message.Chat.CID.Messenger = cmd.ArgsTxt[0]
		cmd.ArgsTxt = cmd.ArgsTxt[1:]
	}
	if len(cmd.ArgsTxt) > 0 && cmd.ArgsTxt[0] == "--" {
		cmd.ArgsTxt = cmd.ArgsTxt[1:]
	}
	cmd.Message.Type = "rich_text"
	cmd.Message.RichText = &ubm_api.RichText{
		{
			Type: "text",
			Text: strings.Join(cmd.ArgsTxt, " "),
		},
	}
	ubm := ubm_api.UBM{
		Type:    "message",
		Message: cmd.Message,
	}
	b, _ := json.Marshal(ubm)
	C.MakeRequest(jsonrpc_any.ChannelProduceRequest{
		UUID: C.UUID,
		Packet: types.Packet{
			Head: types.Head{
				UUID: utils.GenerateUUID(),
				From: cmd.Message.Chat.CID.Messenger,
				Format: types.Format{
					API:     "ubm-api",
					Version: "1.0",
					Method:  "receive",
				},
			},
			Body: b,
		},
	})
}
