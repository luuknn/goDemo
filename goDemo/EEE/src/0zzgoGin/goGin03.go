package main

//import (
//	"database/sql"
//	//"fmt"
//	_ "github.com/go-sql-driver/mysql"
//	"gopkg.in/gin-gonic/gin.v1"
//	"log"
//	"net/http"
//)
//
//type Person struct {
//	Id        int    `json:"id" form:"id"`
//	FirstName string `json:"first_name" form:"first_name"`
//	LastName  string `json:"lastI_name" form:"last_name"`
//}
//
//func main() {
//	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
//	defer db.Close()
//	if err != nil {
//		log.Fatalln(err)
//	}
//	db.SetMaxIdleConns(20)
//	db.SetMaxOpenConns(20)
//	if err := db.Ping(); err != nil {
//		log.Fatalln(err)
//	}
//	router := gin.Default()
//	router.GET("/", func(c *gin.Context) {
//		c.String(http.StatusOK, "it works..")
//	})
//
//	router.GET("/persons", func(c *gin.Context) {
//		rows, err := db.Query("select id ,first_name,last_name from person")
//		defer rows.Close()
//		if err != nil {
//			log.Fatalln(err)
//		}
//		persons := make([]Person, 0)
//		for rows.Next() {
//			var person Person
//			rows.Scan(&person.Id, &person.FirstName, &person.LastName)
//			persons = append(persons, person)
//		}
//		if err = rows.Err(); err != nil {
//			log.Fatalln(err)
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"persons": persons,
//		})
//	})
//
//	router.GET("/person/:id", func(c *gin.Context) {
//		id := c.Param("id")
//		var person Person
//		err := db.QueryRow("select id,first_name,last_name from person where id =?", id).Scan(
//			&person.Id, &person.FirstName, &person.LastName,
//		)
//		if err != nil {
//			//log.Fatalln(err)
//			c.JSON(http.StatusOK, gin.H{
//				"hello": nil,
//			})
//			return
//
//		}
//		c.JSON(http.StatusOK, gin.H{
//			"person": person,
//		})
//
//	})
//
//	router.Run(":8000")
//}

