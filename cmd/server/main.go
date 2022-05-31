package main

import "metrics/internal/server/server"

func main() {
	var httpServer server.Server
	httpServer.Run()
}
