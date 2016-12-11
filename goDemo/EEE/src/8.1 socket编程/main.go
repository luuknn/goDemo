package main

import (
	"fmt"
)

//8 web服务
//web服务可以让你在http协议的基础上通过xml或者json来交换信息 如果你想知道上海的天气预报 中国石油的股价
//或者淘宝商家的一个商品信息 你可以编写一段简短的代码 通过抓取这些信息 然后通过 标准的接口放出来
//就如同你调用一个本地函数 并返回一个值

//web服务背后的关键 在于平台的无关性 你可以运行你的服务在linux系统
//可以与其他windows的asp.net 程序交互 同样的 也可以通过接口和运行在FreeBSD上面的JSP无障碍的通信

//目前主流的有如下几种web服务 rest soap

//rest请求是很直观的  因为rest是基于 HTTP协议的一个补充,它的每次请求都是一个HTTP请求
//然后根据不同的method来处理不同的逻辑 很对web开发者都熟悉 HTTP协议 所以学REST是一件比较容易的事情
//所以我们在8.3小节将详细的讲解如何在Go语言中实现REST方式
//SOAP是W3C在跨网络信息传递和远程计算机函数调用方面的一个标准
//但soap非常复杂 其完整的规范篇幅很长 而且内容仍在增加 Go语言是以简单著称
//所以我们不会介绍soap这样复杂的东西 而Go语言提供了一种天生性能很不错 开发起来很方便的PRC机制
//我们将会在8.4小节 详细的介绍如何使用Go语言来实现RPC
//Go 语言是21世纪的C语言  我们追求的是性能 简单 所以我们在 8.1小节里面介绍如何使用socket编程
//很多游戏服务都是采用socket来编写服务端  因为HTTP协议相对而言 比较耗费性能
//让我们看看GO语言如何来socket编程 目前 随着HTML5的发展
//websockets也逐渐成为很多页游公司接下来开发的一些手段 我们将在 8.2讲解GO语言如何编写websockets的代码

//SOCKET编程

//在很多底层网络应用开发者眼里 一切编程都是socket 话虽然有点夸张 但却也几乎如此看
//现在的网络编程几乎都是用socket来编程 你想过这些情景么
//我们每天打开浏览器浏览网页时  浏览器进程怎么和 web服务器进行通信呢
//当你用QQ聊天时  QQ进程 怎么和服务器或者你好有所在的QQ进行通信的呢
//当你打开PPstream 观看视频时  PPstream进程如何与视频服务器进行通信的呢
//如此种种 都是靠socket来进行通信的  可见 socket编程在现代编程中占据多么重要的地位 这一节 我们将介绍Go语言中如何进行socket编程

//什么是socket
//socket起源于unix 而unix基本哲学之一 就是一切皆文件
//都可以用 打开open 读写write read 关闭close模式 来操作
//socket就是该模式的一个实现  网络socket数据传输是一种特殊的I/O socket也是一种文件描述符
//socket也具有一个类似于 打开文件的函数调用 socket()
//该函数返回一个整型的socket描述符 随后的连接建立 数据传输 等操作都是通过该socket实现的

//常用的socket类型有两种 流式socket (SOCK_STREAM) 和数据报式 (SOCK_DGRAM)
//流式是一种面向连接的socket 针对于面向连接的TCP服务应用
//数据报式 socket是一种无连接的socket 对应于 无连接的UDP服务应用

//socket如何通信
//网络中的进程 之间 如何通过 socket通信呢  首要解决的问题是如何标识一个进程
//否则通信无从谈起 在本地可以通过进程PID来唯一标识一个进程 但是在网络中 是行不通的 其实TCP/ip
//协议已经帮我们解决了这个问题 网络层的IP地址可以唯一标识网络中的主机
//而传输层的 协议+端口 可以唯一标识组集中的应用程序(进程)
//利用这三元组 ip地址 协议 端口 就可以标识网络的进程了  网络中 需要互相通信的进程
//就可以利用这个标志 在他们之间进行交互 请看下面这个TCP IP协议结构图

