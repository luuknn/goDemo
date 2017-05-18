package main

import (
	"fmt"
)

func main() {
	//定义一个int  slice初始值为 1 2 3 4
	var test []int = []int{1, 2, 3, 4}
	//如果我想让 值变成 1 2 3 4 5 6 7 可以使用append 内置函数
	test = append(test, 5)
	test = append(test, 6)
	test = append(test, 7)
	fmt.Println(test)
	//还有一种更加简便的方法 可以添加 5 6 7 如下所示
	test = append(test, []int{5, 6, 7}...)
	fmt.Println(test)

}
