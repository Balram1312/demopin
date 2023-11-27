package main

import (
	"demo_pinpoint/feature"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.Use(cors.Default())
	r.POST("/send-email", feature.SendEmail)

	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
