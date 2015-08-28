package main

import "github.com/chamot1111/dokpi"

func runAppServer() {
	appServer := dokpi.NewAppServer()
	appServer.StartServer()
}

func main() {
	runAppServer()
}
