package parser

import (
  "SimplePAM/models"
  "gorm.io/driver/sqlite"
  "gorm.io/gorm"
  "fmt"
)


func OpenCon() (*gorm.DB, error) {
    db, err := gorm.Open(sqlite.Open("pam.db"), &gorm.Config{})
    if err != nil {
        return nil, fmt.Errorf("failed to open db: %w", err)
    }

    return db,nil
}

func InitDB[T any](db *gorm.DB, model T) {
    db.AutoMigrate(model)
}

func ReadUsernameDB(db *gorm.DB, username string) (*models.User, error) {
    var inUsers *models.User
    result := db.Where("username = ?", username).First(&inUsers)
    if result.Error == gorm.ErrRecordNotFound {
        return nil, result.Error
    } else if result.Error != nil {
        return nil, result.Error
    }
    return inUsers, nil
}

func WriteDB[T any](db *gorm.DB, model T) {
    db.Create(&model)
}

func CheckDB(db *gorm.DB, server_name string) models.Server {
    var inServer models.Server
    db.Where("server = ?", server_name).First(&inServer)
    return inServer
}
