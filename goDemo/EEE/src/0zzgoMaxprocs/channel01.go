package main

//import (
//	"fmt"
//)
//
//type Person struct {
//	Name    string
//	Age     uint8
//	Address Addr
//}
//type Addr struct {
//	city     string
//	district string
//}
//
//type PersonHandler interface {
//	Batch(origs <-chan Person) <-chan Person
//	Handle(orig *Person)
//}
//type PersonHandlerImpl struct{}
//
//func (handler PersonHandlerImpl) Batch(origs <-chan Person) <-chan Person {
//	//初始化通道dests
//	dests := make(chan Person, 10)
//	go func() {
//		//需要被更改的人员信息会通过origs单向通道传递进来 那么我们就应该不断地试图从该通道中接收他们
//		for p := range origs {
//			//变更人员信息
//			handler.Handle(&p)
//			//把人员信息发送给通道dests
//			dests <- p
//		}
//		fmt.Println("All the information has been handled.")
//		//关闭通道dests
//		close(dests)
//
//	}()
//	return dests
//}
//func (handler PersonHandlerImpl) Handle(orig *Person) {
//	//处理人员信息
//	if orig.Address.district == "Baoshan" {
//		orig.Address.district = "Xuhui"
//	}
//}
//func getPersonHandler() PersonHandler {
//	return PersonHandlerImpl{}
//}
//
//var personTotal = 20
//var persons []Person = make([]Person, personTotal)
//var personCount int
//
//func init() {
//	//初始化人员信息
//	for i := 0; i < personTotal; i++ {
//		name := fmt.Sprintf("%s%d", "P", i)
//		p := Person{name, 24, Addr{"Shanghai", "Baoshan"}}
//		persons[i] = p
//	}
//
//}
//func main() {
//	handler := getPersonHandler()
//	//初始化通道origs
//	origs := make(chan Person, 10)
//	//启用G2来处理人员信息
//	dests := handler.Batch(origs)
//	//启用G3来获取人员信息
//	fecthPerson(origs)
//	//启用G4以存储人员信息
//	sign := savePerson(dests)
//	<-sign
//}
//
////接受一个参数 是只允许写入origs通道
//func fecthPerson(origs chan<- Person) {
//	go func() {
//		for _, p := range persons {
//			origs <- p
//		}
//		fmt.Println("All the information has been fetched.")
//		close(origs)
//	}()
//
//}
//
////接受一个参数 是只允许读取dest通道  除非直接强制转换 要么你只能从channel中读取数据
//func savePerson(dest <-chan Person) <-chan byte {
//	sign := make(chan byte, 1)
//	go func() {
//		for {
//			p, ok := <-dest
//			if !ok {
//				fmt.Println("All the information has been saved.")
//				sign <- 0
//				break
//			}
//			savePerson1(p)
//		}
//	}()
//	return sign
//}
//
//func savePerson1(p Person) bool {
//	fmt.Println(p)
//	return true
//}

/*
channel是什么
在Go语言中 Channel即指通道类型 有时也用它来直接指代可以传递某种类型的值的通道

类型的表示法
chan T
关键字chan代表了通道类型的关键字 T则代表了该通道类型的元素类型
例如 type IntChan chan int别名类型IntChan代表了元素类型为int的通道类型  我们可以直接声明一个chan int类型的变量
var IntChan chan int 在被初始化后 变量IntChan就可以被用来传递int类型的元素值了
chan<-T
只能用来发送值  <-表示发送操作符
<-chan
接收通道值 <-表示接收操作符
值表示法
属性和基本操作
基于通道的通讯是在多个goroutine之间进行同步的重要手段 而针对通道本身也是同步的
在同一时刻 仅有一个goroutine 能向一个通道发送元素值
同时也仅有一个goroutine能从它那里接收元素值
通道相当于FIFO先进先出的消息队列
通道中元素都具有原子性 它们是不可被分割的通道中的每一个元素 只可能被某一个goroutine接收 已被接受的
元素值会立刻被从通道中删除

初始化通道
make(chan int,10)
表达式初始化了一个通道类型的值 传递给make函数的第一个参数表明其值的具体类型是元素类型为int的通道类型
而第二个参数指的是 在同一时刻 最多可以容纳10个元素值
import (
	"fmt"
)

type Person struct {
	Name    string
	Age     uint8
	Address Addr
}
type Addr struct {
	city     string
	district string
}

func main() {
	personChan := make(chan Person, 1)
	p1 := Person{"Harry", 32, Addr{"Qidong", "Xian"}}
	fmt.Printf("P1(1):%v\n", p1)
	personChan <- p1
	p1.Address.district = "Shi"
	fmt.Printf("P1(2):%v\n", p1)
	p1_copy := <-personChan
	fmt.Printf("P1(3):%v\n", p1_copy)

}
通道中的元素值丝毫没有受到外界的影响 这说明了 在发送过程中 进行的元素值 属于完全复制 这也保证了我们使用通道传递的值的不变性

单向通道
单向channel只能用于发送或者接收数据
var ch1 chan int//ch1是一个正常的channel 不是单向的
var ch2 chan<- floact64//ch2是单向的channel 只用于写float64数据
var ch3 <-chan  int//ch3是单向channel 只用于读取int数据
channel是一个原生类型 因此不仅支持被传递 还支持 类型转换 只有在介绍了单向channel
的概念后 我们才会明白 类型转换对于channel的意义 就是在单向channel和双向channel之间进行转换
ch4 :=make(chan int)
ch5:=<-chan int(ch4) //ch5 就是一个单向的读取channel
ch6:=chan<- int(ch4)//ch6 是一个单向的写入channel
基于ch4 我们通过类型转换 初始化了两个单向channel 单向读的ch5 和单向写的ch6
func Parse(ch <-chan int)
{
for value:=range ch{
fmt.Println("Parsing value",value)
}
}
关闭通道
close(strChan)
我们应该先明确一点 无论怎么样都不应该在接收端 关闭通道 因为在那里我们无法判断发送端是否还会向通道发送元素值
如何判断一个channel是否已经被关闭 我们可以在读取的时候使用多重返回值的方式
str,ok:=strChan 只需要判断第二个bool的返回值即可 false 表示strChan已经被关闭

import (
	"fmt"
	"time"
)

func main() {
	ch := make(chan int, 5)
	sign := make(chan int, 2)
	go func() {
		for i := 0; i < 5; i++ {
			ch <- i
			time.Sleep(1 * time.Second)
		}
		close(ch)
		fmt.Println("The channel is closed.")
		sign <- 0
	}()
	go func() {
		for {
			e, ok := <-ch
			fmt.Printf("%d (%v)\n", e, ok)
			if !ok {
				break
			}
			time.Sleep(2 * time.Second)
		}
		fmt.Println("Done.")
		sign <- 1

	}()
	<-sign
	<-sign

}

运行时系统并没有在通道ch被关闭之后立即把false作为相应接收操作的第二个结果 而是等到接收端把已在chan通道中的所有元素的值
都接收到了之后才这样做 这确保了在发送端关闭通道的安全性

























*/
