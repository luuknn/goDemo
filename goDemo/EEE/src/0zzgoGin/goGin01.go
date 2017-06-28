package main

//import (
//	"gopkg.in/gin-gonic/gin.v1"
//	"net/http"
//)
//
//func main() {
//
//	router := gin.Default()
//	router.POST("/formpost", func(c *gin.Context) {
//		message := c.PostForm("message")
//		nick := c.DefaultPostForm("nick", "anonymous")
//		c.JSON(http.StatusOK, gin.H{
//			"status": gin.H{
//				"status_code": http.StatusOK,
//				"status":      "ok",
//			},
//			"message": message,
//			"nick":    nick,
//		})
//	})
//
//	router.Run(":8000")
//}

/*
golang
微框架gin  框架一直是敏捷开发中的利器  能让开发者很快的上手并做出应用
Gin是一个golang的微框架 封装比较优雅 API友好 源码注释比较明确  具有快速灵活 容错方便等特点 其实对于 golang
而言 web框架的依赖远比Python java之类要小  自身的 net/http足够简单 性能也非常不错 框架 更像是一些常用函数 或者工具的集合
借助框架开发 不仅可以省去很多常用的封装来的时间 也有助于团队的编码风格 和形成规范
HelloWorld
使用Gin实现helloworld 非常简单  创建一个router 然后使用其Run方法
import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func main() {

	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello World!")
	})
	router.Run(":8000")
}
简单的几行代码 就能实现一个web服务 使用gin的Default方法创建一个路由handler
然后通过HTTP方法绑定路由规则和路由函数  不同于net/http库的路由函数 gin进行了封装
把request和response都封装到了gin.Context的上下文环境 最后是启动路由的Run方法 监听端口 麻雀虽小
五脏俱全  当然 处了 GET方法 gin也支持POST PUT DELETE OPTION等常用的restful方法
*/
/*
restful路由
gin的路由来自httprouter库  因此httprouter具有的功能 gin也具有 不过gin不支持路由正则表达式
router.GET("/user/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s", name)
	})
冒号: 加上一个参数名 组成路由参数 可以使用 c.Params的方法读取值
router.GET("/user/:name/*action", func(c *gin.Context) {
		name := c.Param("name")
		action := c.Param("action")
		message := name + " is " + action
		c.String(http.StatusOK, message)
	})//http://localhost:8000/user/hello/gqy
*/

/*
query string参数 与body参数
web提供的服务通常是 client与server的交互 其中客户端向服务器发送请求 除了路由参数 其他的参数无非两种
查询字符串 query string 和报文体body参数  所谓query string  即 路由用? 后面连接的key1=value1,key2=value2
等形式的参数 当然这个key-value是经过url encode 编码
对于参数的处理 经常会出现参数不存在的情况 对于是否提供默认值 gin也考虑了 并且给出了一个优雅的方案
router.GET("/welcome", func(c *gin.Context) {
		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname")
		c.String(http.StatusOK, "hello %s %s ", firstname, lastname)
	})//http://localhost:8000/welcome?firstname=s&lastname=jl
使用c.DefaultQuery方法读取参数 其中当参数不存在的时候 提供一个默认值 使用Query方法读取正常参数
当参数不存在的时候  返回空字串

body
http的报文体传输数据就比query string稍微复杂一点 常见的格式就有四种  例如
application/json application/x-www-form-urlencoded application/xml multipart/form-data
后面一个主要用于图片上传  json格式很好理解 urlencode 其实也不难 无非就是把querystring的内容
放到了body体里  同样需要urlencode 默认情况下 c.PostFORM解析的是x-www-form-urlcode 或者form-data的参数
func main(){
	router := gin.Default()
	router.POST("/form_post", func(c *gin.Context) {
		message := c.PostForm("message")
		nick := c.DefaultPostForm("nick", "anonymous")
		c.JSON(http.StatusOK, gin.H{
			"status":  gin.H{
				"status_code": http.StatusOK,
				"status":      "ok",
			},
			"message": message,
			"nick":    nick,
		})
	})
}
与get处理query参数一样 post方法也提供了处理默认参数的情况 同理 如果参数不存在 将会得到空字串
前面我们使用c.String返回响应 顾名思义返回string类型 content-type 是plain或者text
调用c.JSON则返回json数据 其中gin.H封装了生成json的方式  是一个强大的工具
使用golang可以像动态语言一样写字面量的json 对于嵌套json的实现 嵌套gin.H即可
发送给服务端 并不是post方法才行 put方法一样可以
func main(){
	router := gin.Default()

	router.PUT("/post", func(c *gin.Context) {
		id := c.Query("id")
		page := c.DefaultQuery("page", "0")
		name := c.PostForm("name")
		message := c.PostForm("message")
		fmt.Printf("id: %s; page: %s; name: %s; message: %s \n", id, page, name, message)
		c.JSON(http.StatusOK, gin.H{
			"status_code": http.StatusOK,
		})
	})
}
*/

