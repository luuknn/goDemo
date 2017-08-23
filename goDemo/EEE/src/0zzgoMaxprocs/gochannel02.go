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
func chansend1(c *hchan,elem unsafe.Pointer){
chansend(c,elem,true,getchallerpc(unsfae.Pointer(&c)))
}
func chansend(c *hchan,ep unsafe.Pointer,block bool,callerpc uintptr) bool{
if c ==nil{
if !block{
return false
}
gopark(nil,nil,"chan send (nil chan)",traceEvGoStop,2)
throw("unreachable")
}
...
if !block && c.closed==0 &&((c.dataqsiz==0&&c.recvq.first==nil)||(c.dataqsiz>0&&c.qcount==c.dataqsiz)){
return false
}
var t0 int64
if blockprofilerate >0{
t0=cputicks()
}
lock(&c.lock)
if c.closed!=0{
unlock(&c.lock)
panic(plainError("send on closed channel"))
}
if sg:=c.recvq.dequeue();sg!=nil{
send(c,sg,ep,func(){unlock(&c.lock)},3)
return true
}
if c.qcount<c.dataqsiz{
qp:=chanbuf(c,c.sendx)
if raceenabled{
raceacquire(qp)
racerelease(qp)
}
typedmemmove(c.elemtype,qp,ep)
c.sendx++
if c.sendx==c.dataqsiz{
c.sendx=0
}
c.qcount++
unlock(&c.lock)
return true
}
if !block{
unlock(&c.lock)
return false
}
gp:=getg()
mysg:=acquireSudog()
mysg.releasetime=0
if t0!=0{
mysg.releasetime=-1
}
mysg.elem=ep
mysg.waitlink=nil
mysg.g=gp
mysg.selectdone=nil
mysg.c=c
gp.waiting=mysg
gp.param=nil
c.sendq.enqueue(mysg)
goparkunlock(&c.lock,"chan send",traceEvGoBlockSend,3)
if mysg!=gp.waiting{
throwe("G waiting list is corrupted")
}
go.waiting=nil
if gp.param==nil{
if c.closed==0{
throw("chansend: spurious wakeup")
}
panic(plainError("send on closed channel"))

}
gp.param =nil
if mysg.releasetime>0{
blockevent(mysg.releasetime-t0,2)
}
mysg.c=ni;
releaseSudog(mysg)
return true

}

== nill channel
先来看下nil channel的情况 也就是向没有make的channel发送数据 上一章讲到 向nil channel发送数据会报deadlock错误
if c==nil{
gopark(nil,nil,"chan send (nil chan)",traceEvGoStop,2)
throw("unreachable")
}
traceEvGoStop =16
func gopark(unlockf func(*g,unsafe.Pointer) bool,lock unsafe,Pointer,reason string,traceEv byte,traceskip int){}
gopark 会将当前goroutine休眠 然后通过unlockf来唤醒 注意我们上面传入的unlockf是nil 也就是向nil channel 发送数据
的goroutine会一直休眠 同理从nil channel读取数据也是一样的处理  回顾上节
func main(){
var x chan int
go func(){
x<-1
}()
<-x
}
这里是一个main goroutine 从nil读取数据 进入休眠 gofunc() 向nil channel发送数据 也进入休眠
然后Go语言启动的时候还有一个goroutine sysmon会一直检测系统运行的情况 比如 checkdead()
func checkdead(){
...
throw("all goroutines are asleep -- deadlock")
}

==closed channel
向close的channel发送数据 直接panic
lock(&c.lock)
if c.closed!=0{
unlock(&c.lock)
panic(plainError("send on closed channel"))
}

==发送数据处理
发送数据分为三种情况
有goroutine阻塞在channel上 此时hchan.buf 为空 直接将数据发送给该goroutine
当前hchan.buf 还有可用空间 将数据放到buffer里面
当前hchan.buf已满 阻塞当前goroutine
第一种情况如下 从当前channel的等待队列中取出等待的goroutine 然后调用send goready负责唤醒goroutine
lock(&c.lock)
if sg:=c.recvq.dequeue;sg!=nil{
//Found a waiting receivver We pass the value we want to send
//directly to the receiver bypassing the channel buffer if any
send(c,sg,ep,func(){unlock(&c.lock)},3)
return true
}
//send processes a send operation on an empty channel c
//the value ep sent by the sender is copied to the receiver sg
//the receiver is then woken up to go on its merry way
//channel c must be empty and locked send unlocks c with unlockf
//sg must already be dequeued from c
//ep must be  non-nil and point to the head or the caller's stack
func send(c *hchan,sg *sudog,ep unsafe.Pointer,unlockf func(),skip int){
...
if sg.elem!=nil{
sendDirect(c.elemtypem,sg,ep)
sg.elem=nil
}
gp:=sg.g
unlockf()
gp.param=unsafe.Pointer(sg)
if sg.releasetime !=0{
sg.releasetime=cputicks()
}
goready(gp,skip+1)
}
第二种情况比较简单 通过比较qcount和dataqsiz来判断hchan.buf是否还有可用空间 除此之后 还需要调整一下sendx和qcount
lock(&c.lock)
if c.qcount<c.dataqsize{
//space is available in the channel buffer ENqueue the element to send
qp:=chanbuf(c,c.sendx)
if raceenabled{
raceacquire(qp)
racerelease(qp)
}
typedmemove(c.elemtype,qp,ep)
c.sendx++
if c.sendx==c.dataqsiz{
c.sendx=0
}
c.qcount++
unlock(&c.lock)
return true
}
第三种情况如下
//Blocj on the channel some receiver will conmplete our op=eration for us
gp:=getg()
mysg:=acquireSudog()
mysg.releasetime=0
if t0!=0{
mysg.releasetime=-1
}
mysg.elem =ep
mysg.waitlink=nil
mysg,g=gp
mysg.selectdone=nil
mysg.c=c
gp.waiting=mysg
gp.param=nil
c.sendq.enqueue(mysg)//当前goroutine如等待队列
goparkunlock(&c.lock,"chan send",traceEvGoBlockSend,3)//休眠

