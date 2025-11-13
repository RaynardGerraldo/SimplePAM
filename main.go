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
    //internal.InitServers()
    //os.Exit(0)
    username := ""
    admin_option := ""
    if len(os.Args) > 1 {
        arg1 := os.Args[1]
        if arg1 == "user" {
           if len(os.Args) > 2 {
                username = os.Args[2]
                if len(username) == 0 {
                    log.Fatal("No username given, try again.")
                }
                fmt.Print("Enter your password: ")
                password, err := terminal.ReadPassword(int(syscall.Stdin))
                if err != nil {
                    log.Fatal(err)
                }
                internal.Auth(username, password)
           } else {
                log.Fatal("Not enough arguments, try again.")
           } 
        }

        if arg1 == "admin" {
            if len(os.Args) > 2 {
                admin_option = os.Args[2]
                if admin_option == "init" {
                    //internal.Init() initializes admin and servers
                    // init cant be called again if theres already an admin present (check admin.json)
                    fmt.Printf("Call internal.Init() here")
                } else if admin_option == "add-user" {
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        // internal.Register(username) registers a new user
                        // Register can only run after admin is authenticated
                        fmt.Printf("adding user: %s", username)
                    } else {
                        fmt.Printf("Not enough arguments, try again")
                    }
                }
                if len(admin_option) == 0 {
                    fmt.Printf("No option given, try again")
                }
            } else {
                fmt.Printf("Not enough arguments, try again")
            }
        }
    }
    // PAM shouldnt have registration, but for now concept for admin input of users.
    //internal.Register(username, password)
}
