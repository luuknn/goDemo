package main

import (
	"fmt"
	"time"
)

func main() {
	//初始化通道
	ch11 := make(chan int, 10000)
	sign := make(chan int, 1)
	//给ch11 通道写入数据
	for i := 0; i < 10000; i++ {
		ch11 <- i
	}
	//关闭ch11通道
	//close(ch11) //为了看效果先注释掉
	//单独起一个goroutine 执行select
	go func() {
		var e int
		ok := true
		for {
			select {
			case e, ok = <-ch11:
				if !ok {
					fmt.Println("End.")
					break
				}
				fmt.Printf("ch11 ->%d\n", e)
			case ok = <-func() chan bool {
				//经过大约1ms后 该接收语句会从timeout接收到一个新元素 并赋值给ok
				//从而恰当地执行了针对单个操作的超时子流程 恰当的结束了当前的for循环
				timeout := make(chan bool, 1)
				go func() {
					time.Sleep(time.Millisecond)
					timeout <- false
				}()
				return timeout
			}():
				fmt.Println("timeout..")
				break
			}
			//通道关闭后退出for循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()
	//惯用手法 读取sign通道数据 为了等待select的goroutine的执行
	<-sign
}

/*
select
golang的select 就是监听IO操作 当IO操作发生时 触发相应的动作
在执行select语句的时候  运行时 系统会自上而下的判断每一个case中发送或接收操作是否可以被立即执行
立即执行:意思是当前的goroutine不会因此操作而被阻塞 还需要依据通道的具体特性 缓存或非缓存

每一个case语句里必须是一个IO操作
所有channel表达式都会被求值 所有被发送的表达式都会被求值
如果任意某个case可以进行 它就执行
如果有多个case都可以运行 select会随机公平的选出一个执行
如果有default 子句 case不满足条件时执行
如果没有default字句 select将阻塞 直到某个case可以运行 go不会重新对channel或值进行求值

select 语句用法
注意到select的代码形式和switch 非常相似 不过select的case里的操作只能是 IO操作
此示例里面 select会一直等待等到某个case语句完成  也就是等到成功从ch1或者ch2中读到数据
如果都不满足且存在default case 那么defaultcase 会被执行 则 selec 语句结束
import (
	"fmt"
)

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	ch1 <- 1
	select {
	case e1 := <-ch1:
		//如果ch1通道成功读取数据 则执行该case处理语句
		fmt.Printf("1th case is selected . e1=%v", e1)
	case e2 := <-ch2:
		//如果ch2通道成功读取数据，则执行该case处理语句
		fmt.Printf("2th case is selected. e2=%v", e2)
	default:
		//如果上面case都没有成功，则进入default处理流程
		fmt.Println("default!.")
	}

}
select 分支选择规则
所有跟在case关键字右边的发送语句 或者接收语句中的通道表达式和元素表达式都会先被求值
无论它们所在的case是否有可能被选择 都会这样
求值顺序 自上而下 从左到右
import (
	"fmt"
)

//定义几个变量 其中chs和numbers分别代表了包含了有限元素和整数列表
var ch1 chan int
var ch2 chan int
var chs = []chan int{ch1, ch2}
var numbers = []int{1, 2, 3, 4, 5}

func getNumber(i int) int {
	fmt.Printf("numbers[%d]\n", i)
	return numbers[i]

}
func getChan(i int) chan int {
	fmt.Printf("chs[%d]\n", i)
	return chs[i]
}
func main() {
	select {
	case getChan(0) <- getNumber(2):
		fmt.Println("1th case is selected.")
	case getChan(1) <- getNumber(3):
		fmt.Println("2th case is selected")
	default:
		fmt.Println("default!")

	}
}

随机执行case
如果同时有多个case满足条件 通过一个伪随机的算法决定哪个case会被执行
func main() {
	chanCap := 5
	ch7 := make(chan int, chanCap)
	for i := 0; i < chanCap; i++ {
		select {
		case ch7 <- 1:
		case ch7 <- 2:
		case ch7 <- 3:
		}
	}
	for i := 0; i < chanCap; i++ {
		fmt.Printf("%v\n", <-ch7)
	}
}

####一些惯用手法示例
示例一  单独启用一个goroutine执行select 等待通道关闭后结束循环
import (
	"fmt"
)

func main() {
	//初始化通道
	ch11 := make(chan int, 10)
	sign := make(chan int, 1)
	//给ch11 通道写入数据
	for i := 0; i < 10; i++ {
		ch11 <- i
	}
	//关闭ch11通道
	close(ch11)
	//单独起一个goroutine 执行select
	go func() {
		var e int
		ok := true
		for {
			select {
			case e, ok = <-ch11:
				if !ok {
					fmt.Println("End.")
					break
				}
				fmt.Printf("ch11 ->%d\n", e)
			}
			//通道关闭后退出for循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()
	//惯用手法 读取sign通道数据 为了等待select的goroutine的执行
	<-sign
}

###示例二
加以改进 我们不想等到通道被关闭后再退出循环 利用一个辅助通道模拟出 操作超时

import (
	"fmt"
	"time"
)

func main() {
	//初始化通道
	ch11 := make(chan int, 1000)
	sign := make(chan int, 1)
	//给ch11 通道写入数据
	for i := 0; i < 1000; i++ {
		ch11 <- i
	}
	//关闭ch11通道
	close(ch11)
	//我们不想等到通道被关闭之后再退出循环 我们创建并初始化一个辅助的通道
	//利用它模拟出操作超时行为
	timeout := make(chan bool, 1)
	go func() {
		time.Sleep(time.Millisecond)
		timeout <- false

	}()
	//单独起一个goroutine 执行select
	go func() {
		var e int
		ok := true
		for {
			select {
			case e, ok = <-ch11:
				if !ok {
					fmt.Println("End.")
					break
				}
				fmt.Printf("ch11 ->%d\n", e)
			case ok = <-timeout:
				//向timeout 通道发送元素false后 该case几乎马上就会被执行 ok=false
				fmt.Println("timeout..")
				break
			}
			//通道关闭后退出for循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()
	//惯用手法 读取sign通道数据 为了等待select的goroutine的执行
	<-sign
}

###示例三
上面实现了单个操作的超时 但是那个超时触发器开始计时有点早

import (
	"fmt"
	"time"
)

func main() {
	//初始化通道
	ch11 := make(chan int, 10000)
	sign := make(chan int, 1)
	//给ch11 通道写入数据
	for i := 0; i < 10000; i++ {
		ch11 <- i
	}
	//关闭ch11通道
	//close(ch11) //为了看效果先注释掉
	//单独起一个goroutine 执行select
	go func() {
		var e int
		ok := true
		for {
			select {
			case e, ok = <-ch11:
				if !ok {
					fmt.Println("End.")
					break
				}
				fmt.Printf("ch11 ->%d\n", e)
			case ok = <-func() chan bool {
				//经过大约1ms后 该接收语句会从timeout接收到一个新元素 并赋值给ok
				//从而恰当地执行了针对单个操作的超时子流程 恰当的结束了当前的for循环
				timeout := make(chan bool, 1)
				go func() {
					time.Sleep(time.Millisecond)
					timeout <- false
				}()
				return timeout
			}():
				fmt.Println("timeout..")
				break
			}
			//通道关闭后退出for循环
			if !ok {
				sign <- 0
				break
			}
		}
	}()
	//惯用手法 读取sign通道数据 为了等待select的goroutine的执行
	<-sign
}

####非缓冲的Channel























*/
