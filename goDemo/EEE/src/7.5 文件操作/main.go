package main

import (
	"fmt"
)

//在任何计算机设备中 文件都是必须的对象
//而在 web编程中 文件的操作 一直是 WEB程序员经常遇到的问题  文件操作 在web应用中 是必须的
//我们经常会遇到生成文件目录 文件(夹)编辑等操作 现在把Go中的这些操作 做个简单的总结并实例示范如何使用

//目录操作
//文件操作的大多数函数都在os包里面
//下面列举了几个目录操作的
//func Mkdir(name string,perm FileMode) error 创建名称为name的目录 权限设置是perm 例如0777
//func MkdirAll(path string,perm FileMode)error 根据path创建多级子目录 例如 jianling/text1/test2
//func Remove(name string) error 删除  名称为name的目录  当目录下有文件或者其他目录会出错
//func RemoveAll(path string) error
//根据path删除多级子目录  如果path是单个名称 那么该目录不删除
package main

import (
	"fmt"
	"os"
)

func main() {
	os.Mkdir("jianling", 0777)
	os.MkdirAll("jianling/test1/test2", 0777)
	err := os.Remove("jianling")
	if err != nil {
		fmt.Println(err)
	}
	os.RemoveAll("jianling")
}

//文件操作
//建立与打开文件
//新建文件可以通过如下两种方法
func Create(name string)(file *File,err Error)
//根据提供的文件名创建新的文件 返回一个文件对象 默认权限 是0666的文件  返回的文件对象是可读写的
func NewFile(fd uintptr,name string) *File
//根据文件描述符创建相应的文件 返回一个文件对象
//通过如下两个方法来打开文件
func Open(name string)(filr *File,err Error)
//该方法打开一个名称为name的文件 但是是只读的方式  内部实现其实调用了OpenFile
func OpenFile(name string,flag int,perm uint32)(file *File,err Error)
//打开名称为name的文件  flag是打开的方式 只读 读写 perm是权限

//写文件
//写文件函数
func (file *File) Write(b []byte)(n int,err Error)
//写入byte类型的信息到文件
func (file *File)WriteAt(b []byte,off int64)(n int,err Error)
//在指定位置开始写入byte类型的信息
func(file *File)WriteString(s string)(ret int,err Error)
//写入string信息到文件
package main

import (
	"fmt"
	"os"
)

func main() {
	userFile := "jianling.txt"
	fout, err := os.Create(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fout.Close()
	for i := 0; i < 10; i++ {
		fout.WriteString("Just a test!\r\n")
		fout.Write([]byte("Just a test~\r\n"))
	}
}
//读文件函数
func (file *File)Read(b[]byte)(n int,err Error)
//读取数据到b中
func (file *File)ReadAt(b[]byte,off int64)(n int,err Error)
//从off开始读取数据到b中
package main

import (
	"fmt"
	"os"
)

func main() {
	userFile := "jianling.txt"
	fl, err := os.Open(userFile)
	if err != nil {
		fmt.Println(userFile, err)
		return
	}
	defer fl.Close()
	buf := make([]byte, 1024)
	for {
		n, _ := fl.Read(buf)
		if 0 == n {
			break
		}
		os.Stdout.Write(buf[:n])

	}
}

//删除文件
//Go语言里面删除文件和删除文件夹是同一个函数
func Remove(name string)Error
//调用该函数就可以 删除文件名为name的文件















