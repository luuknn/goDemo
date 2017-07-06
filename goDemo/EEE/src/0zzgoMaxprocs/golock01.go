// test for Go
//
// Copyright (c) 2015 - Batu <1235355@qq.com>
//
// 创建一个文件存放数据,在同一时刻,可能会有多个Goroutine分别进行对此文件的写操作和读操作.
// 每一次写操作都应该向这个文件写入若干个字节的数据,作为一个独立的数据块存在,这意味着写操作之间不能彼此干扰,写入的内容之间也不能出现穿插和混淆的情况
// 每一次读操作都应该从这个文件中读取一个独立完整的数据块.它们读取的数据块不能重复,且需要按顺序读取.
// 例如: 第一个读操作读取了数据块1,第二个操作就应该读取数据块2,第三个读操作则应该读取数据块3,以此类推
// 对于这些读操作是否可以被同时执行,不做要求. 即使同时进行,也应该保持先后顺序.
//package main
//
//import (
//	"errors"
//	"fmt"
//	"io"
//	"os"
//	"sync"
//	"time"
//)
//
////数据文件的接口类型
//type DataFile interface {
//	// 读取一个数据块
//	Read() (rsn int64, d Data, err error)
//	// 写入一个数据块
//	Write(d Data) (wsn int64, err error)
//	// 获取最后读取的数据块的序列号
//	Rsn() int64
//	// 获取最后写入的数据块的序列号
//	Wsn() int64
//	// 获取数据块的长度
//	DataLen() uint32
//}
//
////数据类型
//type Data []byte
//
////数据文件的实现类型
//type myDataFile struct {
//	f       *os.File     //文件
//	fmutex  sync.RWMutex //被用于文件的读写锁
//	woffset int64        // 写操作需要用到的偏移量
//	roffset int64        // 读操作需要用到的偏移量
//	wmutex  sync.Mutex   // 写操作需要用到的互斥锁
//	rmutex  sync.Mutex   // 读操作需要用到的互斥锁
//	dataLen uint32       //数据块长度
//}
//
////初始化DataFile类型值的函数,返回一个DataFile类型的值
//func NewDataFile(path string, dataLen uint32) (DataFile, error) {
//	f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
//	//f,err := os.Create(path)
//	if err != nil {
//		fmt.Println("Fail to find", f, "cServer start Failed")
//		return nil, err
//	}
//
//	if dataLen == 0 {
//		return nil, errors.New("Invalid data length!")
//	}
//
//	df := &myDataFile{
//		f:       f,
//		dataLen: dataLen,
//	}
//
//	return df, nil
//}
//
////获取并更新读偏移量,根据读偏移量从文件中读取一块数据,把该数据块封装成一个Data类型值并将其作为结果值返回
//
//func (df *myDataFile) Read() (rsn int64, d Data, err error) {
//	// 读取并更新读偏移量
//	var offset int64
//	// 读互斥锁定
//	df.rmutex.Lock()
//	offset = df.roffset
//	// 更改偏移量, 当前偏移量+数据块长度
//	df.roffset += int64(df.dataLen)
//	// 读互斥解锁
//	df.rmutex.Unlock()
//
//	//读取一个数据块,最后读取的数据块序列号
//	rsn = offset / int64(df.dataLen)
//	bytes := make([]byte, df.dataLen)
//	for {
//		//读写锁:读锁定
//		df.fmutex.RLock()
//		_, err = df.f.ReadAt(bytes, offset)
//		if err != nil {
//			//由于进行写操作的Goroutine比进行读操作的Goroutine少,所以过不了多久读偏移量roffset的值就会大于写偏移量woffset的值
//			// 也就是说,读操作很快就没有数据块可读了,这种情况会让df.f.ReadAt方法返回的第二个结果值为代表的非nil且会与io.EOF相等的值
//			// 因此不应该把EOF看成错误的边界情况
//			// so 在读操作读完数据块,EOF时解锁读操作,并继续循环,尝试获取同一个数据块,直到获取成功为止.
//			if err == io.EOF {
//				//注意,如果在该for代码块被执行期间,一直让读写所fmutex处于读锁定状态,那么针对它的写操作将永远不会成功.
//				//切相应的Goroutine也会被一直阻塞.因为它们是互斥的.
//				// so 在每条return & continue 语句的前面加入一个针对该读写锁的读解锁操作
//				df.fmutex.RUnlock()
//				//注意,出现EOF时可能是很多意外情况,如文件被删除,文件损坏等
//				//这里可以考虑把逻辑提交给上层处理.
//				continue
//			}
//		}
//		break
//	}
//	d = bytes
//	df.fmutex.RUnlock()
//	return
//}
//
//func (df *myDataFile) Write(d Data) (wsn int64, err error) {
//	//读取并更新写的偏移量
//	var offset int64
//	df.wmutex.Lock()
//	offset = df.woffset
//	df.woffset += int64(df.dataLen)
//	df.wmutex.Unlock()
//
//	//写入一个数据块,最后写入数据块的序号
//	wsn = offset / int64(df.dataLen)
//	var bytes []byte
//	if len(d) > int(df.dataLen) {
//		bytes = d[0:df.dataLen]
//	} else {
//		bytes = d
//	}
//	df.fmutex.Lock()
//	df.fmutex.Unlock()
//	_, err = df.f.Write(bytes)
//
//	return
//}
//
//func (df *myDataFile) Rsn() int64 {
//	df.rmutex.Lock()
//	defer df.rmutex.Unlock()
//	return df.roffset / int64(df.dataLen)
//}
//
//func (df *myDataFile) Wsn() int64 {
//	df.wmutex.Lock()
//	defer df.wmutex.Unlock()
//	return df.woffset / int64(df.dataLen)
//}
//
//func (df *myDataFile) DataLen() uint32 {
//	return df.dataLen
//}
//
//func main() {
//	//简单测试下结果
//	var dataFile DataFile
//	dataFile, _ = NewDataFile("./mutex_2015_1.dat", 10)
//
//	var d = map[int]Data{
//		1: []byte("gqy"),
//		2: []byte("dxy"),
//		3: []byte("qpf"),
//	}
//
//	//写入数据
//	for i := 1; i < 4; i++ {
//		go func(i int) {
//			wsn, _ := dataFile.Write(d[i])
//			fmt.Println("write i=", i, ",wsn=", wsn, ",success.")
//		}(i)
//	}
//
//	//读取数据
//	for i := 1; i < 4; i++ {
//		go func(i int) {
//			rsn, d, _ := dataFile.Read()
//			fmt.Println("Read i=", i, ",rsn=", rsn, ",data=", d, ",success.")
//		}(i)
//	}
//
//	time.Sleep(10 * time.Second)
//}

