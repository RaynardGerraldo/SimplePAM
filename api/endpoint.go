package main

import (
    "SimplePAM/internal"
    "SimplePAM/parser"
    "github.com/gin-gonic/gin"
    "encoding/base64"
    "gorm.io/gorm"
    "net/http"
    "fmt"
)

type LoginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
}

type RegReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
    Key string `json:"key"`
}

type ServerReq struct {
    Username string `json:"username"`
    Key string `json:"key"`
}

type StatusReq struct {
    Username string `json:"username"`
}

func Login(c *gin.Context, db *gorm.DB) {
    var loginreq LoginReq

    err := c.BindJSON(&loginreq)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    _, err = parser.ReadUsernameDB(db, loginreq.Username)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User doesnt exist."})
        return
    }

    key, valid, err := internal.ReadCred(db, loginreq.Username, []byte(loginreq.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Login failed: %v", err)})
        return
    }

    if valid {
        c.JSON(http.StatusOK, gin.H{"token": base64.StdEncoding.EncodeToString(key)})
    }
}

func Register(c *gin.Context, db *gorm.DB) {
    var regreq RegReq

    err := c.BindJSON(&regreq)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    decodedKey, err := base64.StdEncoding.DecodeString(regreq.Key)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key format"})
        return
    }

    err = internal.Register(db, regreq.Username, []byte(regreq.Password), decodedKey)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Register failed: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": "Account registered"})
}

func InitAdmin(c *gin.Context, db *gorm.DB) {
    var adminreq LoginReq

    err := c.BindJSON(&adminreq)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    key, err := internal.Admin(db, adminreq.Username, []byte(adminreq.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Admin init failed: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"token": base64.StdEncoding.EncodeToString(key)})
}

func InitServer(c *gin.Context, db *gorm.DB) {
    var serverreq RegReq

    err := c.BindJSON(&serverreq)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    
    decodedKey, err := base64.StdEncoding.DecodeString(serverreq.Key)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid key format"})
        return
    }
   
    err = internal.Server(db, serverreq.Username, []byte(serverreq.Password), decodedKey)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Server init failed: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": "Server initialized."})
}

func Status(c *gin.Context, db *gorm.DB) {
    var status StatusReq
    err := c.BindJSON(&status)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }
    _, err = parser.ReadUsernameDB(db, status.Username)
   
    if status.Username == "admin" {
        if err == nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Cant run init, admin already exists"})
            return 
        }
    } else {
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "User doesnt exist."})
            return
        }
    }
    c.JSON(http.StatusOK, gin.H{"error": ""})
    return
}

func AllowedServers(c *gin.Context, db *gorm.DB) {
    var allowed_username StatusReq
    err := c.BindJSON(&allowed_username)

    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    allowed, err := internal.Allowed(db, allowed_username.Username)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Failed to get allowed servers: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"allowed": allowed})
}

func ServersList(c *gin.Context, db *gorm.DB) {
    list, err := internal.ServersList(db)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Failed to get server list: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"list": list})   
}
