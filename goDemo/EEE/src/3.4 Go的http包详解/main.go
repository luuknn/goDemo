package mian

import (
	"fmt"
)
//Go 的http包详解
//前面小节 介绍了 Go怎么样实现了 web工作模式的一个流程 这一小节 我们将详细的解剖一下http包 看它到底是怎么样实现整个过程的
//go的http有两个核心功能 Conn ServeMux
//Conn 的goroutine
//与我们一般编写的http服务器 不同  Go为了 实现高并发和高性能 使用了 Goroutines来处理 Conn的读写事件
//这样每个请求都能 保持独立 相互不会阻塞  可以高效的响应网络事件 这是Go高效的保证
//Go在等待客户端请求里面是这样写的
c,err :=srv.newConn(rw)
if err!=mil {
continue
}
go c.serve()
//这里我们可以看到客户端的每次请求都会创建一个Conn 这个Conn里面保存了该次请求的信息
//然后再传递到对应的handler 该handler中便可以读取到相应的header信息 这样保证了每个请求的独立性

//ServeMux的自定义
//我们前面小节 讲述conn.server的时候 其实内部是调用了 http包默认的路由器
//通过路由器 把本次请求的信息 传递到了后端的处理函数 那么这个路由器是怎么实现的呢
type ServeMux struct{
mu sync.RWMutex//锁  由于请求涉及处理并发 因此 这里需要一个锁机制
m map[string]muxEntry//路由规则  一个string对应一个mux实体 这里的string就是注册 的路由表达式
}
//下面看一下 muxEntry
type muxEntry struct{
explicit bool //是否精确匹配
h Handler //这个路由表达式对应哪个handler
}
//接着看一下Handler的定义
type Handler interfacer{
ServeHTTP(ResponceWriter,*Request)//路由实现器
}
/*
Handler 是一个接口 但是前一小节中 sayHelloName函数并没有实现ServeHTTP这个接口 为什么能添加呢
原来http包里还定义了一个类型 HandlerFunc 我们定义的函数 sayHelloName就是这个HandlerFunc调用后的结果
这个类型默认 就实现了ServeHTPP这个接口 即我们调用了 HandleFunc(f) 强制类型转换f成为 HandlerFunc类型
这样f就拥有了ServeHTTP方法了
*/
type HandlerFunc func (ResponseWriter,*Request){
//ServeHttp calls f(w,r)
func (f HandlerFunc) ServeHTTP(w ResponseWriter,r *Request)
{f(w,r)}
}
//路由器里面存储好了相应的路由规则后 那么 具体的请求又是怎么分发的呢
//路由器 接收到请求之后调用mux.handler(r).ServeHTTP(w,r)
//也就是 调用 对应路由的handler的servehttp 接口 那么mux.handler(r)怎么处理呢
func(mux *ServeMux) handler(r *Request) Handler{
mux.mu.RLock()
defer mux.mu.RUnlock()
//host-specifix pattern takes precedence over genneric ones
h:=mux.match(r.Host +r.URL.Path)
if h==nil{
h=mux.match(r.URL.PATH)
}
if h==nil{
h=NotFoundHandler()
}
return h
}
//原来他是根据用户请求的URL和路由器里面存储的map去匹配的
//当匹配到之后返回存储的handler
 //调用这个handler的ServeHttp接口就可以执行到相应的函数了
 //通过上面这个介绍 我们了解了整个路由过程 Go其实支持外部实现的路由器ListenerAndServer的第二个参数就是用以配置外部路由器的
 //它是一个Handler接口 即 外部路由器 只要实现了Handler接口就可以 我们可以 在自己实现的路由器的ServeHTTP里面实现自定义路由功能
 package main

import (
	"fmt"
	"net/http"
)

type MyMux struct{}

func (p *MyMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path == "/123" {
		sayhelloName(w, r)
		return
	}
	http.NotFound(w, r)
	return
}

func sayhelloName(w http.ResponseWriter, r *http.Request) {

	fmt.Fprintf(w, "Hello myroute!") //这个写入到w的是输出到客户端的
}

func main() {
	mux := &MyMux{}
	http.ListenAndServe(":9090", mux)
}
//GO代码 的执行流程
//通过对http包的分析之后 现在让我们来梳理一下 整个代码的执行过程
//首先调用Http.HandleFunc 
	//	按顺序做了几件事 1 调用了DefaultServerMux的HandleFunc 2 调用了DefaultServeMux的Handle 3 往 DefaultServeMux的map[string]muxEntry中增加对应的handler和路由规则
	
//其次调用了http.ListenAndServr(":9090",nil)
	//按顺序做了几件事情
	//实例化Server
	//调用server的ListenAndServer()
	//调用net.Listen("tcp",addr) 监听端口
	//启动一个for循环 在循环体中 Accept请求
	//对每一个请求实例化一个Conn并且开启了一个goroutine为这个请求进行服务 go c.serve()
	//读取每个请求的内容 w,err :=c.readRequest()
	//判断handler是否为空 如果没有设置handler 这个例子就是没有设置handler handler就设置为DefaultServeMux
	//调用handler的ServeHttp
	//在这个例子中就进入到DefaultServeMux.ServeHTTP
	//根据request选择handler 并且进入到 handler的ServeHTTP mux.handler(r).ServeHTTP(W,R)
	//选择handler 判断是否有路由能满足这个request 循环遍历 ServeMux的muxEntry
	//如果有路由满足 就调用这个路由handler的ServeHt
	//如果没有路由能满足 就调用 NotFoundhandler的ServerHttp



