/*
上面介绍了基本的发送数据 其中multipart/form-data 专用于文件上传 gin文件上传也很方便
和原生net/http方法类似 不同在于 gin把原生的request封装到c.request中了
func main(){
	router := gin.Default()

	router.POST("/upload", func(c *gin.Context) {
		name := c.PostForm("name")
		fmt.Println(name)
		file, header, err := c.Request.FormFile("upload")
		if err != nil {
			c.String(http.StatusBadRequest, "Bad request")
			return
		}
		filename := header.Filename
		fmt.Println(file, err, filename)
		out, err := os.Create(filename)
		if err != nil {
			log.Fatal(err)
		}
		defer out.Close()
		_, err = io.Copy(out, file)
		if err != nil {
			log.Fatal(err)
		}
		c.String(http.StatusCreated, "upload successful")
	})
	router.Run(":8000")
}
使用c.Request.FormFile解析客户端文件name属性 如果不传文件 则会抛错 因此需要处理这个错误
一种方式是直接返回 然后使用os的操作 把文件数据复制到硬盘上
使用下面的命令可以测试上传 注意upload为c.Request.FormFile指定的参数 其值 必须要是绝对路径
curl -X POST http://127.0.0.1:8000/upload -F "upload=@/Users/ghost/Desktop/pic.jpg" -H "Content-Type: multipart/form-data"

################
上传多个文件
单个文件上传很简单 别以为多个文件就会很麻烦 所谓多个文件 无非就是多一次遍历文件 然后一次copy数据存储即可
下面 只写handler 省略 main函数的初始化路由和监听服务器了
router.POST("/multi/upload", func(c *gin.Context) {
		err := c.Request.ParseMultipartForm(200000)
		if err != nil {
			log.Fatal(err)
		}
		formdata := c.Request.MultipartForm
		files := formdata.File["upload"]
		for i, _ := range files { /
			file, err := files[i].Open()
			defer file.Close()
			if err != nil {
				log.Fatal(err)
			}
			out, err := os.Create(files[i].Filename)
			defer out.Close()
			if err != nil {
				log.Fatal(err)
			}
			_, err = io.Copy(out, file)
			if err != nil {
				log.Fatal(err)
			}
			c.String(http.StatusCreated, "upload successful")
		}
	})
与单个文件上传类似 只不过使用了c.Request.MultipartForm 得到文件句柄 再获取文件数据 然后遍历读取

表单上传
上面我们使用的都是curl上传  实际上 用户上传图片 更多是通过表单 ajax和一些requests的请求完成 下面展示一下
web form表单如何上传
我们首先要写一个表单页面 因此需要引入gin如何render模板 前面我们见识了c.String 和c.JSON
下面就来看看 c.HTML
首先需要定义一个模板文件夹 然后调用c.HTML渲染模板 可以通过gin.H给模板传值 至此 无论是string JSON还是HTML
以及后面的XML和YAML 都可以看到Gin封装的接口简明易用
创建一个文件夹 templates 然后创建html文件 upload.Html
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <title>upload</title>
</head>
<body>
<h3>Single Upload</h3>
<form action="/upload", method="post" enctype="multipart/form-data">
    <input type="text" value="hello gin" />
    <input type="file" name="upload" />
    <input type="submit" value="upload" />
</form>
<h3>Multi Upload</h3>
<form action="/multi/upload", method="post" enctype="multipart/form-data">
    <input type="text" value="hello gin" />
    <input type="file" name="upload" />
    <input type="file" name="upload" />
    <input type="submit" value="upload" />
</form>
</body>
</html>
upload 很简单 没有参数  一个用于单个文件上传 一个用于多个文件上传
router.LoadHTMLGlob("templates/*")
router.GET("/upload", func(c *gin.Context) {
	c.HTML(http.StatusOK, "upload.html", gin.H{})
})
使用LoadHTMLClob定义模板文件路径
*/

