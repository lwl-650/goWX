package main

import (
	"fmt"
	"net/http"
	"reflect"

	wechat "github.com/silenceper/wechat/v2"
	"github.com/silenceper/wechat/v2/cache"
	offConfig "github.com/silenceper/wechat/v2/officialaccount/config"
	"github.com/silenceper/wechat/v2/officialaccount/message"
)

func serveWechat(rw http.ResponseWriter, req *http.Request) {
	wc := wechat.NewWechat()
	//这里本地内存保存access_token，也可选择redis memcache 或者自定cache
	memory := cache.NewMemory()
	cfg := &offConfig.Config{
		AppID:     "wxf3b2ce464e28d596",
		AppSecret: "9e04055c8e08d80828fb286ea972488c",
		Token:     "lwl",
		// EncodingAESKey: "",
		Cache: memory,
	}
	officialAccount := wc.GetOfficialAccount(cfg)

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	//设置接收消息的处理方法
	server.SetMessageHandler(func(mixMessage *message.MixMessage) *message.Reply {
		//TODO
		//回复消息：演示回复用户发送的消息
		switch mixMessage.MsgType {
		//自动回复消息
		case message.MsgTypeText:
			text := message.NewText(mixMessage.Content)
			return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
		}
		switch mixMessage.Event {
		//如果是初次订阅回复图文消息
		case message.EventSubscribe:
			//获取token
			miniprogrampagemessage := message.NewCustomerMiniprogrampageMessage(reflect.ValueOf(server.RequestMsg.FromUserName).String(), "title", "appId", "pages/index/index", "缩略图id")
			officialAccount.GetCustomerMessageManager().Send(miniprogrampagemessage)
			return &message.Reply{MsgType: message.MsgTypeMiniprogrampage, MsgData: miniprogrampagemessage}
			//根据token 发送客服消息

		}
		text := message.NewText(mixMessage.Content)

		return &message.Reply{MsgType: message.MsgTypeText, MsgData: text}
	})
	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		fmt.Println(err)
		return
	}
	//发送回复的消息
	server.Send()
}

func main() {
	http.HandleFunc("/", serveWechat)
	fmt.Println("wechat server listener at", ":8001")
	err := http.ListenAndServe(":8001", nil)
	if err != nil {
		fmt.Printf("start server error , err=%v", err)
	}
}
