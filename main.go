package main

import (
	"log"
	"ws_comparator/interfaces/middleware/server"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("no env gotten")
	}
}

func main() {

	server.InitServer().RunServer()

}
