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
