package main

import (
	"fmt"
)

//RPC
//前面几个小节 我们介绍了如何基于socket和HTTP来编写 网络应用
//通过学习 我们了解了 socket和http采用的是 类似 信息交换 模式
//即客户端发送一条消息到服务端 然后 服务器返回一定的信息 以表示 响应
//客户端和服务端之间 约定了 交互信息的格式
//以便双方 都能够 解析交互所产生的信息
//但是很多独立的应用 并没有采用这种模式 而是采用类似常规的函数调用 的方式来完成想要的功能

//RPC就是想实现函数调用模式的网络化 客户端就像 调用本地函数一样 然后客户端把这些参数打包之后
//通过网络传递 到服务器端 服务器端 解包到处理过程中执行 然后执行的结果反馈给客户端

//RPC remote procedure call protocol远程过程调用协议
//是一种通过网络从远程计算机程序请求服务
//而不需要了解底层网络技术的协议 它假定某些传输协议的存在
//如TCP 或UDP 以便为通信程序 之间 携带信息数据
//通过它 可以使 函数调用模式网络化 在OSI网络通信模型中 RPC跨越了传输层和应用层
//RPC使得开发包括网络分布式多程序在内的应用程序更加容易

//RPC工作原理
//运行时 一次客户端对服务器的RPC调用 其内部操作大致有 如下十步
//1调用客户端句柄 执行传送参数
//2调用本地系统内核发送网络消息
//3消息传送到远程主机
//4服务器句柄得到消息 并取得参数
//5执行远程过程
//6执行的过程将结果返回服务器句柄
//7服务器句柄 返回结果 调用远程系统内核
//8消息传回本地主机
//9客户句柄又内核接收消息
//10客户接收句柄返回的数据

//GO RPC
//GO标准包中 已经提供了对RPC的支持 而且支持三个级别的RPC TCP HTTP JSONRPC
//但Go的rpc包是独一无二的rpc 它和传统的rpc系统不同 它只支持Go开发的服务器和客户端之间的交互 因为在内部 他们采用了gob来编码4

//函数必须是导出的 首字母大写
//必须有两个导出类型的参数
//第一个参数是接收的参数 第二个是返回给客户端的参数 第二个参数必须是指针类型的
//函数还要有一个返回值error

//举个例子 正确的RPC函数格式如下
func (t *T) MethodName(argType T1, replyType *T2) error

//T T1 T2 类型必须被 encoding/gob包编解码
//任何的RPC都需要 通过网络来传递数据 GO RPC可以利用http
//和tcp来传递数据 利用http的好处是 可以直接复用net/http里面的一些函数 详细的例子请看下面的实现

//HTTP RPC
//http 的服务端代码实现如下
package main

import (
    "errors"
    "fmt"
    "net/http"
    "net/rpc"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}

func main() {

    arith := new(Arith)
    rpc.Register(arith)
    rpc.HandleHTTP()

    err := http.ListenAndServe(":1234", nil)
    if err != nil {
        fmt.Println(err.Error())
    }
}
//通过上面的例子可以看到 我们注册了一个Arith的RPC服务 然后通过rpc.HandleHTTP
//函数把该服务注册到HTTP协议上 然后我们就可以利用 http的方式来传递数据了
//请看下面的客户端代码
package main

import (
    "fmt"
    "log"
    "net/rpc"
    "os"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "server")
        os.Exit(1)
    }
    serverAddress := os.Args[1]

    client, err := rpc.DialHTTP("tcp", serverAddress+":1234")
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // Synchronous call
    args := Args{17, 8}
    var reply int
    err = client.Call("Arith.Multiply", args, &reply)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

    var quot Quotient
    err = client.Call("Arith.Divide", args, &quot)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}
我们把上面的服务端和客户端的代码分别编译，然后先把服务端开启，然后开启客户端，输入代码，就会输出如下信息：

$ ./http_c localhost
Arith: 17*8=136
Arith: 17/8=2 remainder 1
//通过上面的调用 可以看到参数和返回值 是我们定义的struct类型  在服务端 我们把它们当做调用函数的参数的类型
//在客户端作为 client.call的第2 3 两个参数的类型 客户端最重要的就是这个call函数 它由三个参数
//第一个调用的函数的名字 第二个是要传递的参数 第三个是要返回的参数值 注意是指针类型
//通过上面的代码例子我们可以发现 使用GO的rpc实现相当的简单方便

