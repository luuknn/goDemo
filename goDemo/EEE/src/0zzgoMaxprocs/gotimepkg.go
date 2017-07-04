package main

//import (
//	"fmt"
//	"time"
//)
//
//func main() {
//	//初始化断续器 间隔1s
//	var ticker *time.Ticker = time.NewTicker(2 * time.Second)
//	//num为指定的执行次数
//	num := 3
//	c := make(chan int, num)
//	go func() {
//		for t := range ticker.C {
//			c <- 1
//			fmt.Println("Tick at", t)
//		}
//	}()
//	time.Sleep(10 * time.Second)
//	ticker.Stop()
//	fmt.Println("Ticker stop")
//
//}

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

####time.After函数
time.After函数 表示多少时间后 但是在取出channel内容之前不阻塞 后续程序可以继续执行
鉴于After特性 其通常用来处理程序超时问题
import (
	"fmt"
	"time"
)

func main() {
	ch1 := make(chan int, 1)
	ch2 := make(chan int, 1)
	select {
	case e1 := <-ch1:
		//如果ch1通道成功读取数据 则执行该case处理语句
		fmt.Printf("1th case is selected.e2=%v", e1)
	case e2 := <-ch2:
		//如果ch2通道成功读取数据 则执行该case处理语句
		fmt.Printf("2th case is selected.e2=%v", e2)
	case <-time.After(10 * time.Second):
		fmt.Println("TIME OUT...")

	}
}
time.Sleep 函数 表示休眠多少时间 休眠时处于阻塞状态 后续程序无法执行

####time.Afterfunc函数  示例三  自定义定时器
import (
	"fmt"
	"time"
)

func main() {
	var t *time.Timer
	f := func() {
		fmt.Printf("Expiration time: %v.\n", time.Now())
		fmt.Printf("C's len:%d \n", len(t.C))
	}
	t = time.AfterFunc(1*time.Second, f)
	//让当前goroutine 睡眠2s 确保大于内容的完整
	//这样做的原因是 time.AfterFunc的调用 不会被阻塞 它会以一部的方式在到期事件来临时执行我们的自定义函数
	time.Sleep(2 * time.Second)
	fmt.Println(time.Now(), "Hello")

}
第二行打印内容说明 定时器的字段C并没有缓冲任何元素值 这也说明了 在给定了自定义元素之后 默认的处理方法
(向C发送代表绝对到期时间的元素值)就不会被执行了

####断续器
结构体类型time.Ticker 表示了断续器的静态结构
就是周期性的传达到期时间的装置 这种装置的行为方式与仅有秒针的钟表有些类似 只不过间隔时间可以不是1s
初始化一个断续器
var ticker *timeTicker=time.NewTicker(time.Second)
//示例一: 使用时间控制停止ticker
import (
	"fmt"
	"time"
)

func main() {
	//初始化断续器 间隔2s
	var ticker *time.Ticker = time.NewTicker(2 * time.Second)
	go func() {
		for t := range ticker.C {
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(time.Second * 5) //阻塞 则执行次数是sleep的时间/ticker的时间
	ticker.Stop()
	fmt.Println("Ticker stopped.")
}

###示例二  使用channel 控制停止ticker
import (
	"fmt"
	"time"
)

func main() {
	//初始化断续器 间隔1s
	var ticker *time.Ticker = time.NewTicker(2 * time.Second)
	//num为指定的执行次数
	num := 3
	c := make(chan int, num)
	go func() {
		for t := range ticker.C {
			c <- 1
			fmt.Println("Tick at", t)
		}
	}()
	time.Sleep(10 * time.Second)
	ticker.Stop()
	fmt.Println("Ticker stop")

}


















*/
