package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "gorm.io/gorm"
    "fmt"
)

func toAdmin(db *gorm.DB) error {
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
    
    parser.InitDB(db, &models.User{})

    admin.Hashed = hashed
    admin.Salt = salt
    admin.Master_Key = master_key
    parser.WriteDB(db, admin)

    return toServer(db, key)
}

func toServer(db *gorm.DB, key []byte) error {
    var server models.Server
    var name string
    fmt.Println("\nTry it out with your localhost")
    fmt.Printf("Server username? ")
    fmt.Scan(&name)

    password,err := parser.Prompt("server " + name)
    if err != nil {
        return err
    }

    parser.InitDB(db, &models.Server{})

    server.Server = "server-prod"
    server.Name = name
    server.IP = "localhost"
    // encrypt with DEK
    password, err = crypto.Encrypt(password, key)
    if err != nil {
        return err
    }
    server.Password = password
    
    parser.WriteDB(db, server)

    fmt.Println("admin and server initialized.")
    return nil
}

func Init(db *gorm.DB) error {
    return toAdmin(db)
}
