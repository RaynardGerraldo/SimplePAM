package internal

import (
    "SimplePAM/service"
    "SimplePAM/crypto"
    "SimplePAM/parser"
    "SimplePAM/models"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/scrypt"
    "fmt"
)

func CheckHash(hash []byte, password []byte) bool{
    valid := bcrypt.CompareHashAndPassword(hash, password)
    return valid == nil
}

func ReadCred(username string, password []byte, filename string) ([]byte, bool, error){
    raw, err := parser.Unmarshal(filename)
    if err != nil {
        return nil, false, err
    }
    users, ok := raw.([]models.User)
    if !ok {
        return nil, false, fmt.Errorf("Invalid user format")
    }

    for _, u := range users {
        if u.Username == username {
            if CheckHash(u.Hashed, password) {
                // generate udk
                udk, err := scrypt.Key(password, u.Salt, 32768, 8, 1, 32)
                if err != nil {
                    return nil, false, err
                }
                // get DEK
                DEK,err := crypto.Decrypt(u.Master_Key, udk)
                if err != nil {
                    return nil, false, err
                }

                if username == "admin" {
                    return DEK, true, nil
                } else {
                    return nil, false, service.SSH(DEK, username)
                }
            } else {
                return nil, false, fmt.Errorf("Wrong credentials, try again.")
            }
        } else {
            return nil, false, fmt.Errorf("User doesnt exist.")
        }
    }
    return nil, false, nil
}

func Auth(username string) ([]byte, bool, error){
    password,err := parser.Prompt(username)
    if err != nil {
        return nil, false, err
    }
    if username == "admin" {
        return ReadCred(username, password, "admin.json")
    }
    return ReadCred(username, password, "users.json")
}
