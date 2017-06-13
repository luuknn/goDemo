package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
)

func templateHandler(w http.ResponseWriter, r *http.Request) {
	//	t, _ := template.ParseFiles("src/0zzgohttp/templates/layout.html")
	//	fmt.Println(t.Name())
	//	t.Execute(w, "Hello World")
	/*tmpl := `<!DOCTYPE html>
	<html>
		<head>
		</head>
		<body>
		{{ . }}
		</body>

	</html>`
	t := template.New("hello.html")
	t, _ = t.Parse(tmpl)
	fmt.Println(t.Name())
	t.Execute(w, "Hello world!")*/
	t, _ := template.ParseFiles("src/0zzgohttp/templates/layout.html", "src/0zzgohttp/templates/index.html")
	fmt.Println(t.Name())
	t.ExecuteTemplate(w, "layout", "Hello Monday")
}
func checkErr(err error) {
	if err != nil {
		err.Error()
	}

}
func main() {
	//绑定路由 如果访问 /template 调用Handler方法
	http.HandleFunc("/template", templateHandler)
	//使用tcp协议监听 8888端口号
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}

}

/*
Golang Template 简明笔记
前后端分离的Restful架构大行其道 传统的模板技术已经不多见了 实际上只是渲染的地方由后端转移到了前端 模板的渲染技术本质上还是一样的
简而言之 就是 字串模板和数据的结合

golang提供了两个标准库 用来处理模板 text/template 和 html/template 我们使用 html/template格式化html字符
模板引擎
模板引擎很多 Python的jinja nodejs的jade等都很好 所谓模板引擎 则将 模板和数据进行渲染的输出格式化后的字符程序 对于go
执行这个流程 大概需要三步
1 创建模板对象   2  加载模板字串  3 执行渲染模板
其中最后一步 就是把加载的字符和数据进行格式化 其过程可以总结为
html    -     Template engine    Data
					Template
我们打印了t模板对象的Name方法 实际上 每一个模板 都有一个名字 如果不显示指定这个名字
go将会把文件名(包括扩展名当成名字)本例则是layout.html

go不仅可以解析模板文件 也可以直接解析模板字串 这就是标准的处理 新建 加载 执行三部曲
实际开发中 最终页面很可能是多个模板文件的嵌套结果 go的ParseFiles也支持加载多个模板文件的嵌套结果
go 的ParseFiles也支持加载多个模板文件 不过模板对象的名字则是第一个模板文件名

func templateHandler(w http.ResponseWriter,r *http.Request){
t,_:=template.ParseFiles("templates/layout.html","templates/index.html")
fmt.Println(t.Name())
t.Execute(w,"Hello world")
}
可以看到 打印的还是layout.html的名字 执行模板的时候 并没有index.html的模板内容 此外还有ParseGlob方法 可以通过glob通配符加载模板

模板命名与嵌套
模板命名
前文已经提及 模板对象是有名字的 可以在创建模板对象的时候显示命名 也可以让go自动命名 可是涉及到嵌套模板的时候
该如何命名模板呢 毕竟模板文件有好几个
go 提供了ExecuteTemplate方法 用于执行指定名字的模板 例如加载layout.html模板的时候 可以指定layout.html
*/

/*
func templateHandler(w http.ResponseWriter,r *http.Request){
t,_:=template.ParseFiles("src/0zzgohttp/templates/layout.html")
fmt.Println(t.Name())
t.ExecuteTemplate(w,"layout","Hello world")
}

似乎和Execute方法没有太大的区别 下面修改一下layout.html文件:
{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
<meta charset="UTF-8">
<title>layout</title>
</head>
<body>
	<h3>This is layout 233</h3>
	template data:{{ . }}
</body>
</html>
{{ end }}
在模板文件中 使用了define这个action给模板文件命名了 虽然我们ParseFiles方法返回的模板对象t的名字还是
layout.html 但是ExecuteTemplate执行的模板却是html文件中定义的layout
*/

/*
不仅可以通过define定义模板 还可以通过template action引入模板 类似jinja的include特性
修改layout.html 和index.html
{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>layout</title>
  </head>
  <body>
    <h3>This is layout</h3>
    template data: {{ . }}
    {{ template "index" }}
  </body>
</html>
{{ end }}

{{ define "index" }}
<div style="background: yellow">
    this is index.html
</div>
{{ end }}

go的代码也需要修改 使用ParseFiles加载需要渲染的模板文件
func templateHandler(w http.ResponseWriter, r *http.Request){
    t, _ :=template.ParseFiles("templates/layout.html", "templates/index.html")
    t.ExecuteTemplate(w, "layout", "Hello world")
}
*/

