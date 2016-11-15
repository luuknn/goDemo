package main

import ("fmt")
//定义一个名为"variableName",类型为"type"的变量
var variableName type
//定义多个变量，定义三个都是type类型的变量
var vname1,vname2,vname3 type
//定义变量并初始化其值
var variableNmae type=value
//同时初始化多个变量(3种)
var vname1,vname2,vname3 type=v1,v2,v3
var vname1,vname2,vname3 =v1,v2,v3
vname1,vname2,vname3 :=v1,v2,v3
//下划线是个特殊变量名,任何赋予它的值都会被丢弃
_,b:= 34,35
//常量-所谓常量，也就是在程序编译阶段就确定下来的值，而程序在运行时则无法改变该值。
//在Go程序中，常量可定义为数值、布尔值或字符串等类型。
const constantname=value
//如果需要,也可以明确指定常量的类型
const Pi float32 =3.1314926
//以下是一些常量声明的例子
const Pi=3.1415926
const i=10000
const MaxThread=10

//内置基础类型
Boolean//Go中，布尔值的类型为bool,值为true或false,默认为false
//示例代码
var isActive bool//全局变量声明
var enabled,disabled=true,false //忽略类型的声明
func test(){
var available bool//一般声明
valid:=false//简短声明
available=true //赋值操作
}

//数值类型
//整数类型有带符号和无符号两种 Go同时支持 int和uint
//rune, int8, int16, int32, int64和byte, uint8, uint16, uint32, uint64。其中rune是int32的别称，byte是uint8的别称。
//Go还支持复数默认类型是complex128(64位实数+64位虚数)cpmplex64
var c complex64=5+5i
fmt.Printf("value is : %v", c)

//字符串
//Go的字符串都是采用UTF-8的字符集编码 字符串是用一对双引号或反引号定义 类型是string
//示例代码
var frenchHello string//声明变量为字符串的一般方法
var emptyString string=""//声明了一个字符串变量,初始化为空字符串
func test(){
no,yes,maybe :="no","yes","maybe"//简短声明,同时声明多个变量
frenchHello="Bonjour"//常规赋值
}
//在go中字符串是不可变的，例如以下的代码编译时会报错
var s string ="hello"
s[0]='c'
//但是如果真的想要修改怎么办 且看下面代码
s:="hello"
c:=[]byte(s)//将字符串s转换为[]byte类型
c[0]='c'
s2:=string(c)
fmt.Printf("%s\n", s2)
//Go中可以使用+操作来连接两个字符串
s:="hello"
m:="world"
a:=s+m
fmt.Printf("%s\n", a)
//修改字符串也可以写为：
s:="hello"
s="c"+s[1:]//字符串虽不能更改,但可以进行切片操作
fmt.Printf("%s\n", a)
//如果要声明一个多行的字符串怎么办 可以通过`来声明
m:=`hello
	world`
//被括的字符串为raw字符串,即字符串在代码中的形式就是打印时的形式没有字符转换,换行也将按原样输出

//错误类型
//Go内置有一个error类型,专门用来处理错误信息,Go的package里面还专门有一个包errors来处理错误
err :=errors.New("emit macho dwarf:elf haeder corrupted")
if err!=nil{
fmt.Print(err)
}

//分组声明
//在Go语言中 同时声明多个变量 常量或者导入多个包时 可采用分组的方式进行声明
//例如以下的代码
import "fmt"
import	"os"
const i=100
const pi=3.1415
var i int
var pi float32 
var prefix string
//可以分组写成如下形式
import(
	"fmt"
	"os"
)
const(
	i=100
	pi=3.1415
)
var(
	i int
	pi float32
	prefix string
)

//iota枚举
//Go里面有个关键字 iota 这个关键字用来声明enum的时候采用,它默认开始值是0,每次调用+1
const(
	x	=	iota//x == 0
	y	=	iota//y == 1
	z	=	iota//z == 2
	w	//常量声明省略时，默认和前面一个字的字面相同w=iota 因此w为3
)
const v	iota//每遇到一个const关键字,itoa就会重置此时v==0

