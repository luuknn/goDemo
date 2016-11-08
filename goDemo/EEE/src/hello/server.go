package main

import (
	"io"
	"log"
	"net"
	"os"
	"strconv"
)

func handler(conn net.Conn) {
	defer conn.Close()
	remoteAddr := conn.RemoteAddr().String()
	log.Println("远程IP地址：", remoteAddr)

	//获取文件头信息,例如,文件名
	p := make([]byte, 1024)
	n, err := conn.Read(p) //n表示实际读取到的字节数
	if err != nil {
		log.Printf("读取头文件失败(%s):%v", remoteAddr, err)
		return
	} else if n == 0 {
		log.Printf("空头文件(%s)", remoteAddr)
		return
	}
	fileName := string(p[:n])
	log.Printf("文件:(%s)", fileName)

	//回复确认信息,避免客户端过早发送将文件内容和头信息混杂
	conn.Write([]byte("ok"))

	//打开一个本地文件流
	os.Mkdir("receive", os.ModePerm)
	f, err := os.Create("receive/" + fileName)
	if err != nil {
		log.Printf("无法创建文件(%s):%v", remoteAddr, err)
		return
	}
	defer f.Close()

	_, err = io.Copy(f, conn)
	for {
		buffer := make([]byte, 1024*200) //  每次读取200字节,可自定义
		_, err := conn.Read(buffer)
		if err != nil && err != io.EOF {
			log.Printf("读取失败(%s):%v", remoteAddr, err)
		} else if err == io.EOF {
			break
		}
	}
	if err != nil {
		log.Printf("文件接收失败(%s):%v", remoteAddr, err)
		return
	}
	log.Printf("文件接收成功(%s):%s", remoteAddr, fileName)
}

func runServer(port int) {
	//启动监听
	l, err := net.Listen("tcp", "localhost:"+strconv.Itoa(port))
	if err != nil {
		log.Fatalf("服务端启动监听失败:%v", err)
	}
	log.Println("服务端已启动！")
	//循环接受请求
	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); !ok || !ne.Temporary() {
				log.Printf("接受请求失败:%v", err)
			}
			continue
		}
		log.Println("请求接受成功!")
		go handler(conn)
	}
}
