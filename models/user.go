package models

type User struct {
    Username string `json:"username"`
    Password []byte `json:"password"`
    Servers []string `json:servers"`
}

type Server struct {
    Server string `json:"server"`
    Name string `json:"name"`
    IP string `json:"ip"`
    Password string `json:"password"`
}
