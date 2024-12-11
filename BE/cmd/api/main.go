package main

import (
	"BE-chat-app/cmd/api/db"
	"BE-chat-app/cmd/api/routes"

	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
)



func main(){
	db.Init()

	server := gin.Default()

	routes.RegisterRoutes(server)

	server.Run(":8080")
}