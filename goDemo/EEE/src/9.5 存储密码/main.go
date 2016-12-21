package main

import (
	"fmt"
)

//过去一段时间以来 许多的网站遭遇用户密码数据泄露事件  这其中包括顶级的互联网企业 Linkin 国内诸如CSDN
//该事件横扫整个国内互联网 随后又爆出多玩游戏800万用户资料被泄露 另有传言 人人网 开心网 天涯社区 世纪佳缘 百合网等社区都有可能成为黑客的下一个目标
//层出不穷的类似事件给用户的网上生活造成巨大的影响 人人自危
//那么我们作为一个WEB开发者 在密码存储方案选择时  容易掉入哪些陷阱 以及如何避免这些陷阱
//普通方案
//目前使用最多的密码存储方案是将明文密码做单向哈希后存储 单向哈希算法 有一个特征 无法通过哈希后的摘要 digest恢复原始数据
//这也是单向 两字的来源  常用的单向哈希算法包括 SHA-256 SHA-1 MD5等
//go语言对三种加密算法的实现如下所示
package main

import (
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
)

func main() {
	h := sha256.New()
	io.WriteString(h, "gqy")
	fmt.Printf("% x \n", h.Sum(nil))
	//86 d1 f0 82 1a 7e 4c be 9d 95 2f 22 ae 6d 93 a5 38 26 b7 5e bb 05 ce 18 c5 a9 59 90 e3 ed e6 29

	hh := sha1.New()
	io.WriteString(hh, "gqy")
	fmt.Printf("% x \n", hh.Sum(nil))
	//c5 e7 b1 f8 be b8 79 be 91 17 9c 9d e9 45 62 36 8c c7 f6 ab

	hhh := md5.New()
	io.WriteString(hhh, "gqy")
	fmt.Printf("%x", hhh.Sum(nil))
	//b15d909c980dc7feadecdcc226c4f3e2

}
//单向的哈希算法有两个特性
//同一个密码进行单向哈希 得到的总是唯一确定的摘要
//计算速度快 随着技术的进步 一秒钟就能完成 数十亿次的单向哈希计算

//综合上面两个特点 考虑到多数人使用的密码为常见的组合 攻击者可以将所有的密码的常见组合单向哈希
//得到一个摘要组合 然后与数据库中的摘要进行对比即可获得对应的密码 这个摘要组合也被称为 rainbow table
//因此通过单向加密之后存储的数据 和明文存储没有多大区别 因此 一旦网站的数据库泄露 所有用户的密码本身就大白于天下

//进阶方案
//通过上面的介绍 我们知道黑客可以用rainbow table来破解哈希后的密码 很大程度上是因为 加密时使用的哈希算法是公开的
//如果黑客不知道加密的哈希算法是什么 那他就无从下手了
/*
一个直接的解决办法是，自己设计一个哈希算法。然而，一个好的哈希算法是很难设计的——既要避免碰撞，又不能有明显的规律，做到这两点要比想象中的要困难很多。因此实际应用中更多的是利用已有的哈希算法进行多次哈希。

但是单纯的多次哈希，依然阻挡不住黑客。两次 MD5、三次 MD5之类的方法，我们能想到，黑客自然也能想到。特别是对于一些开源代码，这样哈希更是相当于直接把算法告诉了黑客。

没有攻不破的盾，但也没有折不断的矛。现在安全性比较好的网站，都会用一种叫做“加盐”的方式来存储密码，也就是常说的 “salt”。他们通常的做法是，先将用户输入的密码进行一次MD5（或其它哈希算法）加密；将得到的 MD5 值前后加上一些只有管理员自己知道的随机串，再进行一次MD5加密。这个随机串中可以包括某些固定的串，也可以包括用户名（用来保证每个用户加密使用的密钥都不一样）。
*/
package main

import (
	"crypto/md5"

	"fmt"
	"io"
)

func main() {
	//import "crypto/md5"
	//假设用户名abc，密码123456
	h := md5.New()
	io.WriteString(h, "gqy")

	//pwmd5等于e10adc3949ba59abbe56e057f20f883e
	pwmd5 := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(pwmd5)

	//指定两个 salt： salt1 = @#$%   salt2 = ^&*()
	salt1 := "@#$%"
	salt2 := "^&*()"

	//salt1+用户名+salt2+MD5拼接
	io.WriteString(h, salt1)
	io.WriteString(h, "abc")
	io.WriteString(h, salt2)
	io.WriteString(h, pwmd5)

	last := fmt.Sprintf("%x", h.Sum(nil))
	fmt.Println(last)

}
//在两个salt没有泄露的情况下 黑客如果拿到的是最后这个加密串 就几乎不可能推算出原始的密码是什么了

//专家方案
/*上面的进阶方案在几年前也许是足够安全的方案，因为攻击者没有足够的资源建立这么多的rainbow table。 但是，时至今日，因为并行计算能力的提升，这种攻击已经完全可行。

怎么解决这个问题呢？只要时间与资源允许，没有破译不了的密码，所以方案是:故意增加密码计算所需耗费的资源和时间，使得任何人都不可获得足够的资源建立所需的rainbow table。

这类方案有一个特点，算法中都有个因子，用于指明计算密码摘要所需要的资源和时间，也就是计算强度。计算强度越大，攻击者建立rainbow table越困难，以至于不可继续。

这里推荐scrypt方案，scrypt是由著名的FreeBSD黑客Colin Percival为他的备份服务Tarsnap开发的。

目前Go语言里面支持的库http://code.google.com/p/go/source/browse?repo=crypto#hg%2Fscrypt

dk := scrypt.Key([]byte("some password"), []byte(salt), 16384, 8, 1, 32)
通过上面的的方法可以获取唯一的相应的密码值，这是目前为止最难破解的。

总结

看到这里，如果你产生了危机感，那么就行动起来：

1）如果你是普通用户，那么我们建议使用LastPass进行密码存储和生成，对不同的网站使用不同的密码；
2）如果你是开发人员， 那么我们强烈建议你采用专家方案进行密码存储。

*/




















