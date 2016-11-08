package main

import (
	"fmt"
	"runtime"
	"sync"
)

//从1至1亿循环叠加,并打印结果
func print(wg *sync.WaitGroup, n int) {
	x := 0
	for i := 1; i <= 100000000; i++ {
		x += i
	}
	fmt.Println(n, "sum:", x)
	//标识一次任务完成
	wg.Done()

}

func main() {

	//使用多核运行程序
	runtime.GOMAXPROCS(runtime.NumCPU())
	wg := sync.WaitGroup{}
	//设置等待任务数
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go print(&wg, i)

	}
	//等待任务完成
	wg.Wait()
	fmt.Println("Done.")

}