/*
Gin实战 Gin+MySQL简单的restful风格的API
我们已经了解了golang的gin框架  对于webservice服务  restful风格几乎一统天下 gin也天然的支持restful
下面就使用 gin写一个简单的服务 麻雀虽小 五脏俱全  我们先以一个单文件开始 然后再逐步 分解模块成包 组织代码
import (
	"gopkg.in/gin-gonic/gin.v1"
	"net/http"
)

func main() {
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "it works")
	})
	router.Run(":8000")
}
数据库
安装完毕框架  完成一次请求响应之后 接下来就是安装数据库驱动和初始化数据相关的操作了  首先我们需要新建数据表
一个极其简单的数据表
CREATE TABLE `person` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `first_name` varchar(40) NOT NULL DEFAULT '',
  `last_name` varchar(40) NOT NULL DEFAULT '',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
创建数据表之后 初始化数据库连接池
func main() {
	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "it works")
	})
	router.Run(":8000")
}
使用sql.Open方法会创建一个数据库连接池db 这个db不是数据库连接 它是一个连接池
只有当真正数据库通信的时候才创建连接 例如这里的db.Ping的操作 db.SetMaxIdleConns
db.SetMaxOpenConns(20) 分别设置数据库的空闲连接和最大打开连接 即向mysql服务端发出的所有连接的最大数目
如果不设置 默认都是0 表示打开的连接没有限制 压测会发现 大量的TIME_WAIT状态的连接  虽然mysql的连接数没有上升
设置了这两个参数之后 不存在大量TIME_WAIT状态的连接了 而且qps也没有明显的变化 出于对数据库的保护 最好设置这个连接参数

CURD增删改查
增
func main() {
 ...

 router.POST("/person", func(c *gin.Context) {
  firstName := c.Request.FormValue("first_name")
  lastName := c.Request.FormValue("last_name")
  rs, err := db.Exec("INSERT INTO person(first_name, last_name) VALUES (?, ?)", firstName, lastName)
  if err != nil {
   log.Fatalln(err)
  }
  id, err := rs.LastInsertId()
  if err != nil {
   log.Fatalln(err)
  }
  fmt.Println("insert person Id {}", id)
  msg := fmt.Sprintf("insert successful %d", id)
  c.JSON(http.StatusOK, gin.H{
   "msg": msg,
  })
 })

 ...
}
执行费query操作 使用db的Exec方法 在mysql中使用？做占位符 最后我们把插入后的id返回给客户端

查
查询列表
上面我们增加了一条记录 下面就获取记录 查一般有两个操作 一个是查询列表 其次就是查询具体的某一条记录
两种大同小异
为了给查询结果绑定到golang的变量或对象 我们需要先定义一个结构来绑定对象 在main函数的上方定义Person结构
然后查询我们的数据列表

读取mysql的数据需要一个绑定的过程 db.Query方法返回一个rows对象 这个数据库连接随即也转移到这个对象
因此我们需要定义row.Close操作  然后创建一个[]Person的切片
使用make 而不是直接使用 var persons []Person 的声明方式  还是有所差别的  使用make的方式 当数组切片
没有元素的时候 Json会返回[] 如果直接声明 json会返回null
接下来就是使用rows对象的next方法 遍历所查询的数据 一个个绑定到person对象上  最后append到persons切片

查询单条记录
查询列表需要迭代rows对象 查询单个记录 就没有这么麻烦了 虽然也可以迭代一条记录的结果集
因为查询单个记录的操作实在太常用了  因此golang的database/sql也专门提供了查询方法
import (
	"database/sql"
	//"fmt"
	_ "github.com/go-sql-driver/mysql"
	"gopkg.in/gin-gonic/gin.v1"
	"log"
	"net/http"
)

type Person struct {
	Id        int    `json:"id" form:"id"`
	FirstName string `json:"first_name" form:"first_name"`
	LastName  string `json:"lastI_name" form:"last_name"`
}

func main() {
	db, err := sql.Open("mysql", "root:zxcvb110test@@tcp(rm-uf680nxer55d4wnm4o.mysql.rds.aliyuncs.com:3306)/cmistest?charset=utf8")
	defer db.Close()
	if err != nil {
		log.Fatalln(err)
	}
	db.SetMaxIdleConns(20)
	db.SetMaxOpenConns(20)
	if err := db.Ping(); err != nil {
		log.Fatalln(err)
	}
	router := gin.Default()
	router.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "it works..")
	})

	router.GET("/persons", func(c *gin.Context) {
		rows, err := db.Query("select id ,first_name,last_name from person")
		defer rows.Close()
		if err != nil {
			log.Fatalln(err)
		}
		persons := make([]Person, 0)
		for rows.Next() {
			var person Person
			rows.Scan(&person.Id, &person.FirstName, &person.LastName)
			persons = append(persons, person)
		}
		if err = rows.Err(); err != nil {
			log.Fatalln(err)
		}
		c.JSON(http.StatusOK, gin.H{
			"persons": persons,
		})
	})

	router.GET("/person/:id", func(c *gin.Context) {
		id := c.Param("id")
		var person Person
		err := db.QueryRow("select id,first_name,last_name from person where id =?", id).Scan(
			&person.Id, &person.FirstName, &person.LastName,
		)
		if err != nil {
			//log.Fatalln(err)
			c.JSON(http.StatusOK, gin.H{
				"hello": nil,
			})
			return

		}
		c.JSON(http.StatusOK, gin.H{
			"person": person,
		})

	})

	router.Run(":8000")
}
查询单个记录有一个小问题 当数据不存在的时候  同样也会跑出一个错误 粗暴的使用log退出有点不妥
返回一个nil的时候 万一真的是因为错误  比如sql错误 这种情况如何解决 还需要具体场景设计程序

改
增删改查 下面进行更新的操作了  前面增加记录 我们使用了urlencode的方式提交 更新的api我们自动匹配绑定content-type
router.PUT("/person/:id", func(c *gin.Context) {
 cid := c.Param("id")
 id, err := strconv.Atoi(cid)
 person := Person{Id: id}
 err = c.Bind(&person)
 if err != nil {
  log.Fatalln(err)
 }
 stmt, err := db.Prepare("UPDATE person SET first_name=?, last_name=? WHERE id=?")
 defer stmt.Close()
 if err != nil {
  log.Fatalln(err)
 }
 rs, err := stmt.Exec(person.FirstName, person.LastName, person.Id)
 if err != nil {
  log.Fatalln(err)
 }
 ra, err := rs.RowsAffected()
 if err != nil {
  log.Fatalln(err)
 }
 msg := fmt.Sprintf("Update person %d successful %d", person.Id, ra)
 c.JSON(http.StatusOK, gin.H{
  "msg": msg,
 })
})

删
最后一个操作 就是删除了  删除所需要的功能特性 上面的例子都覆盖了 实现删除也就特别简单了
router.DELETE("/person/:id", func(c *gin.Context) {
 cid := c.Param("id")
 id, err := strconv.Atoi(cid)
 if err != nil {
  log.Fatalln(err)
 }
 rs, err := db.Exec("DELETE FROM person WHERE id=?", id)
 if err != nil {
  log.Fatalln(err)
 }
 ra, err := rs.RowsAffected()
 if err != nil {
  log.Fatalln(err)
 }
 msg := fmt.Sprintf("Delete person %d successful %d", id, ra)
 c.JSON(http.StatusOK, gin.H{
  "msg": msg,
 })
})
我们可以使用删除接口 把数据都删除了 再来验证上面post接口获取列表的时候 当记录没有的时候
切片被json序列化[]还是null

至此 CURD操作的restful风格api 已经完成





































*/
