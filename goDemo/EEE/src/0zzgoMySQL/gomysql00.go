package main

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

func store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}

}
func load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}

}
func main() {
	post := Post{Id: 1, Content: "Hello World!", Author: "Vanyarpy"}
	store(post, "post5")
	var postRead Post
	load(&postRead, "post5")
	fmt.Println(postRead)
}

/*
持久化
程序可以定义为算法+数据 算法是我们的代码逻辑 代码逻辑处理数据 数据的存在形式并不单一 可以存在数据库 文件 无论存在什么地方
处理数据的时候都需要把数据读入内存 如果直接存在内存中 不就可以直接读了么 的确 数据可以存在内存中 涉及数据存储的过程称之为
持久化 下面golang中的数据持久化做简单的介绍 主要包括内存存储 文件存储 和数据库存储
*/

/*
内存存储
所谓内存存储 即定义一些数据结构数组切片 图或者其他自定义结构 把需要持久化的数据存储在这些数据库结构中 使用数据
的时候可以直接操作这些结构
type Post struct {
	Id      int
	Content string
	Author  string
}

var PostById map[int]*Post
var PostsByAuthor map[string][]*Post

func store(post Post) {
	PostById[post.Id] = &post
	PostsByAuthor[post.Author] = append(PostsByAuthor[post.Author], &post)
}
func main() {
	PostById = make(map[int]*Post)
	PostsByAuthor = make(map[string][]*Post)
	post1 := Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"}
	post2 := Post{Id: 2, Content: "Bonjour Monde!", Author: "Pierre"}
	post3 := Post{Id: 3, Content: "Hola Mundo!", Author: "Pedro"}
	post4 := Post{Id: 4, Content: "Greetings Earthlings!", Author: "Sau Sheong"}

	store(post1)
	store(post2)
	store(post3)
	store(post4)
	fmt.Println(PostById[1])
	fmt.Println(PostById[2])
	for _, post := range PostsByAuthor["Sau Sheong"] {
		fmt.Println(post)
	}
	for _, post := range PostsByAuthor["Pedro"] {
		fmt.Println(post)
	}
}
我们定义了两个map的结构PostById PostByAuthor stor方法会把post数据存入两个结构中
当需要数据的时候 再从这两个内存结构中读取即可

内存持久化比较简单 严格来说 这也不算是持久化 比较程序退出会清空内存 所保存的数据也会消失
这种持久化只是相对程序运行时而言 想要程序退出 重启还能读取所存储的数据 这时就得依赖文件或者数据库(非内存数据库)
*/

