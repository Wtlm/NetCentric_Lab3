package main

import (
	"user-management/database"
	"user-management/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	database.InitDB()
	r := gin.Default()
	routes.RegisterRoutes(r)
	r.Run(":8080")
}
