package internal

import (
    "SimplePAM/models"
    "SimplePAM/internal"
    "golang.org/x/crypto/bcrypt"
    "log"
    "encoding/json"
    "os"
)

func Register(username string, password []byte){
    var user models.User

    // encrypt here with master key
    bytes, err := bcrypt.GenerateFromPassword(password, 14)
    if err != nil{
        log.Fatal("Couldnt generate password")
    }
 
    user.Username = username
    user.Password = bytes
    user.Servers = []string{"server-prod", "server-test"}

    users := []models.User{user}
    internal.Writer(users, "users.json") 
}
