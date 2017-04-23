package main

import (
    "fmt"
    "github.com/gin-gonic/gin"
    "os"
)

func StartWebService() {
    port := os.Getenv("PORT")

    if port == "" {
        fmt.Println("$PORT must be set")
    }

    router := gin.New()
    router.Use(gin.Logger())

    router.Static("/", "web")

    router.Run(":" + port)
}
