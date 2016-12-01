package mian

import (
	"fmt"
	"sync"
)
//通过上一小节 我们知道session是在服务器端实现的一种用户和服务器之间认证的解决方案
//目前Go标准包没有为session提供任何支持 这小节我们将会自己动手来实现go版本的session管理和创建4

//session创建过程
//session的基本原理是由服务器为每一个会话维护一份信息数据  客户端和服务端依靠一个全局唯一的标识来访问这份数据
//已达到交互的目的  当用户访问web应用时  服务器端会随需要创建session这个过程概括分为三个步骤
//生成全局唯一的标识符 sessionid
//开辟数据存储空间 一般会在内存中创建 相应的数据结构  但这种情况下 系统一旦掉电
//所有会话数据就会丢失 如果是电子商务类网站 这将造成严重的后果  所以为了解决这类问题 你可以将会话数据写到文件里 或者存储在数据库中
//当然这样会增加I/O开销 但是它可以实现某种程度的session持久化 也更有利于 session的共享
//session的全局唯一标识符发送给客户端
//以上三个步骤中 最关键的是如何发送这个session的唯一标识这一步上 考虑到http协议的定义 数据无非
//可以放到请求行 头域或者body里 所以一般来说会有两种常用的方式  cookie和 URL重写
//1 Cookie服务端通过设置set-cookie头就可以将session的标识符传送到客户端 而客户端此后的每一次请求都会带上这个标识符
//另外一般包含session信息的cookie会将失效时间设置为0 会话cookie 即浏览器进程有效时间
//至于浏览器怎么处理这个0 每个浏览器都有自己的方案 但差别都不会太大
//2 URL重写  所谓url重写 就是在返回给用户的页面里 所有的url后面追加session标识符
//这样 用户在收到响应之后 无论点击响应页面里的哪个链接或提交表单 都会自动带上 session标识符
//从而就实现了 会话的保持  虽然这种做法比较麻烦 但是 如果客户端禁用了cookie的话 这种方案会是首选

//Go实现session管理
//通过上面session的创建过程的讲解 读者应该对 session有了一个大体的认识  但是具体到动态页面技术里面 
//又是怎么实现session的呢  下面我们将结合session的生命周期 lifecycle 来实现go语言版本的session管理
//session管理设计
//我们知道session管理涉及到如下几个因素
//全局session管理器 保证sessionid的全局唯一性 为每个客户关联一个session
//session的存储(可以存储到内存 文件 数据库里) session过期处理
//接下来  我将讲解一下我关于session管理的整个设计思路以及相应的go代码
//定义一个全局session管理器
type Managger struct{
cookieName string//private cookiename
lock sync.Mutex//protects session
provider Provider
maxlifetime int64
}
func NewManager(provideName,cookieName string,naxlifetime int64){*Manager,error}{
	provider,ok:=provides[provideName]
	if !ok{
	return nil,fmt.Errorf("session:unknown provide %q(forgetten import?)", provideName)
	}
return &Manager{provider:provider,cookoeName:cookieName,maxlifetime:maxlifetime},nil
}
//Go 实现整个的流程应该是这样的  在main包中创建一个全局的session管理器
var globalSessions *session.Manage
//然后在init函数中初始化
func init(){
globalSessions=NewManager("memory","gosessionid",3600)

}
//我们知道session是保存在服务器端的数据 它可以以任何方式存储 比如存储在内存 数据库 或者文件中 因此
//我们抽象出 一个Provider接口 用以表征session管理器底层的存储结构
type Provider interface{
SessionInit(sid string)(Session ,error)//实现session初始化 操作成功返回此新的session变量
SessionRead(sid string)(Session,error)//返回sid所代表的session变量 如果不存在 那么将以sid作为参数调用init函数创建并返回一个新的session变量
SessionDestory(sid string) error//用来销毁sid对应的session变量
SessionGC(maxLifeTime int64)//根据maxlifttime 来删除过期的数据
}
//那么session接口需要实现什么样的功能呢  有过web开发经验的读者应该知道 对session 的处理基本就是设置值 读取值 删除值 以及获取sessionid这四个操作 所以我们的session接口也实现这四个操作
type Session interface{
Set(key,value interface{}) error //set session value
Get(key interface{})interface{} //get session value
Delete(key interface{}) error //delete session value
SessionId() string //back current sessionID
}
//以上设计思路来源于database/sql/driver 先定义好接口 然后具体的存储session的结构实现相应的接口并注册后 相应功能这样就可以实现了 下面是用来随意注册session的结构Register函数的实现
var providers=make(map[string]Provider)
// Register makes a session provide available by the provided name.
// If Register is called twice with the same name or if driver is nil,
// it panics.
func Register(name string,provider Provider){
if provider ==nil{
panic ("session:Register provide is nil")
}
if _,dup:=provides[name];dup{
panic("session:register called twice for provide"+name)
}
provides[name]=provider
}
//全局唯一的sessionid
//sessionid是用来访问web应用的每一个用户 因此必须保证它是全局唯一的 GUID 下面代码展示了如何满足这一需求
func (manager *Manager) sessionId() string{
b:=make([byte],32)
if _,err:=io.ReadFull(rand.Reader,b);err!=nil{
return ""
}
return base64.URLEncoding.EncodeToString(b)

}

