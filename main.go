package main

import (
    "fmt"
    "SimplePAM/internal"
    "os"
    "syscall"
    "golang.org/x/crypto/ssh/terminal"
    "log"
)

func main() {
    useroption := ""
    username := ""
    if len(os.Args) > 2 {
        useroption = os.Args[1]
        username = os.Args[2]
    }

    if (useroption != "--u") || (username == "") {
        log.Fatal("No --u (username) given, try again.")
    }
    username = os.Args[2]
    fmt.Print("Enter your password: ")
    password, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        log.Fatal(err)
    }
    
    // PAM shouldnt have registration, but for now concept for admin input of users.
    //internal.Register(username, password)
    
    internal.Auth(username,password)
    //fmt.Println(username)
    //fmt.Println(passwd)
}
