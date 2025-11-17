package crypto

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"fmt"
    "log"
	"io"
	"golang.org/x/crypto/scrypt"
	"golang.org/x/crypto/bcrypt"
)

func Encrypt(plaintext []byte, key []byte) ([]byte, error) {
    if len(key) != 32 {
        return nil, fmt.Errorf("invalid key size: must be 32 bytes")
    }
    c, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }
    gcm, err := cipher.NewGCM(c)
    if err != nil {
        return nil, err
    }
    nonce := make([]byte, gcm.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return nil, err
    }
    return gcm.Seal(nonce, nonce, plaintext, nil), nil
}

func Decrypt(ciphertext []byte, key []byte) ([]byte, error) {
    c, err := aes.NewCipher(key)
    if err != nil {
        return nil, err
    }

    gcm, err := cipher.NewGCM(c)
    if err != nil {
        return nil, err
    }

    nonceSize := gcm.NonceSize()
    if len(ciphertext) < nonceSize {
        return nil, fmt.Errorf("ciphertext too short")
    }

    nonce, actualCiphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

    plaintext, err := gcm.Open(nil, nonce, actualCiphertext, nil)
    if err != nil {
        return nil, err
    }

    return plaintext, nil
}

func AddUser(password []byte, key []byte) ([]byte, []byte, []byte){
    // salt
    salt := make([]byte, 16)
    _, err := rand.Read(salt)
    if err != nil {
        log.Fatal("Cant generate salt: %v", err)
    }

    // udk
    udk, err := scrypt.Key(password, salt, 32768, 8, 1, 32)
    if err != nil {
        log.Fatal("Failed to generate udk: %v", err)
    }

    // hashed pass
    hashed, err := bcrypt.GenerateFromPassword(password, 14)
    if err != nil{
        log.Fatal("Couldnt generate password: %v", err)
    }

    // master key
    master_key,err := Encrypt(key, udk)
    if err != nil{
        log.Fatal("Couldnt generate master key: %v", err)
    }
   
    return hashed, salt, master_key
}


// output salt, udk, hashed, master key > users, admin
// encrypt, hashes password to bcrypt, generate UDK from original password, then use DEK (key) + UDK to generate encrypted key
func Init(password []byte) ([]byte, []byte, []byte, []byte){
    // DEK random
    key := make([]byte, 32)
    _, err := rand.Read(key)
    if err != nil {
        log.Fatal("Failed to generate random key: %v", err)
    }
    
    // salt
    salt := make([]byte, 16)
    _, err = rand.Read(salt)
    if err != nil {
        log.Fatal("Cant generate salt: %v", err)
    }

    // udk
    udk, err := scrypt.Key(password, salt, 32768, 8, 1, 32)
    if err != nil {
        log.Fatal("Failed to generate udk: %v", err)
    }

    // hashed pass
    hashed, err := bcrypt.GenerateFromPassword(password, 14)
    if err != nil{
        log.Fatal("Couldnt generate password: %v", err)
    }

    // master key
    master_key,err := Encrypt(key, udk)
    if err != nil{
        log.Fatal("Couldnt generate master key: %v", err)
    }
   
    return hashed, salt, master_key, key
}
