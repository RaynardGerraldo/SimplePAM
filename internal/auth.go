package internal

import (
    "SimplePAM/service"
    "SimplePAM/crypto"
    "SimplePAM/parser"
    "SimplePAM/models"
    "golang.org/x/crypto/bcrypt"
    "golang.org/x/crypto/scrypt"
    "gorm.io/gorm"
    "fmt"
)

func CheckHash(hash []byte, password []byte) bool{
    valid := bcrypt.CompareHashAndPassword(hash, password)
    return valid == nil
}

func ReadCred(db *gorm.DB, username string, password []byte) ([]byte, bool, error){
    var users []models.User
    db.Find(&users)
 
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
                    return nil, false, service.SSH(db, DEK, username)
                }
            } else {
                return nil, false, fmt.Errorf("Wrong credentials, try again.")
            }
        }     
    }
    return nil, false, fmt.Errorf("User doesnt exist.")
}

func Auth(db *gorm.DB, username string) ([]byte, bool, error){
    password,err := parser.Prompt(username)
    if err != nil {
        return nil, false, err
    }
    return ReadCred(db, username, password)
}
