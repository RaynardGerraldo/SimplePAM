package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "gorm.io/gorm"
)

func Admin(db *gorm.DB, username string, password []byte) ([]byte, error) {
    var admin models.User
    admin.Username = username

    hashed, salt, master_key, key, err := crypto.Init(password)
 
    if err != nil {
        return nil, err
    }
    
    parser.InitDB(db, &models.User{})

    admin.Hashed = hashed
    admin.Salt = salt
    admin.Master_Key = master_key
    parser.WriteDB(db, admin)

    return key, nil
}

func Server(db *gorm.DB, name string, password []byte, key []byte) error {
    var server models.Server
    parser.InitDB(db, &models.Server{})

    server.Server = "server-prod"
    server.Name = name
    server.IP = "localhost"
    // encrypt with DEK
    password, err := crypto.Encrypt(password, key)
    if err != nil {
        return err
    }
    server.Password = password
    
    parser.WriteDB(db, server)
    return nil
}
