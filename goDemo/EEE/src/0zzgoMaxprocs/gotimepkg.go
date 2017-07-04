package main

import (
	"fmt"
	"time"
)

func main() {
	//初始化通道
	ch11 := make(chan int, 1000)
	sign := make(chan byte, 1)
	//给ch11通道写入数据
	for i := 0; i < 1000; i++ {
		ch11 <- i
	}
	//单独起一个goroutine执行select
	go func() {
		var e int
		ok := true
		//首先声明一个 *time.Timer 类型的值 然后在相关case之后声明的匿名函数中尽可能的复用它
		var timer *time.Timer
		for {
			select {
			case e = <-ch11:
				fmt.Printf("ch11 -> %d\n", e)
			case <-func() <-chan time.Time {
				if timer == nil {
					//初始化到期时间据此间隔1ms的定时器
					timer = time.NewTimer(time.Millisecond)
				} else {
					//复用 通过Reset方法重置定时器
					timer.Reset(time.Millisecond)
				}
				//得知定时器到期事件来临时 返回结果
				return timer.C
			}():
				fmt.Println("Timeout..")
				ok = false
				break
			}

			//终止for循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()
	//惯用手法 读取sign通道数据 为了等待select的goroutine执行
	<-sign
}

/*
####time包的定时器/断续器
定时器
在time包中有两个函数可以帮助我们初始化time.Timer
time.Newtimer函数
初始化一个到期时间据此的间隔为3小时30分的定时器
t:=time.Newtimer(3*time.Hour+30*timeMinute)
注意 这里的变量t是*time.NewTimer类型的 这个指针类型的方法集合包含两个方法
Rest 用于重置定时器 该方法返回一个bool类型的值
Stop 用来停止定时器  该方法返回一个bool类型的值  如果返回false 说明该定时器在之前已经到期或者已经被停止了 反之 返回true

通过定时器的字段C 我们可以及时得知定时器到期的这个事件来临 C是一个chan time.Time 类型的缓冲通道 一旦触及到期时间
定时器就会向自己的C字段发送一个time.Time类型的元素值

##示例一 一个简单的定时器
import (
	"fmt"
	"time"
)

func main() {
	//初始化定时器
	t := time.NewTimer(16 * time.Second)
	//当前时间
	now := time.Now()
	fmt.Printf("Now time : %v.\n", now)
	expire := <-t.C
	fmt.Printf("Expiration time: %v.\n", expire)
}

##示例二 我们改造一下之前那个简单的超时操作
import (
	"fmt"
	"time"
)

func main() {
	//初始化通道
	ch11 := make(chan int, 1000)
	sign := make(chan byte, 1)
	//给ch11通道写入数据
	for i := 0; i < 1000; i++ {
		ch11 <- i
	}
	//单独起一个goroutine执行select
	go func() {
		var e int
		ok := true
		//首先声明一个 *time.Timer 类型的值 然后在相关case之后声明的匿名函数中尽可能的复用它
		var timer *time.Timer
		for {
			select {
			case e = <-ch11:
				fmt.Printf("ch11 -> %d\n", e)
			case <-func() <-chan time.Time {
				if timer == nil {
					//初始化到期时间据此间隔1ms的定时器
					timer = time.NewTimer(time.Millisecond)
				} else {
					//复用 通过Reset方法重置定时器
					timer.Reset(time.Millisecond)
				}
				//得知定时器到期事件来临时 返回结果
				return timer.C
			}():
				fmt.Println("Timeout..")
				ok = false
				break
			}

			//终止for循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()
	//惯用手法 读取sign通道数据 为了等待select的goroutine执行
	<-sign
}






















*/
