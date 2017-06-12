package main

import ()

/*
请求的结构
HTTP的交互以请求和响应的应答模式 go的请求我们早就见过了 handler函数的第二个参数http.Requests 其结构为
type Request struct{
Method string

URL *url.URL
Proto string//"HTTP/1.0"
ProtoMajor int //1
ProtoMinor int //0
Header Header
Body io.ReadCloser
ContentLength int64
TransferEncoding[] string
Close bool
Host string
Form url.Values
PostForm url.Values
MultipartForm *multipart.Form
...
ctx context.Context
}
从request结构可以看到 http请求的基本信息都囊括了 对于请求而言主要关注一下请求的URL Method Header Body这些结构
*/

/*
URL
HTTP的url请求格式为 scheme://[userinfo@]host/path[?query][#fragment] go提供了一个URL结构
用来映射HTTP的请求URL

type URL struct{
Scheme string
Opaque string
User *Userinfo
Host string
Path string
RawQuery string
Fragment string
}
URL的格式比较明确 其实更好的名词应该是URI 统一资源定位 url中比较重要的是查询字符串query
通常作为get请求的参数 query是一些使用& 符号分割的key1=value1&key2=value2键值对 由于url编码是ASSIC码
因此query需要进行 urlencode go可以通过 request.URI.RawQuery 读取query
func indexHandler(w http.ResponseWriter, r *http.Request) {
    info := fmt.Sprintln("URL", r.URL, "HOST", r.Host, "Method", r.Method, "RequestURL", r.RequestURI, "RawQuery", r.URL.RawQuery)
    fmt.Fprintln(w, info)
}
*/

/*
header
header 也是HTTP中重要的组成部分 Request结构中就有Header结构 Header本质上是一个
map(map[string] []string) 将http协议的header的key-value进行映射成一个图

Host:example.com
accept-encoding:gzip,deflate
accept_language:en-us
fOO :Bar
foo:two

Header=map[string][]string{
"Accept-Encoding":{"gzip,default"},
"Accept-Language":{"en-us"},
"Foo":{"Bar","Two"},

}
hrader 中的字段包含了很多通信的设置 很多请求 都需要指定Content-Type
func indexHandler(w http.ResponseWriter,r *http.Request){
info :=fmt.Sprintln(r.Header.Get("Content-Type"))
fmt.Fprintln(w,info)
}

☁  ~  curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'name=vanyar&age=27' "http://127.0.0.1:8000?lang=zh&version=1.1.0"
application/x-www-form-urlencoded

Golang提供了不少打印函数 基本上分为三类三种 即Print Println和Printf
Print 比较简单 打印输出到标准输出流 Println也一样 不同在于多打印一个换行符 至于Printf则是打印格式化字符串
三个方法都返回打印的bytes数 Sprint Sprintln Sprintf 则返回打印的字符串 不会输出到标准流中
Fprint Fprintf Fprintln 则把输出的结果打印输出到io.Writer接口中
http中则是http.ResponseWriter这个对象中返回打印的byte数
*/

/*
Body
http中数据通信 主要通过body传输 go把body封装成Request的Body 它是一个ReadCloser接口
接口方法Reader也是一个接口 后者有一个 Read(p []byte)(n int,err error) 方法 因此bidy可以通过byte数组请求的数据

func indexHandler(w http.ResponseWriter,r *http.Request){
info :=fmt.Sprintln(r.Header.Get("Content-Type"))
len :=r.ContentLentgth
body:=make([]byte,len)
r.Body.Read(body)
fmt.Fprintln(w,info,string(body))
}
☁  ~  curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'name=vanyar&age=27' "http://127.0.0.1:8000?lang=zh&version=1.1.0"
application/x-www-form-urlencoded
 name=vanyar&age=27
*/

