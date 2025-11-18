package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "fmt"
)

func Register(username string, DEK []byte, file_exist bool) error {
    // check for dupe usernames
    if file_exist {
        raw, err := parser.Unmarshal("users.json")
        if err != nil {
            return err
        }
        users, ok := raw.([]models.User)
        if !ok {
            return fmt.Errorf("invalid users.json format")
        }
        for _, u := range users {
            if u.Username == username {
                return fmt.Errorf("User already exists.")
            }
        }
    }

    var user models.User
    user.Username = username
    //fmt.Printf("\n%s's password ", username)
    password,err := parser.Prompt(username)
    if err != nil {
        return err
    }

    hashed, salt, master_key, error_msg := crypto.AddUser(password,DEK)
    if error_msg != nil {
        return error_msg
    }
    user.Hashed = hashed
    user.Salt = salt
    user.Master_Key = master_key
    
    user.Servers = []string{"server-prod"}

    users := []models.User{user}
    err = parser.Writer(users, "users.json")
    if err != nil {
        return err
    }
    return nil
}
