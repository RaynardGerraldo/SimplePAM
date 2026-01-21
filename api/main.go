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

    r := gin.New()

    r.Use(gin.Recovery())

    r.Use(Audit())

    r.POST("/login", func(c *gin.Context) {
        Login(c, db)
    })
    r.POST("/status", func(c *gin.Context) {
        Status(c, db)
    })
    r.POST("/initadmin", func(c *gin.Context) {
        InitAdmin(c, db)
    })

    protected := r.Group("/")
    protected.Use(AuthMiddleware())

    {
        protected.POST("/initserver", func(c *gin.Context) {
            InitServer(c, db)
        })

        protected.POST("/register", func(c *gin.Context) {
            Register(c, db)
        })
      
        protected.POST("/allowedservers", func(c *gin.Context) {
            AllowedServers(c, db)
        })

        protected.GET("/serverslist", func(c *gin.Context) {
            ServersList(c, db)
        })
    }
    fmt.Println("PAM Server is running on localhost:8080...")
    r.Run(":8080") 
}
