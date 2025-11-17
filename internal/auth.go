package internal

import (
    "SimplePAM/service"
    "SimplePAM/crypto"
    "SimplePAM/parser"
    "SimplePAM/models"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/scrypt"
    "log"
    "fmt"
)

func CheckHash(hash []byte, password []byte) bool{
    valid := bcrypt.CompareHashAndPassword(hash, password)
    return valid == nil
}

func ReadCred(username string, password []byte, filename string) ([]byte, bool){
    users := parser.Unmarshal(filename).([]models.User)
    for _, u := range users {
        if u.Username == username {
            if CheckHash(u.Hashed, password) {
                // generate udk
                udk, err := scrypt.Key(password, u.Salt, 32768, 8, 1, 32)
                if err != nil {
                    log.Fatal("Failed to generate udk: %v", err)
                }
                // get DEK
                DEK,err := crypto.Decrypt(u.Master_Key, udk)
                if err != nil {
                    log.Fatal("Failed to decrypt to DEK: %v", err)
                }

                if username == "admin" {
                    return DEK, true
                    //service.addUser(DEK)
                } else {
                    service.SSH(DEK, username)
                }
            } else {
                log.Fatal("\nWrong credentials, try again.")
            }
        } else {
            fmt.Println("\nUser doesnt exist.")
        }
    }
    return nil, false
}

func Auth(username string) ([]byte, bool){
    password := parser.Prompt()
    if username == "admin" {
        return ReadCred(username, password, "admin.json")
    } else {
        ReadCred(username, password, "users.json")
    }
    return nil, false
}
