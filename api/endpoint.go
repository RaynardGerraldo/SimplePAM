package main

import (
    "SimplePAM/internal"
    "github.com/gin-gonic/gin"
    "gorm.io/gorm"
    "net/http"
)

type LoginReq struct {
    Username string `json:"username"`
    Password string `json:"password"`
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
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Login failed"})
        return
    }

    if valid {
        c.JSON(http.StatusOK, gin.H{"token": string(key)})
    }
}
