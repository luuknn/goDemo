package main

import (
	"fmt"
	"net/http"
)

func IndexHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello world")
}

func main() {
	http.HandleFunc("/", IndexHandler)
	http.ListenAndServe(":8080", nil)
}

/*
Client -> Requests -> Multiplexer(router) -> handler ->Response ->Client
about hander
hander函数 : 具有func(w http.ResponseWriter,r *http.Requests)签名的函数
handler 处理 函数: 经过 handlerfunc结构包装的handler函数 它实现了ServerHTTP接口方法的函数
	调用handler处理器的ServerHTTP方法时 即调用handler函数本身
handler对象:实现了Handler接口ServeHTTP方法的结构

handler处理器 和handler对象的差别在于 一个是函数 另外一个是结构 它们都有实现了ServeHTTP方法 很多情况下
它们的功能相似 下文就使用统称 handler 这算是 golang通过接口实现的类动态类型吧

handler函数 签名为func(w http.Resp...) handler处理函数 经过handlerfunc结构包装 具有serveHTTP方法
handler对象 实现了serveHTTP接口方法
serveHTTP方法 调用handler处理器的serveHTTP方法等于调用未包装的handler函数
*/

/*
golang没有继承 类多态的方式可以通过接口实现 所谓接口 则是定义声明函数签名 任何结构 只要实现了与接口函数签名相同的方法
就等同于实现了接口 gohttp 服务都是基于handler进行处理的
type Handler interface{
ServeHTTP(ResponseWriter,*Requesr)
}
任何结构体 只要实现了ServeHTTP方法 这个结构 就可以称之为 handler对象 ServeMux会使用handler并调用其serveHttp方法处理请求并返回响应
*/

/*
ServeMux
了解了Handler之后 再看ServeMux ServeMux的源码很简单
type ServeMux struct{
mu sync.RWMutex
m map[string]muxEntry
hosts bool
}
type muxEntry struct{
explicit bool
h handler
pattern string
}
ServeMux结构中最重要的字段为m 这是一个map key是一些url模式 value是一个muxEntry结构
后者里面定义存储了具体的url模式 和handler
当然所谓的ServeMux也实现了ServeHTTP接口 也算是一个handler 不过不是用来处理request和response 而是用来找到路由注册的handler
*/

/*
Server
除了ServeMux和Handler 还有一个结构Server需要了解 从http.ListernAndServe的源码可以看出 它创建了一个server对象
并调用了server对象的 ListenAndServe方法
func ListenAndServer(addr string,handler Handler) error{
server:=&Server{Addr:addr,Handler:handler}
return server.ListenAndServe()
}

查看server的结构如下
type Server struct{
Addr string
Handler Handler
ReadTimeout time.Duration
WriteTimeout time.Duration
TLSConfig *tls.Config

MaxHeaderBytes int

TLSNextProto map[string]func(*Server, *tls.Conn, Handler)
	ConnState func(net.Conn, ConnState)
	ErrorLog *log.Logger
	disableKeepAlives int32    	nextProtoOnce     sync.Once
	nextProtoErr      error
}
serve结构存储了服务器处理请求常见的字段 其中Handler字段也保留Handler接口
如果Serve接口没有提供Handler结构对象 那么会使用DefaultServeMux做multiplexer后面再作分析

*/

/*
创建HTTP服务
创建一个http服务 大致需要经历两个过程 首先需要注册路由 即提供url模式 和handler函数的映射 其次就是实例化一个serve对象 并开启对客户端的监听
http.HandleFunc("/", IndexHandler)
net/http包暴露的注册路由的api很简单 http.HandleFunc 选取了DefaultServeMux作为 multiplexer
func HandleFunc (pattern string,handler func(Response,*Request)){
		DefaultServeMux.HandleFunc(pattern,handler)
}
*/

/*
那么什么是DefaultServeMux呢 实际上 DefaultServeMux是serveMux的一个实例 当然http包也提供了NewServeMux
方法创建了一个ServeMux实例 默认则创建一个DefaultServeMux

//NewServeMux allocates and returns a new ServeMux
func NewServeMux() *ServeMux {return new(ServeMux)}
//DefaultServeMux is the  default ServeMux used by serve.
var DefaultServeMux=&defaultServeMux
var defaultServeMux ServeMux
注意 go创建实例的过程中 也可以使用指针方式 即
type Server struct{}
server :=Server{}
和下面的一样都可以创建Server实例
var DefaultServer Server
var server =&DefalutServer










*/
