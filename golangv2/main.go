package main

import (
	"log"
)

func main() {
	conn := NewServer("localhost:3000")
	log.Fatal(conn.StartServer())
}
