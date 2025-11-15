package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "fmt"
)

func Register(username string, DEK []byte){
    var user models.User

    user.Username = username
    fmt.Printf("\n%s's password ", username)
    password := parser.Prompt()
    hashed, salt, master_key := crypto.AddUser(password,DEK)
    user.Hashed = hashed
    user.Salt = salt
    user.Master_Key = master_key
    
    user.Servers = []string{"server-prod"}

    users := []models.User{user}
    parser.Writer(users, "users.json") 
}
