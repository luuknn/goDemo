package main

import (
	"fmt"
	"reflect"
	"sort"
)

//Go语言里设计最精妙的应该算interface 它让面向对象 内容组织实现非常方便 当你看完这一套 你就会被interface的巧妙设计所折服
//什么是interface
//简单的说 interface是一组method的组合  我们通过interface来定义对象的一组行为
//我们前面一章 最后一个例子中 Student和Employee 都能实现SayHi 虽然他们的内部实现不一样
//但是那不重要 重要的是 他们都能say hi
//让我们来继续做更多的扩展  学生和雇员实现另一个方法 sing  然后学生实现方法借钱  雇员实现方法花钱
//这样 Student实现了三个方法 SayHi Sing BorrowMoney
//Employee 实现了SayHi Sing SpendSalary
//上面这些方法的组合被称为interface(被对象 Student和Employee实现)
//例如 Student和Employee 都实现了interface:SayHi和Sing
//也就是 这两个对象是该interface类型 而employee 没有实现这个interface
//SayHi Sing BorrowMoney 因为employee没有实现BorrowMoney这个方法

//interface类型
//interface 类型定义了一组方法 如果某个对象实现了某个接口的所有方法 则该对象就实现了此接口
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
	loan   float32
}
type Employee struct {
	Human   //匿名字段
	company string
	money   float32
}

//在Human上面定义了一个method
func (h *Human) SayHi() {
	fmt.Printf("Hi ,I am %s ,you can call me on %s\n", h.name, h.phone)
}

//Human对象实现了Sing 方法
func (h *Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

//Human对象实现了Guzzle方法
func (h *Human) Guzzle(beerStein string) {
	fmt.Println("Guzzle Guzzle Guzzle...", beerStein)
}

//Employee 的method 重写 Human的SayHi方法
func (e *Employee) SayHi() {
	fmt.Printf("Hi I am %s,I work at %s.Call me on %s \n", e.name, e.company, e.phone)
}

//Student实现BorroeMoney方法
func (s *Student) BorrowMoney(amount float32) {
	s.loan += amount
}

//Employee实现SpendSalary方法
func (e *Employee) SpendSalary(amount float32) {
	e.money -= amount
}

//定义interface
type Men interface {
	SayHi()
	Sing(lyrics string)
	Guzzle(beerStein string)
}
type YoungChap interface {
	SayHi()
	Sing(song string)
	BorrowMoney(amount float32)
}
type ElderlyGent interface {
	SayHi()
	Sing(song string)
	SpendSalary(amount float32)
}

func main() {
	mark := Student{Human: Human{"Mark", 25, "222-333-YYY"}, school: "MIT"}
	sam := Employee{Human: Human{"sam", 18, "111111-xxx"}, company: "Golang Inc"}
	mark.SayHi()
	(&sam).SayHi()
}
//通过上面的代码 我们可以知道 interface可以被任意的对象实现 我们看到
//上面的Men interface被Human Student Employee 实现
//同理 一个对象可以实现任意多个interface 例如Student实现了Men 和YongChap 两个interface

//interface值
//那么interface里面到底能存什么值呢  如果我们定义一个interface变量,
//那么这个变量里面可以存实现这个interface的任意类型的对象 例如上面例子中
//我们定义了Men interface 类型的变量m 那么m里面可以存 Human Student 或者Employee的值
//因为m能够持有这三种类型的对象 所以我们可以定义一个包含Man类型元素的slice
//这个slice可以被赋予实现了 Men接口的任意结构的对象 这个和我们传统意义上的slice有所不同
//让我们来看一下 下面这个例子
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
	loan   float32
}

type Employee struct {
	Human   //匿名字段
	company string
	money   float32
}

//Human实现Sayhi方法
func (h Human) SayHi() {
	fmt.Printf("Hi, I am %s you can call me on %s\n", h.name, h.phone)
}

//Human实现Sing方法
func (h Human) Sing(lyrics string) {
	fmt.Println("La la la la...", lyrics)
}

//Employee重载Human的SayHi方法
func (e Employee) SayHi() {
	fmt.Printf("Hi, I am %s, I work at %s. Call me on %s\n", e.name,
		e.company, e.phone) //Yes you can split into 2 lines here.
}

// Interface Men被Human,Student和Employee实现
// 因为这三个类型都实现了这两个方法
type Men interface {
	SayHi()
	Sing(lyrics string)
}

