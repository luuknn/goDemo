package main

//import (
//	"encoding/json"
//	"fmt"
//	"log"
//)
//
//type Account struct {
//	Email    string  `json:"email"`
//	Password string  `json:"-"`
//	Money    float64 `json:"money,string"`
//}
//type User struct {
//	Name    string
//	Age     string
//	Roles   []string
//	Skill   map[string]float64
//	Account Account
//	Extra   []interface{}
//	Level   map[string]interface{}
//}
//
//func main() {
//	level := make(map[string]interface{})
//	level["web"] = "Good"
//	level["server"] = 90
//	level["tool"] = nil
//	extra := []interface{}{123, "hello"}
//	account := Account{
//		Email:    "jianling.shih@gmail.com",
//		Password: "111111",
//		Money:    100.5,
//	}
//	skill := make(map[string]float64)
//	skill["java"] = 88
//	skill["js"] = 80.0
//	skill["shell"] = 80
//	user := User{
//		Name:    "JL",
//		Age:     "25",
//		Roles:   []string{"Owner", "Master"},
//		Skill:   skill,
//		Account: account,
//		Extra:   extra,
//		Level:   level,
//	}
//	rs, err := json.Marshal(user)
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(rs)
//	fmt.Println(string(rs))
//}

/*
JSON
http交互生命周期 包含请求和响应 前面我们介绍了很多关于发起 请求 处理请求的内容 现在该聊一聊返回的响应内容了 对于web服务的响应
以前常见的响应是返回服务器渲染的模板  浏览器只要展示模板即可  随着Restful风格的api出现 已经前后端分离 更多的
返回格式是json字串 本节 我们将讨论在golang中如何编码和解码json

JSON是一种数据格式描述语言 以key value构成的哈希结构  类似JavaScript中的对象 Python中的字典 通常json格式的key是字串 其值可以是
任意类型 字串 数字 数组 或者对象结构

数据结构map
json源于JavaScript的对象结构 golang中直接对应的数据结构 可是golang的map也是key value结构 同时struct结构体也可以描述json
当然对于json数据类型 go也会有对象的结构所匹配
数据类型 JSON GOLANG
字串 string   string
整数 number   int64
浮点数    number  float64
数组  array   slice
对象 object  struct
布尔 bool    bool
空值  null  nil
*/
/*
基本结构编码
golang提供了encoding/json的标准库 用于编码json  大致需要两步
首先定义json结构体  使用Marshal方法序列化
定义结构体的时候  只有字段名是大写的 才会被编码到json中
import (
	"encoding/json"
	"fmt"
	"log"
)

type Account struct {
	Email    string
	password string
	Money    float64
}

func main() {
	account := Account{
		Email:    "jianling.shih@gmail.com",
		password: "111111",
		Money:    100.5,
	}
	rs, err := json.Marshal(account)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rs)
	fmt.Println(string(rs))
}
Marshal 方法接受一个空接口的参数 返回一个[]byte 结构 小写命名的password字段没有被编码到json中 生成的json结构字段和Account结构一致
相比字串 数字等基本数据结构 slice切片 map图 是复合结构 这些结构的编码也类似
import (
	"encoding/json"
	"fmt"
	"log"
)

type User struct {
	Name  string
	Age   string
	Roles []string
	Skill map[string]float64
}

func main() {
	skill := make(map[string]float64)
	skill["java"] = 88
	skill["js"] = 80.0
	skill["shell"] = 80
	user := User{
		Name:  "JL",
		Age:   "25",
		Roles: []string{"Owner", "Master"},
		Skill: skill,
	}
	rs, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rs)
	fmt.Println(string(rs))
}
嵌套编码 slice和map可以匹配json的数组和对象  当前提是对象的value是同类型的情况  更通常的做法 对象的key可以是string
但是其值可以是多种结构 golang通过定义结构体来实现这种构造
import (
	"encoding/json"
	"fmt"
	"log"
)

type Account struct {
	Email    string
	password string
	Money    float64
}
type User struct {
	Name    string
	Age     string
	Roles   []string
	Skill   map[string]float64
	Account Account
}

func main() {
	account := Account{
		Email:    "jianling.shih@gmail.com",
		password: "111111",
		Money:    100.5,
	}
	skill := make(map[string]float64)
	skill["java"] = 88
	skill["js"] = 80.0
	skill["shell"] = 80
	user := User{
		Name:    "JL",
		Age:     "25",
		Roles:   []string{"Owner", "Master"},
		Skill:   skill,
		Account: account,
	}
	rs, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rs)
	fmt.Println(string(rs))
}
*/
/*
通过定义嵌套的结构体Account 实现了key与value不一样的结构 golang的数组或切片 其类型也是一样的 如果遇到不同数据类型的数组
则需要借助空结构来实现
import (
	"encoding/json"
	"fmt"
	"log"
)

type Account struct {
	Email    string
	password string
	Money    float64
}
type User struct {
	Name    string
	Age     string
	Roles   []string
	Skill   map[string]float64
	Account Account
	Extra   []interface{}
}

func main() {
	extra := []interface{}{123, "hello"}
	account := Account{
		Email:    "jianling.shih@gmail.com",
		password: "111111",
		Money:    100.5,
	}
	skill := make(map[string]float64)
	skill["java"] = 88
	skill["js"] = 80.0
	skill["shell"] = 80
	user := User{
		Name:    "JL",
		Age:     "25",
		Roles:   []string{"Owner", "Master"},
		Skill:   skill,
		Account: account,
		Extra:   extra,
	}
	rs, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rs)
	fmt.Println(string(rs))
}
*/
/*
使用空接口 也可以定义像结构体实现那种不同value类型的字典结构
当空接口 没有初始化其值的时候 零值是nil 编码成json就是null
import (
	"encoding/json"
	"fmt"
	"log"
)

type Account struct {
	Email    string
	password string
	Money    float64
}
type User struct {
	Name    string
	Age     string
	Roles   []string
	Skill   map[string]float64
	Account Account
	Extra   []interface{}
	Level   map[string]interface{}
}

func main() {
	level := make(map[string]interface{})
	level["web"] = "Good"
	level["server"] = 90
	level["tool"] = nil
	extra := []interface{}{123, "hello"}
	account := Account{
		Email:    "jianling.shih@gmail.com",
		password: "111111",
		Money:    100.5,
	}
	skill := make(map[string]float64)
	skill["java"] = 88
	skill["js"] = 80.0
	skill["shell"] = 80
	user := User{
		Name:    "JL",
		Age:     "25",
		Roles:   []string{"Owner", "Master"},
		Skill:   skill,
		Account: account,
		Extra:   extra,
		Level:   level,
	}
	rs, err := json.Marshal(user)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rs)
	fmt.Println(string(rs))
}
可以看到Extra 返回的并不是一个空的切片 而是null 同时 Level字段实现了先字典的嵌套结构
*/
/*
Struct tag 字段重名   通过上面的例子 我们看到了level字段中key server等是字母小写
其他的都是大写 因为我们在定义结构的时候 只有使用大写字母开头的字段才会被导出 而通常json世界中
更盛行小写字母的方式 看起来就成了一个矛盾 其实不然 golang 提供了 struct tag的方式可以重命名结构字段的输出形式
type Account struct {
	Email    string `json:"email"`
	password string
	Money    float64 `json:"money"`
}
我们使用struct taf重新给Account结构的字段进行了重命名 其中email和money小写了  运行输出正常
忽略字段 -
重命名是一个利器  这个利器还提供了更高级的选项 通常使用marshal的时候 会把结构体的所有除了私有字段都编码到json
中 而实际开发中 我们定义的结构可能更通用 我们需要某个字段可以导出 但是又不能编码到json中
此时使用  struct tag `-`符号 就能完美解决   我们已经知道`_`常用于忽略字段的占位，在tag中则使用短横线`-`
type Account struct {
	Email    string  `json:"email"`
	Password string  `json:"-"`
	Money    float64 `json:"money"`
}
*/
/*
可见即使Password不是私有字段，因为`-`忽略了它，因此没有被编码到json输出。
##### `omitempty`可选字段
对于另外一种字段，当其有值的时候就输出，而没有值(零值)的时候就不输出，则可以使用另外一种选项`omitempty`。
type Account struct {
Email string json:"email"
Password string json:"password,omitempty"
Money float64 json:"money"
}
此时password不会被编码到json输出中。
##### `string`选项
golang是静态类型语言，对于类型定义的是不能动态修改。在json处理当中，struct tag的string可以起到部分动态类型的效果。有时候输出的json希望是数字的字符串，
而定义的字段是数字类型，那么就可以使用string选项。
type Account struct {
	Email    string  `json:"email"`
	Password string  `json:"-"`
	Money    float64 `json:"money,string"`
}
*/
/*
总结
上面所介绍的大致覆盖了golang的json编码处理  总体原则 分两步  首先定义需要编码的结构
然后调用encoding/json标准库的Marshal方法 生成json byte 数组 转换成string类型即可

golang和json的大部分数据结构匹配 对于复合结构 go可以借助结构体和空接口实现json的数组和对象结构
通过struct tag可以灵活的修改json编码字段名和输出控制

既然有JSON的编码 当然会有json的解码 解析json对于golang则需要更多的技巧。
*/
