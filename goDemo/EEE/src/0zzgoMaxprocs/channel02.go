package main

//import (
//	"fmt"
//	"time"
//)
//
//func main() {
//	unbufChan := make(chan int)
//	sign := make(chan byte, 2)
//	go func() {
//		for i := 0; i < 10; i++ {
//			select {
//			case unbufChan <- i:
//			case unbufChan <- i + 10:
//			}
//			fmt.Printf("The %d select is selected\n", i)
//			time.Sleep(time.Second)
//		}
//		close(unbufChan)
//		fmt.Println("The channel is closed.")
//		sign <- 0
//
//	}()
//	go func() {
//	loop:
//		for {
//			select {
//			case e, ok := <-unbufChan:
//				if !ok {
//					fmt.Println("Closed channel.")
//					break loop
//				}
//				fmt.Printf("e: %d\n", e)
//				time.Sleep(2 * time.Second)
//			}
//		}
//		sign <- 1
//	}()
//	<-sign
//	<-sign
//}

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
我们在初始化一个通道时  将其容量设置成0  或者直接忽略对容量的设置 那么就称之为非缓冲通道
ch1:=make(chan int,1)//缓冲通道
ch2:=make(chan int,0)//非缓冲通道
ch3:=make(chan int)//非缓冲通道
向此类通道发送元素值的操作会被阻塞 直到至少有一个针对该通道的接收操作开始进行为止
从此类通道接收元素值的操作会被阻塞 直到至少有一个针对该通道的发送操作开始进行为止
针对非缓冲通道的接收操作会在与之相应的发送操作完成之前完成

对于第三条要特别注意 发送操作在向非缓冲通道发送元素值的时候 会等待能够接收元素值的那个接收操作
并且确保该元素值被成功接收 它才会真正的完成执行 而缓冲通道中 刚好相反 由于元素值的传递是异步的 所以
发送操作在成功向通道发送元素值后就会立即结束 它不会关心是否有接收操作
##示例一
实现多个goroutine之间的同步
import (
	"fmt"
	"time"
)

func main() {
	unbufChan := make(chan int)
	//unbufChan:=make(chan int,1)有缓冲容量
	//启用一个goroutine接收元素值操作
	go func() {
		fmt.Println("Sleep a second...")
		time.Sleep(time.Second) //休息1s
		num := <-unbufChan      //接收unbufChan通道元素值
		fmt.Printf("Received a integer %d .\n", num)

	}()
	num := 1
	fmt.Printf("Send integer %d ...\n", num)
	//发送元素值
	unbufChan <- num
	fmt.Println("Done.")
}
在非缓冲channel中 从打印数据可以看出主goroutine的发送操作在等待一个能够与之配对的接收操作
配对成功后 元素值1 才得以经由unbufChan通道被从主goroutine传递至那个新的goroutine

####select与非缓冲通道
与操作缓冲通道的select相比 它被阻塞的概率一般会大很多  只有存在可配对的操作的时候 传递元素值的动作才能
真正的开始
示例
发送操作间隔1s 接收操作间隔2s
分别向unbufChan通道发送小于10和大于等于10的整数  这样更容易从打印结果分辨出配对的时候哪个case被选中了
下列案例两个case是被随机选择的
import (
	"fmt"
	"time"
)

func main() {
	unbufChan := make(chan int)
	sign := make(chan byte, 2)
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case unbufChan <- i:
			case unbufChan <- i + 10:
			default:
				fmt.Println("default!")
			}
			time.Sleep(time.Second)
		}
		close(unbufChan)
		fmt.Println("The channel is closed.")
		sign <- 0

	}()
	go func() {
	loop:
		for {
			select {
			case e, ok := <-unbufChan:
				if !ok {
					fmt.Println("Closed channel.")
					break loop
				}
				fmt.Printf("e: %d\n", e)
				time.Sleep(2 * time.Second)
			}
		}
		sign <- 1
	}()
	<-sign
	<-sign
}

default case 会在收发操作无法配对的情况下被选中并执行 在这里它被选中的概率是50%
上面的示例给予了我们这样一个启发 使用非缓冲通道能够让我们非常方便地在接收端对发送端的操作频率实施控制
可以尝试去掉defaultcase 看看打印结果 代码稍作修改如下
import (
	"fmt"
	"time"
)

func main() {
	unbufChan := make(chan int)
	sign := make(chan byte, 2)
	go func() {
		for i := 0; i < 10; i++ {
			select {
			case unbufChan <- i:
			case unbufChan <- i + 10:
			}
			fmt.Printf("The %d select is selected\n", i)
			time.Sleep(time.Second)
		}
		close(unbufChan)
		fmt.Println("The channel is closed.")
		sign <- 0

	}()
	go func() {
	loop:
		for {
			select {
			case e, ok := <-unbufChan:
				if !ok {
					fmt.Println("Closed channel.")
					break loop
				}
				fmt.Printf("e: %d\n", e)
				time.Sleep(2 * time.Second)
			}
		}
		sign <- 1
	}()
	<-sign
	<-sign
}
####总结
上面两个例子 第一个有default case 无法配对时执行该语句 而第二个没有default case 无法配对case时
select将阻塞 直到某个case可以运行 上述示例 是直到unbufChan数据被读取操作 不会重新对channel或值进行求值
*/
