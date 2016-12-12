package main

import (
	"fmt"
)

//websoclet 是html5的重要特性 它实现了基于浏览器的远程socket
//它使浏览器和服务器可以进行全双工通信 许多浏览器都已对此做了支持

//在websocket出现之前 为了实现即时通信 采用的技术都是 轮询
//即在特定的时间间隔内 由浏览器对服务器发出HTTP request 服务器在
//收集到请求后 返回最新的数据给浏览器刷新 "轮询"使得浏览器需要对服务器不断发送请求
//这样会占用 大量带宽

//websocket 采用了一些特殊的报文头 使得 浏览器和服务器只需要做一个握手的动作 就可以在浏览器 和服务器之间 建立一条连接通信
//它解决了web实时化的问题 相比传统的http有如下好处
//一个web客户端只建立一个TCP连接
//websocket服务端可以推送push 数据到web客户端
//有更加轻量级的头 减少数据传送量

//websocket url的起始输入是ws:// 或是 wss://在ssl上
//一个带有特定报文头的HTTP握手被发送到了服务器端 接着在服务器端 或是客户端就可以通过
//JavaScript 来使用 某种 套接口 这一套接口可以被用来通过事件句柄 异步接收数据

//websocket的协议颇为简单 在第一次handshake通过以后 连接便建立成功
//其后的通讯数据都是以 \x00开头 \xFF结尾
//在客户端  这个是透明的  websocket组件会自动将原始数据 掐头去尾

//浏览器发出websocket连接请求 然后 服务器发出回应 然后连接建立成功 这个过程通常称为握手 handshakinh
/*在请求中的"Sec-WebSocket-Key"是随机的，对于整天跟编码打交到的程序员，一眼就可以看出来：这个是一个经过base64编码后的数据。服务器端接收到这个请求之后需要把这个字符串连接上一个固定的字符串：

258EAFA5-E914-47DA-95CA-C5AB0DC85B11
即：f7cb4ezEAl6C3wRaU6JORA==连接上那一串固定字符串，生成一个这样的字符串：

f7cb4ezEAl6C3wRaU6JORA==258EAFA5-E914-47DA-95CA-C5AB0DC85B11
对该字符串先用 sha1安全散列算法计算出二进制的值，然后用base64对其进行编码，即可以得到握手后的字符串：

rE91AJhfC+6JdVcVXOGJEADEJdQ=
将之作为响应头Sec-WebSocket-Accept的值反馈给客户端。*/

//GO实现websocket
<html>
<head></head>
<body>
    <script type="text/javascript">
        var sock = null;
        var wsuri = "ws://127.0.0.1:1234";

        window.onload = function() {

            console.log("onload");

            sock = new WebSocket(wsuri);

            sock.onopen = function() {
                console.log("connected to " + wsuri);
            }

            sock.onclose = function(e) {
                console.log("connection closed (" + e.code + ")");
            }

            sock.onmessage = function(e) {
                console.log("message received: " + e.data);
            }
        };

        function send() {
            var msg = document.getElementById('message').value;
            sock.send(msg);
        };
    </script>
    <h1>WebSocket Echo Test</h1>
    <form>
        <p>
            Message: <input id="message" type="text" value="Hello, world!">
        </p>
    </form>
    <button onclick="send();">Send Message</button>
</body>
</html>
//websocket分為客戶端和服務端 客户端通过websocket将信息发送给服务端
//服务端收到信息之后 主动push信息到客户端 然后客户端将输出其 收到的信息 客户端的代码 如上所示
//可以看到客户端js 很容易的就通过 websocket函数建立了一个与服务器 的连接sock
//当握手成功后 会触发 websocket 对象的onopen事件
//告诉客户端连接已经建立成功 客户端一共绑定了 四个事件
//onopen 建立连接后 触发
//onmessage 收到消息后触发
//onerror 发生错误时 触发
//onclose 关闭连接时 触发

//我们的服务器端的实现如下
package main

import (
	"fmt"
	"log"
	"net/http"
	"websocket"
)

func Echo(ws *websocket.Conn) {
	var err error

	for {
		var reply string

		if err = websocket.Message.Receive(ws, &reply); err != nil {
			fmt.Println("Can't receive")
			break
		}

		fmt.Println("Received back from client: " + reply)

		msg := "Received:  " + reply
		fmt.Println("Sending to client: " + msg)

		if err = websocket.Message.Send(ws, msg); err != nil {
			fmt.Println("Can't send")
			break
		}
	}
}

func main() {
	http.Handle("/", websocket.Handler(Echo))

	if err := http.ListenAndServe(":1234", nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
//当客户端 将用户输入的信息send之后 服务器端通过receive 接收到了相应的信息 然后通过send发送了 应答信息

//通过上面的例子 我们看到 客户端和服务器端实现了websocket 非常的方便 GO的源码net分支中已经实现了这个协议 可以直接拿来用
//目前随着 HTML5的发展 websocket会是web开发的一个重点



















