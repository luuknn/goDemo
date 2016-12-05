package main

import (
	"fmt"
)

//JSON (Javascript Object Notation) 是一种轻量级 的数据交换语言 以文字为基础 具有自我描述且易于让人阅读
//尽管JSON是JavaScript的一个子集 但json是独立于语言的文本格式 采用了类似于C语言家族的一些习惯
//JSON与xml最大的不同在于xml是一个完整的标记语言  而json不是
//json由于比xml更小 更快 更易解析 以及 浏览器的内建快速解析支持 使得其更适用于网络数据传输领域
//目前 我们看到很多开放平台  基本上都是采用了JSON作为他们数据交互的接口
//既然json在web开发中如此重要  那么 Go语言对json的支持怎么样呢  Go语言的标准库已经非常好的支持了JSON
//很容易对JSON数据进行 编解码工作
//前一小节的运维例子 用json来表示 结果描述如下
//{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}
//本小节余下内容将以此JSON数据为基础 来介绍go语言的json包对json数据的编解码
//解析json 解析到结构体
//假如有了上面的json串 那么 我们如何来解析这个json串呢  Go的json包中有如下函数
//func Unmarshall(data []byte,v interface{}) error
//通过这个函数我们就可以实现解析的目的  详细的解析例子请看如下代码
package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string
	ServerIP   string
}
type Serverslice struct {
	Servers []Server
}

func main() {
	var s Serverslice
	str := `{"servers":[{"serverName":"Shanghai_VPN","serverIP":"127.0.0.1"},{"serverName":"Beijing_VPN","serverIP":"127.0.0.2"}]}`
	json.Unmarshal([]byte(str), &s)
	fmt.Println(s)

}
//在上面的示例代码中 我们首先定义了与json数据对应的结构体 数组对应slice
//字段名对应JSON里面的KEY 在解析的时候 如何将json数据与struct字段相匹配呢 例如JSON的key是Foo那么 怎么找到对应的字段呢
//首先查找tag含有 Foo的可导出的struct字段（首字母大写）
//其次查找字段名是Foo的导出字段
//最后查找类似FOO或者FoO这样除了首字母之外其他大小写不敏感的导出字段
//聪明的你一定注意到了这一点 能够被赋值的字段必须是可导出字段 即首字母大写
//同时json解析的时候 只会解析能找到的字段 找不到的字段会被忽略  这样的一个好处 是
//当你接收到一个很大的JSON数据结构而你却只想获取其中的部分数据的时候 你只需要将你想要的数据对应的字段名大写 即可轻松解决问题


//解析到interface
//上面那种解析方式是我们知晓被解析的JSON数据的结构的前提下 采取的方案 如果我们不知道被解析的数据的数据格式4
//又应该如何来解析呢
//我们知道interface{} 可以用来存储任意数据类型的对象 这种数据结构 正好用于存储解析的未知结构的json数据的结果
//json包中采用 map[string]interface{}和[]interface{}结构来存储任意的JSON对象和数组
//Go类型和JSON类型的对应关系如下 
//bool代表 JSON booleans flost64代表JSON numbers string 代表JSON strings nil代表JSON null
//现在我们假设有如下的数据
b := []byte(`{"Name":"Wednesday","Age":6,"Parents":["Gomez","Morticia"]}`)
//如果我们在不知道他结构的情况下  我们把它解析到interface{}里面
var f interface{}
err:=json.Unmarshall(b,&f)
//这个时候f里面存储了一个map类似 他们的key是string 值存储在空的interface{}里面
f =map[string] interface{}{
 "Name": "Wednesday",
    "Age":  6,
    "Parents": []interface{}{
        "Gomez",
        "Morticia",
    },
}
//那么如何来访问这些数据呢 通过断言的方式
m:=f.(map[string]interface{})
//通过断言之后 你就可以 通过如下方式来访问里面的数据了
for k,v :range m{
switch vv:=v.(type){
case string :
fmt.Println(k,"is string",vv)
 case int:
        fmt.Println(k, "is int", vv)
    case float64:
        fmt.Println(k,"is float64",vv)
    case []interface{}:
        fmt.Println(k, "is an array:")
        for i, u := range vv {
            fmt.Println(i, u)
        }
    default:
        fmt.Println(k, "is of a type I don't know how to handle")

}
}
//通过上面的示例  可以看到 通过interface{}与tape assert的配合 我们就可以解析未知结构的JSON数了
//上面这个是 官方提供的解决方案  其实我们很多时候 通过类型断言 操作起来 不是很方便 目前bitly公司
//开源了一个叫做 simplejson的包 在处理未知结构体的json时相当方便 详细例子如下所示
js, err := NewJson([]byte(`{
    "test": {
        "array": [1, "2", 3],
        "int": 10,
        "float": 5.150,
        "bignum": 9223372036854775807,
        "string": "simplejson",
        "bool": true
    }
}`))

