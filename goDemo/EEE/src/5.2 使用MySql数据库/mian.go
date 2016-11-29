package main

import (
	"database/sql"
	"fmt"
	_ "mysql"
	//"time"
)

func main() {
	db, err := sql.Open("mysql", "root:111111@/test?charset=utf8")
	checkErr(err)
	//插入数据
	/*
		stmt, err := db.Prepare("INSERT userinfo SET username=?,departname=?,created=?")
		checkErr(err)

		res, err := stmt.Exec("守望子", "研发部门", "2016-11-30")
		checkErr(err)

		id, err := res.LastInsertId()
		checkErr(err)

		fmt.Println(id)*/
	//更新数据
	/*
		stmt, err := db.Prepare("update userinfo set username=? where uid=?")
		checkErr(err)
		res, err := stmt.Exec("gqy", 2)
		checkErr(err)
		affect, err := res.RowsAffected()
		checkErr(err)
		fmt.Println(affect)*/
	//查询数据
	/*rows, err := db.Query("SELECT * FROM userinfo")
	checkErr(err)
	for rows.Next() {
		var uid int
		var username string
		var departname string
		var created string
		err = rows.Scan(&uid, &username, &departname, &created)
		checkErr(err)
		fmt.Println(uid, username, departname, created)

	}*/

	//删除数据
	stmt, err := db.Prepare("delete from userinfo where uid=?")
	checkErr(err)
	res, err := stmt.Exec(2)
	checkErr(err)
	affect, err := res.RowsAffected()
	checkErr(err)
	fmt.Println(affect)
	db.Close()

}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

/*
通过上面的代码 我们可以看出 Go操作mysql数据库是很方便的
关键的几个函数 简单解释
sql.Open() 函数用来打开一个注册过的数据库驱动  Go-MySQL-Driver中注册了mysql这个数据库驱动，第二个参数是DSN(Data Source Name)，它是Go-MySQL-Driver定义的一些数据库链接和配置信息。它支持如下格式：
user@unix(/path/to/socket)/dbname?charset=utf8
user:password@tcp(localhost:5555)/dbname?charset=utf8
user:password@/dbname
user:password@tcp([de:ad:be:ef::ca:fe]:80)/dbname
db.Prepare()函数用来返回准备要执行的sql操作 然后返回准备完毕的执行状态
db.Query()函数用来执行sql返回的Rows结果
stmt.Exec()函数用来执行stmt准备好的sql语句
我们可以看到我们传入的参数 都是=?对应的数据  这样做 的方式 可以一定程度上防止sql注入
*/
