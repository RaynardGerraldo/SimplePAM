package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "log"
    "fmt"
    "golang.org/x/crypto/ssh/terminal"
    "syscall"
)

func Register(username string, DEK []byte){
    var user models.User

    user.Username = username
    fmt.Printf("\nEnter %s's password: ", username)
    password, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatal(err)
    }
    hashed, salt, master_key := crypto.AddUser(password,DEK)
    user.Hashed = hashed
    user.Salt = salt
    user.Master_Key = master_key
    
    user.Servers = []string{"server-prod"}

    users := []models.User{user}
    parser.Writer(users, "users.json") 
}
