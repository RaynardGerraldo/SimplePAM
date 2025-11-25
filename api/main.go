package main

import (
    "SimplePAM/parser"
    "github.com/gin-gonic/gin"
    "fmt"
    "os"
)

func main() {
    db, err := parser.OpenCon()
    if err != nil {
        fmt.Println("Failed to connect to database:", err)
        os.Exit(1)
    }

    r := gin.Default()

    r.POST("/login", func(c *gin.Context) {
        Login(c, db)
    })

    fmt.Println("PAM Server is running on localhost:8080...")
    r.Run(":8080") 
}
