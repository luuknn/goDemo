package main

import (
	"fmt"
)

/*
Web开发中一个很重要的议题就是如何做好用户的整个浏览过程的控制，因为HTTP协议是无状态的，所以用户的每一次请求都是无状态的，我们不知道在整个Web操作过程中哪些连接与该用户有关，我们应该如何来解决这个问题呢？Web里面经典的解决方案是cookie和session，
cookie机制是一种客户端机制 把用户数据保存在客户端 而session机制 是一种服务器端的机制
服务器使用一种类似于散列表的结构来保存信息 每一个网站访客都会被分配给一个唯一的标志符 即 sessionID
它的存放形式 无非两种 要么经过url传递 要么保存在客户端的cookies里  当然 你也可以将session保存到数据库里 这样会更安全 但效率方面会有所下降

6.1小节里面介绍session机制和cookie机制的关系和区别
6.2讲解Go语言 如何来实现session 里面讲实现一个简易的session管理器
6.3小节讲解 如何防止session被劫持的情况 如何有效的保护session 我们知道session其实可以存储在任何地方
6.3小节里面实现的session是存储在内存中的  但是如果我们的应用进一步扩展了  要实现应用的session共享 那么
我们可以把session存储在数据库中  memcache 或者redis 6.4小节 将详细讲解 如何实现这些功能


session和cookie是网站浏览中较为常见的两个概念 也是比较难以辨析的两个概念
但是它们在浏览器需要认证的服务页面 以及页面统计中却相当关键 我们先来了解一下 session和cookie怎么来的 考虑这样一个问题
如何抓取一个访问受限的网页 比如新浪微博好友的主页 个人微博页面等
显然  通过 浏览器 我们可以手动输入用户名和密码来访问页面 而所谓的抓取 就是用程序来模拟完成同样的工作
因此 我们需要了解登录过程中到底发生了什么

当用户来到微博登录页面 输入用户名和密码之后点击登录后 浏览器将认证信息post给远端的服务器
服务器执行验证逻辑 如果验证通过 浏览器会跳转到登录用户的微博首页 在登录成功后 服务器如何验证我们对
其他受限页面的访问呢  因为http协议是无状态的  所以很显然 服务器不可能知道我们已经在上一次http请求中通过了验证
当然 最简单解决方案 就是所有的请求里面都带上用户名和密码 这样虽然可行 但大大加重了 服务器的负担
对于每个request都需要到数据库验证 也大大降低了用户体验 既然直接在请求中带上用户名和密码不可行
那么只有在服务器或者客户端保存一些类似的可以代表身份的信息了  所以就有了cookie和session

cookie 简而言之 就是在本地计算机 保存一些用户操作的历史信息 当然包括登录信息 并在用户再次访问该站点时 浏览器通过http协议将本地
cookie内容发送给服务器 从而完成验证 或继续上一步 操作
session 简而言之就是在服务器上保存用户操作的信息历史
服务器使用session id来标识session
session id 由服务器负责产生 保证随机性与唯一性
相当于一个随机秘钥 避免在握手或者传输中 暴露用户的真实密码  但该方式 仍然需要将发送请求的客户端与session进行对应
所以可以借助cookie机制来获取客户端的标识 session id 也可以通过get方式将id提交给服务器

cookie是由浏览器维持的 存储在客户端的一小段文本信息 伴随着用户请求和页面在web服务器和浏览器之间传递
用户每次访问站点 web应用程序都可以读取cookie包含的信息 浏览器设置里面有cookie隐私数据选项 打开它 可以看到很多已访问网站的cookies
cookie是有时间限制的  根据生命期不同分为两种 会话cookie和持久cookie
如果不设置 过期时间 则表示这个cookie生命周期为创建到浏览器关闭为止 只要关闭浏览器窗口
cookie就消失了 这种生命期为浏览器会话期的cookie 被称为 会话cookie  cookie一般不保存在硬盘上而是保存在内存里
如果设置了过期时间 setMaxAge(606024)浏览器就会把cookie保存在硬盘上 关闭后再次打开浏览器 这些cookie依然有效
直到超过设定的过期时间 存储在硬盘上的cookie可以在不同的浏览器进程间共享 比如两个ie窗口
而对于保存在内存cookie 不用浏览器有不同的处理方式

GO设置cookie
Go语言中通过net/http包中的setcookie来设置
http.SetCookie(w ResponseWriter,cookie *Cookie)
w 表示需要写入的response cookie是一个struct 让我们来看一下 cookie对象是怎么样的
type Cookie struct{
Name string
Value string
Path string
Domain string
Expires time.Time
RawExpires string
// MaxAge=0 means no 'Max-Age' attribute specified.
// MaxAge<0 means delete cookie now, equivalently 'Max-Age: 0'
// MaxAge>0 means Max-Age attribute present and given in seconds
MaxAge int
Secure bool
HttpOnly bool
Raw string
Unparsed []string//Raw text of unparsed attribute-value pairs
}
我们来看一个例子 如何设置cookie
expiration:=*time.LocalTime()
expiration.Year+=1
cookie :=http.Cookie{Name:"username",Value:"jianling",Expires:expiration}
http.SetCookie(w,&cookie)

Go 读取cookie
上面的例子演示了如何设置cookie数据 我们这里来演示一下如何读取cookie
cookie,_:=r.Cookie("username")
fmt.Fprint(w,cookie)

还有另外一种读取方式
for _,cookie :=range r.Cookies(){
fmt.Fprint(w,cookie.Name)
可以看到通过request获取cookie非常方便

}

session 指有始有终的一系列动作 消息 比如打电话 是从拿起电话拨号到挂断电话这个中间的一系列过程 可以称为一个session
然而当session一词 与网络协议相关联时  它又往往隐含了 面向连接 或保持状态 两个含义

session在web开发环境下语义又有了新的扩展 它的含义 是指一类用户用来在客户端和服务器端之间保持状态的解决方案
有时session也用来指 这种解决方案的存储结构
ession机制是一种服务器端的机制，服务器使用一种类似于散列表的结构(也可能就是使用散列表)来保存息。

但程序需要为某个客户端的请求创建一个session的时候，服务器首先检查这个客户端的请求里是否包含了一个session标识－称为session id，如果已经包含一个session id则说明以前已经为此客户创建过session，服务器就按照session id把这个session检索出来使用(如果检索不到，可能会新建一个，这种情况可能出现在服务端已经删除了该用户对应的session对象，但用户人为地在请求的URL后面附加上一个JSESSION的参数)。如果客户请求不包含session id，则为此客户创建一个session并且同时生成一个与此session相关联的session id，这个session id将在本次响应中返回给客户端保存。

session机制本身并不复杂，然而其实现和配置上的灵活性却使得具体情况复杂多变。这也要求我们不能把仅仅某一次的经验或者某一个浏览器，服务器的经验当作普遍适用的。
如上文所述，session和cookie的目的相同，都是为了克服http协议无状态的缺陷，但完成的方法不同。session通过cookie，在客户端保存session id，而将用户的其他会话消息保存在服务端的session对象中，与此相对的，cookie需要将所有信息都保存在客户端。因此cookie存在着一定的安全隐患，例如本地cookie中保存的用户名密码被破译，或cookie被其他网站收集（例如：1. appA主动设置域B cookie，让域B cookie获取；2. XSS，在appA上通过javascript获取document.cookie，并传递给自己的appB）。

通过上面的一些简单介绍我们了解了cookie和session的一些基础知识，知道他们之间的联系和区别，做web开发之前，有必要将一些必要知识了解清楚，才不会在用到时捉襟见肘，或是在调bug时候如无头苍蝇乱转。接下来的几小节我们将详细介绍session相关的知识。
*/