/*
单文件嵌套
总而言之 创建模板对象后和加载多个模板文件 执行模板文件的时候需要指定base模板（layout）
在base模板中 可以include其他命名的模板 无论. define template 这些花括号包裹的东西都是go的action(模板标签)

Action
action是go模板中用于动态执行一些逻辑和战士数据的形式大致分为下面几种
条件语句 迭代 封装 引用
*/

/*
条件判断
条件判断的语法很简单
{{ if arg }}
some content
{{ end }}

{{ if arg }}
some content
{{ else }}
other content
{{ end }}
arg可以是基本数据结构 也可以是表达式 if-end 包裹的内容条件为真的时候展示 与if
语句一样 模板也可以有else语句
func templateHandler(w http.ResponseWriter,r *http.Request){
t,_:=template.ParseFiles("templates/layout.html")
rand.Seed(time.Now().Unix())
t.Executetemplate(w,"layout",rand.Intn(10)>5)

}

{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8">
<title>layout</title>
</head>
<body>
	<h3>This is layout</h3>
	template data: {{ . }}
	{{ if . }}
	Number is greater than 5!
	{{ else }}
	Number is 5 or less!
	{{ end }}

</body>
</html>
{{ end }}

此时 就能看见 当.值为true的时候显示if的逻辑 否则显示else的逻辑
*/
/*
迭代
对于一些数组 切片 或者是map 可以使用迭代的action 与go的迭代类似 使用range进行处理
func templateHandler(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles("templates/layout.html"))
	 daysOfWeek := []string{"Mon", "Tue", "Wed", "Ths", "Fri", "Sat", "Sun"}
	 t.ExecuteTemplate(w,"layout",daysOfWeek)
}

{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>layout</title>
  </head>
  <body>
    <h3>This is layout</h3>
    template data: {{ . }}
    {{ range . }}
    <li>{{ . }}</li>
    {{ end }}
  </body>
</html>
{{ end }}
k=可以看到输出了一堆li列表 迭代的时候还可以使用￥设置循环变量:
{{ range $key,$value := .}}
<li>key :{{ $key }},value:{{ $value }}</li>
{{ else }}
empty
{{ end }}
可以看到与迭代切片很像 range也可以使用else语句
func templateHandler(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles("templates/layout.html"))
    daysOfWeek := []string{}
    t.ExecuteTemplate(w, "layout", daysOfWeek)
}
{{ range . }}
    <li>{{ . }}</li>
{{ else }}
 empty
{{ end }}
当range的结构为空的时候 则会执行else分支的逻辑
*/

/*
with封装
with语言在Python中可以开启一个上下文环境 对于go模板 with语句类似
其含义就是创建一个封闭的作用域 在其范围内 可以使用 .action 而与外面的.无关
只与with的参数有关
{{ with arg }}
此时的点. 就是arg
{{ end }}

{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>layout</title>
  </head>
  <body>
    <h3>This is layout</h3>
    template data: {{ . }}
     {{ with "world"}}
        Now the dot is set to {{ . }}
     {{ end }}
  </body>
</html>
{{ end }}
可见with语句. 与其外面的.是两个不相关的对象 with语句也可以有else
else中.则和with外面的.一样 毕竟只有with语句内才有封闭的上下文

{{ with ""}}
 Now the dot is set to {{ . }}
{{ else }}
 {{ . }}
{{ end }}

*/

/*
引用
我们已经介绍了模板嵌套引用的技巧 引用除了模板的include还包括参数的传递
func templateHandler(w http.ResponseWriter, r *http.Request) {
    t := template.Must(template.ParseFiles("templates/layout.html", "templates/index.html"))
    daysOfWeek := []string{"Mon", "Tue", "Wed", "Ths", "Fri", "Sat", "Sun"}
    t.ExecuteTemplate(w, "layout", daysOfWeek)
}
{{ define "layout" }}
<!DOCTYPE html>
<html>
<head>
    <meta http-equiv="Content-Type" content="text/html; charset=utf-8">
    <title>layout</title>
  </head>
  <body>
    <h3>This is layout</h3>
    layout template data: ({{ . }})
    {{ template "index" }}
  </body>
</html>
{{ end }}

{{ define "index" }}
<div style="background: yellow">
    this is index.html ({{ . }})
</div>
{{ end }}
我们可以修改引用语句  {{ template "index" . }}
把参数传给字模板 再次访问 就能看到index.html模板也有数据了

### 参数，变量和管道
模板的参数可以是go中的基本数据类型，如字串，数字，布尔值，数组切片或者一个结构体。在模板中设置变量可以使用

$variable := value
```
我们在range迭代的过程使用了设置变量的方式。
go还有一个特性就是模板的管道函数，熟悉django和jinja的开发者应该很熟悉这种手法。通过定义函数过滤器，实现模板的一些简单格式化处理。并且通过管道哲学，这样的处理方式可以连成一起。p1 | p2 | p3
例如 模板内置了一些函数，比如格式化输出：12.3456 | printf "%.2f"
*/
