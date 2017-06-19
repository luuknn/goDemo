package main

/*
Prepared 剖析
prepare
前面我们已经学习了sql的基本curd操作 总体而言  有两类操作 Query和Exec 前者返回数据库记录 后者返回数据库影响或者插入相关的结果
上面两种操作 多数是针对单次操作的查询 如果需要批量插入一堆数据 就可以使用Prepared语句 golang处理prepared 语句有其独特的行为 了解其底层的实现 对用好它十分重要

查询
我们可以使用 Query方式查询记录 Query函数提供了两种选择 第一种情况下参数是拼接好的sql 另外一种情况 第一参数是带有占位符的sql
第二个参数是sql的实际参数
rows,err:=db.Query("SELECT * FROM user Where gid=1")
rows,err:=db.Query("SELECT * FROM user Where gid= ?",1)
上面两种方式都能获取数据 那么他们的底层实现是一样的吗 实际上 上面两种方式的底层通信不完全一样
一种是plaintext方式 另外一种是prepared方式

prepared
所谓prepared 即带有占位符的sql语句 客户端将该语句和参数发给mysql服务器 mysql服务器编译成一个prepared语句
这个语句可以根据不同的参数多次调用  prepared语句执行的方式如下
1 准备prepare语句  2执行prepared语句和参数  3 关闭prepared语句
之所以会出现prepare语句方式 主要因为 这样有下面的两个好处
1 避免通过引号组装拼接sql语句 避免sql注入带来的安全风险  2 可以多次执行的sql语句
单纯的看prepare语句的好处 会下意识的觉得既然如此  都使用prepared语句查询不就好了么
其实不然 关于prepared语句注意事项  稍后讨论
*/