//使用TCP/IP协议的应用程序通常采用应用编程接口 UNXI BSD的套接字（socket）和UNIX System V的TLI（已经被淘汰）
//来实现网络进程之间的通信 就目前而言 几乎所有的应用程序都是采用socket 而现在又是网络时代 网络中进程通信是无处不在
//这就是 为什么说 一切皆socket

//socket基础知识
//通过上面的介绍 我们知道socket有两种 TCPsocket和UDPsocket TCP和UDP是协议 而要确定一个进程的三元组 需要ip地址和端口

//IPv4地址
//
//目前的全球因特网所采用的协议族是TCP/IP协议。IP是TCP/IP协议中网络层的协议，是TCP/IP协议族的核心协议。目前主要采用的IP协议的版本号是4(简称为IPv4)，发展至今已经使用了30多年。
//
//IPv4的地址位数为32位，也就是最多有2的32次方的网络设备可以联到Internet上。近十年来由于互联网的蓬勃发展，IP位址的需求量愈来愈大，使得IP位址的发放愈趋紧张，前一段时间，据报道IPV4的地址已经发放完毕，我们公司目前很多服务器的IP都是一个宝贵的资源。
//
//地址格式类似这样：127.0.0.1 172.122.121.111
//
//IPv6地址
//
//IPv6是下一版本的互联网协议，也可以说是下一代互联网的协议，它是为了解决IPv4在实施过程中遇到的各种问题而被提出的，IPv6采用128位地址长度，几乎可以不受限制地提供地址。按保守方法估算IPv6实际可分配的地址，整个地球的每平方米面积上仍可分配1000多个地址。在IPv6的设计过程中除了一劳永逸地解决了地址短缺问题以外，还考虑了在IPv4中解决不好的其它问题，主要有端到端IP连接、服务质量（QoS）、安全性、多播、移动性、即插即用等。
//
//地址格式类似这样：2002:c0e8:82e7:0:0:0:c0e8:82e7

//Go支持的ip类型
//在Go net包中 定义了很多类型 函数 方法 用来网络编程 其中ip定义如下
//type  IP []byte
//在net包中很多函数来操作ip 但是其中比较有用的也就几个 其中ParseIp(s string) ip函数
//会把一个ipv4 或者ipv6的地址 转化成ip类型 请看下面的例子
package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage:%s ip-addr\n", os.Args[0])
		os.Exit(1)

	}
	name := os.Args[1]
	addr := net.ParseIP(name)
	if addr == nil {
		fmt.Println("Invalid addres")
	} else {
		fmt.Println("The addres is ", addr.String())

	}
	os.Exit(0)
}
//执行之后就会发现只要你输入一个IP地址 就会给出相应的IP格式

//TCP Socket
//当我们知道如何通过网络端口访问一个服务时 那么我们能够做什么呢 作为客户端来说
//我们可以通过向远端某台机器的某个网络端口发送一个请求 然后得到在机器的此端口上监听的服务反馈信息
//作为服务端 我们需要把服务绑定到某个指定的端口 并且在此端口上监听 当有客户端来访问的时候能够读取信息并且写入反馈信息

//在Go语言的net包中 有一个类型TCPConn 这个类型可以用来作为客户端和服务端交互的通道
//他有两个主要的函数
func(c *TCPConn) Write(b[]byte)(n int,err os.Error)
func(c *TCPConn) Read(b []byte)(n int,err os.Error)

//TCPConn可以用在客户端和服务器端来读写数据
//还有我们需要知道一个TCPAddr类型 它表示一个TCP的地址信息 它的定义如下
type TCPAddr struct{
IP IP
Port int
}
//在Go语言中 通过ResolveTCPAddr获取一个 TCPAddr
func ResolveTCPAddr(net,addr string)(*TCPAddr,os.Error)
//net参数是 tcp4 tcp6 tcp中任意一个 分别表示 iPV4only 6only 任意一个
//addr表示域名或者IP地址 例如 www.google.com:80 或者 127.0.0.1:22