/*
表单
form
go提供了ParseForm方法用来解析表单提供的数据 即content-type为x-www-form-urlencode的数据
func indexHandler(w http.ResponseWriter,r *http.Request){
contentType:=fmt.Sprintln(r.Header.Get("Content-Type"))
r.ParseForm()
formData:=fmt.Sprintf("%#v",r.Form)
fmt.Fprintf(w,contentType,formData)
}
☁  ~  curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'name=vanyar&age=27' "http://127.0.0.1:8000?lang=zh&version=1.1.0"
application/x-www-form-urlencoded
%!(EXTRA string=url.Values{"name":[]string{"vanyar"}, "age":[]string{"27"}, "lang":[]string{"zh"}, "version":[]string{"1.1.0"}})%

用来读取数据的结构和方法大致有下面几个
fmt.Println(r.Form["lang"])
fmt.Println(r.PostForm["lang"])
fmt.Println(r.FormValue("lang"))
fmt.Println(r.PostFormValue("lang"))
其中r.Form和r.PostForm必须在调用parseForm之后 才会有数据 否则是空数组 后两 无需调用ParseForm的调用就能读取数据
☁  ~  curl -X POST -H "Content-Type: application/x-www-form-urlencoded" -d 'name=vanyar&age=27&lang=en' "http://127.0.0.1:8000?lang=zh&version=1.1.0"
application/x-www-form-urlencoded
%!(EXTRA string=url.Values{"version":[]string{"1.1.0"}, "name":[]string{"vanyar"}, "age":[]string{"27"}, "lang":[]string{"en", "zh"}})%

此时可以看到lang参数 不仅url的query提供了 post的body也提供了 go默认以body数据优先 两者的数据都有 并不会覆盖
如果不想读取url参数 调用PostForm或者PostFormValue读取字段的值即可
r.PostForm["lang"][0]
r.PostFormValue["lang"]
对于form-data的格式的数据 ParseForm的方法只会解析url中的参数 并不会解析body中的参数
☁  ~  curl -X POST -H "Content-Type: multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW" -F "name=vanyar" -F "age=27" -F "lang=en" "http://127.0.0.1:8000?lang=zh&version=1.1.0"
multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW; boundary=------------------------5f87d5bfa764488d
%!(EXTRA string=url.Values{"lang":[]string{"zh"}, "version":[]string{"1.1.0"}}
)%
因此当请求的cont-type 为form-data的术后ParseForm则需要改成MutipartForm
否则r.Form是读取不到body的内容 只能读取到query string中的内容


MutipartForm
ParseMutipartForm 方法需要提供一个读取数据长度的参数 然后使用同样的方法读取表单数据 MutipartForm只会读取body的数据
不会读取url的query数据

func indexHandler(w http.ResponseWriter, r *http.Request) {
    r.ParseMultipartForm(1024)
    fmt.Println(r.Form["lang"])
    fmt.Println(r.PostForm["lang"])
    fmt.Println(r.FormValue("lang"))
    fmt.Println(r.PostFormValue("lang"))
    fmt.Println(r.MultipartForm.Value["lang"])
    fmt.Fprintln(w, r.MultipartForm.Value)
}
可以看到请求之后返回 map[name:[vanyar] age:[27] lang:[en]]。即r.MultipartForm.Value并没有url中的参数。

总结一下，读取urlencode的编码方式，只需要ParseForm即可，读取form-data编码需要使用ParseMultipartForm方法。如果参数中既有url，
又有body，From和FromValue方法都能读取。而带Post前缀的方法，只能读取body的数据内容。其中MultipartForm的数据通过r.MultipartForm.Value访问得到。
*/

