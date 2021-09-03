package main

import (
	"etender/api/property"

	"github.com/gin-gonic/gin"
	cors "github.com/rs/cors/wrapper/gin"
)

const PORT = ":3000"

func main() {

	router := gin.Default()
	router.Use(cors.Default())

	property.Routes(router)

	router.Run(PORT)
}
