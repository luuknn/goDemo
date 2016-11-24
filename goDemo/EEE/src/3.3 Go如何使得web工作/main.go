package main

import (
	"fmt"
)

//前面小节介绍了如何通过Go搭建一个web服务 我们可以看到 简单应用一个net/http包就可以方便的搭建起来了
//那么go在底层 到底怎么做的呢 万变不离其宗 GO的web服务工作也离不开web工作方式
//web工作方式的几个概念
//Request :用户请求信息 用来解析用户的请求信息 包括 post get cookie url等信息
//Response: 服务器需要反馈给客户端的信息
//Conn :用户的每次请求连接
//Handler:处理请求和生成返回信息的处理逻辑

//http包执行流程
/*
1	创建Listen Socket 监听指定的端口 等待客户端请求的到来
2	Listen Socket接受客户端请求 得到Client Socket接下来 Client Socket与客户端通信
3	处理客户端的请求 首先从Client Socket读取HTTP请求的协议头 如果是POST方法 还可能要读取 客户端提交的数据 然后交给相应的handler处理请求 handler 处理完毕准备好客户端需要的数据 通过 client socket写给客户端
*/
//这整个过程里面 我们只要了解 清楚 下面三个问题 也就知道Go是如何让Web运行起来了
//如何监听端口
//如何接收客户端请求
//如何分配handler

//前面一个小节的代码里 我们可以看到Go是通过 一个函数ListenAndServer来处理这些事情的 这个底层 其实是这样处理的
//初始化一个server对象 然后调用了 net.Listen("tcp",addr)也就是底层用TCP协议搭建了一个服务 然后监控我们设置的端口
//下面代码 来自GO的http包源码 通过下面代码 我们可以看到整个的http处理过程
func (srv *Server) Serve(l net.Listener) error {
	defer l.Close()
	var tempDelay time.Duration // how long to sleep on accept failure
	for {
		rw, e := l.Accept()
		if e != nil {
			if ne, ok := e.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				log.Printf("http: Accept error: %v; retrying in %v", e, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return e
		}
		tempDelay = 0
		if srv.ReadTimeout != 0 {
			rw.SetReadDeadline(time.Now().Add(srv.ReadTimeout))
		}
		if srv.WriteTimeout != 0 {
			rw.SetWriteDeadline(time.Now().Add(srv.WriteTimeout))
		}
		c, err := srv.newConn(rw)
		if err != nil {
			continue
		}
		go c.serve()
	}
	panic("not reached")
}

//监控之后如何接收客户端的请求呢 上面代码执行监控端口之后 调用了srv.Serve(nrt.Listener)函数
//这个函数就是处理接收客户端的请求信息 在这个函数里面 起了一个for{} 首先通过listener 接收请求 其次创建一个Conn
//最后单独开了一个goroutine 把这个请求的数据当作参数扔给conn去服务 go c.server()这个就是高并发体现了
//用户的每次请求都是在一个新的goroutine去服务 相互不影响
/*
那么如何具体分配到相应的函数处理请求呢 conn首先会解析request :c.readRequest()
然后获取相应的handler:handler:=c.server.Handler 也就是我们刚才调用函数
listenAndServer时候的第二个参数 我们前面例子传递的是nil 也就是为空
那么默认获取handler =DefaultServerMux 那么这个变量用来做什么的呢
对 这个变量 就是一个路由器 它用来匹配url跳转sayhelloName嘛
这个作用就是注册了请求/的路由规则 当请求uri为 /路由就会 转到函数sayhelloName
DefaultServerMux会调用ServeHTTP方法 这个方法 内部其实就是调用sayhelloName本身
最后通过写入response的信息反馈到客户端

ListenAndServe 监听端口addr
		net.Listen("tcp",addr)
	接收到请求并转交

	接收到用户请求并创建连接conn
srv.Serve(I net.Listener)
	进入for{}
			rw:=I.Accept()
			c:=srv.NewConn()
			go c.serve()
	处理连接
	处理连接*conn.serve()

	分析请求 c.readRequest()
	取出并分析 resp req

	映射url与handlefunc()
	handler =DefaultServeMux
	handler.ServeHTTP(resp,req)
*/
