package main

import (
	"fmt"
)

//有人把Go比作21世纪的C语言 第一因为Go语言设计简单 第二,21世纪最重要的就是并行程序设计 而Go从语言层面就支持了并行
//goroutine是Go并行设计的核心 goroutine 说到底 其实就是线程 但是它比线程更小 十几个 goroutine 可能体现在底层就是5,6个线程
//Go语言 内部帮你实现了这些goroutine之间的内存共享 执行goroutine 只需极少的栈内存 大概是4~5kb 当然会根据相应的数据伸缩 也正因为如此 可以同时运行成千上万个并发任务
//goroutine比thread更易用 更高效 更轻便
//goroutine是通过go的runtime管理的一个线程管理器 goroutine通过关键字go实现 其实就是一个普通的函数
//go hello(a,b,c)
//通过关键字go就启动了一个goroutine 我们来看一个例子
package main

import (
	"fmt"
	"runtime"
)

func say(s string) {
	for i := 0; i < 5; i++ {
		runtime.Gosched()
		fmt.Println(s)
	}
}
func main() {
	go say("world") //开一个新的goroutine 执行
	say("hello")
}
//我们可以看到go关键字很方便的就实现了 并发编程 上面的多个goroutine运行在同一个进程里面 共享内存数据 
//不过设计上我们要遵循 不要通过共享来通信 而要通过通信来共享
//runtime.Gosched()表示 让CPU把时间片让给别人 下次某个时候 继续恢复执行该goroutine
//默认情况下调度器仅使用单线程  也就是说 只实现了并发 想要发挥多核处理器的并行
//需要我们在程序中显示的调用runtime.GOMAXPROCS(n)告诉调度器同时使用 多个线程
//GOMAXPROCS设置了同时运行逻辑代码的系统线程的最大数量 并返回之前的设置

//channels
//goroutine 运行在相同的地址空间 因此访问共享内存必须做好同步 那么goroutine之间如何进行数据的通信呢 Go
//提供了很好的通信机制 channel channel可以与unix shell中的双向管道做类比 可以通过它发送或者接收值
//这些值 只能是特定的类型 channel 类型 定义一个channel时 也需要定义发送到channel的值的类型
//必须使用make创建channel
ci :=make(chan int)
cs :=make(chan string)
cf :=make(chan interface{})
//channel通过操作符<-来接收和发送数据
ch<-v//发送v到channel ch
v:=<-ch//从ch中接收数据 并赋值给v
//我们把这些应用到我们的例子中来:
package main

import "fmt"

func sum(a []int, c chan int) {
	sum := 0
	for _, v := range a {
		sum += v
	}
	c <- sum //send sum to c
}

func main() {
	a := []int{7, 2, 8, -9, 4, 0}
	c := make(chan int)
	go sum(a[:len(a)/2], c)
	go sum(a[len(a)/2:], c)
	x, y := <-c, <-c //receive from c
	fmt.Println(x, y, x+y)
}
//默认情况下channel接收和发送数据都是阻塞的 除非另一端已经准备好 这样就使得goroutines同步变得更加的简单 而不需要显示的lock
//所谓阻塞 也就是 如果读取value :=<ch 它将会被阻塞  直到有数据接收 其次 任何发送(ch<-5)将会阻塞 直到数据被读出 无缓冲channel是在多个goroutine之间同步很棒的工具

//Buffered Channels
//上面我们介绍了 默认的非缓存类型的channel  不过Go也允许 指定channel的缓冲大小 很简单 就是channel 可以存储多少元素
//ch:=make(chan bool ,4)创建了可以存储4个元素的bool型的channel 在这个channel中 前4个可以无阻塞的写入 当写入第5个元素的时候 代码会被阻塞 直到其他goroutine从channel中读取一些元素 腾出空间
ch:=make(chan type,value)
value ==0 // 无缓冲（阻塞）
value >0 //缓冲 (非阻塞 直到value个元素)
//我们看一下下面这个例子  测试一下 修改相应的value值
package main

import "fmt"

func main() {
	c := make(chan int, 2)
	c <- 1
	c <- 2
	fmt.Println(<-c)
	fmt.Println(<-c)
}
//range 和close
//上面这个例子中 我们需要 读取两次c 这样不是很方便  Go考虑到这一点 所以也可以通过range 像操作slice或者map一样操作缓存类型的channel 请看下面的例子
package main

import "fmt"

func fibonacci(n int, c chan int) {
	x, y := 1, 1
	for i := 0; i < n; i++ {
		c <- x
		x, y = y, x+y

	}
	close(c)

}
func main() {
	c := make(chan int, 10)
	go fibonacci(cap(c), c)
	for i := range c {
		fmt.Println(i)
	}

}
//for i:=range c 能够不断的读取channel里面的数据 直到该channel被显示的关闭 上面代码 我们看到可以显示的关闭channel
//生产者 通过 close函数 关闭 channel 关闭channel 之后 就无法再发送任何数据了 在消费方可以通过语法v
//ok:=<ch 测试channel是否被关闭 如果ok返回false 那么说明channel已经没有任何数据并且已经被关闭
//记住 应该在生产者的地方关闭channel  而不是消费的地方去关闭它 这样容易引起panic
//另外记住一点的就是 channel不像文件之类的  不需要经常去关闭 只有当你确实没有任何发送数据了  或者 你想显示的结束range循环之类的

//select
//我们上面介绍的都是 只有一个channel的情况 如果 存在多个channel的时候 我们该如何操作呢 Go里面提供了一个关键字select 通过select 可以监听 channel上的数据流动
//select默认是阻塞的 只有当监听的channel中有发送或者接收可以进行时才会运行 当多个channel都准备好的时候 select是随机的选择一个执行的
package main

import "fmt"

func fibonacci(c, quit chan int) {
	x, y := 1, 1
	for {
		select {
		case c <- x:
			x, y = y, x+y
		case <-quit:
			fmt.Println("quit")
			return
		}
	}
}
func main() {
	c := make(chan int)
	quit := make(chan int)
	go func() {
		for i := 0; i < 10; i++ {
			fmt.Println(<-c)
		}
		quit <- 0
	}()
	fibonacci(c, quit)
}
//在select里面还有default语法 select 其实就是类似switch的功能 
//defaul 就是当监听的channel都没有准备好的时候 默认执行的 (select 不再阻塞等待channel)
select {
case i:=<-c
//use i
default:
//当c阻塞的时候执行这里
}

//超时
 //有时候会出现goroutine阻塞的情况 那么我们如何避免整个程序进入阻塞的情况呢 我们可以用select来设置超时 通过 如下的方式 实现
 func main(){
 c:=make(chan int)
 o:=make(chan bool)
 go func(){
 for {
 select {
 case v:=<-c
 println(v)
 case <-time.After(5*time.Second):
 println("timeout")
 o<-true
break 
 }
 }
 }()
 <-o
 }

//runtime goroutine
//runtime 包中有几个处理goroutine的函数
//Goexit
//退出当前执行的goroutine 但是defer函数还会继续调用

//Gosched
//让出当前goroutine的执行权限 调度器安排其他等待的任务运行 并在 下次某个时刻 从该位置 恢复执行

//NumCPU
//返回CPU核数量

//NumGoroutine
//返回正在执行和排队的任务总数

//GOMAXPROCS
//用来设置时可以运行的cpu核数
























































