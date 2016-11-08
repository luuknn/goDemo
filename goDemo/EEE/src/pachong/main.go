package main

import (
	"fmt"
	"log"
	"net/http"
	"njut/com"
	"path"
	"regexp"
	"runtime"
	"strings"
)

var imgPattern = regexp.MustCompile(`src=\"http(.*).jpg`)

func download(url string, num chan bool, numy chan bool) {
	url = strings.TrimPrefix(url, "src=\"")
	log.Printf("正在下载:%s", url)
	err := com.HttpGetToFile(&http.Client{}, url, nil, "pics/"+path.Base(url))
	if err != nil {
		log.Printf("图片下载失败(%s):%v", url, err)
	}
	numy <- true
	<-num

}

func main() {
	//网址  http://www.jikexueyuan.com/course/dev/?pageNum=18
	runtime.GOMAXPROCS(runtime.NumCPU())
	num := make(chan bool, 10)
	numy := make(chan bool, 192)
	baseurl := "http://www.jikexueyuan.com/course/dev/?pageNum=%d"
	fmt.Println("Hello")

	for i := 2; i < 9; i++ {
		data, err := com.HttpGetBytes(&http.Client{}, fmt.Sprintf(baseurl, i), nil)
		if err != nil {
			log.Fatalf("获取页面失败(%d):%v", 0, err)
		}
		//log.Println(string(data))
		matches := imgPattern.FindAll(data, -1)
		for _, match := range matches {
			log.Println(string(match))
			num <- true
			go download(string(match), num, numy)

		}
	}
	for i := 0; i < 192; i++ {
		<-numy
	}

	fmt.Println("Done")

}
