package main

import (
	"errors"
	"fmt"
)

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("math: square root of negative number")
	}
	// implementation
	return 110 * f, nil
}
func main() {
	f, err := Sqrt(2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(f)
	}
}

//错误处理 调试 和测试
//我们经常看到很多程序员大部分的编程时间都花费在检查bug和修复bug上 无论你是在编写修改代码还是重构系统
//几乎都是花费大量的时间在进行故障排除 和测试 外界都觉得我们程序员是设计师  能够把一个系统从无到有
//是一项很伟大的工作 而且是相当有趣的工作 但事实上 我们每天都是徘徊在排错 调试 测试之间
//当然如果你有良好的习惯和技术方案来直面这些问题 那么你就有可能将排错时间降到最少 而尽可能的将时间花费在更有价值的事情上

//但遗憾的是 很多程序员不愿意在错误处理 调试 和测试能力上下功夫 导致后面应用上线之后查找错误 定位问题 花费更多的时间
//所以我们在设计应用之前 就做好错误处理规划 测试用例等  那么将来修改代码 升级系统都将变得简单
//开发web应用过程中 错误自然难免  那么如何更好的找到错误原因 解决问题呢
//11.1小节将介绍Go语言中如何处理错误 如何设计自己的包 函数的错误处理
//11.2小节将介绍如何使用GDB来调试我们的程序 动态运行情况下的各种变量信息 运行情况的监控和调试
//11.3小节将对Go语言中的单元测试进行深入的探讨 并示例如何来编写单元测试 Go的单元测试规则规范如何定义
//以保证升级修改运行相应的测试代码 就可以进行最小化的测试
//长期以来 培养良好的调试 测试习惯 一直是很多程序员逃避的事情 所以现在你不要再逃避了 就从你现在的项目开发
//从学习GoWb开发开始养成良好的习惯/

//11.1错误处理
//Go原因主要的设计准则是 简洁明白 简洁是指语法和C类似 相当的简单 明白是指任何语句都是很明显的
//不含有任何隐含的东西 在错误处理方案的设计中 也贯彻了这一思想 我们知道在C语言里面是通过返回-1或者null之类的信息来表示错误
//但是对于使用者来说 不查看相应的API说明文档 根本搞不清楚 这个返回值究竟代表什么意思 比如返回0是成功
//还是失败 而Go定义了一个叫做error的类型 来显示的表达错误 在使用时 通过把返回的error变量和nil的比较
//来判定操作是否成功 例如os.Open函数在打开文件失败时将发返回一个部位nil的变量
//func Open(name string)(file *File,err,error)
//下面这个例子通过 调用os Open 打开一个文件如果 错误 那么就会调用 log.Fatal来输出错误的信息
//f,err :=os.Open("filename.txt")
//if err!=nil{
//log.Fatal(err)
//}
//类似于 os。Open函数 标准包中所有可能出现错误的API都会返回一个error变量 以方便错误处理 这个小节将详细的
//介绍error类型的设计 和讨论web应用中如何更好地处理error

