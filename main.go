package main

import (
	"bbs-go/internal/server"
	_ "bbs-go/internal/services/eventhandler"
)

func main() {
	server.Init()
	server.NewServer()
}
