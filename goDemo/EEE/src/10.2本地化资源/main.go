package main

import (
	"fmt"
)

var locales map[string]map[string]string

func main() {
	locales = make(map[string]map[string]string, 2)
	en := make(map[string]string, 10)
	en["pea"] = "pea"
	en["bean"] = "bean"
	locales["en"] = en
	cn := make(map[string]string, 10)
	cn["pea"] = "豌豆"
	cn["bean"] = "毛豆"
	locales["zh-CN"] = cn
	lang := "zh-CN"
	fmt.Println(msg(lang, "pea"))
	fmt.Println(msg(lang, "bean"))

}
func msg(locale, key string) string {
	if v, ok := locales[locale]; ok {
		if v2, ok := v[key]; ok {
			return v2
		}

	}

	return ""

}

/*10.2 本地化资源
前面小节我们介绍了如何设置locale 设置好locale之后我们需要解决的问题 就是如何存储相应的locale对应的信息呢
这里面的信息包括 文本信息 时间和日期 货币值 图片 包含文件以及视图等资源 那么接下来 我们讲对这些信息 进行介绍
go语言中我们把这些格式信息存储在json中 然后通过合适的方式展现出来 接下来中文和英文两种语言对比举例
存储格式文件 en.json 和 zh-CN.json

本地化文本消息
本信息是编写web应用中最常见的 也是本地化资源中最多的信息想要以适合本地语言的方式来显示文本信息
可行的一种方案 建立需要的语言相应的map来维护一个key-value的关系 在输出之前需要从适合的map中去获取
相应的文本 如下是一个简单的示例
package main

import (
	"fmt"
)

var locales map[string]map[string]string

func main() {
	locales = make(map[string]map[string]string, 2)
	en := make(map[string]string, 10)
	en["pea"] = "pea"
	en["bean"] = "bean"
	locales["en"] = en
	cn := make(map[string]string, 10)
	cn["pea"] = "豌豆"
	cn["bean"] = "毛豆"
	locales["zh-CN"] = cn
	lang := "zh-CN"
	fmt.Println(msg(lang, "pea"))
	fmt.Println(msg(lang, "bean"))

}
func msg(locale, key string) string {
	if v, ok := locales[locale]; ok {
		if v2, ok := v[key]; ok {
			return v2
		}

	}

	return ""

}


上面示例演示了不同的locale的文本翻译 实现了中文和英文对于同一个key显示不同语言的实现 上面实现了中文的文本消息
如果想切换到英文版本 只需要把lang设置为en即可

有些时候仅是key value 替换是不能满足需要的例如 I am 30 years old 中文表达 我今年20岁了
而此处 30是一个变量 该怎么办呢 这时候 我们可以结合printf函数来实现 请看下面的代码
en["how old"]="I am %d years old"
cn["how old"]="我今年%d岁了"
fmt.Printf(msg(lang,"how old"),30)
上面的示例代码仅用以演示内部的实现方案 而实际数据是存储在json里面的
所以我们可以通过json.Unmarshal来为相应的map填充数据

本地化日期和时间
因为时区的关系 同一时刻 在不同的地区 表示是不一样的 而且因为locale的关系 时间格式也不尽相同
例如中文环境 下可能显示2012年10月24日 星期三 23时11分13秒 CST，而在英文环境下可能显示:Wed Oct 24 23:11:13 CST 2012。
这里面我们需要解决两点
时区问题 格式问题
en["time_zone"]="America/Chicago"
cn["time_zone"]="Asia/Shanghai"

loc,_:=time.LoadLocation(msg(lang,"time_zone"))
t:=time.Now()
t = t.In(loc)
fmt.Println(t.Format(time.RFC3339))
我们可以通过类似处理文本格式的方式来解决时间格式的问题，举例如下:

en["date_format"]="%Y-%m-%d %H:%M:%S"
cn["date_format"]="%Y年%m月%d日 %H时%M分%S秒"

fmt.Println(date(msg(lang,"date_format"),t))

func date(fomate string,t time.Time) string{
    year, month, day = t.Date()
    hour, min, sec = t.Clock()
    //解析相应的%Y %m %d %H %M %S然后返回信息
    //%Y 替换成2012
    //%m 替换成10
    //%d 替换成24
}
本地化货币值

各个地区的货币表示也不一样，处理方式也与日期差不多，细节请看下面代码:

en["money"] ="USD %d"
cn["money"] ="￥%d元"

fmt.Println(date(msg(lang,"date_format"),100))

func money_format(fomate string,money int64) string{
    return fmt.Sprintf(fomate,money)
}

本地化视图和资源

我们可能会根据Locale的不同来展示视图，这些视图包含不同的图片、css、js等各种静态资源。那么应如何来处理这些信息呢？首先我们应按locale来组织文件信息，请看下面的文件目录安排：

views
|--en  //英文模板
    |--images     //存储图片信息
    |--js         //存储JS文件
    |--css        //存储css文件
    index.tpl     //用户首页
    login.tpl     //登陆首页
|--zh-CN //中文模板
    |--images
    |--js
    |--css
    index.tpl
    login.tpl
有了这个目录结构后我们就可以在渲染的地方这样来实现代码：

s1, _ := template.ParseFiles("views"+lang+"index.tpl")
VV.Lang=lang
s1.Execute(os.Stdout, VV)
而对于里面的index.tpl里面的资源设置如下：

// js文件
<script type="text/javascript" src="views/{{.VV.Lang}}/js/jquery/jquery-1.8.0.min.js"></script>
// css文件
<link href="views/{{.VV.Lang}}/css/bootstrap-responsive.min.css" rel="stylesheet">
// 图片文件
<img src="views/{{.VV.Lang}}/images/btn.png">
采用这种方式来本地化视图以及资源时，我们就可以很容易的进行扩展了。

总结
本小节 介绍了如何使用及存储本地资源 有时需要通过转化函数来实现 有时通过lang来设置 但是
最终都是通过key-value的方式来存储locale对应的数据 在需要时取出相应的locale信息后
如果是文本信息就直接输出 如果是日期时间 或者货币 则需要
通过 printf 或其它格式化函数来处理 而对于不同locale的视图和资源则是最简单的 只要在路径里面加上lang就可以实现了




















*/
