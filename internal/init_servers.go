package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "fmt"
)

func toAdmin() error {
    var admin models.User
    fmt.Println("Your admin username is 'admin' by default")
    admin.Username = "admin"
    password,err := parser.Prompt(admin.Username)
    if err != nil {
        return err
    }

    hashed, salt, master_key, key, err := crypto.Init(password)
 
    if err != nil {
        return err
    }

    admin.Hashed = hashed
    admin.Salt = salt
    admin.Master_Key = master_key
    admin.Servers = []string{}
    
    admin_ins := []models.User{admin}
    err = parser.Writer(admin_ins, "admin.json")
    if err != nil {
        return err
    }
    return toServer(key)
}

func toServer(key []byte) error {
    var server models.Server
    var name string

    fmt.Println("\nTry it out with your localhost")
    fmt.Printf("Server username? ")
    fmt.Scan(&name)

    password,err := parser.Prompt("server " + name)
    if err != nil {
        return err
    }

    server.Server = "server-prod"
    server.Name = name
    server.IP = "localhost"
    // encrypt with DEK
    password, err = crypto.Encrypt(password, key)
    if err != nil {
        return err
    }
    server.Password = password

    servers := []models.Server{server}
    err = parser.Writer(servers, "servers.json")
    if err != nil {
        return err
    }
    fmt.Println("admin and server initialized.")
    return nil
}

func Init() error {
    return toAdmin()
}
