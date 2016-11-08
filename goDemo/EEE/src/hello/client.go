package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func runClient(port int, file string) {
	//连接服务器
	conn, err := net.Dial("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Printf("无法建立连接:%v", err)
		return
	}
	defer conn.Close()
	log.Println("链接建立成功！")

	//打开要传送的文件
	f, err := os.Open(file)
	if err != nil {
		log.Printf("无法打开文件:%v", err)
		return
	}
	defer f.Close()
	//写入头文件信息并等待确认
	conn.Write([]byte(file))

	p := make([]byte, 2)
	_, err = conn.Read(p)
	if err != nil {
		log.Printf("无法获得服务器端信息:%v", err)
		return
	} else if string(p) != "ok" {
		log.Printf("无效的服务器端相应:%s", string(p))
		return
	}
	log.Println("头信息发送成功！")
	_, err = io.Copy(conn, f)
	if err != nil {
		log.Print("发送文件失败(%s):%v", file, err)
		return
	}
	log.Println("文件发送成功！")
}
