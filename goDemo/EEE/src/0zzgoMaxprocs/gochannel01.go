package main

/*
import (
	"fmt"
)

func main() {
	var x = make(chan int)
	go func() {
		x <- 1
	}()
	c := <-x
	fmt.Println(c)

}


channel是Go语言中的一个非常重要的特性
要想理解channel要先知道CSP模型 CSP Communicating Sequential Process的简称 中文可以叫做 通信顺序进程 是一种并发编程模型
简单来讲 CSP模型 由并发执行的实体 线程或进程所组成 实体之间通过发送消息 进行通信 这里发送消息时使用的就是通道 或者叫做channel
CSP模型的关键是关注channel 而不关注发送消息的实体 Go语言实现了CSP部分理论 goroutine对应CSP中并发执行的实体 channel也就对应着CSP中的
channel
channel的基础知识 创建channel
unBufferChan:=make(chan int)//1
创建的是无缓冲的channel
buffChan:=make(chan int,N)//2
创建的是有缓冲的channel
如果使用channel之前没有make 会出现dead lock错误
func main() {
	var x =chan int
	go func() {
		x <- 1
	}()
      <-x
}
channel 读写操作
ch:=make(chan int,10)
//读操作
x<-ch
//写操作
ch<-x
channel的种类
channel分为无缓冲channel 和有缓冲channel 两者的区别如下
无缓冲 发送和接收动作是同时发生的 如果没有goroutine读取channel 发送者会一直阻塞
缓冲 缓冲channel类似一个有容量的队列 当队列满的时候 发送者会阻塞 当队列空的时候 接受者会阻塞

关闭channel
channel可以通过built-in 函数close()来关闭
ch:=make(chan int)
close(ch)
关于关闭channel有几点需要注意的是
重复关闭channel会导致panic
向关闭的channel发送数据会panic
从关闭的channel读取数据不会panic 读出channel中已有的数据之后再读就是channel类似的默认值 比如chan int 类型的channel关闭之后读取到的值为0

对于上面的的第三点 我们需要区分一下 channel 的值是默认值 还是channel关闭了  使用ok-idiom方式 这种方式在map中比较常用
ch :=make(chan int,10)
...
close(ch)
val,ok:=<-ch
if ok==false{
//channel closed
}
channel的典型用法
func main(){
x:=make(chan int)
go func(){
x<-1
}()
<-x
}
select
select 一定程度上可以类比于linux中的IO多路复用中的select
后者相当于提供了多个IO事件的统一管理 而Golang中的select相当于提供了对多个channel的统一管理
当然这只是select在channel上的一种使用方法
select{
case e,ok:=<-ch1:
...
case e,ok:=<-ch2:
...
default:
}
值得注意的是select中的break只能跳到select这一层
select使用的时候一般配合for循环使用 像下面这样  因为正常select里面的流程 也就执行一遍 这么看来select中的break就稍显鸡肋了
所以使用break的时候一般配置label使用 label定义在for循环这一层
for {
    select {
        ...
    }
}

range channel
range channel 可以直接取得channel中的值 当我们使用range来操作channel的时候 一旦 channel关闭 channel内部数据读完之后循环自动结束
func consumer(ch chan int){
for x:=range ch{
fmt.Println(x)
...
}
}
func producer(ch chan int){
for _,v:=range values{
ch <-v
}
}
超时控制
在很多操作情况下 都需要超时控制 利用select实现超时控制 下面是一个简单的示例
select{
case <-ch:
//get data from ch
case <- time.After(2*time.Second)
//read data from ch timeout
}
类似的 上面的time.After 可以换成其他任何的异常控制流
生产者 消费者模型
利用缓冲channel 可以很轻松的实现生产者 消费者模型 上面的range 示例其实就是一个简单的生产者 消费者模型

单向channel
单向channel 顾名思义就是只能读或者写的channel 但是仔细一想 只能写的channel 如果不读其中的值有什么用呢
其实单向channel主要用在函数声明中
func foo(ch chan<- int)<-chan int{...}
foo的形参是一个只能写的channel 那么就表示函数foo只会对ch进行写 当然你传入的参数可以是一个普通channel foo返回值是一个
只能读取的channel 那么表示foo返回值规范用法 就是只能读取 这种写法 在Golang原生代码库中有非常多的示例
// Done returns a channel which is closed if and when this pipe is closed
// with CloseWithError.
func (p *http2pipe) Done() <-chan struct{} {
    p.mu.Lock()
    defer p.mu.Unlock()
    if p.donec == nil {
        p.donec = make(chan struct{})
        if p.err != nil || p.breakErr != nil {

            p.closeDoneLocked()
        }
    }
    return p.donec
}
也许你会说这么写在功能上和使用普通channel并不会有什么差别 确实是这样的 但是使用单向channel编程体现了一种非常优秀的编程规范
convention over configuration 约定优于配置 这种编程范式在Ruby中体现的尤为明显

总结
Golang的channel将goroutine隔离开 并发编程的时候可以将注意力放在channel上 在一定程度上 这个和消息队列的解耦功能还是挺像的






















*/