/*
文件上传
form-data格式用的最多的方式就是在图片上传的时候r.MultipartForm.Value 是post的body字段数据
r.MultipartForm.File则包含了图片数据
func indexHandler(w http.ResponseWriter,r *http.Request){
r.ParseMultipartForm(1024)
fileheader:=r.Multipart.File["file"][0]
fmt.Println(fileHeader)
file,err:=fileHeader.Open()
if err ==nil{
data,err :=ioutil.ReadAll(file)
if err==nil{
fmt.Println(len(date))
fmt.Fprintln(w,string(data))
}

}
fmt.Println(err)
}
发出请求之后 可以看见返回了图片当然go 提供了更好的工具函数 r.FormFile 直接读取上传文件的数据 而不需要再使用ParseMultipartForm方法

file,_,err:=r.FormFile("file")
if err == nil{
data,err:=ioutil.ReadAll(file)
if err==nil{
fmt.Println(len(data))
fmt.Fprintln(w,string(date))
}
}
fmt.Println(err)
这种情况 只适用于除了文件字段 没有其他字段的时候 如果仍然需要读取lang参数 还是需要加上ParseMultipartForm调用 读取到了上传文件 接下来就是很普通的写文件的io操作了
*/

/*
Response

请求和响应是http的孪生兄弟，不仅它们的报文格式类似，相关的处理和构造也类似。go构造响应的结构是ResponseWriter接口。
type ResponseWriter interface {
    Header() Header
    Write([]byte) (int, error)
    WriteHeader(int)
}
里面的方法也很简单，Header方法返回一个header的map结构。WriteHeader则会返回响应的状态码。Write返回给客户端的数据。

我们已经使用了fmt.Fprintln 方法，直接向w写入响应的数据。也可以调用Write方法返回的字符。
func indexHandler(w http.ResponseWriter, r *http.Request) {
    str := `<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>`
    w.Write([]byte(str))
}
☁  ~  curl -i http://127.0.0.1:8000/
HTTP/1.1 200 OK
Date: Wed, 07 Dec 2016 09:13:04 GMT
Content-Length: 95
Content-Type: text/html; charset=utf-8
<html>
<head><title>Go Web Programming</title></head>
<body><h1>Hello World</h1></body>
</html>%                                                                                  ☁  ~
go根据返回的字符，自动修改成了text/html的Content-Type格式。返回数据自定义通常需要修改header相关信息。
func indexHandler(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(501)
    fmt.Fprintln(w, "No such service, try next door")
}
☁  ~  curl -i http://127.0.0.1:8000/
HTTP/1.1 501 Not Implemented
Date: Wed, 07 Dec 2016 09:14:58 GMT
Content-Length: 31
Content-Type: text/plain; charset=utf-8
No such service, try next door
重定向

重定向的功能可以更加设置header的location和http状态码实现。
func indexHandler(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Location", "https://google.com")
    w.WriteHeader(302)
}
☁  ~  curl -i http://127.0.0.1:8000/
HTTP/1.1 302 Found
Location: https://google.com
Date: Wed, 07 Dec 2016 09:20:19 GMT
Content-Length: 31
Content-Type: text/plain; charset=utf-8
重定向是常用的功能，因此go也提供了工具方法，http.Redirect(w, r, "https://google.com", http.StatusFound)。

与请求的Header结构一样，w.Header也有几个方法用来设置headers
func (h Header) Add(key, value string) {
    textproto.MIMEHeader(h).Add(key, value)
}
func (h Header) Set(key, value string) {
    textproto.MIMEHeader(h).Set(key, value)
}
func (h MIMEHeader) Add(key, value string) {
    key = CanonicalMIMEHeaderKey(key)
    h[key] = append(h[key], value)
}
func (h MIMEHeader) Set(key, value string) {
    h[CanonicalMIMEHeaderKey(key)] = []string{value}
}
Set和Add方法都可以设置headers，对于已经存在的key，Add会追加一个值value的数组中，，set则是直接替换value的值。即 append和赋值的差别。

*/
/*
总结
对于web应用程式 处理请求 返回响应是基本的内容 golang很好的封装了Request和ResponseWriter给开发者
无论请求还是响应 都是针对url header 和body相关数据的处理 也是http协议的基本内容
除了body的数据处理 有时候也需要处理header中的数据  一个常见的例子就是处理cookie 这将会在cookie的话题中讨论



*/
