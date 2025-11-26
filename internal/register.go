package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "gorm.io/gorm"
    "fmt"
)

func Register(db *gorm.DB, username string, password []byte, DEK []byte) error {
    // check for dupe usernames
    dupe, err := parser.ReadUsernameDB(db, username)
    if err == nil {
        return fmt.Errorf("User already exists")
    }

    // handle error other than record not found
    if dupe == nil {
        if err != gorm.ErrRecordNotFound {
            return fmt.Errorf("Error in reading database: %w", err)
        }
    }


    var user models.User
    user.Username = username

    hashed, salt, master_key, error_msg := crypto.AddUser(password,DEK)
    if error_msg != nil {
        return error_msg
    }

    // write db
    user.Hashed = hashed
    user.Salt = salt
    user.Master_Key = master_key
    server, err := parser.CheckDB(db, "server-prod")
    if err != nil {
        return fmt.Errorf("server-prod not found: %w", err)
    }
    user.Servers = append(user.Servers, server)

    err = parser.WriteDB(db, user)
    if err != nil {
        return fmt.Errorf("Write failed: %w", err)
    }

    return nil
}
