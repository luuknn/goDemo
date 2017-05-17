package main

import (
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
)

var buf []byte

func upload(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "GET" {
		t, err := template.ParseFiles("upload.html")
		checkErr(err)
		t.Execute(w, nil)
	} else {
		//解析form中file上传的名字
		file, handle, err := r.FormFile("file")
		checkErr(err)
		//打开 已只读 文件不存在创建 方式打开 要存放的路径
		f, err := os.OpenFile("test/"+handle.Filename, os.O_WRONLY|os.O_CREATE, 0666)
		//文件拷贝
		io.Copy(f, file)
		checkErr(err)
		//关闭对应打开的文件
		defer f.Close()
		defer file.Close()
		fmt.Println("upload success")
		fmt.Println("test/" + handle.Filename)
	}

}
func checkErr(err error) {
	if err != nil {
		err.Error()
	}

}
func main() {
	//绑定路由 如果访问 /upload 调用Handler方法
	http.HandleFunc("/upload", upload)
	//使用tcp协议监听 8888端口号
	err := http.ListenAndServe(":8888", nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}

}
