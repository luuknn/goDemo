package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

/*
表单 是我们平时编写web应用常用的工具 通过表单我们可以方便的让客户端 和服务器 进行数据的交互
表单是一个包含表单元素的区域  表单元素 是允许用户在表单中 输入信息的元素（比如 文本域 下拉列表 单选框 复选框等）
表单使用表单标签 <form>定义
Go里面对于form处理已经有很方便的方法了  在request里面的有专门的form处理 可以很方便的整合到web开发里面来
4.1 小节 里面将间接Go如何处理表单的输入 由于不能信任任何用户的输入 所以我们需要对这些 输入进行有效验证
4.2 小节 将就如何进行一些普通的验证进行详细的演示
HTTP协议是一种无状态的协议 那么如何才能 辨别是否是同一个用户呢 同时又如何保证 一个表单 不出现多次递交的情况呢
4.3 4.4 小节里面将对cookie(cookie是存储在客户端的信息 能够每次通过header和服务器进行交互的数据)等详细讲解
表单还有一个很大的功能就是能够上传文件 那么Go是如何处理文件上传的呢 针对大文件上传我们如何有效的处理呢
4.5 小节 我们将一起处理文件上传的知识
*/
/*
4.1  上面递交表单到服务器/login 当用户输入信息点击登录 之后 会跳转到服务器 路由 login里面
我们首先要判断 这个是什么方式传递过来的 post还是get
http包里面有个很简单的方式就可以获取 我们在前面web的例子的基础上来 看看怎么处理 login页面的forum数据
*/
func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm() //解析url传递的参数，对于POST则解析响应包的主体（request body）
	//注意:如果没有调用ParseForm方法，下面无法获取表单的数据
	fmt.Println(r.Form) //这些信息是输出到服务器端的打印信息
	fmt.Println("path", r.URL.Path)
	fmt.Println("scheme", r.URL.Scheme)
	fmt.Println(r.Form["url_long"])
	for k, v := range r.Form {
		fmt.Println("key:", k)
		fmt.Println("val:", strings.Join(v, ""))
	}
	fmt.Fprintf(w, "Hello 守望子!") //这个写入到w的是输出到客户端的
}

func login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //获取请求的方法
	if r.Method == "GET" {
		t, _ := template.ParseFiles("login.html")
		t.Execute(w, nil)
	} else {
		//请求的是登陆数据，那么执行登陆的逻辑判断
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
	}
}

func main() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	http.HandleFunc("/login", login)         //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

//request.Form 是一个url.Values类型 里面存储的类似 key=value的信息 下面展示了 可以对form数据进行的一些操作
/*
v:=url.Vaules{}
v.Set("name","Ava")
v.Add("friend","Jess")
v.Add("friend","Sarah")
fmt.Println(v.Get("name"))
fmt.Println(v.Get("friend"))
fmt.Println(v["friend"])
Tips: Request本身也提供了FormValue()函数来获取用户提交的参数。如r.Form["username"]也可写成r.FormValue("username")。调用r.FormValue时会自动调用r.ParseForm，所以不必提前调用。r.FormValue只会返回同名参数中的第一个，若参数不存在则返回空字符串。

*/
