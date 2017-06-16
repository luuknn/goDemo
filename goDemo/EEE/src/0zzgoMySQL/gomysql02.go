package main

import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	Db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
}

/*
数据库查询
我们了解了数据库连接与连接池 拿到了连接当然就是为了跟数据库交互 对于数据库交互 无怪乎两类操作
读和写   其中怎么读 怎么写 读和写的过程糅合在一起就会遇到复杂的事务  本章主要关注数据库的读写操作

读取数据
database/sql提供了Query和QueryRow方法进行查询数据库 对于Query方法 主要分为三步
1 从连接池中请求一个连接 2 执行查询的sql语句  3 将数据库连接所属权传递给Result结果集
Query返回的结果集是sql.Rows 类型  它有一个Next方法 可以迭代数据库的游标 进而获取每一行的数据
大概的使用方式如下
rows, err := db.Query("SELECT world FROM test.hello")
if err != nil{
    log.Fatalln(err)
}
for rows.Next(){
    var s string
    err = rows.Scan(&s)
    if err !=nil{
        log.Fatalln(err)
    }
    log.Printf("found row containing %q", s)
}
rows.Close()

上面的代码 我们已经见过好多次了 rows.Next方法设计用来迭代  当它迭代到最后一行数据之后 会触发一个
io.EOF的信号 即引发一个错误 同时go会自动调用rows.Close方法 释放连接  然后返回false此时循环将会结束退出

通常你会正常迭代完数据然后退出循环  可是如果并没有正常的循环 而因其他的错误导致退出了循环
此时rows.Next处理结果集的过程并没有完成  归属于rows的连接不会被释放到连接池 因此十分有必要正确的处理rows.Close事件
如果没有关闭rows连接  将导致大量的连接并且不会被其他函数重用 就像溢出了一样 最终将导致数据库无法使用

那么如何阻止这样的行为呢？上述代码已经展示，无论循环是否完成或者因为其他原因退出，都显示的调用rows.Close方法，确保连接释放。
又或者使用defer指令在函数退出的时候释放连接，即使连接已经释放了，rows.Close仍然可以调用多次，是无害的。
使用defer的时候需要注意，如果一个函数执行很长的逻辑，例如main函数，那么rows的连接释放就会也很长，好的实践方案是尽可能的越早释放连接。
rows.Next循环迭代的时候，因为触发了io.EOF而退出循环。为了检查是否是迭代正常退出还是异常退出，需要检查rows.Err。例如上面的代码应该改成：
rows, err := db.Query("SELECT world FROM test.hello")
if err != nil{
    log.Fatalln(err)
}
defer rows.Close()
for rows.Next(){
    var s string
    err = rows.Scan(&s)
    if err !=nil{
        log.Fatalln(err)
    }
    log.Printf("found row containing %q", s)
}
rows.Close()
if err = rows.Err(); err != nil {
    log.Fatal(err)
}

读取单条记录
Query方法是读取多行结果集  实际开发中 很多查询 只需要单条记录 不需要再通过Next迭代
golang提供了QueryRow方法用于查询单条记录的结果集

var s string
err = db.QueryRow("SELECT world FROM test.hello LIMIT 1").Scan(&s)
if err != nil{
    if err == sql.ErrNoRows{
        log.Println("There is not row")
    }else {
        log.Fatalln(err)
    }
}
log.Println("found a row", s)

QueryRow方法的使用很简单 它要么返回sql.Row类型 要么返回一个error
如果是发送了错误 则会延迟到Scan调用结束后返回 如果没有错误 则 Scan正常执行 只有当查询的结果为空的时候
会触发一个sql.ErrNoRows错误 你可以选择先检查错误 再调用 Scan 方法 或者 先调用Scan再检查错误












*/
