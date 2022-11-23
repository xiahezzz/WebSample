package main

import (
	"web/server"
)

func main() {
	MyServer := server.NewServer()
	//gin.SetMode(gin.ReleaseMode)

	MyServer.Start() // listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
}