func main() {
	mike := Student{Human{"Mike", 25, "222-222-XXX"}, "MIT", 0.00}
	paul := Student{Human{"Paul", 26, "111-222-XXX"}, "Harvard", 100}
	sam := Employee{Human{"Sam", 36, "444-222-XXX"}, "Golang Inc.", 1000}
	Tom := Employee{Human{"Tom", 36, "444-222-XXX"}, "Things Ltd.", 5000}

	//定义Men类型的变量i
	var i Men

	//i能存储Student
	i = mike
	fmt.Println("This is Mike, a Student:")
	i.SayHi()
	i.Sing("November rain")

	//i也能存储Employee
	i = Tom
	fmt.Println("This is Tom, an Employee:")
	i.SayHi()
	i.Sing("Born to be wild")

	//定义了slice Men
	fmt.Println("Let's use a slice of Men and see what happens")
	x := make([]Men, 3)
	//T这三个都是不同类型的元素，但是他们实现了interface同一个接口
	x[0], x[1], x[2] = paul, sam, mike

	for _, value := range x {
		value.SayHi()
	}
}
//通过上面的代码你会发现 interface 就是一组 抽象方法的集合 它必须由其他非interface类型
//而不能自我实现 go通过interface实现了 duck-typing 
//即 当一只鸟走起来像鸭子 游起来像鸭子 叫起来像鸭子 那么这只鸟就可以被称为鸭子

//空interface
//空interface(interface{})不包含任何的method 正因为如此 所有类型都实现了 空interface
//空interface对于描述起不到任何的作用 因为它不包含任何的method
//但 空interface在我们需要存储任意类型的数值的时候相当有用 因为 它可以存储任意类型的数值
//定义a为空接口
var a interface{}
var i int=5
s:="Hello World"
//a可以存储任意类型的数值
a=i
a=s
//一个函数把interface{}作为参数 那么他可以接受任意类型的值 作为参数
//如果一个函数 返回interface{}那么也就可以返回类型的值

//interface 函数参数
//interface 的变量可以持有任意实现该interface类型的对象 这给我们编写函数method
//提供了一些额外的思考 我们是不是可以通过定义interface参数 让函数接受各种类型的参数
//举个例子 fmt.Printlan 是我们常用的一个函数 但是你是否注意到它可以接受任意类型的数据
//打开fmt源码文件 你会看到 这样的一个定义
type Stringer interface{
String() string
}
//也就是说 任何实现了string方法的类型都能作为参数被fmt.println调用 让我们来试一试
package main

import (
	"fmt"
	"strconv"
)

type Human struct {
	name  string
	age   int
	phone string
}

//通过这个方法Human实现了fmt.Stringer
func (h Human) String() string {
	return "<" + h.name + "-" + strconv.Itoa(h.age) + "years -  ✆ " + h.phone + ">"

}
func main() {
	Bob := Human{"Bob", 23, "000-777-2222"}
	fmt.Println("This Human is :", Bob)
}
//我们回顾一下 前面的Box例子 你会发现Color也定义了一个method String其实也实现了fmt.Stringer这个interface
//如果需要某个类型被fmt包以特殊的格式输出  你就必须实现这个Stringer这个接口 如果没有实现这个接口 fmt将以默认的方式输出
//实现同样的功能
fmt.Println("The biggest one is", boxes.BiggestsColor().String())
fmt.Println("The biggest one is", boxes.BiggestsColor())
//注: 实现了error接口的对象  即实现了Error() string的对象 使用fmt输出时 会调用Error方法 因此不必再定义String()方法了

//interface变量存储的类型
//我们知道interface的变量里面可以存储任意类型的数值 该类型实现了interface
//那么我们怎么方向知道这个变量里面实际保存了的是哪个类型的对象呢  目前常用的有两种方法
//Comma-ok 断言
//Go 语言里面有一个语法 可以直接判断 是否是该类型的变量 value,ok =element.(t)
//这里value就是变量的值 ok是一个bool类型 element是interface变量 T是断言的类型
//如果element里面确实存储了T类型的数值 那么ok返回true 否则返回 false
//让我们通过一个例子 来更加深入的理解
package main

import (
	"fmt"
	"strconv"
)

type Element interface{}
type List []Element
type Person struct {
	name string
	age  int
}

