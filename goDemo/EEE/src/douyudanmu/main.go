package main

import (
	"douyudanmu/douyu"
	"flag"
	"fmt"
)

var (
	roomid = flag.Int("roomid", 258090, "房间号")
)

// 默认弹幕服务器
const (
	DefaultDouyuDanmuHost = "openbarrage.douyutv.com"
	DefaultDouyuDanmuPort = 8601
)

// DanmuHandle 为自定义的弹幕处理
func DanmuHandle(message *douyu.Message) {
	contentType, ok := message.Field("type")
	if !ok {
		return
	}
	switch contentType {
	case douyu.TypeChatMsg:
		// 默认全部为string
		nick, _ := message.Field("nn")
		level, _ := message.Field("level")
		text, _ := message.Field("txt")
		fmt.Printf("<level %s> - %s >>>》 %s\n", level, nick, text)
	case douyu.TypeUserEnter:
		nick, _ := message.Field("nn")
		level, _ := message.Field("level")
		fmt.Printf("!!!!!欢迎<lv %s> %s 进入房间\n", level, nick)
	}
}

func main() {
	flag.Parse()
	fmt.Println("===========房间号:", *roomid)

	douyuClient := douyu.New()
	douyuClient.Connect(DefaultDouyuDanmuHost, DefaultDouyuDanmuPort)

	douyuClient.JoinRoom(*roomid)
	douyuClient.HandleFunc(DanmuHandle)
	douyuClient.Watch()
}
