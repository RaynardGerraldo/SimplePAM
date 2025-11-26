package models

import "gorm.io/gorm"

type User struct {
    gorm.Model 
    Username   string    `gorm:"uniqueIndex;not null" json:"username"`
    Hashed     []byte    `gorm:"not null" json:"hashed"`
    Salt       []byte    `gorm:"not null" json:"salt"`
    Master_Key []byte    `gorm:"not null" json:"master_key"`
    // third table many2many
    Servers    []*Server `gorm:"many2many:user_servers;" json:"servers"`
}

type Server struct {
    gorm.Model
    Server   string `gorm:"uniqueIndex;not null" json:"server"`
    Name     string `gorm:"not null" json:"name"`
    IP       string `gorm:"not null" json:"ip"`
    Password []byte `gorm:"not null" json:"password"`
}
