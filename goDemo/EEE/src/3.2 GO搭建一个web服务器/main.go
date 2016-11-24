package main

import (
	"fmt"
)

//前面小节已经介绍了web是基于http协议的一个服务 Go语言里面提供了一个完整的net/http包 通过http包
//就可以很方便的搭建一个可以运行的web服务 同时使用这个包 能很简单地对web路由 静态文件 模板 cookie等数据进行设置和操作
//http包建立web服务器
package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
)

func sayhelloName(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()       //解析参数，默认是不会解析的
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

func main() {
	http.HandleFunc("/", sayhelloName)       //设置访问的路由
	err := http.ListenAndServe(":9090", nil) //设置监听的端口
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}

/*我们看到上面的代码，要编写一个web服务器很简单，只要调用http包的两个函数就可以了。

如果你以前是PHP程序员，那你也许就会问，我们的nginx、apache服务器不需要吗？Go就是不需要这些，因为他直接就监听tcp端口了，做了nginx做的事情，然后sayhelloName这个其实就是我们写的逻辑函数了，跟php里面的控制层（controller）函数类似。

如果你以前是python程序员，那么你一定听说过tornado，这个代码和他是不是很像，对，没错，go就是拥有类似python这样动态语言的特性，写web应用很方便。

如果你以前是ruby程序员，会发现和ROR的/script/server启动有点类似。
我们看到Go通过简单的几行代码就已经运行起来一个web服务了，而且这个Web服务内部有支持高并发的特性，我将会在接下来的两个小节里面详细的讲解一下go是如何实现Web高并发的。
*/
