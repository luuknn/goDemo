package main

import (
	"fmt"
	"sort"
)

type Person struct {
	Name string
	Age  int
}

func (p Person) String() string {
	return fmt.Sprintf("%s:%d", p.Name, p.Age)
}

//ByAge implements sort.Interface for []Person based on the Age field
type ByAge []Person

func (a ByAge) Len() int           { return len(a) }
func (a ByAge) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByAge) Less(i, j int) bool { return a[i].Age < a[j].Age }
func main() {
	people := []Person{
		{"Jack", 12},
		{"Frank", 28},
		{"Lily", 18},
		{"John", 7},
	}
	fmt.Println(people)
	sort.Sort(ByAge(people))
	fmt.Println(people)
}
