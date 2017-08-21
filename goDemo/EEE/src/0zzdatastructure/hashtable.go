package main

import ()

/*
hash table 中文可以称作哈希表
简单来说就是提供(key,value)的存取 当然只存key也是可以的
一般来说 哈希表的插入和查找的平均时间复杂度都是O(1)因此在日常工作中使用也较为广泛
1General Hash Table
实现
为了保证插入和查找的平均复杂度为O(1) hash table 底层一般都是使用数组来实现
对于给定的key 一般先进行hash操作 然后相对哈希表的长度取模 将key映射到指定的地方
index=hash(key) % hash_table_size //index就是key存储位置的索引
这里的核心在于如何选取合适的hash函数 如果提前知道key的一些相关信息 往往可以选取一个不错的hash函数
常用的hash函数有SHA-1 SHA-256 SHA-512,
冲突处理
冲突 也叫作碰撞 意思是两个或者多个key映射到了哈希表的同一个位置 冲突处理一般有两种方法
开放定址 open addressing 和开链 separate chaining
开放定址 的意思是当发生冲突时  我们从当前位置向后按某种策略遍历哈希表 当发现可用的空间的时候 则插入元素
开放地址有一次探测、二次探测和双重哈希。一次探测是指我们的遍历策略是一个线性函数，比如依次遍历冲突位置之后的第 1，2，3…N 位置。如果直接遍历 1，4（=2^2），9 (=3^2)，这就是二次探测的一个例子。双重哈希就是遍历策略间隔由另一个哈希函数来确定。
key “John Smith” 和 “Sandra Dee” 在 index = 152 位置出现冲突，使用开放地址的方法将 “Sandra Dee” 存放在 index = 153 的位置。之后 key “Ted Baker” 的映射位置为 index = 153，又出现冲突，则将其存放在 index = 154 的位置。由这个例子我们可以看出这种处理方法的一个缺点：解决旧问题的同时会引入新的问题。

开链的思想是哈希表中的每一个元素都是一个类似链表或者其他数据结构的head 当出现冲突时 我们就在链表后面
添加元素 这也就意味着 如果某一个位置冲突过多的话 插入的时间复杂度将退化为O(1)补充一点
如果哈希表的每个元素都是一个链表头的 那么又可以分为头存储元素和不存储元素两种
简单比较一下这两种处理方法的优劣：开放定址在解决当前冲突的情况下同时可能会导致新的冲突，而开链不会有这种问题。同时开链相比于开放定址局部性较差，在程序运行过程中可能引起操作系统的缺页中断，从而导致系统颠簸。















*/