//定义了string方法 实现了fmt.stringer
func (p Person) String() string {
	return "(name: " + p.name + " - age: " + strconv.Itoa(p.age) + " years)"
}
func main() {
	list := make(List, 3)
	list[0] = 1       //an int
	list[1] = "Hello" //a string
	list[2] = Person{"Dennis", 70}
	for index, element := range list {
		if value, ok := element.(int); ok {
			fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
		} else if value, ok := element.(string); ok {
			fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
		} else if value, ok := element.(Person); ok {
			fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
		} else {
			fmt.Println("list[%d] is of a different type", index)
		}

	}

}
//是不是很简单啊 同时你是否注意到 多个ifelse里面 还记得我前面介绍流程里面讲过 if里面允许初始化变量
//也许你注意到了 我们断言的类型越多 那么ifelse 也就越多 所以 才引出下面要介绍的switch
//switch 测试
//最好的讲解就是代码例子 现在让我们重写上面的这个实现
package main

import (
	"fmt"
	"strconv"
)

type Element interface{}
type List []Element
type Person struct {
	name string
	age  int
}

//定义了string方法 实现了fmt.stringer
func (p Person) String() string {
	return "(name: " + p.name + " - age: " + strconv.Itoa(p.age) + " years)"
}
func main() {
	list := make(List, 3)
	list[0] = 1       //an int
	list[1] = "Hello" //a string
	list[2] = Person{"Dennis", 70}
	for index, element := range list {
		switch value := element.(type) {
		case int:
			fmt.Printf("list[%d] is an int and its value is %d\n", index, value)
		case string:
			fmt.Printf("list[%d] is a string and its value is %s\n", index, value)
		case Person:
			fmt.Printf("list[%d] is a Person and its value is %s\n", index, value)
		default:
			fmt.Println("list[%d] is of a different type", index)
		}
	}
}
//这里有一点需要强调的是 element.(type) 语法只能在switch测试的时候进行使用 如果你要在
//switch外面判断一个类型就要使用comma-ok

//嵌入interface
//Go里面真正吸引人的地方是他的内置逻辑语法 就像我们学习struct时学习的匿名字段
//多么优雅啊 那么相同的逻辑引入到interface里面 那不是更加完美了  如果一个interface1作为
//interface2的一个嵌入字段 那么interface2 隐式的包含了interface1里面的method
//我们可以看到源码包 container/heap 里面有这样的一个定义
type Interface interface{
sort.Interface//嵌入字段sort.Interface
Push(x interface{})//a push method to push elements into the heap
Pop() interface{}//a Pop elements that pops elements from the heap
}
//我们看到sort.interface其实就是嵌入字段 把sort.interface所有的method都包含进来了 也就是下面的三个方法
type Interface interface{
//Len is the number of elements in the collection
Len() int
//less returns whether the element with index i should  sort
//before the element with index i
Less(i,j int)bool
//Swap swaps the elements with indexs i j
Swap(i,j int)
}
//另外一个例子 就是io包下面的io.ReadWriter 他包含了io包下面的reader和write两个interface
//io.ReadWriter
type ReadWriter interface{
Reader
Writer
}

//反射
//Go语言实现了反射 所谓反射 就是动态运行时的状态 我们一般用到的包是reflect包
//如何运用reflect包 官方的这篇文章详细的讲解了reflect的实现原理
//使用 reflect一般分为三步 下面 简要的讲解一下  要去反射一个类型的值(这些值都实现了空interface)
//首先需要把它转化成reflect对象  (reflect.Type 或者reflect.Value)根据不同的情况调用不同的函数
//这两种获取方式如下
t:=reflect.TypeOf(i)//得到类型的元数据 通过t 我们呢获取类型定义里面所有元素
v:=reflect.ValueOf(i)//得到实际的值  通过v我们获取存储在里面的值 还可以去改变值
//转化为reflect对象之后 我们就可以进行一些操作了  也就是将 reflect对象转化成相应的值 例如
tag:=t.Elem().Field(0).Tag//获取定义在struct里面的标签
name:=v.Elem().Field(0).String()//获取存储在第一个字段里面的值
//获取反射值能返回相应的类型和数值

package main

import (
	"fmt"
	"reflect"
)

