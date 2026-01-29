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

func Server(db *gorm.DB, serverName string, name string, password []byte, key []byte, ip string, port uint16) error {
    var server models.Server
    parser.InitDB(db, &models.Server{})

    server.Server = serverName
    server.Name = name
    if len(ip) < 0 {
        server.IP = "localhost"
    } else {
        server.IP = ip
    }
    
    if port <= 0 {
        server.Port = 22
    } else {
        server.Port = port
    }

    // encrypt with DEK
    password, err := crypto.Encrypt(password, key)
    if err != nil {
        return err
    }
    server.Password = password
    
    parser.WriteDB(db, server)
    return nil
}

/*func AddtoUser(db *gorm.DB) error {
    var username string
    var serverName string

    fmt.Printf("Username? ")
    fmt.Scanf(&username)

    var user models.User
    result := db.Where("username = ?", username).First(&user)
    if result.Error != nil {
        return fmt.Errorf("Username %s not found: %w", username, err)
    }

    fmt.Printf("Server name to add to %s? ", username)
    fmt.Scanf(&serverName)
    
    server, err := parser.CheckDB(db, serverName)
    if err != nil {
        return fmt.Errorf("%s not found: %w", serverName, err)
    }
    user.Servers = append(user.Servers, server)
    return nil
}*/
