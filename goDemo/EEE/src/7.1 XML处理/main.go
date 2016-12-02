package main

import (
	"fmt"
)

//文本处理
//web开发中对于文本处理是非常重要的一部分
//我们往往需要对输出或者输入的内容进行处理 这里的文本包括字符串 数字 Json XML等等
//Go语言作为一门高性能语言  对这些文本的处理都有官方的标准库来支持 而且在你使用中 你会发现
//go的标准库一些设计相当的巧妙 而且对于使用者来说 也很方便就能处理这些文本  本章我们将通过四小节的介绍 让用户对Go语言的处理文本有一个很好的认识

//xml是目前很多标准接口的交互语言  很多时候 和一些java编写 的webserver进行交互都是基于xml标准进行交互的
//7.1小节将介绍如何处理xml文本 我们使用xml之后发现它太复杂了
//很多互联网企业对外的API大多数采用了JSON格式 这种格式描述简单 但是又很好的表达意思
//7.2 小节 我们将讲述如何来处理这样的json格式数据  正则是一个让人又爱又恨的工具
//它处理文本的能力非常强大 我们在前面表单验证里面已经有所领略 它的强大
//7.3小节 将详细的更深入的讲解如何利用好Go的正则 web开发中一个很重要的部分就是MVC分离 在Go语言web开发中 V
//有一个专门的包来支持 template 7.4小节将详细的讲解如何使用模板 来进行输出内容
//7.5小节 将详细的介绍如何进行文本和文件夹的操作
//7.6小节介绍了字符串相关的操作

//XML处理
//xml作为一种数据交换和信息传递的格式已经十分普及  而随着web服务日益广泛的应用
//现在xml在日常的开发工作中 也扮演了愈发重要的角色 这一小节 我们将就Go语言标准中的xml相关处理的包进行介绍
//这个小节不会涉及xml规范相关的内容 而是介绍如何用Go语言来编解码xml文件相关的知识
//假如你是一名运维人员  你为你所管理的服务器生成了如下内容的xml配置文件
<?xml version="1.0" encoding="utf-8">
<servers version="1">
<server>
<serverName>shanghai_vpn</serverName>
<serverIP>127.0.0.1</serverIp>
</server>
<serverName>beijing_vpn</serverName>
<serverIP>127.0.0.2</serverIp>
</servers>
//上面的xml文档 描述了两个服务器的信息  包含了服务器名和服务器的IP信息 接下来Go例子 以此xml描述的信息进行操作
//解析xml
//如何解析如上这个xml文件呢 我们可以通过xml包的Unmarshal函数来达到我们的目的
func Unmarshal(data []byte,v interface{}) error
//data接收的是xml数据流 v是要输出的结构定义为interface 也就是可以把xml转化为任意的格式
//我们这里  主要介绍 struct的转换  因为struct和xml都有类似树结构的特征
//示例代码如下
package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Recurlyservers struct {
	XMLName     xml.Name `xml:"servers"`
	Version     string   `xml:"version,attr"`
	Svs         []server `xml:"server"`
	Description string   `xml:",innerxml"`
}

type server struct {
	XMLName    xml.Name `xml:"server"`
	ServerName string   `xml:"serverName"`
	ServerIP   string   `xml:"serverIP"`
}

func main() {
	file, err := os.Open("D:/server.xml")
	if err != nil {
		fmt.Printf("erroe:%v", err)

	}
	defer file.Close()
	data, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Printf("erroe:%v", err)
	}
	v := Recurlyservers{}
	err = xml.Unmarshal(data, &v)
	if err != nil {
		fmt.Printf("error:%v", err)
		return
	}
	fmt.Println(v)
}
//如上例子 输出如下结果
{{ servers} 1 [{{ server} Shanghai_VPN 127.0.0.1} {{ server} Beijing_VPN 127.0.0.2}] 
    <server>
        <serverName>Shanghai_VPN</serverName>
        <serverIP>127.0.0.1</serverIP>
    </server>
    <server>
        <serverName>Beijing_VPN</serverName>
        <serverIP>127.0.0.2</serverIP>
    </server>
}
//上面例子中 将xml文件解析成对应struct对象是通过xml.Unmarshal来完成的
//这个过程如何实现的  可以看到我们的struct后面多了一些类似于 xml:"serverName"这样的内容 
//这个是struct的一个特性  它们被称为 struct tag 它们是用来辅助反射的  我们来看一下Unmarshal的定义
func Unmarshal(data []byte ,v interface{}) error
//我们看到函数定义了两个参数 第一个是xml数据流 第二个是存储的对应类型
//目前支持 struct slice和string xml包内部采用了反射来进行 数据的映射 所以 v里面的字段必须是导出的
//Unmarshal解析的时候xml元素和字段怎么对应起来的呢  这是有一个优先级读取流程的  首先会 读取struct tag
//如果没有  那么就会对应字段名  必须注意一点的是解析的时候tag 字段名 xml 元素都是大小写敏感的  所以必须一一对应字段

