package internal

import (
    "SimplePAM/models"
    "gorm.io/gorm"
)

func Allowed(db *gorm.DB, username string) ([]string, error) {
    var user models.User
    var all []string
    result := db.Preload("Servers").Where("username = ?", username).First(&user)
    if result.Error != nil {
        return nil, result.Error
    }
    for _, s := range user.Servers {
        all = append(all, s.Server)
    }
    return all, nil
}

func ServersList(db *gorm.DB) ([]models.Server, error) {
    var servers []models.Server
    result := db.Find(&servers)
    if result.Error != nil {
        return nil, result.Error
    }
    
    return servers, nil
}
