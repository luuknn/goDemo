package main

//import (
//	"database/sql"
//	_ "github.com/go-sql-driver/mysql"
//	"log"
//)
//
//func main() {
//	var err error
//	Db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
//	if err != nil {
//		log.Fatal(err)
//	}
//
//	defer Db.Close()
//	_, err = Db.Exec("CREATE TABLE IF NOT EXISTS cmistest.hello(word varchar(50))")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	rs, err := Db.Exec("insert into cmistest.hello(word) values ('hellp ')")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	rowCount, err := rs.RowsAffected()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	log.Printf("inserted %d rows", rowCount)
//
//	rows, err := Db.Query("select word from cmistest.hello")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	for rows.Next() {
//		var s string
//		err = rows.Scan(&s)
//		if err != nil {
//			log.Fatalln(err)
//		}
//		log.Printf("found row containing %q", s)
//	}
//	rows.Close()
//}

/*
Golang mysql 连接与连接池
database/sql 是golang的标准库之一 它提供了一系列接口方法 用于访问数据库 它并不会提供数据库特有的方法
那些特有的方法交给数据库驱动去实现

database/sql库提供了一些type 这些类型对掌握它的用法非常重要

DB数据库对象。 sql.DB类型代表了数据库 和其他语言不一样 它不是数据库连接 golang中的连接 来自内部实现的连接池
连接的建立是惰性的  当你需要连接的时候 连接池会自动帮你创建  通常你不需要操作连接池 一切都有go来帮你完成

Results结果集。数据库查询的时候都会有结果集 sql.Rows类型表示查询返回多行数据的结果集
sql.Row则表示单行查询结果的结果集 当然 对于插入更新和删除 返回的结果集类型为sql.Result

Statement语句  sql.Stmt类型表示sql查询语句 例如DDL DML等类似sql语句 可以把当成prepare语句构造查询 也可以直接使用
sql.DB的函数对其操作

warming up
下面就开始我们的sql数据库之旅 mysql  驱动 使用go-sql-driver/mysql

对于其他语言 查询数据的时候需要 创建一个连接 对于go而言 则是需要创建一个数据库抽象对象 连接将会在查询需要的时候 由连接池创建并维护
使用 sql.Open函数创建 数据库对象  它的第一个参数是数据库驱动名  第二个参数是一个连接 字串 (符合DSN风格 可以是一个tcp连接 一个unix socket等)
package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	var err error
	Db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	if err != nil {
		log.Fatal(err)
	} else {
		fmt.Println("create DB  :  success")
	}
	defer Db.Close()

}
创建了数据库对象之后 在函数退出的时候 需要 释放连接 即调用 sql.Close方法 例子使用了 defer语句设置 释放连接
接下来进行一些基本的数据库操作 首先 我们使用exec方法执行一条sql 创建一个数据表
import (
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	var err error
	Db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}

	defer Db.Close()
	_, err = Db.Exec("CREATE TABLE IF NOT EXISTS cmistest.hello(word varchar(50))")
	if err != nil {
		log.Fatalln(err)
	}
	rs, err := Db.Exec("insert into cmistest.hello(word) values ('hellp ')")
	if err != nil {
		log.Fatalln(err)
	}
	rowCount, err := rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	log.Printf("inserted %d rows", rowCount)

	rows, err := Db.Query("select word from cmistest.hello")
	if err != nil {
		log.Fatalln(err)
	}
	for rows.Next() {
		var s string
		err = rows.Scan(&s)
		if err != nil {
			log.Fatalln(err)
		}
		log.Printf("found row containing %q", s)
	}
	rows.Close()
}
我们使用Query方法执行select 查询语句 返回的是一个sql.Rows类型的结果集  迭代后者的Next方法
然后使用Scan方法给变量s赋值  以便取出结果  最后再把结果集关闭(释放连接)
通过上面一个简单的例子  介绍了 database/sql的基本数据库查询操作
*/
/*
sql.DB
正如上文所言  sql.DB是数据库的抽象 虽然 它容易被误以为是数据库连接  它提供了一些跟数据库交互的函数
同时管理维护了一个数据库连接池 帮你处理了单调而重复的管理工作  并且在多个goroutines也是十分安全

sql.DB表示的是数据库抽象  因此 你有几个数据库就需要为每一个数据库创建一个sql.DB对象
因为它维护了一个连接池  因此不需要频繁的创建和销毁 它需要长时间保持 因此最好是设置成一个全局变量 以便其它代码可以访问
创建数据库对象需要引入标准库database/sql 同时还需要引入 驱动 go-sql-driver/mysql
使用_表示引入驱动的变量 这样做的目的为了在代码中不至于和标注库的函数变量namespace冲突
*/
/*
连接池
只用sql.Open 函数创建连接池  可是此时只是初始化了连接池 并没有创建任何连接 连接创建都是惰性的 只有当你真正使用到连接的时候 连接池才会创建连接
连接池很重要 它直接影响着你的程序行为
连接池的工作原理却相当简单  当你的函数 如Exec Query 调用需要访问底层数据库的时候 函数首先会向连接池请求一个连接
如果连接池有空闲的连接  则返回给函数 否则连接池将会创建一个新的连接 给函数 一旦连接给了函数 连接则归属于函数
函数执行完毕 后 要要么把连接所属权归还给连接池 要么传递给下一个需要连接的Rows对象 最后使用完连接的对象也会把连接释放到连接池

请求一个连接的函数有好几种 执行完毕处理连接的方式稍有差别 大致如下
db.Ping()调用完毕后 会马上把连接返回给连接池
db.Exec()调用完毕后马上把连接返回给连接池 但它返回的Result对象还保留着连接的引用 当后面的代码需要处理结果集的时候连接将会被重用
db.Query()调用完毕后会将连接传递给sql.Rows类型 但后者迭代完毕 或者显示的调用Close()方法后 连接将会被释放回到连接池
db.QueryRow()调用完毕后会将连接传递给sql.Row类型 当.Scan()方法调用之后 把连接释放回到连接池
db.Begin()调用完毕后将连接传递给sql.Tx类型对象 当 .Commit()或Rollback()方法调用后释放连接

因为每一个连接都是惰性创建的  如何验证sql.Open调用之后 sql.DB对象可用呢 通常使用db.Ping()方法初始化
db, err := sql.Open("driverName", "dataSourceName")
if err != nil{
    log.Fatalln(err)
}
defer db.Close()
err = db.Ping()
if err != nil{
   log.Fatalln(err)
}
调用了ping之后连接池一定会初始化一个数据库连接 当然 实际上对于失败的处理 应该定义一个符合自己需要的
方式 现在为了演示  简单的使用log.Fatalln(err)表示了
*/

