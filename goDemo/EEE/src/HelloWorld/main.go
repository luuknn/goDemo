package main

import (
	"fmt"
	"time"
)

func ready(w string, sec int64) {
	secs := time.Duration(sec) * time.Second
	time.Sleep(secs)
	fmt.Println(w, "is ready")

}
func main() {
	go ready("Tee", 2)
	go ready("Coffee", 1)
	fmt.Println("T'm waiting")
	sec := 1
	secs := time.Duration(sec) * time.Second
	time.Sleep(5 * secs)
}
