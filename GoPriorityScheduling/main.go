package main

import (
	app "goPriorityScheduler/App"
	"goPriorityScheduler/server"
)

func main() {
	app := app.NewApp()
	HttpServer := server.NewHTTPServer(8000, app)
	HttpServer.StartServer()
}