//TCP RPC
//上面我们实现了基于HTTP协议的RPC 接下来我们要实现基于TCP协议的RPC 服务端的实现代码如下所示
package main

import (
    "errors"
    "fmt"
    "net"
    "net/rpc"
    "os"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}

func main() {

    arith := new(Arith)
    rpc.Register(arith)

    tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        rpc.ServeConn(conn)
    }

}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
//上面这个代码 和http服务器相比 不同之处在于 此处我们采用了TCP协议 然后需要我们自己控制连接 当有客户端连接上来后 我们需要把这个连接交给rpc来处理
//如果你留心了 你会发现 他是一个阻塞型的单用户程序 如果想要实现多并发 那么可以使用goroutine来实现 我们前面走在socket小节的时候
//已经介绍过如何处理goroutine 下面展现了TCP实现的RPC客户端
package main

import (
    "fmt"
    "log"
    "net/rpc"
    "os"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "server:port")
        os.Exit(1)
    }
    service := os.Args[1]

    client, err := rpc.Dial("tcp", service)
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // Synchronous call
    args := Args{17, 8}
    var reply int
    err = client.Call("Arith.Multiply", args, &reply)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

    var quot Quotient
    err = client.Call("Arith.Divide", args, &quot)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}
//这个客户端代码和http客户端代码对比 唯一的区别一个是DialHTTP 一个是Dial(tcp)其他处理一模一样

//JSON RPC
//
//JSON RPC是数据编码采用了JSON，而不是gob编码，其他和上面介绍的RPC概念一模一样，下面我们来演示一下，如何使用Go提供的json-rpc标准包，请看服务端代码的实现：

package main

import (
    "errors"
    "fmt"
    "net"
    "net/rpc"
    "net/rpc/jsonrpc"
    "os"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
    *reply = args.A * args.B
    return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
    if args.B == 0 {
        return errors.New("divide by zero")
    }
    quo.Quo = args.A / args.B
    quo.Rem = args.A % args.B
    return nil
}

func main() {

    arith := new(Arith)
    rpc.Register(arith)

    tcpAddr, err := net.ResolveTCPAddr("tcp", ":1234")
    checkError(err)

    listener, err := net.ListenTCP("tcp", tcpAddr)
    checkError(err)

    for {
        conn, err := listener.Accept()
        if err != nil {
            continue
        }
        jsonrpc.ServeConn(conn)
    }

}

func checkError(err error) {
    if err != nil {
        fmt.Println("Fatal error ", err.Error())
        os.Exit(1)
    }
}
//通过示例我们可以看出 json-rpc是基于TCP协议实现的，目前它还不支持HTTP方式。
//
//请看客户端的实现代码：

package main

import (
    "fmt"
    "log"
    "net/rpc/jsonrpc"
    "os"
)

type Args struct {
    A, B int
}

type Quotient struct {
    Quo, Rem int
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage: ", os.Args[0], "server:port")
        log.Fatal(1)
    }
    service := os.Args[1]

    client, err := jsonrpc.Dial("tcp", service)
    if err != nil {
        log.Fatal("dialing:", err)
    }
    // Synchronous call
    args := Args{17, 8}
    var reply int
    err = client.Call("Arith.Multiply", args, &reply)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d*%d=%d\n", args.A, args.B, reply)

    var quot Quotient
    err = client.Call("Arith.Divide", args, &quot)
    if err != nil {
        log.Fatal("arith error:", err)
    }
    fmt.Printf("Arith: %d/%d=%d remainder %d\n", args.A, args.B, quot.Quo, quot.Rem)

}

//总结
//Go已经提供了对RPC的良好支持 通过上面HTTP TCP JSON RPC的实现
//我们就可以很方便的开发很多分布式的web应用
//我想作为读者的你已经领会到了这一点 但遗憾的是目前Go尚未提供对SOQP RPC的支持
//欣慰的是 现在已经有第三方的开源实现了






















