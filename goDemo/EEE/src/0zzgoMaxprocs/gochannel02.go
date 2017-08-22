package main

/*
Go Channel
1hchan struct
channel的底层数据结构是hchan struct
type hchan struct{
qcount uint //队列中的数据个数
dataqsiz uint //channel的大小
buf unsafe.pointer//存放数据的环形数组
elemsize uint16 //channel中数据类型的大小
closed uint32 //表示channel是否关闭
elemtype *_type //元素数据类型
sendx uint//send 的数组索引
recvx uint //recv的数组索引
recvq waitq//<-ch recv行为阻塞在channel上的goroutine队列
sendq waitq //ch<- 阻塞在channel上的goroutine队列
lock mutex
}
type waitq struct{
first *sudog
last *sudog
}
type sudog struct {
g *g
selectdone *uint32
next *sudog
prev *sudog
elem unsafe.Pointer
acquiretime int64
releasetime int64
ticket uint32
parent uint32
parent *sudog
waitlink *sudog
waittail *sudog
c *hchan
}
上面直接对各个字段做了解释 我们可以看到channel其实就是一个队列加一个锁
只不过这个锁是一个轻量级锁 其中recvq是读操作阻塞在channel的goroutine列表
sendq是写操作阻塞在channel的goroutine列表
列表的实现是sudog 其实就是一个对g的结构的封装

2 make
通过make创建channel对应的代码如下
func makechan(t *chantype,size int64) *hchan{
elem :=t.elem
if elem.size>=1<<16{
throw("makechan:invalid channel element type")
}
if hchanSize%maxAlign !=0 ||elem.align>maxAlign{
throw("makechan: bad alignment")
}
 if size < 0 || int64(uintptr(size)) != size || (elem.size > 0 && uintptr(size) > (_MaxMem-hchanSize)/elem.size) {
        panic(plainError("makechan: size out of range"))
    }
var c *hchan
if elem.kind$kondNoPointers!=0||size==0{
 c = (*hchan)(mallocgc(hchanSize+uintptr(size)*elem.size, nil, true))
        if size > 0 && elem.size != 0 {
            c.buf = add(unsafe.Pointer(c), hchanSize)
        } else {
            // race detector uses this location for synchronization
            // Also prevents us from pointing beyond the allocation (see issue 9401).
            c.buf = unsafe.Pointer(c)
        }
    } else {
        c = new(hchan)
        c.buf = newarray(elem, int(size))
    }
    c.elemsize = uint16(elem.size)
    c.elemtype = elem
    c.dataqsiz = uint(size)

    if debugChan {
        print("makechan: chan=", c, "; elemsize=", elem.size, "; elemalg=", elem.alg, "; dataqsiz=", size, "\n")
    }
    return c
}
最前面的两个if是一些异常判断 元素类型大小限制和对齐限制 第三个if也很明显 判断size大小是否小于0 或者过大
int64(uintptr(size))!=size 这句也是判断size是否为负 值得一说的是最后面的判断条件
uintprt(size)>(_MaxMem-hchanSize)/elem.size
_MaxMem 这个是Arena区域的最大值 用来分配给堆  也就是说 channel是在堆上分配的
再往下就可以看到分配的代码了 如果channel内数据类型不含有指针且size>0 则将分配在连续的内存区域
如果size=0 实际上buf是不分配空间的
if elem.kind&kindNopointers!=0||size ==0{
c =(*hchan)(mallocgc(hchanSize+uintptr(size)*elem.size, nil, true))
 if size > 0 && elem.size != 0 {
        c.buf = add(unsafe.Pointer(c), hchanSize)
    } else {
             c.buf = unsafe.Pointer(c)
    }
}
除了上面的情况 剩下的 也就是size>0 channel 和channel.buf是分别进行分配的 剩下的代码是剩下字段的处理
else{
c=new(hchan)
c.buf=newarray(elem,int(size))
}
总结一下 make chan的过程是在堆上进行分配 返回一个hchan的指针

3 send
send 也就是ch<-x 对应的函数如下




















*/
