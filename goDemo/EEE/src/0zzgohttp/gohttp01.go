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

因此DefaultServerMux 的HandleFunc(pattern,handler)方法实际是定义在ServeMux下的

func(mux ServeMux) HandleFunc(pattern string,handler func(ResponseWriter,Request)){
mux.Handle(pattern,HandlerFunc(handler))
}
上述代码中 HanlerFunc是一个函数类型 同时实现了handler接口的servehttp方法 使用handlerfunc类型包装下
路由定义的indexhandler函数 其目的就是为了让这个函数也实现servehttp方法 转变成一个handler处理器(函数)
type HandlerFunc func(RequestWriter,*Request)
func(f HandleFunc)ServeHTTP(w ResponseWriter,r*Request){
f(w,r)
}
一旦这样做了 就意味着我们的indexhandler函数也有了servehttp方法
此外 servemux的handle方法将会对 pattern和handler函数做一个map映射

func(mux *ServeMux) Handle(pattern string,handler Handler){
mux.mu.Lock()
defer mux.mu.Unlock()
if pattern==""{
panic("http:invalid pattern"+pattern)
}
if handler ==nil{
panic("http: nil handler")
if mux.m[pattern].explicit{
panic("http:mutiple registrations for"+pattern)
}
if mux.m==nil{
mux.m=make(map[string]muxEntry)
}
mux.m[pattern] = muxEntry{explicit: true, h: handler, pattern: pattern}

if pattern[0] != '/' {
    mux.hosts = true
}

n := len(pattern)
if n > 0 && pattern[n-1] == '/' && !mux.m[pattern[0:n-1]].explicit {

    path := pattern
    if pattern[0] != '/' {
        path = pattern[strings.Index(pattern, "/"):]
    }
    url := &url.URL{Path: path}
    mux.m[pattern[0:n-1]] = muxEntry{h: RedirectHandler(url.String(), StatusMovedPermanently), pattern: pattern}
}
}
// 由此可见 handle函数的主要目的在于把handle和pattern模式绑定到map[string]muxEntry的map上 其中 muxentry保存了
更多pattern和handler的信息 还记得前面讨论的server结构吗 server的m字段就是map[string]muxEntry这样一个map
此时 pattern和handler的路由注册完成 接下来就是如何开始server的监听 以接收客户端的请求
开启监听
注册好了路由之后 启动web服务还需要开启服务器监听 http的ListenAndServer方法中可以看到创建了一个server对象
并调用了 server对象的同名方法

func ListenAndServe(add string,handler Handler)error{
server:=&Server{Addr:addr,Handler:handler}
return server.ListenAndServe()
}
func(srv Server)ListenAndServe() error{
addr:=srv.Addr
if addr==""{
addr=":http"
}
ln,err:=net.Listen("tcp",addr)
if err!=nil{
return err
}
return srv.Serve(tcpKeepAliveListener{ln.(net.TCPListener)})

}

//server的listenandserver方法中 会初始化监听地址addr同时调用listen方法设置监听 最后将监听的tcp对象传入serve方法
func (srv *Server) Serve(l net.Listener) error{
defer l.Close()
...
baseCtx:=context.Background()
ctx:=context.WithValue(baseCtx,ServerContexKey,srv)
ctx=context.WithValue(ctx,LocalAddrContextKey,l.Addr())
for{
rw,e:=l.Accept()
...
c:=srv.newConn(rw)
c.setState(c.rwc,StateNew)//before Serve can return
go c.serve(ctx)
}
}
//处理请求
监听开启之后 一旦客户端请求到底 go就开启一个协程处理请求 主要逻辑都在serve方法之中
serve方法比较长 其主要职责 就是创建一个上下文对象 然后调用Listener的Accept方法用来 获取连接数据并使用newConn
方法创建连接对象 最后使用goroutine协程的方式处理连接请求 因为每一个连接都开启了一个协程 请求的上下文都不同 同时又保证了go的高并发 serve又是一个长长的方法

 func(c *conn)serve(ctx context.Context){
 c.remoteAddr=c.rwc.RemoteAddr().String()
 defer func(){
 if err:=recover();err!=nil{
 const size=64<<10
 buf:=make([]byte,size)
 buf=buf[:runtime.Stack(buf,false)]
 c.server.logf("httpLpanic serving %v:%v\n%s",c.remoteAddr,err,buf)
 }
 if !c.hijacked(){
 c.close()
 c.setState(c.rwc,StateClosed)
 }
 }()
 
 ...
 
 for{
 w,err:=c.readRequest(ctx)
 if c.r.remain!=c.server.initialReadLimitSize(){
 //if we read any bytes off the wire ,we're active
 c.setState(c.rwc,StateActive)
 }
 ... 
 }
 ...
 
 serverHandler{c.server}.ServeHTTP(w,w.req)
 w.cancelCtx()
 if c.hijacked(){
 return
 }
 w.finishRequest()
 if !w.shouldReuseConnection(){
 if w.requestBodyLimitHit || w.closedRequestBodyEarly(){
 c.closeWriteAndWait()
 }
 return
 }
 c.setState(c.rwc,StateIdle)
 }
//尽管serve很长 里面的结构和逻辑还是很清晰的 使用defer定义函数退出时 连接关闭相关的处理 。 然后就是读取连接的
网络数据 并处理读取完毕时候的状态 接下来就是调用 serveHandler{c.server}.ServeHTTP(w,w.req)方法处理请求了
最后就是请求处理完毕的逻辑 serveHandler是一个重要的结构 它只有一个字段 即serve结构 同时它实现了handler接口方法
ServeHTTP 并在该方法中做了一个重要的事情 初始化multiplexer路由多路复用器 如果server对象 没有指定Handler
则使用DefaultServeMux作为路由Multiplexer并调用初始化Handler 的servehttp方法

type serverHandler struct{
srv *Server
}
func(sh serverHandler)ServeHTTP(rw ResponseWriter,req Request){
handler:=sh.srv.Handler
if handler==nil{
handler=DefaultServeMux
}
if req.RequestURI==""&&req.Method=="OPTIONS"{
handler=globleOptionsHandler{}
}
handler.ServeHTTP(rw,req)
}

这里DefaultServeMux的ServeHTTP方法其实也是定义在ServeMux结构中的，相关代码如下：

func (mux ServeMux) (w ResponseWriter, r Request) {
if r.RequestURI == “*” {
if r.ProtoAtLeast(1, 1) {
w.Header().Set(“Connection”, “close”)
}
w.WriteHeader(StatusBadRequest)
return
}
h, _ := mux.Handler(r)
h.ServeHTTP(w, r)
}

func (mux ServeMux) Handler(r Request) (h Handler, pattern string) {
if r.Method != “CONNECT” {
if p := cleanPath(r.URL.Path); p != r.URL.Path {
_, pattern = mux.handler(r.Host, p)
url := *r.URL
url.Path = p
return RedirectHandler(url.String(), StatusMovedPermanently), pattern
}
}
return mux.handler(r.Host, r.URL.Path)
}

func (mux *ServeMux) handler(host, path string) (h Handler, pattern string) {
mux.mu.RLock()
defer mux.mu.RUnlock()

// Host-specific pattern takes precedence over generic ones
if mux.hosts {
    h, pattern = mux.match(host + path)
}
if h == nil {
    h, pattern = mux.match(path)
}
if h == nil {
    h, pattern = NotFoundHandler(), ""
}
return
}

func (mux *ServeMux) match(path string) (h Handler, pattern string) {
var n = 0
for k, v := range mux.m {
if !pathMatch(k, path) {
continue
}
if h == nil || len(k) > n {
n = len(k)
h = v.h
pattern = v.pattern
}
}
return
}
```

//mux的serveHTTP方法通过调用其handler方法寻找注册到路由上的handler函数 并调用该函数的ServeHTTP方法 本例则是IndexHandler函数

mux的Handler方法对URL简单的处理 然后调用handler方法 后者会创建一个锁 同时调用match方法返回一个handler和pattern

在match方法中 mux的m字段是map[string]muxEntry图 后者存储了pattern和hadler处理器函数 因此通过迭代m寻找出注册路由的
pattern模式与实际url匹配的handler函数并返回

返回的结构一直传递到mux的ServeHTTP方法 接下来调用handler函数的ServeHTTP方法 即IndexHandler函数 然后把response写到http.RequestWriter对象返回给客户端

上述函数运行结束 即serverHandler{c.server}.ServeHTTP(w, w.req) 运行结束 接下来就是对请求处理完毕之后希望和连接断开的相关逻辑

至此 http服务大致介绍完毕 包括注册路由 开启监听 处理连接 路由处理函数
*/

//总结
/*多数的web应用基于HTTP协议 客户端和服务器通过request-response的方式交互 一个server并不可少的两部分莫过于
路由注册和连接处理 go通过一个ServeMux实现了的miltiplexer路由多路复用器来管理路由 同时提供一个Handler接口
提供ServeHTTP用来实现handler处理其函数 后者可以处理实际request并构造response

ServeMux和handler处理器函数的连接桥梁就是handler接口ServeMux的ServeHTTP方法实现了寻找注册路由的handler的函数
并调用该handler的ServeHTTP方法 ServeHTTP方法就是真正处理请求和构造响应的地方


























