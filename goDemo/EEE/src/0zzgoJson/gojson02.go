package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type User struct {
	UserName interface{} `json:"username"`
	Password string      `json:"password"`
	Email    string
	Phone    int64
}

var jsonString string = `{
"username":15895920212,
"password":"111111"
}
`

func Decode(r io.Reader) (u *User, err error) {
	u = new(User)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	switch t := u.UserName.(type) {
	case string:
		u.UserName = t
	case float64:
		u.UserName = int64(t)
	}
	return

}

func main() {
	user, err := Decode(strings.NewReader(jsonString))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", user)
	fmt.Println(user.UserName)

}

/*
golang处理json  解码
定义结构
与编码json的Marshal类似 解析json也提供了Unmarshal方法 对于解析json 也大致分为两步 首先定义结构 然后
调用Unmarshal 方法序列化 我们先从简单的例子着手
import (
	"encoding/json"
	"fmt"
	"log"
)

type Account struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Money    float64 `json:"money"`
}

var jsonString string = `{
"email":"jianling.shih@gmail.com",
"password":"111111",
"money":100.5
}
`

func main() {
	account := Account{}
	err := json.Unmarshal([]byte(jsonString), &account)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", account)

}
unmarshal接受一个byte数组 和空接口指针的参数 和sql中 读取数据类似 先定一个数据实例 然后传指针地址
与编码类似  golang会将json的数据结构和go的数据结构进行匹配 匹配的原则就是寻找tag的相同的字段
然后查找字段 查询的时候是大小写不敏感的
type Account struct {
    Email    string  `json:"email"`
    PassWord string
    Money    float64 `json:"money"`
}
把Password的tag去掉 依然可以把json的password匹配到Pawwword 但是如果结构的字段是私有化的 即使tag符合
也不会被解析
type Account struct {
    Email    string  `json:"email"`
    password string  `json:"password"`
    Money    float64 `json:"money"`
}
上面的password并不会被解析赋值json的password 大小写不敏感只是针对公有字段而言 再寻找tag或字段的
时候匹配不成功 则会抛弃这个json字段的值
*/
/*
string tag
在编码的时候我们使用tag string 可以把结构定义的数字类型以字串形式编码
同样在解码的时候 只有字串类型的数字 才能被正确解析 或者会报错
import (
	"encoding/json"
	"fmt"
	"log"
)

type Account struct {
	Email    string  `json:"email"`
	Password string  `json:"password"`
	Money    float64 `json:"money,string"`
}

var jsonString string = `{
"email":"jianling.shih@gmail.com",
"password":"111111",
"money":"100.5"
}
`

func main() {
	account := Account{}
	err := json.Unmarshal([]byte(jsonString), &account)
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", account)

}
-tag
与编码一样 tag 的- 也不会被解析 但是会初始化其零值
type Account struct {
    Email    string  `json:"email"`
    Password string  `json:"password"`
    Money    float64 `json:"-"`
}
稍微总结一下 解析json最好的方式就是定义与将要被解析的json的结构
*/
/*
动态解析
通常更加json的格式预先定义golang的结构进行解析是最理想的情况 可是实际开发中 理想情况往往都存在理想的
愿望中 很多json非但格式不确定 有的还有可能是动态数据类型

例如通常登录的时候 往往既可以使用手机号做用户名 也可以使用邮件做用户名 客户端传的json可以是字串 也可以是数字
此时服务端解析就需要技巧了

Decode
前面我们使用了简单的方法Unmarshal直接解析json字串 下面我们使用更底层的方法NewDecode和Decode方法
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type User struct {
	UserName string `json:"username"`
	Password string `json:"password"`
}

var jsonString string = `{
"username":"jianling.shih@gmail.com",
"password":"111111"
}
`

func Decode(r io.Reader) (u *User, err error) {
	u = new(User)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	return

}

func main() {
	user, err := Decode(strings.NewReader(jsonString))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", user)

}
我们定义了一个Decode函数 在这个函数进行json字串的解析  然后调用json的NewDecode方法构造一个Decode对象
最后使用这个对象的Decode方法赋值给定义好的结构对象 对于字串 可以使用strings.NewReader方法 让字串变成一个Stream对象
*/
/*
接口
如果客户端传的username的值是一个数字类型的手机号 那么上面的解析方法会失败  正如我们之前所介绍的动态类型
行为一样 使用空接口可以hold住这样的情景
type User struct {
	UserName interface{} `json:"username"`
	Password string      `json:"password"`
}
运行后看到 &{UserName:1.5895920212e+10 Password:111111}
怎么说 貌似是成功了 可是返回的数字是科学计数法 有点奇怪 可以使用golang的断言
然后转型
import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strings"
)

type User struct {
	UserName interface{} `json:"username"`
	Password string      `json:"password"`
}

var jsonString string = `{
"username":15895920212,
"password":"111111"
}
`

func Decode(r io.Reader) (u *User, err error) {
	u = new(User)
	err = json.NewDecoder(r).Decode(u)
	if err != nil {
		return
	}
	switch t := u.UserName.(type) {
	case string:
		u.UserName = t
	case float64:
		u.UserName = int64(t)
	}
	return

}

func main() {
	user, err := Decode(strings.NewReader(jsonString))
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Printf("%+v\n", user)

}
输出&{UserName:15895920212 Password:111111} 看起来听好 可是我们的username字段始终是一个空接口 使用它的时候
还需要转换类型 这样情况来看 解析的时候就应该转换好类型 那么用的时候就省心了
修改定义的结构如下
type User struct {
    UserName interface{} `json:"username"`
    Password string      `json:"password"`
    Email string
    Phone int64
}
这样就能 通过 fmt.Println(user.Email + " add me") 使用字段进行操作了
*/
/*




































 */