//TCP client
//Go语言中通过net包中的DialTCP函数来建立一个TCP连接 并且返回一个TCPConn类型的对象
//当连接建立时 服务器端 也创建一个同类型的对象 此时客户端和服务器端通过各自拥有的TCPConn对象来进行数据交换
//一般而言 客户端 通过TCPConn对象将请求信息 发送到服务器端 读取服务器端的响应信息
//服务器端读取并解析来自客户端的请求 并返回应答信息 这个连接 只有当任一端关闭了连接之后才失效 不然 这连接可以一直在使用 建立连接的函数定义如下
func DialTCP(net string,laddr,raddr *TCPAddr)(c *TCPConn,err os.Error)
//net参数tcp4 tcp6 tcp任意一个 laddr表示本机地址 一般设置为nil raddr 表示 远程服务地址
//接下来我们写一个简单的例子 模拟一个基于HTTP协议的客户端去连接一个WEB服务器
//我们要写一个简单的http请求头 格式类型如下
"HEAD / HTTP/1.0\r\n\r\n"
//从服务端接收到的响应信息格式可能如下
HTTP/1.0 200 OK
Last-Modified: Thu, 25 Mar 2010 17:51:10 GMT
Content-Length: 18074
Connection: close
Date: Sat, 28 Aug 2016 00:43:48 GMT
Server: lighttpd/1.4.23

//我们的客户端代码如下所示
package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"os"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %s host:port ", os.Args[0])
		os.Exit(1)
	}
	service := os.Args[1]
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	conn, err := net.DialTCP("tcp", nil, tcpAddr)
	checkError(err)
	_, err = conn.Write([]byte("HEAD / HTTP/1.0\r\n\r\n"))
	checkError(err)
	result, err := ioutil.ReadAll(conn)
	checkError(err)
	fmt.Println(string(result))
	os.Exit(0)
}
func checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
//通过上面的代码 我们可以看出  首先 程序将用户输入的参数service传入 net.ResolveTCPAddr
//获取一个tcpAddr 然后把 tcpAddr传入 DialTCP后创建了一个连接conn
//通过conn来发送请求信息  最后通过 ioutil.ReadAll 从conn读取全部是文本 也就是 服务端响应的反馈信息

//TCP SERVER
//上面 我们编写了一个TCP客户端 也通过net包 创建了一个服务器端程序
//在服务器端我们需要绑定 服务到指定的非激活 端口 并监听此端口
func  ListenTCP(net string,laddr *TCPAddr)(l *TCPListener,err os.Error)
func (l *TCPListener) Accept()(c conn,err os.Error)

//下面我们实现一个简单的时间同步服务 监听端口7777
package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

func main() {
    service := ":7777"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        daytime := time.Now().String()
        conn.Write([]byte(daytime)) // don't care about return value
        conn.Close()                // we're finished with this client
    }
}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
//上面的服务跑起来后 它将会 一直在那里等待 直到有新的客户端 请求到达
//上面代码有个缺点 执行的时候是单任务的 不能同时接收多个请求
//那么该如何改造 以使他支持并发呢  Go里面有一个goroutine 机制 请看改造后的代码
package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

func main() {
    service := ":1200"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    defer conn.Close()
    daytime := time.Now().String()
    conn.Write([]byte(daytime)) // don't care about return value
    // we're finished with this client
}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
//通过把业务处理分离到函数 handlerClient 我们就可以 进一步实现多并发执行了
//看上去是不是很帅 增加 go关键字 就实现了服务端的多并发 从这个小例子也可以看出goroutine的强大之处
//有的朋友 可能要问 这个服务端没有处理客户端实际请求的内容 如果我们需要通过从客户端发送不同的请求来获取不同的时间格式
//而且需要一个长连接 该怎么做呢 请看:
package main

import (
    "fmt"
    "net"
    "os"
    "time"
    "strconv"
)

func main() {
    service := ":1200"
    tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
    checkError(err)
    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)
    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        go handleClient(conn)
    }
}

func handleClient(conn net.Conn) {
    conn.SetReadDeadline(time.Now().Add(2 * time.Minute)) // set 2 minutes timeout
    request := make([]byte, 128) // set maxium request length to 128KB to prevent flood attack
    defer conn.Close()  // close connection before exit
    for {
        read_len, err := conn.Read(request)

        if err != nil {
            fmt.Println(err)
            break
        }

        if read_len == 0 {
            break // connection already closed by client
        } else if string(request) == "timestamp" {
            daytime := strconv.FormatInt(time.Now().Unix(), 10)
            conn.Write([]byte(daytime))
        } else {
            daytime := time.Now().String()
            conn.Write([]byte(daytime)) 
        }

        request = make([]byte, 128) // clear last read content
    }
}

