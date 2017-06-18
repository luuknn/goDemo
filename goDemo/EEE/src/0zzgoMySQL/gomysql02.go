package main

/*
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM posts WHERE id >0")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cols)
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}
	var results []map[string]string
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			log.Fatalln(err)
		}
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	for k, v := range results {
		fmt.Println(k, v)
	}
}
*/
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

rows.Scan 原理
结果集方法 Scan 可以把数据库取出的字段赋值给指定的数据结构 它的参数是一个空接口的切片 这就意味着可以传入任何值
通常把需要赋值的目标变量的指针当成参数传入 它能将 数据库取出的值赋值到指针值对象上

var var1,var2 string
err=rows.Scan(&var1,&var2)

在一些特殊案例中 如果你不想把值赋值给指定的目标变量 那么需要使用 *sql.RawBytes类型
大多数情况下我们不必这么做 但是还是需要注意 在db.QueryRow().Scan()中不能使用sql.RawBytes

Scan还会帮我们自动推断除数据字段匹配目标变量 比如有个数据库字段的类型是varchar 而他的值是一个数字串1 如果
我们定义的目标变量是string 则scan赋值后目标变量是数字string 如果声明的目标变量是一个数字类型 那么scan
会自动调用 strconv.ParseInt()或者 strconv.ParseInt()方法将字段转换成和声明的目标变量一致的类型 当然
如果有些字段无法转换成功 则会返回错误 因此在调用scan后都需要检查错误
 var world int
 err=stmt.QueryRow(1).Scan(&world)//此时scan会把字段转变成数字整型给world变量

 var world string
 err=stmt.QueryRow(1).Scan(&world)//此时 scan取出的字段就是字串 同样的如果字段是int类型 声明的变量是string类型
 scan也会自动将数字转换成字符串赋值给目标变量

 空值处理
 数据库有一个特殊类型 NULL空值 可是NULL不能通过scan直接跟普通变量赋值 甚至也不能将null赋值给nil对于null必须
 指定特殊的类型 这些类型定义在database/sql库中 例如sql.NullFloat64 如果在标准库中找不到匹配的类型 可以尝试在驱动中寻找
 var (
   s1 string
	s2 sql.Null
	String i1 int
	f1 float64
	f2 float64
)
// 假设数据库的记录为 ["hello", NULL, 12345, "12345.6789", "not-a-float"]
err = rows.Scan(&s1, &s2, &i1, &f1, &f2) if err != nil {
log.Fatal(err) }
 sql: Scan error on column index 4: converting string "not-a- oat" to a  oat64: strconv.ParseFloat: parsing "not-a- oat": invalid syntax

 如果忽略err，强行读取目标变量，可以看到最后一个值转换错误会处理，而不是抛出异常：
 err = rows.Scan(&s1, &s2, &i1, &f1, &f2)
log.Printf("%q %#v %d %f %f", s1, s2, i1, f1, f2)
// 输出
 "hello" sql.NullString{String:"", Valid:false} 12345 12345.678900
0.000000

可以看到，除了最后一个转换失败变成了零值之外，其他都正常的取出了值，尤其是null匹配了NullString类型的目标变量。
对于null的操作，通常仍然需要验证：
var world sql.NullString
err := db.QueryRow("SELECT world FROM hello WHERE id = ?", id).Scan(&world)
...
if world.Valid {
	  wrold.String
} else {
    // 数据库的value是不是null的时候，输出 world的字符串值， 空字符串
    world.String
}
对应的，如果world字段是一个int，那么声明的目标变量类似是sql.NullInt64，读取其值的方法为world.Int64。

但是有时候我们并不关心值是不是Null,我们只需要吧他当一个空字符串来对待就行。这时候我们可以使用[]byte（null byte[]可以转化为空string）
var world []byte
err := db.QueryRow("SELECT world FROM hello WHERE id = ?", id).Scan(&world)
...
log.Println(string(real_name)) // 有值则取出字串值，null则转换成 空字串。

*/
/*
自动匹配字段
在执行查询的时候 我们定义了目标变量 同时查询的时候也写明了字段 如果不指明字段 或者字段的顺序 和查询的不一样
都有可能出错 因此如果能够自动匹配查询的字段值 将会十分节省代码 同时也易于维护
go提供了 Columns方法用获取字段名 与大多数函数一样 读取失败将会 返回一个err因此需要错误检查
cols,err :=rows.Columns()
if err!=nil{
log.Fatalln(err)
}
对于不定字段的查询 我们可以定义一个map的key和value用来表示数据库一条记录的row值
通过rows.Columbs得到col作为map的key值
import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

func main() {
	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	rows, err := db.Query("SELECT * FROM posts WHERE id >0")
	if err != nil {
		log.Fatalln(err)
	}
	defer rows.Close()
	cols, err := rows.Columns()
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println(cols)
	vals := make([][]byte, len(cols))
	scans := make([]interface{}, len(cols))
	for i := range vals {
		scans[i] = &vals[i]
	}
	var results []map[string]string
	for rows.Next() {
		err = rows.Scan(scans...)
		if err != nil {
			log.Fatalln(err)
		}
		row := make(map[string]string)
		for k, v := range vals {
			key := cols[k]
			row[key] = string(v)
		}
		results = append(results, row)
	}
	for k, v := range results {
		fmt.Println(k, v)
	}
}
数据库有三个字段 我们使用*取出所有字段 使用rows.Columns()获取字段名 是一个string的数组
然后创建一个切片vals  用来存放取出来的数据结果 类似是byte的切片 接下来还需要定义一个切片 这个切片用来scan
将数据库的值复制给到它
完成这一步后 vals则得到了scan复制给他的值 因为是byte的切片 因此在循环一次 将其转换成string即可
转换后的row即我们取出的数据行值  最后组装到result切片中

Exec
前面介绍了很多关于查询方面的内容 查询是读方便的内容 对于写 插入更新和删除 这类操作和query不太一样
写的操作 只关系是否写成功了
database/sql 提供了 Exec方法用于执行写的操作
我们也见识到了 Exec返回一个sql.Result类型 它有两个方法 LastInsertId 和RowsAffexted
返回一个数据库自增的id 这是一个int64类型的值
Exec执行完毕之后 连接会立即释放回到连接池中  因此不需要像 query那样再手动调用 row的close方法

总结
目前 我们大致了解了数据库的CRUD操作 对于读的操作 需要定义目标变量才能scan数据记录
scan会智能的帮我们转换一些数据 取决于定义的目标变量类型 对于null的处理
可以使用databa/sql 或驱动提供的类型声明 也可以使用 []byte 将其转换成空字符串 除了读数据之外
对于写的操作 提供了Exec方法

在实际应用中 与数据库交互 往往写的sql语句还带有参数 这类sql可以称之为prepare语句
prepare语句有很多好处 可以防止sql注入 可以批量执行等 但是prepare的连接管理有其自己的机制 也有其使用上的陷阱
关于prepare的使用 将会在以后进行讨论。
*/
