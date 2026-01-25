package main

import (
    "fmt"
    "os"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
)

var upgrader = websocket.Upgrader {
    ReadBufferSize: 1024,
    WriteBufferSize: 1024,
    CheckOrigin: func(r *http.Request) bool { return true },
}

func echoHandler(c *gin.Context) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        fmt.Println("Upgrade failed: ", err)
        os.Exit(1)
    }

    defer conn.Close()

    for {
        msgType, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }

        err = conn.WriteMessage(msgType, msg)
        if err != nil {
            break
        }
    }
}
