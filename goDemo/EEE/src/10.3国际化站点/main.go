package main

import ("fmt")

func main() {
}

/*10.3国际化站点
前面小节介绍了 如何处理本地化资源 即locale一个相应的配置文件 那么如果处理多个本地化资源呢 而对于一些我们经常用到的
例如 简单的文本翻译 时间日期 数字等如何处理呢 本小节 将一一解决这些问题
管理多个本地包
在开发一个应用时 首先我们要决定是只支持一种语言 还是多种语言 如果支持多种语言 我们需要制定一个组织结构
以方便将来更多语言的添加 在此我们设计如下 locale有关的文件放置在
config/locales下 假设你要支持中文和英文 那么你需要 在这个文件夹下放置en.json 和zh.json大概的内容如下所示
# zh.json
{
"zh":{
"submit": "提交",
"create":"创建"
}
}
#en.json
{
"en":{
"submit":"Submit",
"create":"Create"
}
}
为了支持国际化 在此我们使用了一个国际化相关的包 go-i18n首先 我们向 goi18n包注册
config/locales 这个目录 以加载所有的locale文件
Tr:=i18n.NewLocale()
Tr.LoadPath("config/locales")
这个包使用起来很简单 你可以通过下面的方式进行测试
fmt.Println(Tr.Translate("submit"))
//输出Submit
Tr.SetLocale("zn")
fmt.Println(Tr.Translate("submit"))
//输出“递交”

//自动加载本地包
//上面我们介绍了如何自动加载自定义语言包 其实go-i18n 库已经预加载了很多默认的格式信息 例如时间格式 货币格式
//用户可以在自定义配置时改写这些默认配置 请看下面的处理过程

//加载默认配置文件 这些文件都放在 go-i18n/locales 下面
//文件命名zh.json en-json en-US.json 等 可以不断的扩展支持更多的语言*/
func (il *IL) loadDefaultTranslations(dirPath string ) error{
dir,err :=os.Open(dirPath)
if err!=nul {
return err
}
defer dir.Close()
names,err:=dir.Readdirnames(-1)
if err!=nil{
return err}
for _, name := range names {
		fullPath := path.Join(dirPath, name)

		fi, err := os.Stat(fullPath)
		if err != nil {
			return err
		}

		if fi.IsDir() {
			if err := il.loadTranslations(fullPath); err != nil {
				return err
			}
		} else if locale := il.matchingLocaleFromFileName(name); locale != "" {
			file, err := os.Open(fullPath)
			if err != nil {
				return err
			}
			defer file.Close()

			if err := il.loadTranslation(file, locale); err != nil {
				return err
			}
		}
	}

	return nil
}
通过上面的方法加载配置信息到默认的文件这样我们就可以在我们没有自定义时间信息的时候执行如下的代码获取队形的信息
//locale=zh 的情况下 执行如下代码
fmt.Println(Tr.Time(time.Now()))
//输出 :2017年3月27日 星期一 07:17:58 CST

fmt.Println(Tr.Time(time.Now(),"long"))
//输出 2017年 3月 27日
fmt.Println(Tr.Money(11.11))
//输出11.11
//面我们实现了多个语言包的管理和加载，而一些函数的实现是基于逻辑层的，例如："Tr.Translate"、"Tr.Time"、"Tr.Money"等，虽然我们在逻辑层可以利用这些函数把需要的参数进行转换后在模板层渲染的时候直接输出，但是如果我们想在模版层直接使用这些函数该怎么实现呢？不知你是否还记得，在前面介绍模板的时候说过：Go语言的模板支持自定义模板函数，下面是我们实现的方便操作的mapfunc：
//
//文本信息
//文本信息调用Tr.Translate来实现相应的信息转换，mapFunc的实现如下：
//
//func I18nT(args ...interface{}) string {
//	ok := false
//	var s string
//	if len(args) == 1 {
//		s, ok = args[0].(string)
//	}
//	if !ok {
//		s = fmt.Sprint(args...)
//	}
//	return Tr.Translate(s)
//}
//注册函数如下：
//
//t.Funcs(template.FuncMap{"T": I18nT})
//模板中使用如下：
//
//{{.V.Submit | T}}
//时间日期
//时间日期调用Tr.Time函数来实现相应的时间转换，mapFunc的实现如下：
//
//func I18nTimeDate(args ...interface{}) string {
//	ok := false
//	var s string
//	if len(args) == 1 {
//		s, ok = args[0].(string)
//	}
//	if !ok {
//		s = fmt.Sprint(args...)
//	}
//	return Tr.Time(s)
//}
//注册函数如下：
//
//t.Funcs(template.FuncMap{"TD": I18nTimeDate})
//模板中使用如下：
//
//{{.V.Now | TD}}
//货币信息
//货币调用Tr.Money函数来实现相应的时间转换，mapFunc的实现如下：
//
//func I18nMoney(args ...interface{}) string {
//	ok := false
//	var s string
//	if len(args) == 1 {
//		s, ok = args[0].(string)
//	}
//	if !ok {
//		s = fmt.Sprint(args...)
//	}
//	return Tr.Money(s)
//}
//注册函数如下：
//
//t.Funcs(template.FuncMap{"M": I18nMoney})
//模板中使用如下：
//
//{{.V.Money | M}}
//总结
//通过这小节 我们知道了如何实现一个多语言的Web应用 通过自定义语言包我们可以方便的实现多语言 而且
//通过配置文件能够非常方便的扩充多语言 默认情况下 go-i18n会自定加载一些公共的配置信息
//例如时间 货币 等 我们就可以非常方便的使用 同时为了支持在模板中使用这些函数 也实现了 相应的模板函数
//这样就允许我们在开发web应用的时候直接在模板中通过pipeline的方式来操作多语言包


10.4小结
//通过这一章的介绍 读者应该对如何操作i18n有了深入的了解 我也根据这一章介绍的内容 实现了一个开源的
//解决方案 go-18n 通过这个开源库 我们可以很方便的实现多语言版本的web应用 使得我们的应用能够轻松的实现国际化
//如果你发现这个开源库中的错误或者那些缺失的地方 请一起参与到这个开源项目中来 让我们这个库争取成为Go的标准库


















*/
