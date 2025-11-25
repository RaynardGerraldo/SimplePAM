package main

import (
    "SimplePAM/internal"
    "github.com/gin-gonic/gin"
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

func Login(c *gin.Context, db *gorm.DB) {
    var loginreq LoginReq

    err := c.BindJSON(&loginreq)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    key, valid, err := internal.ReadCred(db, loginreq.Username, []byte(loginreq.Password))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Login failed: %v", err)})
        return
    }

    if valid {
        c.JSON(http.StatusOK, gin.H{"token": string(key)})
    }
}

func Register(c *gin.Context, db *gorm.DB) {
    var regreq RegReq

    err := c.BindJSON(&regreq)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request format"})
        return
    }

    err = internal.Register(db, regreq.Username, []byte(regreq.Password), []byte(regreq.Key))
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": fmt.Sprintf("Register failed: %v", err)})
        return
    }
    c.JSON(http.StatusOK, gin.H{"success": "Account registered"})
}
