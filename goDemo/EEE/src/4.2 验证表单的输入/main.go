package main

import (
	"fmt"
)

/*
开发web的一个原则就是 不能信任用户输入的任何信息  所以验证和过滤用户输入的信息就变得非常重要
我们平常编写web应用主要有两方面数据验证 一个是在页面端的js验证 一个是在服务器端的验证
这小节 我们讲解的是 如何在 服务器端验证

必填字段
你想要确保从一个表单元素中得到一个值  例如前面小节里面的用户名 我们如何处理呢 Go有一个内置函数len可以获取字符串的长度 这样 我们就可以通过len来获取数据的长度 例如
if len(r.Form("username")[0]==0){
//为空的处理
}
r.From对不同类型的表单元素的留空有不同的处理
对于空文本框 空文本区域 以及文件上传 元素的值为空值
而如果是未选中的复选框和单选按钮 则根本不会在r.Form中产生相应条目
如果我们用上面例子中的 方式去获取数据时程序就会报错
所以我们需要通过r.Form.Get()来获取值 因为 如果字段不存在 通过该方式获取的是空值
但是通过 r.Form.Get()只能获取单个的值 如果是map值 必须通过上面的方式来获取

数字
你想要 确保一个表单输入框中获取的只能是数字 例如 你想通过表单获取某个人的具体年龄是50岁还是10岁
而不是像一把年纪了 或者年轻着呢 这种描述
如果我们判断是正整数 那么我们需要先转化int类型 然后进行处理
getint,err:=strconv.Atoi(r.Form.Get("age"))
if err!=nil{
//数字转化出错了  那么可能就不是数字
}
//接下来就可以判断这个数字的大小范围了
if getint>100{
//太大了
}
还有一种方式就是正则匹配的方式
if m,_::=regexp.MatchString("^[0-9]+$",r.Form.Get("age"));!m{
return false
}
对于性能要求很高的 用户来说 这是一个老生常谈的问题了  他们认为应该尽量避免使用正则表达式 因为使用正则表达式的速度会变慢
但是在目前机器的性能那么强劲的情况下 对于这种简单的正则表达式效率和类型转换函数是没有什么差别的。如果你对正则表达式很熟悉，而且你在其它语言中也在使用它，那么在Go里面使用正则表达式将是一个便利的方式。
Go实现的正则是RE2，所有的字符都是UTF-8编码的。

中文
有时候我们想通过表单元素获取一个用户的中文名字但是又为了 保证获取的是正确的中文 我们需要进行验证
if m,_:=regexp.MatchString("^[\\x{4e00}-\\x{9fa5}]+$", r.Form.Get("realname"));!m{
return false
}

英文
我们希望通过表单元素获取一个英文值 例如我们想知道一个用户的英文名 我们可以很简单 的通过正则验证数据
if m,_:=regexp.MatchString("^[a-zA-Z]+$",f.Form.Get("engname"));!m{
return false
}

电子邮件地址
你想知道用户输入的一个email地址是否正确 通过如下这个方式可以验证
if m,_:=regexp.MatchString(`^([\w\.\_]{2,10})@(\w{1,}).([a-z]{2,4})$`, r.Form.Get("email")); !m {
fmt.Println("no")
}else{
fmt.Println("yes")
}

手机号码
if m, _ := regexp.MatchString(`^(1[3|4|5|8][0-9]\d{4,8})$`, r.Form.Get("mobile")); !m {
    return false
}

下拉菜单
如果我们想要判断表单里面<select> 元素生成的下拉菜单中是否有 被选中的项目 有些时候 黑客可能会
伪造这个下拉菜单不存在的值发送给你 那么如何判断这个值 是否是我们预设的值呢
<select name="fruit">
<option value="apple">apple</option>
<option value="pear">pear</option>
<option value="banane">banane</option>
</select>
slice:=[]string{"apple","pear","banane"}

for _, v := range slice {
    if v == r.Form.Get("fruit") {
        return true
    }
}
return false
上面这个函数包含在我开源的一个库里面(操作slice和map的库)，https://github.com/astaxie/beeku

单选按钮
如果我们想要判断radio按钮是否有一个被选中了，我们页面的输出可能就是一个男、女性别的选择，但是也可能一个15岁大的无聊小孩，一手拿着http协议的书，另一只手通过telnet客户端向你的程序在发送请求呢，你设定的性别男值是1，女是2，他给你发送一个3，你的程序会出现异常吗？因此我们也需要像下拉菜单的判断方式类似，判断我们获取的值是我们预设的值，而不是额外的值。
<input type="radio" name="gender" value="1">男
<input type="radio" name="gender" value="2">女
那么我们也可以类似下拉菜单的做法一样
slice:=[]int{1,2}
for _,v:=range slice{
if v==r.From.Get("gender"){
return true
}
}
return false

复选框
有一项选择兴趣的复选框 你想确定用户选中的和你提供给用户选择的是同一类型的数据
<input type="checkbox" name="interest" value="football">足球
<input type="checkbox" name="interest" value="basketball">篮球
<input type="checkbox" name="interest" value="tennis">网球
对于复选框我们的验证和单选有点不一样 因为接收到的是一个slice
slice:=[]string{"football","basketball","tennis"}
a:=Slice_diff(r.Form["interest"],slice)
if a == nil{
    return true
}

return false

日期和时间
t := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
fmt.Printf("Go launched at %s\n", t.Local())
获取time之后  我们就可以 进行很多时间函数的操作 具体的判断 就根据自己的需求调整

身份证号码
如果我们想验证表单输入的是否是身份证 通过正则也可以方便的验证 但是身份证有15位 和18位 我们两个都需要进行验证
//验证15位身份证，15位的是全部数字
if m, _ := regexp.MatchString(`^(\d{15})$`, r.Form.Get("usercard")); !m {
    return false
}

//验证18位身份证，18位前17位为数字，最后一位是校验位，可能为数字或字符X。
if m, _ := regexp.MatchString(`^(\d{17})([0-9]|X)$`, r.Form.Get("usercard")); !m {
    return false
}
上面列出了 我们一些常用的服务器端的表单元素验证 希望通过这个引导入门 能够让你对Go的数据验证有所了解
特别是Go里面正则处理
*/