//array.slice.map
//array就是数组,它的定义方式如下
var arr [n]type
//在[n]type中，n表示数组的长度，type表示存储元素的类型。对数组的操作和其它语言类似，都是通过[]来进行读取或赋值：
var arr [10]int //声明了一个int类型的数组
arr[o]=42//数组下标是从0开始的
arr[1]=43//赋值操作
fmt.Printf("The first element is %d\n", arr[0])//获取数据 返回42
fmt.Printf("The last element is %d\n", arr[9])//返回未赋值的最后一个元素，默认返回0
//由于长度也是数组类型的一部分因此[3]int [4]int 是不同的类型 数组也就不能改变长度
//数组之间的赋值是值赋值 即当把一个数组作为参数传入函数的时候,传入的其实是该数组的副本,而不是指针
//同样 数组可以使用另一种:=来声明
a:=[3]int{1,2,3}//声明了一个长度为3的int数组
b:=[10]int{1,2,3}//声明了一个长度为10的int数组前三个元素初始化为1,2,3，其它默认为0
c:=[...]{4,5,6}//可以省略长度。。。的方式,Go会自动根据元素个数来计算长度
//也许你会说 我想数组里面的值还是数组 能实现吗 当然咯Go支持嵌套数组,即多维数组
//声明了一个二维数组,该数组已两个数组作为元素,其中每个数组中又有4个int类型的元素
doubleArray :=[2][4]int{[4]int{1,2,3,4},[4]int{5,6,7,8}}
//如果内部的元素和外部的一样,那么上面的声明可以简化,直接忽略内部的类型
easyArray :=[2][4]int{{1,2,3,4},{5,6,7,8}}
//SLICE
//在很多应用场景中 数组并不能满足我们的需求 在定义数组的时候 我们并不知道需要多大的数组因此我们需要动态数组在Go里面这种数据结构叫做slice
//slice 并不是真正意义上的动态数组 而是一个引用类型 slice总是指向一个底层array slice的声明也和array一样，只是不需要长度
var fslice []int //和声明array一样,只是少了长度
//接下来我们可以声明一个slice 并初始化数据,如下所示
slice :=[]byte{'a','b','c','d'}
//slice可以从一个数组或者一个已经存在的slice中再次声明slice 通过array[i:j]来获取，i是开始位置 j是结束位置 但不包含arr[j]
var ar =[10]byte{'a','b','c','d','e','f','g','h','i','j'}//声明一个含有10个元素类型为byte的数组
var a,b[]byte//声明两个含有byte的slice
a=ar[2:5]//现在a含有的元素 ar[2] ar[3] ar[4]
//slice 有一些简要的操作
var array=[10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
//声明两个slice
var aSlice,bSlice []byte
//演示一些简便操作
aSlice =array[:3]//abc
aSlice=array[5:]//fghij
aSlice =array[:]//包含了全部元素
//从slice中获取slice
aSlice = array[3:7]  // aSlice包含元素: d,e,f,g，len=4，cap=7
bSlice = aSlice[1:3] // bSlice 包含aSlice[1], aSlice[2] 也就是含有: e,f
bSlice = aSlice[:3]  // bSlice 包含 aSlice[0], aSlice[1], aSlice[2] 也就是含有: d,e,f
bSlice = aSlice[0:5] // 对slice的slice可以在cap范围内扩展，此时bSlice包含：d,e,f,g,h
bSlice = aSlice[:]   // bSlice包含所有aSlice的元素: d,e,f,g
//slice 是引用类型所以当改变其中元素的值时，其他的所有引用都会改变该值,aSlice值修改 bSlice中的值也会改变
//从概念上面来说 slice像一个结构体 这个结构体包含了三个元素
//一个指针 长度 最大长度
Array_a := [10]byte{'a', 'b', 'c', 'd', 'e', 'f', 'g', 'h', 'i', 'j'}
Slice_a := Array_a[2:5]
//slice 有几个有用的内置函数
//len获取slice的长度 cap获取slice的最大容量 append像slice里面追加一个或多个元素 然后返回一个和slice一样类型的slice
//copy函数 copy从源slice 的src中复制元素到目标dst 并且返回复制元素的个数
//（即(cap-len) == 0）时，此时将动态分配新的数组空间。返回的slice数组指针将指向这个空间，而原数组的内容将保持不变；其它引用此数组的slice则不受影响。
//MAP
//map[keyType]valueType
//map的读取和设置也和slice一样 通过key来操作 只是slice的index只能是int类型 而map多了很多类型,可以是int 也可以是string以及完全定义了== 和!=操作的类型
var numbers map[string] int 
numbers :=make(map[string] int)
numbers["one"]=1
numbers["ten"]=10
numbers["three"]=3
//使用map需要注意几点
//map是无序的,每次打印出来的map都会不一样
//map的长度是不固定的
//内置的len函数同样适用于map
//map的值可以很方便的修改 通过numbers["one"]=11 可以很容易的把key为one的字典改为11
// 初始化一个字典
rating := map[string]float32 {"C":5, "Go":4.5, "Python":4.5, "C++":2 }
// map有两个返回值，第二个返回值，如果不存在key，那么ok为false，如果存在ok为true

if ok {
    fmt.Println("C# is in the map and its rating is ", csharpRating)
} else {
    fmt.Println("We have no rating associated with C# in the map")
}

delete(rating, "C")  // 删除key为C的元素
//上面说过了 map也是一种引用类型 如果两个map指向同一个底层 那么一个改变另一个也相应改变
m:=make(map[string]string)
m["Hello"]="Bonjour"
m1:=m
m1["Hello"]="Sault"//现在m["Helllo"]的值已经是sault了

//make、new操作
//
//make用于内建类型（map、slice 和channel）的内存分配。new用于各种类型的内存分配。
//
//内建函数new本质上说跟其它语言中的同名函数功能一样：new(T)分配了零值填充的T类型的内存空间，并且返回其地址，即一个*T类型的值。用Go的术语说，它返回了一个指针，指向新分配的类型T的零值。有一点非常重要：
//
//new返回指针。
//内建函数make(T, args)与new(T)有着不同的功能，make只能创建slice、map和channel，并且返回一个有初始值(非零)的T类型，而不是*T。本质来讲，导致这三个类型有所不同的原因是指向数据结构的引用在使用前必须被初始化。例如，一个slice，是一个包含指向数据（内部array）的指针、长度和容量的三项描述符；在这些项目被初始化之前，slice为nil。对于slice、map和channel来说，make初始化了内部的数据结构，填充适当的值。
//
//make返回初始化后的（非零）值。
//下面这个图详细的解释了new和make之间的区别。



























