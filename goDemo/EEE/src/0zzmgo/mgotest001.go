package main

//import (
//	"fmt"
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	"log"
//)
//
//type Person struct {
//	Name  string
//	Phone string
//}
//
//func main() {
//	session, err := mgo.Dial("106.14.181.47")
//	if err != nil {
//		panic(err)
//	}
//	defer session.Close()
//	session.SetMode(mgo.Monotonic, true)
//	c := session.DB("mydb").C("people")
//	err = c.Insert(&Person{"Ale", "+55 55 8118 9639"}, &Person{"Cla", "+55 56 7868 1341"})
//	if err != nil {
//		log.Fatal(err)
//	}
//	result := Person{}
//	err = c.Find(bson.M{"name": "Ale"}).One(&result)
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println("Phone:", result.Phone)
//
//}