//session的创建
//我们需要为每一个来访用户分配或获取与他相关联的session 以便后面根据session信息来验证操作
//sessionStart 这个函数就是用来检测是否已经有某个session与当前来访用户发生了关联 如果没有创建
func (manager *Manager) SessionStart(w http.RespondeWriter,r *http.request)(session Session){
manager.lock.Lock()
defer manager.lock.Unlock()
cookie,err :=r.Cookie(manager.cookieName)
if err!=nil||cookie.value==""{
sid:=manager.provider.SessionInit(sid)
cookie :=http.Cookie{Name: manager.cookieName, Value: url.QueryEscape(sid), Path: "/", HttpOnly: true, MaxAge: int(manager.maxlifetime)}
http.SetCookie(w,&cookie)
}else{
sid,_ :=url.QueryUnescape(cookie.Value)
session,_=manager.provide.SessionRead(sid)

}
return
}
我们用前面login操作来演示session的运用：

func login(w http.ResponseWriter, r *http.Request) {
    sess := globalSessions.SessionStart(w, r)
    r.ParseForm()
    if r.Method == "GET" {
        t, _ := template.ParseFiles("login.gtpl")
        w.Header().Set("Content-Type", "text/html")
        t.Execute(w, sess.Get("username"))
    } else {
        sess.Set("username", r.Form["username"])
        http.Redirect(w, r, "/", 302)
    }
}
操作值：设置、读取和删除

SessionStart函数返回的是一个满足Session接口的变量，那么我们该如何用他来对session数据进行操作呢？

上面的例子中的代码session.Get("uid")已经展示了基本的读取数据的操作，现在我们再来看一下详细的操作:

func count(w http.ResponseWriter, r *http.Request) {
    sess := globalSessions.SessionStart(w, r)
    createtime := sess.Get("createtime")
    if createtime == nil {
        sess.Set("createtime", time.Now().Unix())
    } else if (createtime.(int64) + 360) < (time.Now().Unix()) {
        globalSessions.SessionDestroy(w, r)
        sess = globalSessions.SessionStart(w, r)
    }
    ct := sess.Get("countnum")
    if ct == nil {
        sess.Set("countnum", 1)
    } else {
        sess.Set("countnum", (ct.(int) + 1))
    }
    t, _ := template.ParseFiles("count.gtpl")
    w.Header().Set("Content-Type", "text/html")
    t.Execute(w, sess.Get("countnum"))
}
通过上面的例子可以看到，Session的操作和操作key/value数据库类似:Set、Get、Delete等操作

因为Session有过期的概念，所以我们定义了GC操作，当访问过期时间满足GC的触发条件后将会引起GC，但是当我们进行了任意一个session操作，都会对Session实体进行更新，都会触发对最后访问时间的修改，这样当GC的时候就不会误删除还在使用的Session实体。

session重置

我们知道，Web应用中有用户退出这个操作，那么当用户退出应用的时候，我们需要对该用户的session数据进行销毁操作，上面的代码已经演示了如何使用session重置操作，下面这个函数就是实现了这个功能：

//Destroy sessionid
func (manager *Manager) SessionDestroy(w http.ResponseWriter, r *http.Request){
    cookie, err := r.Cookie(manager.cookieName)
    if err != nil || cookie.Value == "" {
        return
    } else {
        manager.lock.Lock()
        defer manager.lock.Unlock()
        manager.provider.SessionDestroy(cookie.Value)
        expiration := time.Now()
        cookie := http.Cookie{Name: manager.cookieName, Path: "/", HttpOnly: true, Expires: expiration, MaxAge: -1}
        http.SetCookie(w, &cookie)
    }
}
session销毁

我们来看一下Session管理器如何来管理销毁,只要我们在Main启动的时候启动：

func init() {
    go globalSessions.GC()
}

func (manager *Manager) GC() {
    manager.lock.Lock()
    defer manager.lock.Unlock()
    manager.provider.SessionGC(manager.maxlifetime)
    time.AfterFunc(time.Duration(manager.maxlifetime), func() { manager.GC() })
}
我们可以看到GC充分利用了time包中的定时器功能，当超时maxLifeTime之后调用GC函数，这样就可以保证maxLifeTime时间内的session都是可用的，类似的方案也可以用于统计在线用户数之类的。

总结

至此 我们实现了一个用来在Web应用中全局管理Session的SessionManager，定义了用来提供Session存储实现Provider的接口,下一小节，我们将会通过接口定义来实现一些Provider,供大家参考学习。





































