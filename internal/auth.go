package internal

import (
    "SimplePAM/crypto"
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
    var user models.User
    err := db.Where("username = ?", username).First(&user).Error
    if err != nil {
        if err == gorm.ErrRecordNotFound {
            return nil, false, fmt.Errorf("User not found")
        }
        return nil, false, fmt.Errorf("db error: %w", err)
    }

    if CheckHash(user.Hashed, password) {
        // generate udk
        udk, err := scrypt.Key(password, user.Salt, 32768, 8, 1, 32)
        if err != nil {
            return nil, false, err
        }
        // get DEK
        DEK,err := crypto.Decrypt(user.Master_Key, udk)
        if err != nil {
            return nil, false, err
        }
        
        return DEK, true, nil
    }
    return nil, false, fmt.Errorf("Wrong credentials, try again.")
}
