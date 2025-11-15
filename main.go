package main

import (
    "fmt"
    "SimplePAM/internal"
    "os"
    "log"
)

func main() {
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
                internal.Auth(username)
           } else {
                log.Fatal("Not enough arguments, try again.")
           } 
        }

        if arg1 == "admin" {
            if len(os.Args) > 2 {
                admin_option = os.Args[2]
                if admin_option == "init" {
                    internal.Init()
                    //internal.Init() initializes admin and servers
                    // init cant be called again if theres already an admin present (check admin.json)
                } else if admin_option == "add-user" {
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        // Register can only run after admin is authenticated
                        DEK, err := internal.Auth(arg1)
                        if err {
                            internal.Register(username, DEK)
                            fmt.Println("\nadding user: %s", username)
                        } else {
                            log.Fatal("Not authorized.")
                        }
                        // internal.Register(username) registers a new user
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