/*
golang的pliantext和prepare查询方式
现在我们再回顾下上面调用Query函数的两个操作  对于第一个操作 执行plianttext的sql语句
先看db.Quert方法
//Query executes a query that returns rows ,typically a SELECT
//The args are for any placeholder parameters in the query
func(db *DB) Query(query string,args ...interface{})(*Rows,error){
var rows *Rows
var err error
for i:=0;i<maxBadConnRetries;i++{
rows,err=db.query(query,args,cachedOrNewConn)
if err!=driver.ErrBadConn{
break
}
}
if err == driver.ErrBadConn{
return db.query(query,args,alwaysNewConn)
}
return rows,err
}
Query方法我们很熟悉了  它的内部调用db.query方法 并且根据连接重连的状况选择是cachedorNewConn模式还是alwaysNewConn模式
前者会从返回一个cached连接或者等待一个可用连接 甚至也可能建立一个新的连接 后者表示 打开连接时的策略
这就是签名所说的retry10次连接
func (db *DB) query(query string, args []interface{}, strategy connReuseStrategy) (*Rows, error) {
 ci, err := db.conn(strategy)
 if err != nil {
  return nil, err
 }
 return db.queryConn(ci, ci.releaseConn, query, args)
}
query方法的逻辑很简单 通过db.conn 方法返回一个新创建或者缓存的空闲连接 driverConn 随机调用 queryConn方法
//queryConn executes a query on the given connection
//The connection gets released by the releaseConn function
func(db *DB)queryConn(dc *driverConn,releaseConn func(error),query string,args []interface{})(*Rows,error)
 // 判断驱动是否实现了Queryer
 if queryer, ok := dc.ci.(driver.Queryer); ok {
  dargs, err := driverArgs(nil, args)
  if err != nil {
   releaseConn(err)
   return nil, err
  }
  dc.Lock()
  rowsi, err := queryer.Query(query, dargs)  // 调用驱动的查询方法  connection.go 第305行
  dc.Unlock()
  if err != driver.ErrSkip {  // 不带参数的返回
   if err != nil {
    releaseConn(err)
    return nil, err
   }
   // Note: ownership of dc passes to the *Rows, to be freed
   // with releaseConn.
   rows := &Rows{
    dc:          dc,
    releaseConn: releaseConn,
    rowsi:       rowsi,
   }
   return rows, nil
  }
 }
 dc.Lock()
 si, err := dc.ci.Prepare(query)  // 带参数的返回，创建prepare对象
 dc.Unlock()
 if err != nil {
  releaseConn(err)
  return nil, err
 }
 ds := driverStmt{dc, si}
 rowsi, err := rowsiFromStatement(ds, args...)   // 执行语句
 if err != nil {
  dc.Lock()
  si.Close()
  dc.Unlock()
  releaseConn(err)
  return nil, err
 }
 // Note: ownership of ci passes to the *Rows, to be freed
 // with releaseConn.
 rows := &Rows{
  dc:          dc,
  releaseConn: releaseConn,
  rowsi:       rowsi,
  closeStmt:   si,
 }
 return rows, nil
}


queryConn 函数 内容比较多 先判断驱动是否实现了Query 如果实现了即调用其Query方法
方法会针对sql查询语句做查询 例如myslq的驱动如下connection.go 第305行左右 即
func (mc *mysqlConn) Query(query string, args []driver.Value) (driver.Rows, error) {
 if mc.netConn == nil {
  errLog.Print(ErrInvalidConn)
  return nil, driver.ErrBadConn
 }
 if len(args) != 0 {
  if !mc.cfg.InterpolateParams {
   return nil, driver.ErrSkip
  }
  // try client-side prepare to reduce roundtrip
  prepared, err := mc.interpolateParams(query, args)
  if err != nil {
   return nil, err
  }
  query = prepared
  args = nil
 }
 // Send command
 err := mc.writeCommandPacketStr(comQuery, query)
 if err == nil {
  // Read Result
  var resLen int
  resLen, err = mc.readResultSetHeaderPacket()
  if err == nil {
   rows := new(textRows)
   rows.mc = mc
   if resLen == 0 {
    // no columns, no more data
    return emptyRows{}, nil
   }
   // Columns
   rows.columns, err = mc.readColumns(resLen)
   return rows, err
  }
 }
 return nil, err
}
Query 先检查参数是否为0  然后调用writeCommandPacketStr方法执行sql并通过readResultSetHeaderPacket读取数据库服务返回的结果
如果参数不为0  会先判断是否是prepare语句 这里会返回一个driver.ErrSkip的错误
把函数执行权再返回到queryConn函数中 然后再调用 si,err:=dc.ci.Prepare(query)创建一个Stmt对象  接下来调用rowsiFromStatement执行查询
func rowsiFromStatement(ds driverStmt, args ...interface{}) (driver.Rows, error) {
 ds.Lock()
 want := ds.si.NumInput()
 ds.Unlock()
 // -1 means the driver doesn't know how to count the number of
 // placeholders, so we won't sanity check input here and instead let the
 // driver deal with errors.
 if want != -1 && len(args) != want {
  return nil, fmt.Errorf("sql: statement expects %d inputs; got %d", want, len(args))
 }
 dargs, err := driverArgs(&ds, args)
 if err != nil {
  return nil, err
 }
 ds.Lock()
 rowsi, err := ds.si.Query(dargs)
 ds.Unlock()
 if err != nil {
  return nil, err
 }
 return rowsi, nil
}
rowsiFromStatement方法会调用驱动的ds.si.Query(dargs)方法，执行最后的查询。大概再statement.go的第84行

func (stmt *mysqlStmt) Query(args []driver.Value) (driver.Rows, error) {
 if stmt.mc.netConn == nil {
  errLog.Print(ErrInvalidConn)
  return nil, driver.ErrBadConn
 }
 // Send command
 err := stmt.writeExecutePacket(args)
 if err != nil {
  return nil, err
 }
 mc := stmt.mc
 // Read Result
 resLen, err := mc.readResultSetHeaderPacket()
 if err != nil {
  return nil, err
 }
 rows := new(binaryRows)
 if resLen > 0 {
  rows.mc = mc
  // Columns
  // If not cached, read them and cache them
  if stmt.columns == nil {
   rows.columns, err = mc.readColumns(resLen)
   stmt.columns = rows.columns
  } else {
   rows.columns = stmt.columns
   err = mc.readUntilEOF()
  }
 }
 return rows, err
}
调用 stmt和参数执行sql查询。查询完毕之后，返回到queryConn方法中，使用releaseConn释放查询的数据库连接。
*/
/*
自定义prepare查询
从query查询可以看到 对于占位符的prepare语句 go内部通过 dc.ci.Prepare(query)会自动创建一个stmt对象
其实我们也可以自定义stmt语句 使用方法如下
stmt,err:=db.Prepare("SELECT * FROM user where gid= ?")
if err!=nil{
log.Fatalln(err)
}
defer stmt.Close()
rows,err:=stmt.Query(1)
if err!=nil{
log.Fatalln(err)
}
即通过Prepare方法创建一个stmt对象 然后执行stmt对象的Query方法 得到sql.Rows结果集 最后关闭stmt.Close
这个过程就和之前所说的prepare三步骤匹配了
创建stmt的prepare方式是golang的一个设计 其目的是Prepare once execute many times
为了批量执行sql语句 但是通常会造成所谓的三次网络请求(three network round-trips)即prepering executing和closing三次请求

对于大多数数据库 prepared的过程都是 先发送一个带占位符的sql语句到服务器 服务器 返回一个statementid
然后再把这个id和参数发送给服务器执行 最后再发送关闭statement命令

golang的实现了连接池 处理prepare方式也需要特别注意 调用Prepare方法返回stmt的时候 golang会在
某个空闲的连接上进行prepare语句 然后就把连接释放回到连接池  可是golang会记住这个连接 当需要执行参数的时候 就再次找到之前记住的连接进行执行 等到stmt.Close调用的时候 再释放该连接

在执行参数的时候 如果记住的连接正处于忙碌阶段  此时 golang将会从新选一个新的空闲连接进行prepare(re-prepare)
当然  即使是reprepare 同样也会遇到刚才的问题 那么将会一而再再而三的进行reprepare直到找到空闲连接进行查询的时候

这种情况将会导致leak连接的情况 尤其是高并发的情况 将会导致大量的prepare过程
因此使用stmt的情况需要仔细考虑应用场景 通常在应用程序中 多次执行同一个sql语句的情况并不多 因此减少prepare语句的使用

之前有一个疑问 是不是所有sql语句都不能带占位符 因为这是prepare语句 只要看了以便database/sql和驱动的源码 才恍然大悟 对于
query(prepare,args)的方式 正如我们前面所分析的  database/sql会使用 ds.si.Query(dargs)创建stmt 然后就立即执行
prepare和参数 最后关闭stmt  整个过程都是同一个连接上完成 因此不存在reprepare的情况 当然也无法使用所谓的创建一次 执行多次的目的

对于prepare的使用方式 基于其好处和缺点  我们将会在后面再讨论 目前需要注意的大致就是：
1 单次查询不需要使用prepared 每次使用stmt语句都是三次网络请求次数 prepared execute close
2不要循环中创建prepare语句      3  注意关闭stmt

尽管会有reprepare过程  这些操作依然是database/sql帮我们所做的  与连接retry10次一样  开发者无需担心

对于Query操作 如此 同理 Exec操作也一样

总结
目前我们学习database/sql提供两类查询操作 Query和Exec方法 他们都可以使用plaintext和prepare方式查询
对于后者 可以有效的避免数据库注入  而prepare方式又可以有显示的声明stmt对象 也有隐藏的方式
显示的创建stmt会有3次网络请求  创建 -执行 - 关闭  再批量操作 可以考虑这种做法  另外一种创建方式创建prepare后就
执行  因此 不会因为reprepare导致高并发下的leak连接问题

具体使用哪种方式 还得基于应用场景 安全过滤 和连接管理等考虑 至此 关于查询和执行操作已经介绍了很多
关系型数据库的另外一个特性就是关系和事务处理  下一节 我们将会讨论database/sql的数据库事务功能












*/
