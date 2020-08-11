package main

import (
	routes "PortalMGTNIIP/routes"
	"log"
	_ "os"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	router := gin.Default()

	router.LoadHTMLGlob("template/*")
	router.Use(cors.Default())
	routes.Routes(router)
	log.Fatal(router.Run(":4747"))
}