func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
        os.Exit(1)
    }
}
//上面这个例子中 我们用conn.Read()不断读取客户端发来的请求
//由于我们需要保持与客户端的长连接 所以不能在读取完一次请求后就关闭连接
//由于 conn.SetReadDeadline()设置了超时 当一定时间内客户端无请求发送 conn便会自动关闭
//下面的for循环即会因为连接已关闭而跳出
//需要注意的是，request在创建时需要指定一个最大长度以防止flood attack；每次读取到请求处理完毕后，需要清理request，因为conn.Read()会将新读取到的内容append到原内容之后。

//控制TCP连接
//TCP 有很多连接控制函数  我们平常 用到比较多的有如下几个函数
func DialTimeout(net,addr string,timeout time.Duration)(conn,error)
//设置 建立连接的超时时间 客户端与服务端都适用 当超过设置时间时 连接自动关闭
func (c *TCPConn) SetReadDeadline(t time.Time) error
func (c*TCPConn) SetWriteDeadline(t time.Time) error
//设置读取 写入一个连接的超时时间 当超过设置时间时 连接自动关闭
func (c *TCPConn) SetKeepAlive(keepalive bool) os.Error
//设置客户端是否和服务器端保持长连接 可以降低建立TCP连接时的握手开销 对于一些 需要频繁交换数据的应用场景 比较适用

//UDP Socket
//Go语言包中处理UDP Socket和TCP Socket不同的地方就是在服务器端处理多个客户端请求数据包的方式不同,UDP缺少了对客户端连接请求的Accept函数。其他基本几乎一模一样，只有TCP换成了UDP而已。UDP的几个主要函数如下所示：
func ResolveUDPAddr(net, addr string) (*UDPAddr, os.Error)
func DialUDP(net string, laddr, raddr *UDPAddr) (c *UDPConn, err os.Error)
func ListenUDP(net string, laddr *UDPAddr) (c *UDPConn, err os.Error)
func (c *UDPConn) ReadFromUDP(b []byte) (n int, addr *UDPAddr, err os.Error
func (c *UDPConn) WriteToUDP(b []byte, addr *UDPAddr) (n int, err os.Error)
//一个UDP的客户端代码如下所示,我们可以看到不同的就是TCP换成了UDP而已：
//package main

import (
    "fmt"
    "net"
    "os"
)

func main() {
    if len(os.Args) != 2 {
        fmt.Fprintf(os.Stderr, "Usage: %s host:port", os.Args[0])
        os.Exit(1)
    }
    service := os.Args[1]
    udpAddr, err := net.ResolveUDPAddr("udp4", service)
    checkError(err)
    conn, err := net.DialUDP("udp", nil, udpAddr)
    checkError(err)
    _, err = conn.Write([]byte("anything"))
    checkError(err)
    var buf [512]byte
    n, err := conn.Read(buf[0:])
    checkError(err)
    fmt.Println(string(buf[0:n]))
    os.Exit(0)
}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
        os.Exit(1)
    }
}
//我们来看一下UDP服务器端如何来处理：
package main

import (
    "fmt"
    "net"
    "os"
    "time"
)

func main() {
    service := ":1200"
    udpAddr, err := net.ResolveUDPAddr("udp4", service)
    checkError(err)
    conn, err := net.ListenUDP("udp", udpAddr)
    checkError(err)
    for {
        handleClient(conn)
    }
}
func handleClient(conn *net.UDPConn) {
    var buf [512]byte
    _, addr, err := conn.ReadFromUDP(buf[0:])
    if err != nil {
        return
    }
    daytime := time.Now().String()
    conn.WriteToUDP([]byte(daytime), addr)
}
func checkError(err error) {
    if err != nil {
        fmt.Fprintf(os.Stderr, "Fatal error ", err.Error())
        os.Exit(1)
    }
}

//总结 通过对TCP和UDP socket编程的描述和实现 可见go已经很好的支持了socket编程
//而且使用起来相当的方便 go提供了很多函数
//通过这些函数可以很容易就编写出 高性能的socket应用


