/*
文件存储
文件存储 顾名思义 就是将需要存储的数据写入文件中 然后文件保存在硬盘中 需要读取数据的时候
再载入文件 把数据读取到内存中 所写入的数据和创建的文件可以自定义 例如 一个纯文本 格式化文本
甚至是二进制文件都可以  无非就是编码写入 读取解码的两个过程

下面我们介绍三种常用的文件存储方式 纯文本文件 csv文件 或二进制文件

纯文本
纯文本文件是最简单的一种文件存储方式 只需要将保存的字符串写入文本保存即可 golang提供了
ioutil库用于读写文件 也提供了os相关的文件创建 写入 保存的工具函数
import (
	"fmt"
	"io/ioutil"
)

func main() {
	data := []byte("Hello World!\n")
	fmt.Println(data)
	err := ioutil.WriteFile("data1", data, 0644)
	if err != nil {
		panic(err)
	}
	read1, _ := ioutil.ReadFile("data1")
	fmt.Println(string(read1))
}
我们先创建了一个byte类型的数组Hello World!\n 一共13个字符 对应的切片为
[72 101 108 108 111 32 87 111 114 108 100 33 10] 调用ioutil的WriteFile方法
即可创建data1 文件 并且文件存储的是文本字符串 使用ReadFile方法可以读取文本字符串内容
注意 读取的数据也是一个byte类型的切片 因此需要使用string转换成文本

除了ioutil库 还可以使用os库的函数进行文件读写操作
import (
	"fmt"
	"os"
)

func main() {
	data := []byte("Hello World!\n")
	file1, _ := os.Create("data2")
	defer file1.Close()

	bytes, _ := file1.Write(data)
	fmt.Printf("Write %d bytes to file \n", bytes)

	file2, _ := os.Open("data2")
	defer file2.Close()

	read2 := make([]byte, len(data))
	bytes, _ = file2.Read(read2)
	fmt.Printf("Read %d bytes from file\n", bytes)
	fmt.Println(read2, string(read2))
}
使用os的Creadte方法 创建一个文件 返回一个文件句柄结构 对于文件这种资源结构 即使定义defer资源清理是一个
好习惯 使用Write将数据写入文件 文件的写入完毕

读取的时候略显麻烦 使用Open函数打开文件句柄 创建一个空的byte切片 然后使用Read方法读取数据 并赋值给切片
如果想要文本字符 还需要调用string转换格式

csv
csv文件是一种以逗号分隔单元数据的文件 类似表格 但是很轻量 对于存储一些结构化的数据很有用 golang提供了专门处理csv的库

和纯文本文件读写类似 csv文件需要通过os创建一个文件句柄 然后调用相关的csv函数读写数据
import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

func main() {
	csvFile, err := os.Create("posts.csv")
	if err != nil {
		panic(err)
	}
	defer csvFile.Close()
	allPosts := []Post{
		Post{Id: 1, Content: "Hello World!", Author: "Sau Sheong"},
		Post{Id: 2, Content: "Bonjour Monde!", Author: "Pierre"},
		Post{Id: 3, Content: "Hola Mundo!", Author: "Pedro"},
		Post{Id: 4, Content: "Greetings Earthlings!", Author: "Sau Sheong"},
	}
	writer := csv.NewWriter(csvFile)
	for _, post := range allPosts {
		line := []string{strconv.Itoa(post.Id), post.Content, post.Author}
		fmt.Println(line)
		err := writer.Write(line)
		if err != nil {
			panic(err)
		}
	}
	writer.Flush()

	file, err := os.Open("posts.csv")
	if err != nil {
		panic(err)
	}
	defer file.Close()
	reader := csv.NewReader(file)
	reader.FieldsPerRecord = -1
	record, err := reader.ReadAll()
	if err != nil {
		panic(err)
	}
	var posts []Post
	for _, item := range record {
		id, _ := strconv.ParseInt(item[0], 0, 0)
		post := Post{Id: int(id), Content: item[1], Author: item[2]}
		posts = append(posts, post)

	}
	fmt.Println(posts[1].Id)
	fmt.Println(posts[1].Content)
	fmt.Println(posts[1].Author)

}
创建了文件句柄之后 使用csv的函数NewWriter创建一个可写对象 然后依次遍历数据 写入数据 写完的时候 需要调用Flush方法
读取csv文件也类似 创建一个NewReader的可读对象 然后读取内容

gob
无论纯文本 还是csv文件的读写 所存储的数据文件是可以直接用文本工具打开的 对于一些不希望被文件工具打开 需要将数据写成二进制
幸好go提供了gob模板用于创建二进制文件
import (
	"bytes"
	"encoding/gob"
	"fmt"
	"io/ioutil"
)

type Post struct {
	Id      int
	Content string
	Author  string
}

func store(data interface{}, filename string) {
	buffer := new(bytes.Buffer)
	encoder := gob.NewEncoder(buffer)
	err := encoder.Encode(data)
	if err != nil {
		panic(err)
	}
	err = ioutil.WriteFile(filename, buffer.Bytes(), 0600)
	if err != nil {
		panic(err)
	}

}
func load(data interface{}, filename string) {
	raw, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	buffer := bytes.NewBuffer(raw)
	dec := gob.NewDecoder(buffer)
	err = dec.Decode(data)
	if err != nil {
		panic(err)
	}

}
func main() {
	post := Post{Id: 1, Content: "Hello World!", Author: "Vanyarpy"}
	store(post, "post5")
	var postRead Post
	load(&postRead, "post5")
	fmt.Println(postRead)
}
定义一个函数 用于写入数据 使用NewEncoder方法创建一个encoder对象 然后对数据进行二进制编码 最后将数据写入文件中
读文件的内容的过程与之相反 先读取文件的内容 然后把这个二进制内容转换成一个buffer对象 最后再解码 调用的过程也很简单

通过上面的小例子 我们讨论了golang中基本文件读写操作 基本上涉及的都有纯文本 格式化文本 和二进制文本的读写操作
通过文件持久化数据比起 内存才是真正的持久化 然而很多应用的开发 持久化更多还是和数据库打交道

关于数据库又是一个很大的话题 我们先简单的讨论下sql 后续再针对mysql的操作做详细的介绍



*/
