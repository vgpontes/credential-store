package main

import (
	"log"
)

func main() {
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	err = db.Init()
	if err != nil {
		log.Fatal(err)
		panic(err)
	}

	server := NewAPIServer(db)
	server.Run()
}
