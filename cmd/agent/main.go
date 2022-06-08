package main

import (
	"metrics/internal/agent"
)

func main() {
	app := agent.NewHTTPClient()
	app.Run()
}