func main() {
	var x float64 = 3.4
	v := reflect.ValueOf(x)
	fmt.Println("Type:", v.Type())
	fmt.Println("kind is float64:", v.Kind() == reflect.Float64)
	fmt.Println("value:", v.Float())

}
//最后 反射的话 那么反射的字段是可以修改的 我们前面学过传值和传引用
//这里面也是一样的道理 反射的字段必须是可读写的意思
//错误例子
var x float64 =3.4
v:=reflect.ValueOf(x)
v.SetFloat(7,1)
//如果要修改相应的值,必须这样写
var x float64=3.4
p:=reflect.ValueOf(&x)
v:=p.Elem()
v.SetFloat(7.1)
//上面只是对反射的简单介绍 更深入理解还需要自己在编程中的不断实践
//Go 反射的规则
//1 反射的规则
//在运行时反射是程序检查其所拥有的结构 尤其是类型的一种能力 这是元编程的一种形式
//它同时也是造成混淆的重要来源,//接下来的 反射一词  是表示 在Go中的反射

//2 类型Types和接口interfaces
//由于反射是建于 类型系统之上的 就从复习一下Go中的类型开始吧
//Go是静态类型的 每一个变量有一个静态的类型 也就是说 有一个已知类型并且在编译的时候就确定下来
//int float32 MyType []byte 等等 如果定义
type MyInt int
var i int
var j MyInt
//那么i的类型是int j的类型是MyInt 即使变量i和j有相同的底层类型 它们仍然是有不同的静态类型的
//未经转换是不能相互直接赋值的
//在类型中有一个重要的类别就是接口类型 表达了固定的一个方法集合 
//一个接口变量可以存储任意实际值(非接口) 只要这个值 实现了接口的方法
//众所周知 的一个例子就是is io.Reader 和io.Writer 来自 io包的类型Reader和Writer
type Reader interface{
Read(p []byte)(n int,err error)
}
type Writer interface{
Write(p []byte)(n int,err error)
}
//任何用这个申明实现了read或者weite方法的类型 可以说它实现了io.Reader(或io.Write)
//基本讨论来说 这意味着io.reader类型的变量可以保存任意值 只要这个值的类型实现了read方法
var r io.Reader
r=os.Stdin
r=bufio.NewReader(r)
r=new(bytes.Buffer)
//and so on
//有一个事情是一定要明确的  不论r保存了什么值 r的类型 总是io.reader Go是静态类型
//r的静态类型是io.reader
//接口类型的一个极端重要的例子是空接口 interface{}
//它表示空的方法集合 由于任何值都有另个或者多个方法 所以任何值都可以满足它

//3 接口的表现representation
//接口类型的变量存储了两个内容  赋值给变量实际的值和这个值的类型描述 更准确的说
//值是底层实现了接口的数据项目 而类型描述了这个项目完整的类型 例如下面
var r io.Reader 
tty,err:=os.OpenFile("/dev/tty",os.O_RDWR,0)
if err !=nil{
return nil,err
}
r=tty
//r包含的是(value ,type)对 如，如 (tty, os.File)。注意类型 os.File 除了 Read 方法还实现了其他方法：尽管接口值仅仅提供了访问 Read* 方法的可能，但是内部包含了这个值的完整的类型信息。这也就是为什么可以这样做：
var w io.Writer
w=r.(io.Writer)
//在这个赋值中表达式是一个类型断言 它断言r内部项同时实现了io.Writer 因此可以赋值它到w
//在赋值之后，w 将会包含 (tty, os.File)。跟在 r* 中保存的一致。接口的静态类型决定了哪个方法可以通过接口变量调用，即便内部实际的值可能有一个更大的方法集。
//接下来 可以这样做
var empty interface{}
empty=w
//而空接口值e 也将同样包含(tty,os.File)这很方便 空接口可以保存任何值 同时保留关于那个值的所有信息
//(这里无需类型断言 因为w肯定满足空接口在这个例子中，将一个值从 Reader 变为 Writer，由于 Writer 的方法不是 Reader 的子集，所以就必须明确使用类型断言。)
//一个很重要的细节是 接口内部的对 总是 value 实际值 的格式  而不会有 value 接口类型的格式 接口不能保存接口值
//现在准备好 来进行反射了

//4 反射的第一条规则
//从接口值到反射对象的反射
//在基本层面上 反射只是一个检查存储在接口变量中的类型和值的 算法。
//从头说起 reflect包中有两个类型需要了解 type和value

































