package main

//import (
//	"fmt"
//	"net/http"
//	"time"
//)
//
//func text(w http.ResponseWriter, r *http.Request) {
//	fmt.Fprintln(w, "hello tomorrow")
//}
//func index(w http.ResponseWriter, r *http.Request) {
//	w.Header().Set("Content-Type", "text/html")
//	html := `<doctype html>
//        <html>
//        <head>
//          <title>Hello World</title>
//        </head>
//        <body>
//        <p>
//          <a href="/welcome">Welcome</a> |  <a href="/message">Message</a>
//        </p>
//        </body>
//</html>`
//	fmt.Fprintln(w, html)
//}
//func main() {
//	mux := http.NewServeMux()
//	mux.Handle("/", http.HandlerFunc(index))
//	mux.HandleFunc("/text", text)
//	server := &http.Server{
//		Addr:         ":8000",
//		ReadTimeout:  60 * time.Second,
//		WriteTimeout: 60 * time.Second,
//		Handler:      mux,
//	}
//	server.ListenAndServe()
//
//}

/*
Golang构建HTTP服务 Handler ServeMux与中间件 Golang标准库http包提供了基础的http服务 这个服务又基于
Handler接口和ServeMux结构的做Mutipexer 实际上 go的作者设计handler这样的接口 不仅提供了默认的ServeMux对象
开发者也可以自定义 ServeMux对象

本质上ServeMux只是一个路由管理器 而它本身也实现了Handler接口的ServeHTTP方法 因此 围绕Handler接口的方法
ServeHTTP 可以轻松的写出go中的中间件

在go的http路由原理讨论中 追本溯源还是讨论Handler接口和ServeMux结构 下面就基于这两个对象开始更多关于go中http的故事吧

介绍http库源码的时候 创建http服务的代码很简单 实际上代码隐藏了很多细节 才有了后来的流程介绍 本文的主要目的
是把这些细节暴露 从更底层的方式开始 一步步隐藏细节 完成样例代码一样的逻辑 了解更多http包的原理后 才能基于此构建中间件

自定义的Handler
标准库http提供了Handler接口 用于开发者实现自己的handler只要实现接口的ServeHTTP方法即可
*/
/*
package main

import (
	"fmt"
	"net/http"
)

type textHandler struct {
	responseText string
}

func (th *textHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, th.responseText)
}

type indexHandler struct{}

func (in *indexHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `
<doctype html>
<html>
<head>
<title>Hello World</title>
</head>
<body>
<p>
<a href="/welcome">Welcome</a> | <a href="/message">Message</a>
</p>
</body>
</html>
`
	fmt.Fprintln(w, html)
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", &indexHandler{})

	thWelcome := &textHandler{"TextHandler!"}
	mux.Handle("/text", thWelcome)
	http.ListenAndServe(":8000", mux)

}
上面自定义了两个handler结构都实现了ServeHTTP方法 我们知道 NewServeMux可以创建一个ServeMux实例 ServeMux同时也实现了
ServeHTTP方法 因此代码中的mux也是一种handler 把它当成参数传给http.ListenAndServe方法 后者会把mux传给server实例
因为指定了handler因此整个http服务就不再是DefaultServeMux 而是mux 无论是在注册路由还是提供请求服务的时候

有一点值得注意 这里并没有使用HandleFunc注册路由而是直接使用了mux注册路由。当没有指定mux的时候 系统需要创建一个defaultServeMux,
此时我们已经有了mux 因此 不再需要http.HandleFunc方法 直接使用mux 的Handle方法注册即可

此外 Handle第二个参数是一个handler处理器 并不是HandleFunc的一个handler函数 其原因也是因为
mux.Handle本质上就需要绑定url的patter模式和handler处理器 即可 既然indexHandler是handle处理器
当然就能作为参数 一切请求的处理过程 都交给实现接口方法ServeHTTP就行了
1 http.HandleFunc(pattern,function) 2 defaultServemux.HandleFunc(pattern,function) 3 mux.Handle(pattern,HandlerFunc(handler))
right 4 mux.handle(pattern,handler)
12只是为了创建一个ServeMux实例 然后调用实例的Handle方法 右边的直接调用了mux实例的Handle方法
创建handler处理器
上面费尽口舌啰嗦 不就是 1 2 3 与3的差别么 开发者只需要写函数即可 不用再定义结构 因此 下面将直接创建handler函数
调用go的方法 将函数转变成handler处理器
*/
/*
package main

import (
	"fmt"
	"net/http"
)

func text(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello Monday")
}
func index(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	html := `<doctype html>
        <html>
        <head>
          <title>Hello World</title>
        </head>
        <body>
        <p>
          <a href="/welcome">Welcome</a> |  <a href="/message">Message</a>
        </p>
        </body>
</html>`
	fmt.Fprintln(w, html)
}
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(index))
	//mux.Handle("/text", http.HandlerFunc(text))
	mux.HandleFunc("/text", text)
	http.ListenAndServe(":8000", mux)

}
代码中使用了http.HandlerFunc 方法 直接将一个handler函数转变成实现了handler处理器 等价于图中3的步骤
而mux.HandleFunc('/text',text)就更进一步 与图中2步骤一致
与defaultServemux.HandleFunc(pattern,function)的用法一样

使用默认的DefaultServeMux
经过了上面两个过程的转化 隐藏更多的细节 代码与defaultServeMux的方式越来越像 下面再去掉自定义的serveMux只需要修改main函数的处理逻辑即可

func main() {
	http.Handle("/", http.HandlerFunc(index))
	http.HandleFunc("/text", text)
	http.ListenAndServe(":8000", nil)
}
*/

/*
自定义Server
默认的DefaultServeMux创建的判断来自server对象 如果server对象不提供handler才会使用默认的
serveMux对象 既然ServeMux可以自定义 那么Serve对象 一样可以
使用http.Server 即可创建自定义的serve对象
func main(){
    http.HandleFunc("/", index)
    server := &http.Server{
        Addr: ":8000",
        ReadTimeout: 60 * time.Second,
        WriteTimeout: 60 * time.Second,
    }
    server.ListenAndServe()
}
自定义的serverMux对象也可以传到server对象中
func main() {
	mux := http.NewServeMux()
	mux.Handle("/", http.HandlerFunc(index))
	mux.HandleFunc("/text", text)
	server := &http.Server{
		Addr:         ":8000",
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
		Handler:      mux,
	}
	server.ListenAndServe()

}
可见 go中的路由和处理函数之间的关系非常密切 同时又很灵活 通过巧妙的使用Handler接口可以设计出优雅的中间件程序

*/

/*
中间件Middleware
所谓中间件 就是连接上下级不同功能的函数或者软件 通常进行一些包裹函数的行为 为被包裹函数提供添加一些功能或行为.
前文的HandleFunc就能把签名为 func(w http.ResponseWriter,r *http.Request)的函数包裹成handler这个函数也算是中间件

go的http中间件很简单，只要实现一个函数签名为func(http.Handler) http.Handler的函数即可。
http.Handler是一个接口，接口方法我们熟悉的为serveHTTP。返回也是一个handler。
因为go中的函数也可以当成变量传递或者或者返回，因此也可以在中间件函数中传递定义好的函数，只要这个函数是一个handler即可，
即实现或者被handlerFunc包裹成为handler处理器。


















*/
