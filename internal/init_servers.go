package internal

import (
    "SimplePAM/models"
    "SimplePAM/parser"
    "SimplePAM/crypto"
    "fmt"
    "syscall"
    "golang.org/x/crypto/ssh/terminal"
    "log"
)

func toAdmin() {
    var admin models.User
    fmt.Printf("Your admin username is 'admin' by default")
    admin.Username = "admin"
    fmt.Print("\nEnter your password: ")
    password, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatal(err)
    }
    hashed, salt, master_key := crypto.Init(password)
    admin.Hashed = hashed
    admin.Salt = salt
    admin.Master_Key = master_key
    admin.Servers = []string{}
    
    admin_ins := []models.User{admin}
    parser.Writer(admin_ins, "admin.json")
}

func toServer() {
    var server models.Server
    var name string
    var password string

    fmt.Printf("Try it out with your localhost")
    fmt.Printf("Username ? ")
    fmt.Scan(&name)
    fmt.Printf("\nServer password here: ")
    fmt.Scan(&password)
    
    /*password, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatal(err)
    }*/

    server.Server = "server-prod"
    server.Name = name
    server.IP = "localhost"
    server.Password = password

    servers := []models.Server{server}
    parser.Writer(servers, "servers.json")
}

func Init(){
    toAdmin()
    toServer()
}
