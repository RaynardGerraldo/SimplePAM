package internal

import (
    "SimplePAM/models"
    "golang.org/x/crypto/bcrypt"
    "log"
    "encoding/json"
    "os"
)

func Register(username string, password []byte){
    var user models.User
    bytes, err := bcrypt.GenerateFromPassword(password, 14)
    if err != nil{
        log.Fatal("Couldnt generate password")
    }
 
    user.Username = username
    user.Password = bytes
    user.Servers = []string{"server-prod", "server-test"}

    users := []models.User{user}

    toJson, err := json.MarshalIndent(users, "", " ")
    if err != nil {
        log.Fatal("Couldnt parse to JSON")
    }

    err = os.WriteFile("users.json", toJson, 0644)
    if err != nil{
        log.Fatal("Couldnt write to file")
    }
}
