package internal

import (
    "SimplePAM/models"
    "SimplePAM/service"
    "golang.org/x/crypto/bcrypt"
    "os"
    "io/ioutil"
    "log"
    "encoding/json"
)

func CheckHash(hash []byte, password []byte) bool{
    valid := bcrypt.CompareHashAndPassword(hash, password)
    return valid == nil
}

func Auth(username string, password []byte){
    // read from users.json, match username and password from args.
    jsonfile, err := os.Open("users.json")
    if err != nil {
        log.Fatal("Couldnt open users.json", err)
    }
    defer jsonfile.Close()

    bytes, err := ioutil.ReadAll(jsonfile)
    if err != nil {
        log.Fatal("Couldnt read users.json", err)
    }
    
    var user models.User
    err = json.Unmarshal(bytes, &user)
    if err != nil {
        log.Fatal("Error unmarshalling json", err)
    }
    if user.Username == username {
        service.SSH(CheckHash(user.Password, password))
    } else{
        service.SSH(false)
    }
}
