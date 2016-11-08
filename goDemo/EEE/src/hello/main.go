package main

import (
	"flag"
	"fmt"
	"log"
)

var (
	mode = flag.String("mode", "", "运行模式")
	port = flag.Int("port", 5050, "服务端监听端口")
	file = flag.String("file", "", "文件名称")
)

func main() {
	flag.Parse()

	fmt.Println("Hello World!")
	fmt.Println(*mode, *port, *file)
	//检查运行模式
	switch *mode {
	case "server":
		runServer(*port)
	case "client":
		runClient(*port, *file)
	default:
		log.Fatalf("未知的运行模式:%v", *mode)
	}
}
