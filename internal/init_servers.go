package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "fmt"
    "log"
)

func toAdmin() {
    var admin models.User
    fmt.Println("Your admin username is 'admin' by default")
    admin.Username = "admin"
    password := parser.Prompt()
    hashed, salt, master_key, key := crypto.Init(password)
    admin.Hashed = hashed
    admin.Salt = salt
    admin.Master_Key = master_key
    admin.Servers = []string{}
    
    admin_ins := []models.User{admin}
    parser.Writer(admin_ins, "admin.json")
    toServer(key)
}

func toServer(key []byte) {
    var server models.Server
    var name string

    fmt.Println("\nTry it out with your localhost")
    fmt.Printf("Server username? ")
    fmt.Scan(&name) 
    fmt.Printf("Server password? ")
    password := parser.Prompt()
    server.Server = "server-prod"
    server.Name = name
    server.IP = "localhost"
    // encrypt with DEK
    password, err := crypto.Encrypt(password, key)
    if err != nil {
        log.Fatal(err)
    }
    server.Password = password

    servers := []models.Server{server}
    parser.Writer(servers, "servers.json")
}

func Init(){
    toAdmin()
}