arr, _ := js.Get("test").Get("array").Array()
i, _ := js.Get("test").Get("int").Int()
ms := js.Get("test").Get("string").MustString()

//生成JSON
//我们开发很多应用的时候 最后都是要输出json数据串 那么如何来处理呢  json包里面通过Marshal函数来处理
//函数定义如下 func Marshal(v interface{})([]byte,error)
//假设 我们还是需要生成上面的服务器列表信息 那么如何来处理呢  请看下面的例子
package main

import (
	"encoding/json"
	"fmt"
)

type Server struct {
	ServerName string
	ServerIP   string
}
type Serverslice struct {
	Servers []Server
}

func main() {
	var s Serverslice
	s.Servers = append(s.Servers, Server{ServerName: "Shanghai_VPN", ServerIP: "127.0.0.1"})
	s.Servers = append(s.Servers, Server{ServerName: "Beijing_VPN", ServerIP: "127.0.0.2"})
	b, err := json.Marshal(s)
	if err != nil {
		fmt.Println("json err", err)

	}
	fmt.Println(string(b))
}

//我们看到 上面输出的字段名都是大写的 如果你想用小写的怎么办呢 把结构体的字段名 改成小写的？
//JSON输出的时候必须注意 只有导出的字段才会被输出  如果修改 字段名那么会发现  什么都不会输出  所以必须通过 struct tag定义来实现
type Server struct{
ServerName string `json:"serverName"`
 ServerIP   string `json:"serverIP"`
}
type Serverslice struct {
    Servers []Server `json:"servers"`
}
//通过修改上面的结构体定义 输出的JSON串  就和我们最开始定义的JSON串保持一致了
//针对 JSON输出 我们在定义struct tag的时候需要注意几点是
//字段的tag 是- 那么这个字段不会输出到json
//tag中带有 自定义名称  那么这个自定义名称会出现在JSON的字段名中 例如上面例子中的serverName
//tag 中如果带有 omitempty选项 那么如果该字段值为空 就不会输出到JSON串中
//如果字段类型是bool string int int64 等 而tag中带有",string"选项 那么这个字段输出到json的时候会把该字段相应的值转换成json字符串
package main

import (
	"encoding/json"

	"os"
)

type Server struct {
	// ID 不会导出到JSON中
	ID int `json:"-"`

	// ServerName 的值会进行二次JSON编码
	ServerName  string `json:"serverName"`
	ServerName2 string `json:"serverName2,string"`

	// 如果 ServerIP 为空，则不输出到JSON串中
	ServerIP string `json:"serverIP,omitempty"`
}

func main() {
	s := Server{
		ID:          3,
		ServerName:  `Go "1.0" `,
		ServerName2: `Go "1.0" `,
		ServerIP:    ``,
	}
	b, _ := json.Marshal(s)
	os.Stdout.Write(b)
}
//会输出以下内容
//{"serverName":"Go \"1.0\" ","serverName2":"\"Go \\\"1.0\\\" \""}
//Marshal函数只有在转换成功的时候才会返回数据，在转换的过程中我们需要注意几点：
//
//JSON对象只支持string作为key，所以要编码一个map，那么必须是map[string]T这种类型(T是Go语言中任意的类型)
//Channel, complex和function是不能被编码成JSON的
//嵌套的数据是不能编码的，不然会让JSON编码进入死循环
//指针在编码的时候会输出指针指向的内容，而空指针会输出null
//本小节，我们介绍了如何使用Go语言的json标准包来编解码JSON数据，同时也简要介绍了如何使用第三方包go-simplejson来在一些情况下简化操作，学会并熟练运用它们将对我们接下来的Web开发相当重要。
//


















