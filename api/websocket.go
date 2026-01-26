package main

import (
    "fmt"
    "os"
    "gorm.io/gorm"
    "net/http"
    "github.com/gorilla/websocket"
    "github.com/gin-gonic/gin"
    "encoding/base64"
    "SimplePAM/service"
    "SimplePAM/models"
    "SimplePAM/crypto"
)


type HandshakeReq struct {
    Type string `json:"type"`
    Token string `json:"token"`
    DEK   string `json:"dek"`
    Target string `json:"target"`
}

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

func SSHHandler(c *gin.Context, db *gorm.DB) {
    conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
    if err != nil {
        fmt.Println("Upgrade failed: ", err)
        os.Exit(1)
    }

    defer conn.Close()

    var handshake HandshakeReq

    err = conn.ReadJSON(&handshake)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, []byte("Error reading handshake"))
    }

    if handshake.Type != "handshake" {
        conn.WriteMessage(websocket.TextMessage, []byte("Invalid message type"))
        return
    }
    
    _, err = ValidateToken(handshake.Token)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, []byte("Unauthorized: Invalid token"))
        return
    }

    decodedKey, err := base64.StdEncoding.DecodeString(handshake.DEK)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, []byte("Invalid key format"))
        return
    }

    var server models.Server
    searchServer := db.Where("server = ?", handshake.Target).First(&server)
    if searchServer.Error != nil {
        conn.WriteMessage(websocket.TextMessage, []byte("Server not found: "+handshake.Target))
        return
    }

    password, err := crypto.Decrypt(server.Password, decodedKey)
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, []byte("Decryption Failed"))
        return
    }

    conn.WriteMessage(websocket.TextMessage, []byte("\r\nAuthentication Successful. Connecting...\r\n"))

    rw := &WebSocketRW{Conn: conn}

    err = service.InternalSSH(rw, rw, server.Name, string(password), server.IP)
    
    if err != nil {
        conn.WriteMessage(websocket.TextMessage, []byte("SSH Error: "+err.Error()))
    }
}
