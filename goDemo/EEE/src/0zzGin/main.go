package main

import (
	db "src/0zzGin/database"
)

func main() {
	defer db.SqlDB.Close()
	router := initRouter()
	router.Run(":8000")
}
