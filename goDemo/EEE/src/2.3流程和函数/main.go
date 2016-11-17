package main
import (
"fmt"
)
//流程和函数
//流程控制
//流程控制包含分为三大类 条件判断  循环控制 无条件跳转
//if 不需要括号
if x>10{
fmt.Println("x is greater than 10")
}else{
fmt.Println("x is less than 10")
}
//Go的if还有个强大的地方就是条件判断语句里允许声明一个变量 这个变量的作用域只能在该条件逻辑快内 其他地方不起作用
//计算获取x 然后根据x返回的大小 判断是否大于10
if x:=computeValue();x>10{
fmt.Println("x is greater than 10")
}else {
fmt.Println("x is less than 10")
}
//这个地方如果这样调用就编译出错了 因为x是条件里面的变量
fmt.Println(x)
//多个条件的时候 如下所示
if integer == 3{
fmt.Println("The integer is equal to 3")
}else if integer <3{
fmt.Println("The intager is less than 3")
}else {
fmt.Println("The integer is greater than 3")
}
//goto
//Go有goto语句 请明智的使用它 用goto跳转到必须在当前函数内定义的标签
func myFunc(){
i:=0
Here://这行的第一个词 以冒号结束作为标签
 println(i)
 i++
 goto here//跳转到Here去
}
f//for
//Go里面最强大的一个控制逻辑就是for 它既可以用来循环读取数据 又可以当作while来控制逻辑还能迭代操作
for expression1;expression2;expression3{
//...
}
//ex1 ex3是变量声明或者函数调用返回值之类的 ex2是用来条件判断的 ex1循环开始前调用 ex3 在每轮循环结束调用
package main
import "main"
func main(){
sum:=0;
for index:=0;index<10;index++{
sum+=index
}
fmt.Println("Sum is equal to",sum)
}//输出 sum is equal to 45
//有些时候需要进行多个赋值操作，由于Go里面没有,操作符，那么可以使用平行赋值i, j = i+1, j-1
//有些时候如果我们忽略expression1和expression3：
sum:=1
for ;sum<1000;{
sum+=sum
}
//其中;也可以省略 那么就变成如下的代码了  是不是似曾相识 对 这就是while的功能
sum:=1
for sum<1000{
sum+=sum
}
//在循环里面有两个关键操作break和continue ,break操作是跳出当前循环，continue是跳过本次循环。当嵌套过深的时候，break可以配合标签使用，即跳转至标签所指定的位置，详细参考如下例子：
for　index:=10;index>0;index--{
if index==5{
break//或者continue
}
fmt.Println(index)
}
//break打印出来10 9 8 7 6
//continue 打印出来10 9 8 7 6 5 4 3 2 1
//break和continue还可以跟着标号，用来跳到多重循环中的外层循环
//for配合range 可以用于读取slice 和map的数据
for k,v:=range map{
fmt.Println("map's key:",k)
fmt.Println("map's val",v)
}
//由于Go支持多值返回 而对于声明而未被调用的变量 编译器会报错,在这种情况下,可以使用_来丢弃不需要的返回值
for _,v :=range map{
fmt.Println("map's val:",v)
}
//switch
//sExpr和expr1、expr2、expr3的类型必须一致。Go的switch非常灵活，表达式不必是常量或整数，执行的过程从上至下，直到找到匹配项；而如果switch没有表达式，它会匹配true。
i  :=10
switch i{
case 1:
fmt.Println("i is equal to 1")
case 2,3,4:
fmt.Println("i is equal to 2 ,3,or4")
case 10:
fmt.Println("i is equal to 10")
default:
fmt.Println("ALL I know is that i is an integer")
}
//switch 默认相当于 每个case后面带有break 匹配成功后不会向下执行其他case而是跳出整个switch
//但是fallthrough 强制执行后面的case代码
integer :=6
switch integer{
case 4:
fmt.Println("The integer was <=4")
fallthrough
case 5:
fmt.Println("The integer was <=5")
fallthrough
default:
fmt.Println("default case")
}
//函数 
//函数是Go里面的核心设计,它通过关键字func来声明 它的格式如下
func funcName(input1 type1,inpute2 type2)(output1 type1,output2 type2){
//这里是处理逻辑代码
//返回多个值
return value1,value2
}
//从上面的代码 我们看出
//关键字func用来声明一个函数funcName
//函数可以有一个或者多个参数 每个参数后面带有类型 通过,分隔
//函数可以返回多个值
//上面返回值声明了两个变量output1 output2 如果你不想声明也可以 直接就两个类型
//如果只有一个返回值且不声明返回值变量 那么你可以省略包括返回值的括号
//如果没有返回值 那么就直接省略最后的返回信息
//如果有返回值 那么必须在函数的外层添加return语句
//计算Max的值
package main
import("fmt")
//返回a，b中的最大值
func max(a,b int)int{
if a>b{
return a
}
return b
}
func main(){
x:=3
y:=4
z:=5
max_xy:=max(x,y)//调用函数max(x,y)
max_xz:=max(x,z)//调用函数max(x,z)
fmt.Printf("max(%d,%d)=%d\n", x,y,max_xy)
fmt.Printf("max(%d,%d)=%d\n", y,z,max(y,z))//也可以在这直接调用它
}
//多个返回值
Go语言比c更先进的特性 其中一点就是函数能够返回多个值
package main
import "fmt"
//返回A+B 和A*B
func SumAndProduct(A,B int)(int,int){
return A+B,A*B
}
func main(){
x:=3
y:=4
xPLUSy,xTIMESy:=SumAndProduct(x,y)
fmt.Printf("%d+%d=%d\n", x,y,xPLUSy)
fmt.Printf("%d*%d=%d\n", x,y,xTIMESy)
}
//最好命名返回值 增强可读性
func SumAndProduct(A,B int) (add int,Multipled int){
add =A+B
Multiplied =A*B
return
}
//变参
//Go函数支持变参 接受变参的函数是有着不定数量的参数的 为了做到这点 首先需要定义函数 使其接受变参
func myfunc (arg ...int){}
//arg ...int 告诉Go这个函数接受不定数量的参数 注意 这些参数的类型全部是int
//在参数中 变量arg 是一个int的slice
for _,n:=range arg{
fmt.Printf("And the number is:%d\n", n)
}
//传值与传指针
//当我们传一个参数值到被调用的函数里面 实际上是传了这个值得一份copy当在被调用函数中修改参数值的时候，调用函数中相应的实参 不会发生任何变化 因为数值变化只作用在copy上
//举例
package main
import "fmt"
//简单的一个函数 实现了参数+1的操作
func add1(a int) int{
a=a+1//我们改变了a的值
return a
}
func main(){
x:=3
fmt.Println("x= ",x)
x1 :=add1(x)//调用add1()
fmt.Println("x+1= ",x1)//应该输出 x+1=4
fmt.Println("x=",x)//应该输出x=3
}
//我们知道，变量在内存中是存放于一定地址上的，修改变量实际是修改变量地址处的内存。只有add1函数知道x变量所在的地址，才能修改x变量的值。所以我们需要将x所在地址&x传入函数，并将函数的参数的类型由int改为*int，即改为指针类型，才能在函数中修改x变量的值。此时参数仍然是按copy传递的，只是copy的是一个指针。请看下面的例子
package main
import "fmt"
//简单的一个函数 实现了参数+1的操作
func add1(a *int) int{
*a=*a+1//我们改变了a的值
return *a
}
func main(){
x:=3
fmt.Println("x= ",x)
x1 :=add1(&x)//调用add1()
fmt.Println("x+1= ",x1)//应该输出 x+1=4
fmt.Println("x=",x)//应该输出x=4
}
//这样我们就达到了修改x的目的 那么传指针到底有什么好处呢
//传指针使得多个函数能操作同一个对象
//传指针比较轻量级 8bytes 只是传内存地址 我们可以用指针传递体积大的结构体 如果用参数值传递的话在每次copy上面就会花费相对较多的系统开销(内存和时间)所以当你要传递大的结构体的时候 用指针是一个明智的选择
//Go语言中string slice map这三种类型的实现机制类似指针 所以可以直接传递 而不用取地址后传递指针
//注: 若函数 需改变slice的长度 则 仍需要取地址传递指针