/*
参数绑定
我们已经见识了 x-www-form-urlencoded 类型的参数处理  现在越来越多的应用习惯 使用JSON来通信
也就是无论返回的response还是提交的request 其content-type类型都是application/json的格式
而对于一些旧的web表单页面还是x-www-form-urlencoded的形式 这就需要我们的服务器能hold住
这多种content-type的参数了 golang中要处理并非易事 好在有gin  他们的model bind功能非常强大

type User struct {
Username string `form:"username" json:"username" binding:"required"`
Passwd string `form:"passwd" json:"passwd" binding:"required"`
Age int `form:"age" json:"age"`
}
func main(){
router :=gin.Default()
router.POST("/login",func(c *gin.Context))
var user User
var err error
contentType :=c.Request.Header.Get("Content-Type")
switch contentType{
case "application/json":
 err=c.BindJSON(&user)
 case "application/x-www-form-urlencoded":
 err=c.BindWith(&user,bingding.Form)
}
if err!=nil{
fmt.Println(err)
log.Fatal(err)
}
c.JSON(http.StatusOK,gin.H{
"user":user.Username,
"passwd":user.Passwd,
"age":user.Age
})
})
}
先定义一个User模型结构体 然后针对客户端的content-type 一次使 BindJSON 和BindWith方法
可以看到 结构体中 设置了 binding标签字段 (username和passwd) 如果没传会抛出错误 非binding的字段
age 对于客户端 没有传 user结构会用零值填充 对于User结构没有的参数 会自动被忽略
使用json还需要注意一点 json是有数据类型的  因此 对于 {"passwd":"123"} 和 {"passwd":123}是不同的数据类型
解析需要符合 对应的数据类型 否则会出错
当然 gin还提供了更加高级的方法 c.Bind 它会content-type字段推断 是bind表单还是json参数
router.POST("/login",fun(c *gin.Context){
var user User
err:=c.Bind(&user)
if err!=nil{
fmt.Println(err)
log.Fatal(err)
}
c.JSON(http.StatusOK,gin.H{
"username":user.Username,
"passwd":user.passwd,
"age":user.Age,
})
})
*/

/*
多格式渲染
既然请求可以使用不同的content-type 响应也是如此 通常会有html text plain json和xml等
gin提供了很优雅的渲染方法  到目前为止 我们已经见识了c.String c.JSON c.HTML 下面介绍下c.XML
router.GET("/render",func(c *gin.Context){
contentType :=c.DefaultQuery("content-type","json")
if contentType == "json" {
		c.JSON(http.StatusOK, gin.H{
			"user":   "rsj217",
			"passwd": "123",
		})
	} else if contentType == "xml" {
		c.XML(http.StatusOK, gin.H{
			"user":   "rsj217",
			"passwd": "123",
		})
	}
})

重定向
gin对于重定向 的请求 相当简单 调用上下文的Rediret方法
router.GET("/redict/google", func(c *gin.Context) {
	c.Redirect(http.StatusMovedPermanently, "https://google.com")
})

分组路由
v1 := router.Group("/v1")
v1.GET("/login", func(c *gin.Context) {
	c.String(http.StatusOK, "v1 login")
})
v2 := router.Group("/v2")
v2.GET("/login", func(c *gin.Context) {
	c.String(http.StatusOK, "v2 login")
})

middleware中间件
golang的net/http设计的一大特点就是特别容易构建中间件 gin也提供了类似的中间件
需要注意的是中间件 只对注册过的路由函数起作用  对于分组路由 嵌套使用中间件  可以限定中间件的作用范围
中间件分为全局中间件  单个路由中间件 和群组中间件

全局中间件
先定一个中间件函数
 func MiddleWare() gin.HandlerFunc{
 return func(c *gin.Context){
 fmt.Println("before middleware")
 c.Set("request","client_request")
 c.Next()
 fmt.Println("before middleware")
 }
 }
 该函数很简单  只会给c上下文添加一个属性 并赋值  后面的路由处理器 可以根据 被中间件装饰后提取其值需要注意 虽然名为全局
 中间件 只要注册中间件的过程之前设置的路由 将不会受注册的中间件所影响  只有注册了中间件一下代码的路由规则
 才会被中间件装饰
 router.Use(MiddleWare())
{
	router.GET("/middleware", func(c *gin.Context) {
		request := c.MustGet("request").(string)
		req, _ := c.Get("request")
		c.JSON(http.StatusOK, gin.H{
			"middile_request": request,
			"request": req,
		})
	})
}
使用router装饰中间件，然后在/middlerware即可读取request的值，注意在router.Use(MiddleWare())
代码以上的路由函数，将不会有被中间件装饰的效果。

使用花括号包含被装饰的路由函数只是一个代码规范，即使没有被包含在内的路由函数，只要使用router进行路由
，都等于被装饰了。想要区分权限范围，可以使用组返回的对象注册中间件。
如果没有注册就使用MustGet方法读取c的值将会抛错，可以使用Get方法取而代之。
上面的注册装饰方式，会让所有下面所写的代码都默认使用了router的注册过的中间件。

单个路由中间件
当然，gin也提供了针对指定的路由函数进行注册。
router.GET("/before", MiddleWare(), func(c *gin.Context) {
	request := c.MustGet("request").(string)
	c.JSON(http.StatusOK, gin.H{
		"middile_request": request,
	})
})
把上述代码写在 router.Use(Middleware())之前，同样也能看见/before被装饰了中间件。

群组中间件
群组的中间件也类似，只要在对于的群组路由上注册中间件函数即可：
authorized := router.Group("/", MyMiddelware())
// 或者这样用：
authorized := router.Group("/")
authorized.Use(MyMiddelware())
{
    authorized.POST("/login", loginEndpoint)
}
群组可以嵌套，因为中间件也可以根据群组的嵌套规则嵌套。

*/

