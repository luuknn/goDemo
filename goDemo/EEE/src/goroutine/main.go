package main

import (
	"fmt"
	"runtime"
)

//从1至1亿循环叠加,并打印结果
func print1(c chan bool, n int) {
	x := 0
	for i := 1; i <= 100000000; i++ {
		x += i
	}
	fmt.Println(n, "test:", x)
	c <- true

}

func main1() {

	//使用多核运行程序
	runtime.GOMAXPROCS(runtime.NumCPU())
	c := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go print1(c, i)

	}
	for i := 0; i < 10; i++ {
		<-c
	}
	fmt.Println("Done.")

}