//defer
//Go语言中有不错的设计,即 延迟语句defer 你可以在函数中添加多个defer语句
//当函数执行到最后时 这些defer语句会按照逆序执行 最后该函数返回 特别是 当你在进行一些打开资源的操作时 遇到错误需要提前返回 在返回前你需要关闭相应的资源 不然很容易造成资源泄露等问题
//举例 打开一个资源
func ReadWrite() bool{
file.Open("file")
//做一些工作
if failureX{
file.close()
return false
}
if failureY{
file.close()
return false
}
file.Close()
return true
}
//我们看到上面有很多重复的代码 Go的defer有效解决了这个问题 使用它以后 不但代码量减少了很多 而且程序变得很优雅 在defer后指定的函数会在函数退出前调用
func ReadWrite() bool{
file.Open("file")
defer file.Close()
if failureX{
return false
}
if failureY{
return false
}
return true
}
//如果有很多调用defer 那么defer是采用后进先出的模式 
for i:=0;i<5;i++{
defer fmt.Printf("%d", i)
}//以上代码会输出  4 3 2 1 0

//函数作为值、类型
//在Go中函数 也是一种变量 我们可以通过type来定义它 它的类型就是所有拥有相同的参数，相同的返回值的一种类型
type typeName func (input1 inputType1,input2,inputType2[,...]) (result1 resultType1 [,...]) 
//那么函数作为类型到底有什么好处呢  那就是可以把这个类型的函数当作值来传递 请看下面的例子
package main
import "fmt"
type testInt func(int) bool//声明了一个函数类型
func isOdd(integer int) bool{
if integer%2 ==0{
return false
}
return true
}
func isEven(integer int) bool{
if integer%2 == 0{
return true
}
return false
}
//声明的函数类型在这个地方当了一个参数
func filter(slice []int, f testInt) []int {
    var result []int
    for _, value := range slice {
        if f(value) {
            result = append(result, value)
        }
    }
    return result
}
func main(){
slice :=[]int {1,2,3,4,5,7}
fmt.Println("slice =",slice)
odd :=filter(slice,isOdd)//函数当值来传递
fmt.Println("Odd elements of slice are:",odd)
even:=filter(slice,isEven)//函数但值来进行传递
fmt.Println("Even elements of slice are :",even)
}
//函数当作值和类型在我们写一些通用接口的时候非常有用 通过上面的例子我们看到testInt这个类型是一个函数类型 然后两个filter函数的参数和返回值与testInt类型是一样的 但是我们可以实现很多种逻辑 这样使得我们的程序变得非常灵活

