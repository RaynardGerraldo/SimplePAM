package models

type User struct {
    Username string `json:"username"`
    Hashed []byte `json:"hashed"`
    Salt []byte `json:"salt"`
    Master_Key []byte `json:"master_key"`
    Servers []string `json:"servers"`
}

// ADJUST THIS FOR DEK
type Server struct {
    Server string `json:"server"`
    Name string `json:"name"`
    IP string `json:"ip"`
    Password string `json:"password"`
}