/*
中间件最大的作用 莫过于用于一些记录log 错误handler 还有就是对部分借口的鉴权 下面就实现一个简易的鉴权中间件

router.GET("/auth/signin", func(c *gin.Context) {
	cookie := &http.Cookie{
		Name:     "session_id",
		Value:    "123",
		Path:     "/",
		HttpOnly: true,
	}
	http.SetCookie(c.Writer, cookie)
	c.String(http.StatusOK, "Login successful")
})
router.GET("/home", AuthMiddleWare(), func(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"data": "home"})
})

登录函数会设置一个session_id的cookie，注意这里需要指定path为/，不然gin会自动设置cookie的path为/auth，
一个特别奇怪的问题。/homne的逻辑很简单，使用中间件AuthMiddleWare注册之后，将会先执行AuthMiddleWare的逻辑，然后才到/home的逻辑。
func AuthMiddleWare() gin.HandlerFunc {
	return func(c *gin.Context) {
		if cookie, err := c.Request.Cookie("session_id"); err == nil {
			value := cookie.Value
			fmt.Println(value)
			if value == "123" {
				c.Next()
				return
			}
		}
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "Unauthorized",
		})
		c.Abort()
		return
	}
}
从上下文的请求中读取cookie，然后校对cookie，如果有问题，则终止请求，
直接返回，这里使用了c.Abort()方法。

异步协程
golang的高并发一大利器就是协程  gin里可以借助协程实现异步任务 因为涉及异步 过程 请求的上下文需要copy到
异步的上下文 并且这个上下文是只读的
router.GET("/sync", func(c *gin.Context) {
	time.Sleep(5 * time.Second)
	log.Println("Done! in path" + c.Request.URL.Path)
})
router.GET("/async", func(c *gin.Context) {
	cCp := c.Copy()
	go func() {
		time.Sleep(5 * time.Second)
		log.Println("Done! in path" + cCp.Request.URL.Path)
	}()
})

在请求的时候sleep 五秒钟  同步的逻辑可以看到服务的进程睡眠了 异步的逻辑则看到响应返回了 然后程序还在后台的协程处理

自定义路由
gin不仅可以使用框架本身的router进行run 也可以配合使用 net/http 本身的功能
fun main(){
 router := gin.Default()
 http.ListenAndServe(":8080", router)
}
或者

func main() {
    router := gin.Default()
    s := &http.Server{
        Addr:           ":8000",
        Handler:        router,
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }
    s.ListenAndServe()
}
当然还有一个优雅的重启和结束进程的方案 后面将会探索使用 supervisor管理golang的进程

总结
Gin是一个轻巧而强大的golang web框架  涉及常见开发的功能 我们都做了简单的介绍











*/
