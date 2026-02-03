package main

import (
    "SimplePAM/parser"
    "github.com/gin-contrib/cors"
    "github.com/gin-gonic/gin"
    "fmt"
    "os"
    "time"
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

    r.Use(cors.New(cors.Config{
        AllowOrigins:    []string{"http://localhost:8080"},
        AllowMethods:    []string{"GET", "POST", "PUT", "PATCH", "DELETE", "HEAD"},
        AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
        MaxAge:          12 * time.Hour,
    }))

    r.POST("/login", func(c *gin.Context) {
        Login(c, db)
    })
    r.POST("/status", func(c *gin.Context) {
        Status(c, db)
    })
    r.POST("/initadmin", func(c *gin.Context) {
        InitAdmin(c, db)
    })
    
    r.GET("/ws/echo", func(c *gin.Context) {
        echoHandler(c)
    })

    r.GET("/ws/ssh", func(c *gin.Context) {
        SSHHandler(c, db)
    })

    protected := r.Group("/")
    protected.Use(AuthMiddleware())
    {
        protected.POST("/addtouser", func(c *gin.Context) {
            AddtoUser(c, db)
        })

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

        protected.GET("/userslist", func(c *gin.Context) {
            UsersList(c, db)
        })
    }
    fmt.Println("PAM Server is running on localhost:8080...")
    r.Static("/web", "./frontend")
    r.Run(":8080") 
}