/*
####golang同步 锁的使用案例介绍
互斥锁
互斥锁是传统的并发程序对共享资源进行访问控制的主要手段  它由标准库代码包sync中的Mutex结构体类型代表
只有两个公开的方法
Lock Unlock
类型sync.Mutex的零值表示了未被锁定的互斥量
var mutex sync.Mutex
mutex.Lock()
##示例
import (
	"fmt"
	"sync"
	"time"
)

func main() {
	//声明
	var mutex sync.Mutex
	fmt.Println("Lock the lock.(GO)")
	//加锁mutex
	mutex.Lock()
	fmt.Println("The lock is locked.(GO)")
	for i := 0; i < 4; i++ {
		go func(i int) {
			fmt.Printf("Lock the lock.(G%d)\n", i)
			mutex.Lock()
			fmt.Printf("The lock is locked.(G%d)\n", i)
		}(i)
	}
	//休息一会 等待打印结果
	time.Sleep(time.Second)
	fmt.Println("Unlock the lock.(GO)")
	//解锁mutex
	mutex.Unlock()
	fmt.Println("The lock is unlocked.(GO)")
	//休息一会 等待打印结果
	time.Sleep(time.Second)
}
建议 : 同一个互斥锁的成对锁定和解锁操作放在同一层次的代码块中

####读写锁
针对读写操作的互斥锁 它可以分别针对读操作和写操作 进行锁定和解锁操作 遵循的访问控制与互斥锁有所不同
它允许任意读操作同时进行 同一时刻 只允许有一个写操作进行
=============华丽丽分割线========================
并且一个写操作被进行过程中 读操作的进行也是不被允许的
并且一个写操作被执行过程中  读操作的进行也是不被允许的
读写锁控制下的多个写操作之间都是互斥的
写操作和读操作之间也是互斥的
多个读操作之间却不存在互斥关系
读写锁由结构体类型 sync.RWMutex代表
写操作的锁定和解锁
func(*RWMutex) Lock
func(&RWMutex) Unlock
 读操作的锁定和解锁
 func(*RWMutex) Rlock
 func(*RWMutex) RUnlock
 注意
 写解锁在进行的时候会试图唤醒所有所有想进行读锁定而被阻塞的goroutine
 读解锁在进行的时候只会在已无任何读锁定的情况下试图唤醒一个因欲进行写锁定而被阻塞的goroutine
 若对一个未被写锁定的读写锁进行解锁 会引起一个运行时的panic
 而对一个未被读锁定的读写锁进行读解锁  并不会引起恐慌

 ####锁的完整示例
// test for Go
//
// Copyright (c) 2015 - Batu <1235355@qq.com>
//
// 创建一个文件存放数据,在同一时刻,可能会有多个Goroutine分别进行对此文件的写操作和读操作.
// 每一次写操作都应该向这个文件写入若干个字节的数据,作为一个独立的数据块存在,这意味着写操作之间不能彼此干扰,写入的内容之间也不能出现穿插和混淆的情况
// 每一次读操作都应该从这个文件中读取一个独立完整的数据块.它们读取的数据块不能重复,且需要按顺序读取.
// 例如: 第一个读操作读取了数据块1,第二个操作就应该读取数据块2,第三个读操作则应该读取数据块3,以此类推
// 对于这些读操作是否可以被同时执行,不做要求. 即使同时进行,也应该保持先后顺序.
package main

import (
    "fmt"
    "sync"
    "time"
    "os"
    "errors"
    "io"
)

//数据文件的接口类型
type DataFile interface {
    // 读取一个数据块
    Read() (rsn int64, d Data, err error)
    // 写入一个数据块
    Write(d Data) (wsn int64, err error)
    // 获取最后读取的数据块的序列号
    Rsn() int64
    // 获取最后写入的数据块的序列号
    Wsn() int64
    // 获取数据块的长度
    DataLen() uint32
}

//数据类型
type Data []byte

//数据文件的实现类型
type myDataFile struct {
    f *os.File  //文件
    fmutex sync.RWMutex //被用于文件的读写锁
    woffset int64 // 写操作需要用到的偏移量
    roffset int64 // 读操作需要用到的偏移量
    wmutex sync.Mutex // 写操作需要用到的互斥锁
    rmutex sync.Mutex // 读操作需要用到的互斥锁
    dataLen uint32 //数据块长度
}

//初始化DataFile类型值的函数,返回一个DataFile类型的值
func NewDataFile(path string, dataLen uint32) (DataFile, error){
    f, err := os.OpenFile(path, os.O_APPEND|os.O_RDWR|os.O_CREATE, 0666)
    //f,err := os.Create(path)
    if err != nil {
        fmt.Println("Fail to find", f, "cServer start Failed")
        return nil, err
    }

    if dataLen == 0 {
        return nil, errors.New("Invalid data length!")
    }

    df := &myDataFile{
        f : f,
        dataLen:dataLen,
    }

    return df, nil
}

//获取并更新读偏移量,根据读偏移量从文件中读取一块数据,把该数据块封装成一个Data类型值并将其作为结果值返回

func (df *myDataFile) Read() (rsn int64, d Data, err error){
    // 读取并更新读偏移量
    var offset int64
    // 读互斥锁定
    df.rmutex.Lock()
    offset = df.roffset
    // 更改偏移量, 当前偏移量+数据块长度
    df.roffset += int64(df.dataLen)
    // 读互斥解锁
    df.rmutex.Unlock()

    //读取一个数据块,最后读取的数据块序列号
    rsn = offset / int64(df.dataLen)
    bytes := make([]byte, df.dataLen)
    for {
        //读写锁:读锁定
        df.fmutex.RLock()
        _, err = df.f.ReadAt(bytes, offset)
        if err != nil {
            //由于进行写操作的Goroutine比进行读操作的Goroutine少,所以过不了多久读偏移量roffset的值就会大于写偏移量woffset的值
            // 也就是说,读操作很快就没有数据块可读了,这种情况会让df.f.ReadAt方法返回的第二个结果值为代表的非nil且会与io.EOF相等的值
            // 因此不应该把EOF看成错误的边界情况
            // so 在读操作读完数据块,EOF时解锁读操作,并继续循环,尝试获取同一个数据块,直到获取成功为止.
            if err == io.EOF {
                //注意,如果在该for代码块被执行期间,一直让读写所fmutex处于读锁定状态,那么针对它的写操作将永远不会成功.
                //切相应的Goroutine也会被一直阻塞.因为它们是互斥的.
                // so 在每条return & continue 语句的前面加入一个针对该读写锁的读解锁操作
                df.fmutex.RUnlock()
                //注意,出现EOF时可能是很多意外情况,如文件被删除,文件损坏等
                //这里可以考虑把逻辑提交给上层处理.
                continue
            }
        }
        break
    }
    d = bytes
    df.fmutex.RUnlock()
    return
}

func (df *myDataFile) Write(d Data) (wsn int64, err error){
    //读取并更新写的偏移量
    var offset int64
    df.wmutex.Lock()
    offset = df.woffset
    df.woffset += int64(df.dataLen)
    df.wmutex.Unlock()

    //写入一个数据块,最后写入数据块的序号
    wsn = offset / int64(df.dataLen)
    var bytes []byte
    if len(d) > int(df.dataLen){
        bytes = d[0:df.dataLen]
    }else{
        bytes = d
    }
    df.fmutex.Lock()
    df.fmutex.Unlock()
    _, err = df.f.Write(bytes)

    return
}

func (df *myDataFile) Rsn() int64{
    df.rmutex.Lock()
    defer df.rmutex.Unlock()
    return df.roffset / int64(df.dataLen)
}

func (df *myDataFile) Wsn() int64{
    df.wmutex.Lock()
    defer df.wmutex.Unlock()
    return df.woffset / int64(df.dataLen)
}

func (df *myDataFile) DataLen() uint32 {
    return df.dataLen
}

func main(){
    //简单测试下结果
    var dataFile DataFile
    dataFile,_ = NewDataFile("./mutex_2015_1.dat", 10)

    var d=map[int]Data{
        1:[]byte("batu_test1"),
        2:[]byte("batu_test2"),
        3:[]byte("test1_batu"),
    }

    //写入数据
    for i:= 1; i < 4; i++ {
        go func(i int){
            wsn,_ := dataFile.Write(d[i])
            fmt.Println("write i=", i,",wsn=",wsn, ",success.")
        }(i)
    }

    //读取数据
    for i:= 1; i < 4; i++ {
        go func(i int){
            rsn,d,_ := dataFile.Read()
            fmt.Println("Read i=", i,",rsn=",rsn,",data=",d, ",success.")
        }(i)
    }

    time.Sleep(10 * time.Second)
}





























*/
