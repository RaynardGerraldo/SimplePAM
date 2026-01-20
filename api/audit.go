package main

import (
    "github.com/gin-gonic/gin"
    "fmt"
    "os"
    "time"
)

func Audit() gin.HandlerFunc {
    return func(c *gin.Context) {
        start := time.Now()
        c.Next()
        end := time.Now()

        latency := end.Sub(start)

        status := c.Writer.Status()
        clientIP := c.ClientIP()
        method := c.Request.Method
        path := c.Request.URL.Path

        f, err := os.OpenFile("audit.log", os.O_WRONLY | os.O_CREATE | os.O_APPEND, 0644)
        if err != nil {
            fmt.Println("Error opening log file: ", err)
            return
        }

        defer f.Close()

        log := fmt.Sprintf("[%d] %s %s | %s | %v\n", status, method, path, clientIP, latency)

        _, err = f.WriteString(log)

        if err != nil {
            fmt.Println("Error writing to logfile: ", err)
        }
    }
}