//panic和recover
//Go没有像java那样的异常机制 它不能抛出异常 而是使用了panic和recover机制 一定要记住 你应该把它当做最后的手段来使用 也就是说 你的代码中应当没有或者很少有panic的东西 这是个强大的工具 请明智的使用它 那么我们该如何使用它
//panic
//是一个内建函数 可以中断原有的控制流程 进入一个令人恐慌的流程中 
//当函数F调用panic 函数F的执行被中断 但是F的延迟函数会正常执行 然后F返回到调用它的地方
//在调用的地方 F的行为就像调用了panic 这一过程继续向上 直到发生panic和goroutine
//中所有调用的函数返回 此时 程序退出 恐慌可以直接调用panic产生 也可以运行时错误产生 例如访问越界的数组
//recover
//是一个内建函数 可以让进入恐慌的流程中的goroutine恢复过来
//recover 仅仅在延迟函数中有效 在正常的执行过程中 调用recover会返回nil
//没有其他任何效果 如果当前的goroutine陷入恐慌 调用recover可以捕获到panic的输入值并且恢复正常的执行
//下面这个函数演示了如何在过程中使用panic
var user =os.Getenv("USER")
func init(){
if user ==""{
panic ("no value for $USER")
}
}
//下面这个函数检查作为其参数的函数在执行时是否会产生panic
func throwsPanic(f func()) (b bool){
defer func(){
if x:=recover();x!=nil{
b=true
}
}()
f()//执行函数f 如果d中出现了panic 那么就可以恢复过来
return
}

//main函数和init函数
//Go里面有两个保留函数 init函数 能够应用于所有的package和main函数（只能应用于packagemain）
//这两种函数在定义时 不能有任何的参数和返回值 虽然一个package 里面可以写任意多个init函数 但
//无论是对于可读性还是维护性来说 我们都强烈建议用户在一个package的 每个文件中只写一个init函数
//Go程序 会自动调用 init 和main 
//程序的初始化和执行都起始于main包。如果main包还导入了其它的包，那么就会在编译时将它们依次导入。有时一个包会被多个包同时导入，那么它只会被导入一次（例如很多包可能都会用到fmt包，但它只会被导入一次，因为没有必要导入多次）。当一个包被导入时，如果该包还导入了其它的包，那么会先将其它包导入进来，然后再对这些包中的包级常量和变量进行初始化，接着执行init函数（如果有的话），依次类推。等所有被导入的包都加载完毕了，就会开始对main包中的包级常量和变量进行初始化，然后执行main包中的init函数（如果存在的话），最后执行main函数。

//import
//我们在写Go代码的时候 经常使用到import这个命令来导入包文件 而我们经常看到的方式 参考如下
import("fmt")
//然后我们代码里面可以通过如下的方式进行调用
fmt.Println("hello world")
//上面这个fmt是Go语言的标准库 其实是goroot下去加载该模块
//Go还支持如下两种方式来加载自己写的模块
//1 相对路径 import "./model" //当前文件同一目录的model目录 但是不建议这种方式来import
//2 绝对路径 import "shorturl/model" //加载gopath/src/shorturl/model模块
//上面展示了一些import常用的几种方式 但还有一些特殊的import 让很多新手很费解
//1 点操作
import(
. "fmt"
)//点操作的含义就是 这个包导入之后 调用这个包的函数时可以省略前缀的包名
//fmt.Println("hello world")可以省略的写成 Println("hello world")
//2.别名操作 别名操作顾名思义 我们可以把包名命名成另一个我们用起来很容易记忆的名字
import(
f "fmt"
)
//别名操作的话调用包函数时 前缀变成了我们的前缀 即f.Println("hello world")
//3 操作
//这个操作经常是很多人费解的一个操作符 请看下面这个import
import(
"database/sql"
	_"github.com/ziutek/mymysql/godrv"
)
//_操作其实是引入该包 而不是直接使用包里面的函数 而是调用了该包的init函数














































