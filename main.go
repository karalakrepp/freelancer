package main

import (
	"log"

	"github.com/karalakrepp/Golang/freelancer-project/controller"
	"github.com/karalakrepp/Golang/freelancer-project/database"
	"github.com/karalakrepp/Golang/freelancer-project/token"
)

func main() {

	db, err := database.NewPostgresStore()
	if err := db.Init(); err != nil {
		log.Fatal(err)

	}
	if err != nil {
		log.Fatal(err)
	}
	maker, err := token.NewJWTMaker("afdassssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssssdassssssssssssss")
	if err != nil {
		log.Println(err)
	}
	svc := controller.NewApiService(":8000", db, maker)
	svc.Routes()
}
