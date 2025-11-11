package models

type User struct {
    Username string `json:"username"`
    Password []byte `json:"password"`
    Servers []string `json:servers"`
}