//Error 类型
//error类型是一个接口类型 这是它的定义
//type error interface{
//
//Error() string
//}
//error 是一个内置的接口类型 我们可以在 builtin包下面找到相应的定义 而我们在很多内部包里面用到的
//error都是errors包下面实现的私有结构errorString
//type errorString struct{
//s string
//}
//
//func (e *errorString) Error() string{
//return e.s
//}
//你可以通过errors.New把一个字符串转化为errorString 来得到一个满足erroe接口的对象
// New returns an error that formats as the given text.
//func New(text string) error {
//	return &errorString{text}
//}
/*package main

import (
	"errors"
	"fmt"
)

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, errors.New("math: square root of negative number")
	}
	// implementation
	return 110 * f, nil
}
func main() {
	f, err := Sqrt(2)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Print(f)
	}
}
*/
//自定义Error
//
//通过上面的介绍我们知道error是一个interface，所以在实现自己的包的时候，通过定义实现此接口的结构，我们就可以实现自己的错误定义，请看来自Json包的示例：
//
//type SyntaxError struct {
//	msg    string // 错误描述
//	Offset int64  // 错误发生的位置
//}
//
//func (e *SyntaxError) Error() string { return e.msg }
//Offset字段在调用Error的时候不会被打印，但是我们可以通过类型断言获取错误类型，然后可以打印相应的错误信息，请看下面的例子:
//
//if err := dec.Decode(&val); err != nil {
//	if serr, ok := err.(*json.SyntaxError); ok {
//		line, col := findLine(f, serr.Offset)
//		return fmt.Errorf("%s:%d:%d: %v", f.Name(), line, col, err)
//	}
//	return err
//}
//需要注意的是，函数返回自定义错误时，返回值推荐设置为error类型，而非自定义错误类型，特别需要注意的是不应预声明自定义错误类型的变量。例如：
//
//func Decode() *SyntaxError { // 错误，将可能导致上层调用者err!=nil的判断永远为true。
//    var err *SyntaxError     // 预声明错误变量
//    if 出错条件 {
//        err = &SyntaxError{}
//    }
//    return err               // 错误，err永远等于非nil，导致上层调用者err!=nil的判断始终为true
//}
上面例子简单的演示了如何自定义Error类型。但是如果我们还需要更复杂的错误处理呢？此时，我们来参考一下net包采用的方法：
//
//package net
//
//type Error interface {
//    error
//    Timeout() bool   // Is the error a timeout?
//    Temporary() bool // Is the error temporary?
//}
//在调用的地方，通过类型断言err是不是net.Error,来细化错误的处理，例如下面的例子，如果一个网络发生临时性错误，那么将会sleep 1秒之后重试：
//
//if nerr, ok := err.(net.Error); ok && nerr.Temporary() {
//	time.Sleep(1e9)
//	continue
//}
//if err != nil {
//	log.Fatal(err)
//}
//错误处理
//
//Go在错误处理上采用了与C类似的检查返回值的方式，而不是其他多数主流语言采用的异常方式，这造成了代码编写上的一个很大的缺点:错误处理代码的冗余，对于这种情况是我们通过复用检测函数来减少类似的代码。
//
//请看下面这个例子代码：
//
//func init() {
//	http.HandleFunc("/view", viewRecord)
//}
//
//func viewRecord(w http.ResponseWriter, r *http.Request) {
//	c := appengine.NewContext(r)
//	key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
//	record := new(Record)
//	if err := datastore.Get(c, key, record); err != nil {
//		http.Error(w, err.Error(), 500)
//		return
//	}
//	if err := viewTemplate.Execute(w, record); err != nil {
//		http.Error(w, err.Error(), 500)
//	}
//}
//上面的例子中获取数据和模板展示调用时都有检测错误，当有错误发生时，调用了统一的处理函数http.Error，返回给客户端500错误码，并显示相应的错误数据。但是当越来越多的HandleFunc加入之后，这样的错误处理逻辑代码就会越来越多，其实我们可以通过自定义路由器来缩减代码(实现的思路可以参考第三章的HTTP详解)。
//
//type appHandler func(http.ResponseWriter, *http.Request) error
//
//func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if err := fn(w, r); err != nil {
//		http.Error(w, err.Error(), 500)
//	}
//}
//上面我们定义了自定义的路由器，然后我们可以通过如下方式来注册函数：
//
//func init() {
//	http.Handle("/view", appHandler(viewRecord))
//}
//当请求/view的时候我们的逻辑处理可以变成如下代码，和第一种实现方式相比较已经简单了很多。
//
//func viewRecord(w http.ResponseWriter, r *http.Request) error {
//	c := appengine.NewContext(r)
//	key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
//	record := new(Record)
//	if err := datastore.Get(c, key, record); err != nil {
//		return err
//	}
//	return viewTemplate.Execute(w, record)
//}
//上面的例子错误处理的时候所有的错误返回给用户的都是500错误码，然后打印出来相应的错误代码，其实我们可以把这个错误信息定义的更加友好，调试的时候也方便定位问题，我们可以自定义返回的错误类型：
//
//type appError struct {
//	Error   error
//	Message string
//	Code    int
//}
//这样我们的自定义路由器可以改成如下方式：
//
//type appHandler func(http.ResponseWriter, *http.Request) *appError
//
//func (fn appHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
//	if e := fn(w, r); e != nil { // e is *appError, not os.Error.
//		c := appengine.NewContext(r)
//		c.Errorf("%v", e.Error)
//		http.Error(w, e.Message, e.Code)
//	}
//}
//这样修改完自定义错误之后，我们的逻辑处理可以改成如下方式：
//
//func viewRecord(w http.ResponseWriter, r *http.Request) *appError {
//	c := appengine.NewContext(r)
//	key := datastore.NewKey(c, "Record", r.FormValue("id"), 0, nil)
//	record := new(Record)
//	if err := datastore.Get(c, key, record); err != nil {
//		return &appError{err, "Record not found", 404}
//	}
//	if err := viewTemplate.Execute(w, record); err != nil {
//		return &appError{err, "Can't display record", 500}
//	}
//	return nil
//}
//如上所示，在我们访问view的时候可以根据不同的情况获取不同的错误码和错误信息，虽然这个和第一个版本的代码量差不多，但是这个显示的错误更加明显，提示的错误信息更加友好，扩展性也比第一个更好。

//总结 在程序设计中 容错是相当重要的一部分工作 在Go中它是通过错误处理来实现的 error虽然只是一个接口 但是其变化却可以有很多 
//我们可以根据自己的需求来实现不同的处理 最后介绍的错误处理方案 希望能给大家在如何设计更好的web错误处理方案上带来一点思路















