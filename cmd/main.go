package main

import (
	"clean-architecture-gochat/internal/routers"
	"log"
)

// @title Clean Architecture Chat API
// @version 1.0
// @description This is a chat application API using Clean Architecture
// @host localhost:8080
// @BasePath /
func main() {
	// 加了這一行後，gin就不會輸出debug的訊息了
	//gin.SetMode(gin.ReleaseMode)
	r := routers.SetupRouter()

	log.Println("Starting server on :8080")
	r.Run(":8080")
}