//Go语言的反射机制 利用这些tag信息 来将 来自xml文件中的数据 反射成对应的struct对象  关于反射 如何利用struct tag 更多内容请参阅 reflect中的相关内容
//解析XML到struct的时候遵循如下的规则
// 如果struct的字段是string 或者[]byte类型且 它的tag含有 ",innerxml"Unmarshal会将此字段所对应元素内的
//所有内嵌的原始xml累加到此字段上  如上面的Description定义
//如果struct 中 有一个叫做XMLName的字段 并且类型为xml.Name 那么在解析的时候就会保存这个element的名字到该字段如servers
//如果某个struct字段的tag定义中含有XML结构中element的名称，那么解析的时候就会把相应的element值赋值给该字段，如上servername和serverip定义。
//如果某个struct字段的tag定义了中含有",attr"，那么解析的时候就会将该结构所对应的element的与字段同名的属性的值赋值给该字段，如上version定义。
//如果某个struct字段的tag定义 型如"a>b>c",则解析的时候，会将xml结构a下面的b下面的c元素的值赋值给该字段。
//如果某个struct字段的tag定义了"-",那么不会为该字段解析匹配任何xml数据。
//如果struct字段后面的tag定义了",any"，如果他的子元素在不满足其他的规则的时候就会匹配到这个字段。
//如果某个XML元素包含一条或者多条注释，那么这些注释将被累加到第一个tag含有",comments"的字段上，这个字段的类型可能是[]byte或string,如果没有这样的字段存在，那么注释将会被抛弃。

//上面详细讲述了 如何定义struct 的tag 只要设置对了tag 那么xml解析就如上面的示例般简单
//tag和xml的element是一一对应的关系 如上所示  我们还可以通过slice来表示多个同级元素
//注意 为了正确解析 go语言的xml包要求struct定义中的所有字段必须是可导出的  即首字母大写

//输出XML
//假若 我们不是要 解析如上所示的XML文件 而是生成它  那么在Go语言中 又该如何实现呢
//xml包中  提供了Marshal和 MarshalIndent两个函数 来满足我们的需求 
//这两个函数主要的区别 是第二个函数会增加 前缀 和缩进  函数的定义 如下所示
func Marshal(v interface{})([]byte,error)
func MarshalIndent(v interface{},prefix,indent string)([]byte,error)
//两个函数的第一个参数是用来生成xml的结构定义类型数据 都是返回生成的xml数据流
//下面我们来看一下 如何输出如上的xml

package main

import (
	"encoding/xml"
	"fmt"
	"os"
)

type Servers struct {
	XMLName xml.Name `xml:"servers"`
	Version string   `xml:"version,attr"`
	Svs     []server `xml:"server"`
}

type server struct {
	ServerName string `xml:"serverName"`
	ServerIP   string `xml:"serverIP"`
}

func main() {
	v := &Servers{Version: "1"}
	v.Svs = append(v.Svs, server{"shanghai_vpn", "127.0.0.1"})
	v.Svs = append(v.Svs, server{"beijing_vpn", "127.0.0.2"})
	output, err := xml.MarshalIndent(v, "  ", "  ")
	if err != nil {
		fmt.Printf("error:%v", err)
	}
	os.Stdout.Write([]byte(xml.Header))
	os.Stdout.Write(output)
}
//和我们之前定义文件的格式一模一样  之后会有os.Stdout.Write([]byte(xml.Header))
//这句代码的出现 是因为 xml.MarshalIndent和xml.Narshal输出的信息都是不带xml头的
//为了生成正确的xml文件 我们使用了 xml包 预定义的Header变量

//我们看到Marshal函数接收的参数v是interface{}类型的，即它可以接受任意类型的参数，那么xml包，根据什么规则来生成相应的XML文件呢？
//
//如果v是 array或者slice，那么输出每一个元素，类似value
//如果v是指针，那么会Marshal指针指向的内容，如果指针为空，什么都不输出
//如果v是interface，那么就处理interface所包含的数据
//如果v是其他数据类型，就会输出这个数据类型所拥有的字段信息
//生成的XML文件中的element的名字又是根据什么决定的呢？元素名按照如下优先级从struct中获取：
//
//如果v是struct，XMLName的tag中定义的名称
//类型为xml.Name的名叫XMLName的字段的值
//通过struct中字段的tag来获取
//通过struct的字段名用来获取
//marshall的类型名称
//我们应如何设置struct 中字段的tag信息以控制最终xml文件的生成呢？
//
//XMLName不会被输出
//tag中含有"-"的字段不会输出
//tag中含有"name,attr"，会以name作为属性名，字段值作为值输出为这个XML元素的属性，如上version字段所描述
//tag中含有",attr"，会以这个struct的字段名作为属性名输出为XML元素的属性，类似上一条，只是这个name默认是字段名了。
//tag中含有",chardata"，输出为xml的 character data而非element。
//tag中含有",innerxml"，将会被原样输出，而不会进行常规的编码过程
//tag中含有",comment"，将被当作xml注释来输出，而不会进行常规的编码过程，字段值中不能含有"--"字符串
//tag中含有"omitempty",如果该字段的值为空值那么该字段就不会被输出到XML，空值包括：false、0、nil指针或nil接口，任何长度为0的array, slice, map或者string
//tag中含有"a>b>c"，那么就会循环输出三个元素a包含b，b包含c，例如如下代码就会输出

FirstName string   `xml:"name>first"`
LastName  string   `xml:"name>last"`

<name>
<first>Jianling</first>
<last>Shih</last>
</name>
//上面我们介绍了如何使用Go语言的xml包来编解码xml文件 重要的一点是对xml的所有操作 都是通过struct tag 来实现的
//所以 学会对struct tag的运用变得非常重要 在文章中 我们简要的列举了如何定义tag

