/*
连接失败
关于连接池另外一个知识点 就是你不必检查或者尝试处理连接失败的情况  当你进行 数据库操作的时候
如果连接失败了  database/sql会帮你处理 实际上 当从连接池取出的连接断开的时候  databasesql会自动尝试重连10次
仍然无法重连的情况下会自动从连接池再获取一个或者新建另外一个

连接池配置
配置连接池有两个的方法：

db.SetMaxOpenConns(n int) 设置打开数据库的最大连接数。包含正在使用的连接和连接池的连接。如果你的函数调用需要申请一个连接，并且连接池已经没有了连接或者连接数达到了最大连接数。此时的函数调用将会被block，直到有可用的连接才会返回。设置这个值可以避免并发太高导致连接mysql出现too many connections的错误。该函数的默认设置是0，表示无限制。
db.SetMaxIdleConns(n int) 设置连接池中的保持连接的最大连接数。默认也是0，表示连接池不会保持释放会连接池中的连接的连接状态：即当连接释放回到连接池的时候，连接将会被关闭。这会导致连接再连接池中频繁的关闭和创建。
对于连接池的使用依赖于你是如何配置连接池，如果使用不当会导致下面问题：

大量的连接空闲，导致额外的工作和延迟。
连接数据库的连接过多导致错误。
连接阻塞。
连接池有超过十个或者更多的死连接，限制就是10次重连。
大多数时候，如何使用sql.DB对连接的影响大过连接池配置的影响。这些具体问题我们会再使用sql.DB的时候逐一介绍。

掌握了database/sql关于数据库连接池管理内容，下一步则是使用这些连接，进行数据的交互操作啦。


*/
