package main

import (
	"github.com/sangjinsu/proglog/internal/server"
	"log"
)

func main() {
	svr := server.NewHttpServer(":8080")
	log.Fatal(svr.ListenAndServe())
}
