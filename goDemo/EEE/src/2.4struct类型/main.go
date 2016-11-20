package main

import (
	"fmt"
)
//Go语言中 我们可以声明新的类型 作为其它类型的属性或字段的容器
//例如我们自定义类型 person 代表一个人的实体 这个实体拥有的属性 姓名和年龄
//这样的类型 我们称之为 struct 如下代码所示
type person struct{
name string
age int
}
//看到了么 声明一个struct 如此简单 上面的类型 包含有两个字段 name是用来保存用户名称这个属性 age用来保存用户年龄这个属性
//如何使用 struct呢 请看下面的代码
type person struct{
name string
age int
}
var P person //P现在就是person类型的变量了
P.name="Jianling"//赋值给"Jianling"的name属性
P.age=25//赋值 25给变量P的age属性
fmt.Printf("The Person's name is %s", P.name)//访问P的name属性
//除了上面这种P的声明使用外 还有两种声明使用方式
//1 按照顺序提供初始化值
P:=person{"Tom",25}
//2 通过 field:value 的方式初始化 这样的话可以任意顺序
P:=person{age:25,name:"Tom"}
//下面我们看一个完整使用struct的例子
package main

import "fmt"

//声明一个新的类型
type person struct {
	name string
	age  int
}

//比较两个人的年龄,返回年龄大的那个人 并且返回年龄差 struct也是传值的
func Older(p1, p2 person) (person, int) {
	if p1.age > p2.age {
		return p1, p1.age - p2.age
	}
	return p2, p2.age - p1.age

}
func main() {
	var tom person
	//赋值初始化
	tom.name, tom.age = "Tom", 18
	//两个字段都写清楚的初始化
	bob := person{age: 25, name: "Bob"}
	//按照struct 定义顺序初始化值
	paul := person{"paul", 43}
	tb_Older, tb_diff := Older(tom, bob)
	tp_Older, tp_diff := Older(tom, paul)
	//bp_Older, bp_diff := Older(bob, paul)
	fmt.Printf("Of %s and %s,%s is older by %d years.\n", tom.name, bob.name, tb_Older.name, tb_diff)
	fmt.Printf("Of %s and %s,%s is older by %d years.\n", tom.name, paul.name, tp_Older.name, tp_diff)
}
//struct的匿名字段
//我们上面介绍了 如何定义一个struct 定义的时候是字段名与其类型一一对应 实际上 Go支持只提供类型 而不写字段名的方法 也就是匿名字段 也称为嵌入字段
//当匿名字段是一个struct的时候 那么这个struct所拥有的全部字段被隐式的引入了当前定义的这个struct
//让我们来看一个例子 让上面说的这些更加具体化
package main

import (
	"fmt"
)

type Human struct {
	name   string
	age    int
	weight int
}
type Student struct {
	Human      //匿名字段 那么默认Student 就包含了Human的所有字段
	speciality string
}

func main() {
	//我们初始化一个学生
	mark := Student{Human{"Mark", 25, 120}, "Computer Science"}
	//我们访问相应的字段
	fmt.Println("His name is :", mark.name)
	fmt.Println("His age is :", mark.age)
	fmt.Println("His weight is :", mark.weight)
	fmt.Println("His speciality is", mark.speciality)
	//修改对应的备注信息
	mark.speciality = "AI"
	fmt.Println("Mark changed his speciality~")
	fmt.Println("His speciality is ", mark.speciality)
	//修改他的年龄信息
	fmt.Println("Mark become old")
	mark.age = 40
	fmt.Println("his age is ", mark.age)
	//修改他的体重信息
	fmt.Println("Mark is not an athlet anymore")
	mark.weight += 60
	fmt.Println("His weight is ", mark.weight)
}
//我们看到student 访问属性age和name的时候 就像访问自己所拥有的字段一样
//对 匿名字段就是这样，能够实现字段的继承 是不是很酷啊 还有比这个更酷的呢 那就是student还能访问Human这个字段 作为字段名 
mark.Human =Human{"Marcus",55,220}
mark.Human.age-=1
//通过匿名访问和修改字段相当的有用 但是不仅仅是struct字段哦  所有的内置类型和自定义类型都是可以作为匿名字段的,请看下面的例子
package main

import (
	"fmt"
)

type Skills []string
type Human struct {
	name   string
	age    int
	weight int
}
type Student struct {
	Human      //匿名字段 那么默认Student 就包含了Human的所有字段
	Skills     //匿名字段 自定义类型的 string slice
	int        //内置类型作为匿名字段
	speciality string
}

func main() {
	//初始化学生Jane
	jane := Student{Human: Human{"Jane", 35, 100}, speciality: "Biology"}
	//现在我们来访问相应的字段
	fmt.Println("Her name is", jane.name)
	fmt.Println("Her age is", jane.age)
	fmt.Println("Her weight is", jane.weight)
	fmt.Println("Her speciality is", jane.speciality)
	//我们来修改他的 skill技能字段
	jane.Skills = []string{"anatomy"}
	fmt.Println("Her skills are ", jane.Skills)
	fmt.Println("She acquired two new ones")
	jane.Skills = append(jane.Skills, "physics", "golang")
	fmt.Println("Her skills now are ", jane.Skills)
	//修改匿名内置类型字段
	jane.int = 3
	fmt.Println("Her preferred number is ", jane.int)
}
//从上面的例子我们看出来struct 不仅仅能够将struct 作为匿名字段自定义类型 内置类型都可以作为匿名字段
//而且可以在相应的字段上面进行函数操作 如例子中的append
//这里有个问题 如果human里面有一个字段叫做phone student也有一个字段叫做phone 那么该怎么办呢
//Go里面很简单的解决了这个问题 最外层的优先访问 也就是当你通过 student.phone 访问的时候 是访问student里面的字段 而不是human里面的字段
//这样就允许我们去重载通过匿名字段 继承的一些字段 当然 如果 我们想访问重载后对应匿名类型里面的字段可以通过 匿名字段名来访问 请看下面的例子
package main

import "fmt"

type Human struct {
	name  string
	age   int
	phone string //Human类型拥有的字段
}
type Empoyee struct {
	Human      //匿名字段 Human
	speciality string
	phone      string //雇员的phone字段
}

func main() {
	Bob := Empoyee{Human{"Bob", 34, "777-8888-XXXXX"}, "Designer", "333-333"}
	fmt.Println("Bob's work phone is:", Bob.phone)
	//如果我们要访问Human的phone字段
	fmt.Println("Bob's  personal phone is:", Bob.Human.phone)
}




















