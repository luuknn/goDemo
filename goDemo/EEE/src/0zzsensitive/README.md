# Sensitive

敏感词过滤

#


Usage:

```go
package main

import (
	"fmt"
	"0zzsensitive"
)

func main() {
	filter := sensitive.New()
	filter.LoadWordDict("../dict/dict.txt")
	filter.AddWord("MMP")

	fmt.Println(filter.Filter("MMP 金三胖又开始挑衅了")) //  金又开始挑衅了
	// 42 即 "*"
	fmt.Println(filter.Replace("MMP 金三胖又开始挑衅了", 42)) // *** 金**又开始挑衅了
	fmt.Println(filter.FindIn("MMP 金三胖又开始挑衅了"))      // true MMP

	fmt.Println(filter.FindIn("M|MP 金三|胖又开始挑衅了")) // false
	filter.UpdateNoisePattern(`\|`)
	fmt.Println(filter.FindIn("M|MP 金三|胖又开始挑衅了")) // true MMP
}
```
