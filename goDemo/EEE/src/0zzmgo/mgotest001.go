package main

//import (
//	"fmt"
//	"gopkg.in/mgo.v2"
//	"gopkg.in/mgo.v2/bson"
//	//	"log"
//)
//
//type Mail struct {
//	Id    bson.ObjectId "_id"
//	Name  string
//	Email string
//}
//
//func main() {
//	//连接数据库
//	session, err := mgo.Dial("106.14.181.48")
//	if err != nil {
//		panic(err)
//	}
//	defer session.Close()
//	//获取数据库 获取集合
//	c := session.DB("mydb").C("mail")
//	//存储数据
//	m1 := Mail{bson.NewObjectId(), "user1", "user1@gmail.com"}
//	m2 := Mail{bson.NewObjectId(), "user2", "user2@gmail.com"}
//	m3 := Mail{bson.NewObjectId(), "user2", "user3@gmail.com"}
//	err = c.Insert(&m1, &m2, &m3)
//	if err != nil {
//		panic(err)
//	}
//	//读取数据
//	ms := []Mail{}
//	err = c.Find(bson.M{"name": "user2"}).All(&ms)
//	if err != nil {
//		panic(err)
//	}
//	//显示数据
//	for i, m := range ms {
//		fmt.Printf("%s,%d,%s\n", m.Id.Hex(), i, m.Email)
//	}
//
//}

//
//type Person struct {
//	Name  string
//	Phone string
//}
//
//func main() {
//	session, err := mgo.Dial("106.14.181.48")
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
