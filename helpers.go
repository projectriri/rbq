package main

import (
	"encoding/json"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/jsonrpc-any"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/ubm-api"
	"github.com/projectriri/bot-gateway/utils"
	"io/ioutil"
)

func heyMaster() {
	goshujinsama, err := ioutil.ReadFile("goshujinsama.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(goshujinsama, &Goshujinsama)
	if err != nil {
		panic(err)
	}
}

func isMyMaster(uid ubm_api.UID) (bool, string) {
	for name, master := range Goshujinsama {
		if master.Messenger == uid.Messenger && master.ID == uid.ID {
			return true, name
		}
	}
	return false, ""
}

func sendMessage(cid ubm_api.CID, msg ubm_api.Message) {
	msg.CID = &cid
	ubm := ubm_api.UBM{
		Type:    "message",
		Message: &msg,
	}
	b, _ := json.Marshal(ubm)
	C.MakeRequest(jsonrpc_any.ChannelProduceRequest{
		UUID: C.UUID,
		Packet: types.Packet{
			Head: types.Head{
				UUID: utils.GenerateUUID(),
				From: "rbq",
				To:   cid.Messenger,
				Format: types.Format{
					API:     "ubm-api",
					Version: "1.0",
					Method:  "send",
				},
			},
			Body: b,
		},
	})
}

func sendText(cid ubm_api.CID, text string) {
	msg := ubm_api.Message{
		Type: "rich_text",
		RichText: &ubm_api.RichText{
			{
				Type: "text",
				Text: text,
			},
		},
	}
	sendMessage(cid, msg)
}
