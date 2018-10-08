package main

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/projectriri/bot-gateway/adapters/jsonrpc-server-any/client/golang"
	"github.com/projectriri/bot-gateway/router"
	"github.com/projectriri/bot-gateway/types"
	"github.com/projectriri/bot-gateway/types/cmd"
	"github.com/projectriri/bot-gateway/types/ubm-api"
)

var C jsonrpc_sdk.Client
var Goshujinsama = make(map[string]ubm_api.UID)

func main() {

	_, err := toml.DecodeFile("config.toml", &config)
	if err != nil {
		panic(err)
	}

	heyMaster()

	C = jsonrpc_sdk.Client{}
	C.Init(config.Addr, config.UUID)

	C.Accept = []router.RoutingRule{
		{
			From: ".*",
			To:   ".*",
			Formats: []types.Format{
				{
					API:     "cmd",
					Version: "1.0",
					Method:  "cmd",
				},
				{
					API:     "ubm-api",
					Version: "1.0",
					Method:  "receive",
				},
			},
		},
	}

	C.Dial()
	pkts, _ := C.GetUpdatesChan(0)
	for pkt := range pkts {
		switch pkt.Head.Format.API {
		case "cmd":
			var command cmd.Command
			json.Unmarshal(pkt.Body, &command)
			switch command.CmdStr {
			case "rbq:say":
				if len(command.ArgsStr) == 0 {
					continue
				}
				sendText(command.Message.Chat.CID, command.ArgsStr)
			case "rbq:who":
				if ok, master := isMyMaster(command.Message.From.UID); ok {
					sendText(command.Message.Chat.CID, fmt.Sprintf(
						"%s 是咱的主人呐！", master,
					))
				} else {
					sendText(command.Message.Chat.CID, "汝好像不是咱的主人呐！")
				}
			case "rbq:listen":
				listen(command.Message.Chat.CID, command.Message.From.PrivateChat, command.ArgsStr)
			case "rbq:leave":
				stopListen(command.Message.Chat.CID, command.Message.From.PrivateChat)
			case "rbq:scdo":
				scdo(command)
			case "rbq:ntr":
				addScdoer(command.Message.Chat.CID, command.Message.From.UID, command.ArgsTxt)
			case "rbq:use":
				heyMaster()
				sendText(command.Message.Chat.CID, "啊啊啊不要、不要、啊……才没有感到好舒服")
			case "ping":
				fallthrough
			case "rbq:ping":
				sendText(command.Message.Chat.CID, "主人想要做什么？QvQ")
			}
		case "ubm-api":
			var ubm ubm_api.UBM
			json.Unmarshal(pkt.Body, &ubm)
			if ubm.Type != "message" || ubm.Message == nil {
				continue
			}
			onListen(ubm.Message)
		}
	}
}
