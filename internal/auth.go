package internal

import (
    "SimplePAM/models"
    "SimplePAM/service"
    "SimplePAM/crypto"
    "SimplePAM/parser"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/scrypt"
    "os"
    "io/ioutil"
    "log"
    "encoding/json"
)

func CheckHash(hash []byte, password []byte) bool{
    valid := bcrypt.CompareHashAndPassword(hash, password)
    return valid == nil
}

func ReadCred(username string, password []byte, filename string) ([]byte, bool){
    jsonfile, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer jsonfile.Close()

    bytes, err := ioutil.ReadAll(jsonfile)
    if err != nil {
        log.Fatal("Couldnt read", err)
    }
    
    var users []models.User
    err = json.Unmarshal(bytes, &users)
    if err != nil {
        log.Fatal("Error unmarshalling json", err)
    }

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
