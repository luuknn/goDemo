package main

import (
	"fmt"
)

//面向对象
//前面两章我们介绍了函数和struct 那你是否想过函数当作struct的字段一样处理呢
//今天我们就讲解一下函数的另一种形态 带有接受者的函数 我们称之为 method
//method
//现在假设有这么一个场景 我们定义了一个struct 叫长方形 现在要计算面积 那么 按照我们的一般思路 应该会用下面的方式来进行实现
package main

import "fmt"

type Rectangle struct {
	width, height float64
}

func area(r Rectangle) float64 {

	return r.width * r.height
}
func main() {
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	fmt.Println("Area of r1 is:", area(r1))
	fmt.Println("Area of r2 is:", area(r2))
}
//该段代码可以计算出长方形的面积 但是area 不是Rectangle的方法实现的 而是将Rectangle的对象
//r1 r2作为参数传入函数计算面积的
//这样的实现当然没有问题喽 但是要增加圆形 正方形 五边形的时候 你想计算他们的面积的时候怎么办呢
//那就只能增加新的函数 但是函数名就要跟着换了  变成area_circle
//像下图所表示的那样 椭圆 代表函数 而这些函数并不属于struct 他们是单独存在的
//很显然 这样的实现方式并不优雅  并且从概念上来说 面积是形状的一个属性 它是属于这个特定形状的 就像长和宽一样
//基于上面的原因 就有了 method的概念 method是附属在一个给定的类型上的 它的语法和函数声明几乎一样
//只是func后面增加了一个receiver 也就是method所依从的主体
//用上面提到的形状的例子来说 method area()是依赖于某个形状来发生作用的 Rectangle.area()的发出者
//是rectangle area是rectangle的方法 而非一个外围函数
//更具体的说 Rectangle存在字段length和width 同时存在方法area() 这些字段和方法都属于rectangle
//用Rob Pike的话来说 就是 "A method is a function with an implicit first argument,called a receiver"
//method 的语法如下
func (r ReceiverType) funcName(parameters)(results)
//下面我们用最开始的例子 用method来进行实现
package main

import (
	"fmt"
	"math"
)

type Rectangle struct {
	width, height float64
}
type Circle struct {
	radius float64
}

func (r Rectangle) area() float64 {

	return r.width * r.height
}
func (c Circle) area() float64 {

	return c.radius * c.radius * math.Pi
}
func main() {
	r1 := Rectangle{12, 2}
	r2 := Rectangle{9, 4}
	c1 := Circle{10}
	fmt.Println("Area of r1 is:", r1.area())
	fmt.Println("Area of r2 is:", r2.area())
	fmt.Println("Area of c1 is:", c1.area())
}
//在使用method的时候重要注意几点
//虽然method的名字一模一样 但是如果接受者不一样 那么method就不一样
//method里面可以访问接收者的字段
//调用method通过.访问 就像struct里面访问字段一样
//在上例 method area 分别属于rectangle和Circle 于是他们的receiver就变成了Rectangle和circle
//或者说这个area()方法是由rectangle和circle发出的
//值得注意的一点是 method虚线标出 意思是 以值传递 是的receiver可以是指针 两者差别在于
//指针作为receiver 会对实例对象的内容发生操作 而普通类型作为receiver仅仅是以副本作为操作对象 并不对原实例对象发生操作
//那是不是method只能作用在struct上面 当然不是
//它可以定义 任何你自定义的类型 内置类型 struct等各种类型上面
//struct只是自定义类型里面一种比较特殊的类型而已 还有其他自定义类型声明 可以通过如下这样的申明来实现
type typeName typeLiteral
//请看下面这个申明自定义类型的代码
type ages int 
type money float32
type months map[string]int
m:=month{
"January":31,
"February":28,
...
"December":31
}
//看到了吗 简单的狠吧 这样你就可以在自己的代码里面定义有意义的类型了 实际上 只是定义了一个别名
//例如上面的 ages 代替了int
//好了 让我们回到method
//你可以在任何的自定义类型中定义任意多的 method 接下来让我们看一个复杂一点的例子
package main

import "fmt"

const (
	WHITE = iota
	BLACK
	BLUE
	RED
	YELLOW
)

