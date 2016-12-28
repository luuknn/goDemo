package main

import (
	"ToolRadix/redis"
	"fmt"
	"log"
	"strconv"
)

type Album struct {
	Title  string
	Artist string
	Price  float64
	Likes  int
}

func main() {
	conn, err := redis.Dial("tcp", "localhost:6379")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	/*
		//GET操作
		conn, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		title, err := conn.Cmd("HGET", "album:1", "title").Str()
		if err != nil {
			log.Fatal(err)
		}
		artist, err := conn.Cmd("HGET", "album:1", "artist").Str()
		if err != nil {
			log.Fatal(err)
		}
		price, err := conn.Cmd("HGET", "album:1", "price").Float64()
		if err != nil {
			log.Fatal(err)
		}
		likes, err := conn.Cmd("Hget", "album:1", "likes").Int()
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s by %s:£%.2f [%d likes]\n", title, artist, price, likes)*/

	/*
		//SET操作
		conn, err := redis.Dial("tcp", "localhost:6379")
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()
		resp := conn.Cmd("HMSET", "album:1", "title", "Electric Ladyland", "artist", "Jimi Hendrix", "price", 4.95, "likes", 8)
		// 使用client.Cmd()方法 通过它我们可以发送一个命令到我们的服务器 这里会返回给我们一个Resp对象的指针
		if resp.Err != nil {
			log.Fatal(resp.Err)
		}
		//在这个示例中 我们所关心的不是Redis返回什么 因为所有成功的操作都返回"ok"字符串
		//所以我们不需要对*Resp对象做错误检查 在这样的情况下 我们只需要检查err就好了
		fmt.Println("Electric Ladyland added!")*/
}
