package main

//import (
//	"database/sql"
//	"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"log"
//)
//
//func doSomething() {
//	panic("A panic Running Error")
//
//}
//func clearTransaction(tx *sql.Tx) {
//	err := tx.Rollback()
//	if err != sql.ErrTxDone && err != nil {
//		log.Fatalln(err)
//	}
//
//}
//func main() {
//	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer db.Close()
//
//	tx, err := db.Begin()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	defer clearTransaction(tx)
//	rs, err := tx.Exec("update posts set content='hello world'where author='forever'")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	rowAffected, err := rs.RowsAffected()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println(rowAffected)
//	rs, err = tx.Exec("update posts set content='hello'where author='jl'")
//	if err != nil {
//		log.Fatalln(err)
//	}
//	rowAffected, err = rs.RowsAffected()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	fmt.Println("hhh", rowAffected)
//	//doSomething()
//	if err := tx.Commit(); err != nil {
//		// tx.Rollback() 此时处理错误，会忽略doSomthing的异常
//		log.Fatalln(err)
//	}
//}

/*
Transaction 事务
事务处理  是数据的重要特性  尤其是对于一些支付系统  事务保证性对业务逻辑会有重要影响 golang的mysql驱动也封装好了
事务相关的操作  我们已经学习了db的Query和Exec方法处理查询和修改数据库

tx对象
一般查询使用的是db对象的方法  事务则使用另外一个对象 sql.Tx对象 使用db的Begin方法可以创建tx对象 tx对象也有
数据库交互的Query Exec Prepare方法 用法和db的相关用法类似  查询或修改的操作完毕之后 需要调用tx对象的commit提交或者rollback方法回滚

一旦创建了tx对象 事务处理都依赖与tx对象 这个对象会从连接池中取出一个空闲的连接 接下来的sql执行都基于这个连接
直到commit或rollback调用之后  才会把连接释放到连接池

在事务处理的时候 不能使用db的查询方法 虽然后者可以获取数据 可是这不属于同一个事务处理 将不会接受commit和rollback的改变
一个简单的事务例子如下
tx,err:=db.Begin()
tx.Exec(query1)
tx.Exec(query2)
tx.commit()
在tx中使用db是错误的
tx,err:=db.Begin()
db.Exec(query1)
tx.Exec(query2)
tx.commtit()
上述代码在调用db的Exec方法的时候 tx会绑定连接到事务中 db则是一个额外的连接 两者不是同一个事务
需要注意 Begin和Commit方法 与sql语句中的BEGIN或COMMIT语句没有关系

事务与连接
创建Tx对象的时候 会从连接池中取出连接 然后调用相关的Exec方法的时候 连接仍然会绑定在该事务处理中
在实际的事务处理中 go可能创建不同的连接 但是那些其他连接都不属于该事务 例如上面例子中db创建的连接和tx的连接就不是一回事

事务的连接生命周期从Begin函数调用开始  直到Commit和Rollback函数的调用结束 事务也提供了prepare语句的使用方式
但是需要使用Tx.Stmt方法创建  prepare设计的初衷是多次执行 对于事务 有可能需要多次执行同一条sql然而
然而无论是正常的prepare和事务处理  prepare对于连接的管理都有点小复杂  因此尽量避免在事务中使用prepare的方式
例如下面的例子就容易导致错误
tx,_:=db.Begin()
defer tx.Rollback()
stmt,_ tx.Prepare("INSERT...")
defer stmt.Close()
tx.Commit()
因为stmt.Close 使用defer语句 即函数退出的时候 再清理stmt 可是实际执行过程的时候 tx.Commit就已经释放了连接
当函数退出的时候 再执行stmt.Close的时候 连接可能又被使用了

事务并发
对于sql.Tx对象 因为事务过程只有一个连接  事务内的操作都是顺序执行的  在开始下一个数据库交互之前
必须先完成上一个数据库交互  例如下面的例子
rows,_ :=db.Query("SELECT ID FROM USER")
for rows.Next(){
var mid,did int
rows.Scan(&mid)
db.QueryRow("SELECT idFROM detail_user WHERE master =?",mid).Scan(&did)
}
调用了Query方法之后 在Next方法中取结果的时候 rows是维护了一个连接 再次调用QueryRow的时候 db会再从连接池中取出一个新的连接
rows和db的连接两者可以并存 并且相互不影响
可是 这样的逻辑在事务处理中将会失效
rows, _ := tx.Query("SELECT id FROM user")
for rows.Next() {
   var mid, did int
   rows.Scan(&mid)
   tx.QueryRow("SELECT id FROM detail_user WHERE master = ?", mid).Scan(&did)
}
tx执行了Query方法后 连接转移到rows上  在Next方法中 tx.QueryRow将尝试获取该连接进行数据库操作
因为还没有调用rows.Close 因此底层的连接属于busy状态  tx是无法再进行查询的
上面的例子看起来有点傻 毕竟涉及这样的操作 使用query的join就能规避这个问题 例子只是为了说明tx的使用问题

*/

/*
前面对事务解释了一堆  说了那么多 其实还不动手实践下 下面就事务的使用做简单的介绍
因为事务是单个连接 因此 任何事务处理过程出现了异常  都需要使用rollback 一方面为了保证数据完整一致性
另一方面是释放事务绑定的连接
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func doSomething() {
	panic("A panic Running Error")

}
func clearTransaction(tx *sql.Tx) {
	err := tx.Rollback()
	if err != sql.ErrTxDone && err != nil {
		log.Fatalln(err)
	}

}
func main() {
	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	tx, err := db.Begin()
	if err != nil {
		log.Fatalln(err)
	}
	defer clearTransaction(tx)
	rs, err := tx.Exec("update posts set content='hello world'where author='forever'")
	if err != nil {
		log.Fatalln(err)
	}
	rowAffected, err := rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(rowAffected)
	rs, err = tx.Exec("update posts set content='hello'where author='jl'")
	if err != nil {
		log.Fatalln(err)
	}
	rowAffected, err = rs.RowsAffected()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("hhh", rowAffected)
	//doSomething()
	if err := tx.Commit(); err != nil {
		// tx.Rollback() 此时处理错误，会忽略doSomthing的异常
		log.Fatalln(err)
	}
}
我们定义了一个clearTransaction(tx)函数 该函数会执行rollback操作 因为我们事务处理过程中任何一个错误
都会导致main函数退出  因此在main函数退出执行defer的rollback 操作 回滚事务和释放连接

如果不添加defer最后Commit后check错误err后 再rollback那么当doSomething发生异常的时候 函数就退出了
此时还没有执行到tx.Commit 这样就导致事务的连接没有关闭  事务也没有回滚


总结
database/sql提供了事务处理的功能 通过Tx对象实现 db.Begin 会创建tx对象 后者的Exec和
Query执行事务的数据库操作 最后在tx的Commit和Rollback中完成数据库事务的提交和回滚 同时释放连接

tx事务环境中 只有一个数据库连接  事务内的Exec都是依次执行的  事务中也可以使用db进行查询
但是db查询的过程会新建连接 这个连接的操作不属于该事务

*/
