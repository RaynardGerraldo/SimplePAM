package main

import (
    "SimplePAM/internal"
    "SimplePAM/parser"
    "SimplePAM/models"
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
    ServerName string `json:"servername"`
}

type ServerReq struct {
    ServerName string `json:"servername"`
    Username string `json:"username"`
    Password string `json:"password"`
    Key string `json:"key"`
    IP string `json:"ip"`
    Port uint16 `json:"port"`
}

type StatusReq struct {
    Username string `json:"username"`
}

type AddtoUserReq struct {
    Username string `json:"username"`
    ServerName string `json:"servername"`
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

    jwt, err := GenerateToken(loginreq.Username, "user")
    if err != nil {
        c.JSON(500, gin.H{"error": "Could not generate token"})
        return
    }

    if valid {
        c.JSON(http.StatusOK, gin.H{"token": base64.StdEncoding.EncodeToString(key), "jwt": jwt})
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

    err = internal.Register(db, regreq.Username, []byte(regreq.Password), decodedKey, regreq.ServerName)

    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Register failed: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": "Account registered"})
}

func InitAdmin(c *gin.Context, db *gorm.DB) {
    var adminreq LoginReq
    var count int64
    db.Model(&models.User{}).Count(&count)
    if count > 0 {
        c.JSON(http.StatusBadRequest, gin.H{"error": "System already initialized. Please login."})
        return
    }

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

    jwt, err := GenerateToken(adminreq.Username, "admin")
    if err != nil {
        c.JSON(500, gin.H{"error": "Could not generate token"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"token": base64.StdEncoding.EncodeToString(key), "jwt": jwt})
}

func InitServer(c *gin.Context, db *gorm.DB) {
    var serverreq ServerReq

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
   
    err = internal.Server(db, serverreq.ServerName, serverreq.Username, []byte(serverreq.Password), decodedKey, serverreq.IP, serverreq.Port)

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
   
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "User doesnt exist."})
        return
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

func UsersList(c *gin.Context, db *gorm.DB) {
    list, err := internal.UsersList(db)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Failed to get user list: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"list": list})   
}

func AddtoUser(c *gin.Context, db *gorm.DB) {
    var toadd AddtoUserReq
    err := c.BindJSON(&toadd)

    var user models.User
    result := db.Where("username = ?", toadd.Username).First(&user)
    if result.Error != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Username %s not found: %w", toadd.Username, result.Error)})
        return 
    }

    server, err := parser.CheckDB(db, toadd.ServerName)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Server %s not found: %w", toadd.ServerName, err)})
        return
    }

    err = db.Model(&user).Association("Servers").Append(server)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Write failed: %w", err)})
        return
    }

    c.JSON(http.StatusOK, gin.H{"success": "Server added to user."})
}