type Color byte
type Box struct {
	width, height, depth float64
	color                Color
}
type BoxList []Box //a slice of boxes
func (b Box) Volume() float64 {
	return b.depth * b.height * b.width
}
func (b *Box) SetColor(c Color) {
	b.color = c
}
func (bl BoxList) BiggestsColor() Color {
	v := 0.00
	k := Color(WHITE)
	for _, b := range bl {
		if b.Volume() > v {
			v = b.Volume()
			k = b.color
		}
	}
	return k
}
func (bl BoxList) PaintItBlack() {
	for i, _ := range bl {
		bl[i].SetColor(BLACK)
	}
}
func (c Color) String() string {
	strings := []string{"WHITE", "BLACK", "BLUE", "RED", "YELLOW"}
	return strings[c]
}
func main() {
	boxes := BoxList{
		Box{4, 4, 4, RED},
		Box{10, 10, 1, YELLOW},
		Box{1, 1, 20, BLACK},
		Box{10, 10, 1, BLUE},
		Box{10, 30, 1, WHITE},
		Box{20, 20, 20, YELLOW},
	}

	fmt.Printf("We have %d boxes in our set\n", len(boxes))
	fmt.Println("The volume of the first one is", boxes[0].Volume(), "cm³")
	fmt.Println("The color of the last one is", boxes[len(boxes)-1].color.String())
	fmt.Println("The biggest one is", boxes.BiggestsColor().String())

	fmt.Println("Let's paint them all black")
	boxes.PaintItBlack()
	fmt.Println("The color of the second one is", boxes[1].color.String())

	fmt.Println("Obviously, now, the biggest one is", boxes.BiggestsColor().String())
}
//现在让我们回头来看看setColor这个method 它的receiver是一个指向Box的指针 
//我们定义setcolor的真正目的是想改变 这个Box的颜色 如果不传Box的指针 那么setcoloe接受的
//其实是Box的一个copy 也就是说 method内对于颜色值的修改 其实只作用于copy而不是真正的Box 所以我们需要传入指针
//这里可以把receiver当作method的第一个参数来看 然后结合前面函数讲解的传值和传引用来看就不难理解
//这里你也许会问 那setcolor函数里面应该这样定义*b.color =c 而不是b.color =c 因为我们需要读取到 指针相应的值
//你是对的，其实Go里面这两种方式都是正确的，当你用指针去访问相应的字段时(虽然指针没有任何的字段)，Go知道你要通过指针去获取这个值，看到了吧，Go的设计是不是越来越吸引你了
//也许细心的读者会问 paintitblack里面调用setcolor的时候是不是应该写成(&bl[i]).setcolor(BLACK)
//因为setcolor的receiver是*Box 而不是Box
//你又说对了  这两种方式 都可以 因为Go知道 receiver是指针 它自动帮你转了
//也就是说
//如果一个method的receiver 是*T 你可以在一个T类型的实例变量V上面调用这个method 而不需要&V去调用这个method
//类似的
//如果一个method 的receiver是T 你可以在一个*T类型的变量P上面调用这个method 而不是*P去调用这个method
//所以 你不用担心你调用的指针的method 还是不是指针的method Go知道你要做的一切

//Method 继承
//前面一章我们学习了字段的继承 那么你会发现Go的一个神奇之处 method也是可以继承的
//如果匿名字段 实现了一个method 那么包含这个匿名字段的struct也能调用该method 让我们来看一下 这个例子
package main

import "fmt"

type Human struct {
	name  string
	age   int
	phone string
}
type Student struct {
	Human  //匿名字段
	school string
}
type Employee struct {
	Human   //匿名字段
	company string
}

//在Human上面定义了一个method
func (h *Human) SayHi() {
	fmt.Printf("Hi ,I am %s ,you can call me on %s\n", h.name, h.phone)
}
func main() {
	mark := Student{Human{"Mark", 25, "222-333-YYY"}, "MIT"}
	sam := Employee{Human: Human{"sam", 18, "111111-xxx"}, company: "Golang Inc"}
	mark.SayHi()
	sam.SayHi()
}

//method 重写
//在上面的例子中 如果Employee想要实现自己的SayHi 怎么办 简单 和匿名字段冲突一样的道理
//我们可以在Employee上面定义一个method 重写匿名字段的方法 请看下面的例子
package main

import "fmt"

type Human struct {
	name  string
	age   int
	phone string
}
type Student struct {
	Human  //匿名字段
	school string
}
type Employee struct {
	Human   //匿名字段
	company string
}

//在Human上面定义了一个method
func (h *Human) SayHi() {
	fmt.Printf("Hi ,I am %s ,you can call me on %s\n", h.name, h.phone)
}

//Employee 的method 重写 Human的method
func (e *Employee) SayHi() {
	fmt.Printf("Hi I am %s,I work at %s.Call me on %s \n", e.name, e.company, e.phone)
}
func main() {
	mark := Student{Human{"Mark", 25, "222-333-YYY"}, "MIT"}
	sam := Employee{Human: Human{"sam", 18, "111111-xxx"}, company: "Golang Inc"}
	mark.SayHi()
	(&sam).SayHi()
}


















 








































