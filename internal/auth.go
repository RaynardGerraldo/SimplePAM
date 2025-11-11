package internal

import (
    "SimplePAM/models"
    "golang.org/x/crypto/bcrypt"
    "fmt"
    "os"
    "io/ioutil"
    "log"
    "encoding/json"
)

func CheckHash(hash []byte, password []byte) bool{
    err := bcrypt.CompareHashAndPassword(hash, password)
    return err == nil
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
        fmt.Println(CheckHash(user.Password, password))
    } else{
        fmt.Println("\nNo username match")
    }
}