===recv
读取<-c 和发送的情况非常类似
== nil channel
func chanrecv(c *hchan,ep unsafe.Pointer,block bool)(selected,received bool){
if c==nil{
if !block{
return
}
gopark(nil,nil,"chan receive(nil chan)",traceEvGoStop,2)
throw("unreachable")
}
...
}

==closed channel
从closed channel 接收数据 如果channel中还有数据 接着走下面的流程 如果已经没有数据了 则返回默认值
使用ok-idiom方式读取的时候 第二个参数返回false
lock(&c.lock)
if c.closed!=0&&c.qcount==0{
if raceenabled{
raceacquire(unsafe.Pointer(c))
}
unlock(&c.lock)
if eq!=nil{
typedmemclr(c.elemtype,ep)
}
retuen true,false
}
==接收数据处理
当前有发送goroutine阻塞在channel上 buf已满
lock(&c.lock)
if sg:=c.sendq.dequeue();sg!=nil{
recv(c,sg,ep,func(){unlock(&c.lock)},3)
return  true ,true

}
==buf 中有可用数据
if c.qcount>0{
//receive directly from queue
qp:=chanbuf(c,c.recvx)
if raceenabled{
raceacquire(qp)
racerelease(qp)
}
if ep!=nil{
 typedmemmove(c.elemtype, ep, qp)
}

 typedmemclr(c.elemtype, qp)
    c.recvx++
    if c.recvx == c.dataqsiz {
        c.recvx = 0
    }
    c.qcount--
    unlock(&c.lock)
    return true, true
}
== buf为空 阻塞
// no sender available: block on this channel.
gp := getg()
mysg := acquireSudog()
mysg.releasetime = 0
if t0 != 0 {
    mysg.releasetime = -1
}
// No stack splits between assigning elem and enqueuing mysg
// on gp.waiting where copystack can find it.
mysg.elem = ep
mysg.waitlink = nil
gp.waiting = mysg
mysg.g = gp
mysg.selectdone = nil
mysg.c = c
gp.param = nil
c.recvq.enqueue(mysg)
goparkunlock(&c.lock, "chan receive", traceEvGoBlockRecv, 3)

====close
关闭channel也就是close(ch)对应的代码如下 去掉冗余代码
func closechan(c *hchan) {
    if c == nil {
        panic(plainError("close of nil channel"))
    }

    lock(&c.lock)
    if c.closed != 0 {
        unlock(&c.lock)
        panic(plainError("close of closed channel"))
    }

    c.closed = 1

    var glist *g

    // release all readers
    for {
        sg := c.recvq.dequeue()
        if sg == nil {
            break
        }
        if sg.elem != nil {
            typedmemclr(c.elemtype, sg.elem)
            sg.elem = nil
        }
        if sg.releasetime != 0 {
            sg.releasetime = cputicks()
        }
        gp := sg.g
        gp.param = nil
        if raceenabled {
            raceacquireg(gp, unsafe.Pointer(c))
        }
        gp.schedlink.set(glist)
        glist = gp
    }

    // release all writers (they will panic)
    for {
        sg := c.sendq.dequeue()
        if sg == nil {
            break
        }
        sg.elem = nil
        if sg.releasetime != 0 {
            sg.releasetime = cputicks()
        }
        gp := sg.g
        gp.param = nil
        if raceenabled {
            raceacquireg(gp, unsafe.Pointer(c))
        }
        gp.schedlink.set(glist)
        glist = gp
    }
    unlock(&c.lock)

    // Ready all Gs now that we've dropped the channel lock.
    for glist != nil {
        gp := glist
        glist = glist.schedlink.ptr()
        gp.schedlink = 0
        goready(gp, 3)
    }
}
close channel的工作 除了将c.closed设置为1 还需要
唤醒recvq队列里面的阻塞goroutine
唤醒sendq队列里面的阻塞goroutine

处理方式分别遍历recvq和sendq队列 将所有的goroutine放到glist队列中  最后唤醒glist队列中的goroutine

++select channel
golang中的select语句的实现 在runtime/select.go文件中 我们简单看下select和channel一起用的时候
select{
case c<-x:
...foo
default:
...bar
}
会被编译为
if selectnbsend(c,v){
...foo
}else{
..bar
}
对应 selectnbsend函数如下
func selectnbsend(c *hchan, elem unsafe.Pointer) (selected bool) {
    return chansend(c, elem, false, getcallerpc(unsafe.Pointer(&c)))
}


select {
case v = <-c
    ... foo
default:
    ... bar
}

if selectnbrecv(&v, c) {
    ... foo
} else {
    ... bar
}
func selectnbrecv(elem unsafe.Pointer, c *hchan) (selected bool) {
    selected, _ = chanrecv(c, elem, false)
    return
}

==select {
case v, ok = <-c:
    ... foo
default:
    ... bar
}


if c != nil && selectnbrecv2(&v, &ok, c) {
    ... foo
} else {
    ... bar
}
func selectnbrecv2(elem unsafe.Pointer, received *bool, c *hchan) (selected bool) {
    // TODO(khr): just return 2 values from this function, now that it is in Go.
    selected, *received = chanrecv(c, elem, false)
    return
}

##
golang的channel实现集中在文件runtime/chan.go 参考legendtkl









*/
