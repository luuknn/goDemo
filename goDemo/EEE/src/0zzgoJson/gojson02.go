package main

//import (
//	"encoding/json"
//	"fmt"
//)
//
//type Person struct {
//	Name string `json:"name"`
//	Age  int    `json:"age"`
//}
//type Place struct {
//	City    string `json:"city"`
//	Country string `json:"country"`
//}
//
//var jsonString string = `{
//        "things": [
//            {
//                "name": "Alice",
//                "age": 37
//            },
//            {
//                "city": "Ipoh",
//                "country": "Malaysia"
//            },
//            {
//                "name": "Bob",
//                "age": 36
//            },
//            {
//                "city": "Northampton",
//                "country": "England"
//            }
//        ]
//    }`
//
//func decode(jsonStr []byte) (persons []Person, places []Place) {
//	var data map[string][]json.RawMessage
//	err := json.Unmarshal(jsonStr, &data)
//	if err != nil {
//		fmt.Println(err)
//		return
//	}
//	fmt.Printf("%+v\n", data["things"])
//	for _, item := range data["things"] {
//		persons = addPerson(item, persons)
//		places = addPlace(item, places)
//	}
//	return
//}
//func addPerson(item json.RawMessage, persons []Person) []Person {
//	person := Person{}
//	if err := json.Unmarshal(item, &person); err != nil {
//		fmt.Println(err)
//	} else {
//		if person != *new(Person) {
//			persons = append(persons, person)
//		}
//	}
//	return persons
//}
//func addPlace(item json.RawMessage, places []Place) []Place {
//	place := Place{}
//	if err := json.Unmarshal(item, &place); err != nil {
//		fmt.Println(err)
//	} else {
//		if place != *new(Place) {
//			places = append(places, place)
//		}
//	}
//	return places
//}
//func main() {
//	personsA, placesA := decode([]byte(jsonString))
//	fmt.Printf("%+v\n", personsA)
//	fmt.Printf("%+v\n", placesA)
//
//}

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
延迟解析
因为UserName字段  实际上是在使用的时候 才会用到他的具体类型  因此我们可以延迟解析
使用json.RawMessage方式 将json的字串继续以byte数组方式存在
type User struct {
	UserName json.RawMessage `json:"username"`
	Password string          `json:"password"`
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
	var email string
	if err = json.Unmarshal(u.UserName, &email); err == nil {
		u.Email = email
		return
	}
	var phone int64
	if err = json.Unmarshal(u.UserName, &phone); err == nil {
		u.Phone = phone
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
输出  &{UserName:[49 53 56 57 53 57 50 48 50 49 50] Password:111111 Email: Phone:15895920212}
[49 53 56 57 53 57 50 48 50 49 50]
总体而言 延迟解析和使用空接口的方式类似 需要再次调用Unmarshal方法 对json.RawMessage 进行解析 原理和解析到接口的形式类似
*/
/*
不定字段解析
对于未知json结构的解析 不同的数据类型可以映射到接口 或者使用延迟解析 有时候 会遇到json的数据字段都不一样的情况
例如 需要解析下面一个json字串
var jsonString string = `{
        "things": [
            {
                "name": "Alice",
                "age": 37
            },
            {
                "city": "Ipoh",
                "country": "Malaysia"
            },
            {
                "name": "Bob",
                "age": 36
            },
            {
                "city": "Northampton",
                "country": "England"
            }
        ]
    }`

json字串的是一个对象 其中一个key things的值是一个数组 这个数组的每一个item都未必一样 大致是两种数据
结构 可以抽象为person和place 即 定义下面的结构体
import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type Place struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

var jsonString string = `{
        "things": [
            {
                "name": "Alice",
                "age": 37
            },
            {
                "city": "Ipoh",
                "country": "Malaysia"
            },
            {
                "name": "Bob",
                "age": 36
            },
            {
                "city": "Northampton",
                "country": "England"
            }
        ]
    }`

func decode(jsonStr []byte) (persons []Person, places []Place) {
	var data map[string][]map[string]interface{}
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	for i := range data["things"] {
		item := data["things"][i]
		if item["name"] != nil {
			persons = addPerson(persons, item)
		} else {
			places = addPlace(places, item)
		}
	}
	return
}
func addPerson(persons []Person, item map[string]interface{}) []Person {
	name := item["name"].(string)
	age := item["age"].(float64)
	person := Person{name, int(age)}
	persons = append(persons, person)
	return persons
}
func addPlace(places []Place, item map[string]interface{}) []Place {
	city := item["city"].(string)
	country := item["country"].(string)
	place := Place{City: city, Country: country}
	places = append(places, place)
	return places

}
func main() {
	personsA, placesA := decode([]byte(jsonString))
	fmt.Printf("%+v\n", personsA)
	fmt.Printf("%+v\n", placesA)

}
unmarshal json字串到一个map结构 然后迭代item 并使用type断言的方式 解析数据 迭代的时候会判断item是否是person还是place然后调用对应的解析方法
*/

/*
混合结构
混合结构很好理解 如同我们前面解析username为email和phone两种情况 就在结构中定义好这两种结构即可
import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type Place struct {
	City    string `json:"city"`
	Country string `json:"country"`
}
type Mixed struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	City    string `json:"city"`
	Country string `json:"country"`
}

var jsonString string = `{
        "things": [
            {
                "name": "Alice",
                "age": 37
            },
            {
                "city": "Ipoh",
                "country": "Malaysia"
            },
            {
                "name": "Bob",
                "age": 36
            },
            {
                "city": "Northampton",
                "country": "England"
            }
        ]
    }`

func decode(jsonStr []byte) (persons []Person, places []Place) {
	var data map[string][]Mixed
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", data["things"])
	for i := range data["things"] {
		item := data["things"][i]
		if item.Name != "" {
			persons = append(persons, Person{Name: item.Name, Age: item.Age})
		} else {
			places = append(places, Place{City: item.City, Country: item.Country})
		}
	}
	return
}
func main() {
	personsA, placesA := decode([]byte(jsonString))
	fmt.Printf("%+v\n", personsA)
	fmt.Printf("%+v\n", placesA)

}
混合结构的思路很简单 借助golang会初始化没有匹配的json和抛弃没有匹配的json给特定的字段赋值
比如每一个item都具有四个字段 只不过有的会匹配person的json数据 有的则匹配place 没有匹配的字段则是零值
接下来再根据item的具体情况 分别赋值到对应的Person或者Place结构
混合结构的解析方式也很不错  思路还是借助了解析json中抛弃不要的字段  借助零值的处理
*/
/*
json.RawMessage
json.RawMessage 非常有用  延迟解析也可以使用这个样例  我们已经介绍过类似的技巧  下面就贴代码了
import (
	"encoding/json"
	"fmt"
)

type Person struct {
	Name string `json:"name"`
	Age  int    `json:"age"`
}
type Place struct {
	City    string `json:"city"`
	Country string `json:"country"`
}

var jsonString string = `{
        "things": [
            {
                "name": "Alice",
                "age": 37
            },
            {
                "city": "Ipoh",
                "country": "Malaysia"
            },
            {
                "name": "Bob",
                "age": 36
            },
            {
                "city": "Northampton",
                "country": "England"
            }
        ]
    }`

func decode(jsonStr []byte) (persons []Person, places []Place) {
	var data map[string][]json.RawMessage
	err := json.Unmarshal(jsonStr, &data)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v\n", data["things"])
	for _, item := range data["things"] {
		persons = addPerson(item, persons)
		places = addPlace(item, places)
	}
	return
}
func addPerson(item json.RawMessage, persons []Person) []Person {
	person := Person{}
	if err := json.Unmarshal(item, &person); err != nil {
		fmt.Println(err)
	} else {
		if person != *new(Person) {
			persons = append(persons, person)
		}
	}
	return persons
}
func addPlace(item json.RawMessage, places []Place) []Place {
	place := Place{}
	if err := json.Unmarshal(item, &place); err != nil {
		fmt.Println(err)
	} else {
		if place != *new(Place) {
			places = append(places, place)
		}
	}
	return places
}
func main() {
	personsA, placesA := decode([]byte(jsonString))
	fmt.Printf("%+v\n", personsA)
	fmt.Printf("%+v\n", placesA)

}
把things的item数组解析成一个json.RawMessage 然后在定义其他结构逐步解析 上述例子其实在真实开发环境下
应该尽可能避免 像person或place这样的数据 可以定义两个数组分别存储他们 这样就方便很多

总结
关于golang解析json的基本介绍到此结束 想要解析简单 就需要定义明确的map结构
面对无法确定的数据结构或类型  再动态解析 可以借助接口与断言的方式解析 也可以用json.RawMessage延迟解析
具体使用情况 还得考虑实际的需求和应用场景 总而言之  使用json作为现在spi的数据通信方式 已经很普遍了
*/
