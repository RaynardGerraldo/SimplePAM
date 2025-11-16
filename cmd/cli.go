package cmd

import (
    "fmt"
    "SimplePAM/internal"
    "os"
    "log"
)

func checkCreds(filename string) bool {
    info, err := os.Stat(filename)
    if os.IsNotExist(err) {
        return false
    }
    return !info.IsDir()
}

func Cli() {
    username := ""
    admin_option := ""
    if len(os.Args) > 1 {
        arg1 := os.Args[1]
        if arg1 == "user" {
           if len(os.Args) > 2 {
                username = os.Args[2]
                if len(username) == 0 {
                    fmt.Println("No username given, try again.")
                }
                if !checkCreds("users.json") {
                    fmt.Println("No users exist, run add-user.")
                    os.Exit(1)
                }
                internal.Auth(username)
           } else {
                fmt.Println("Not enough arguments, try again.")
           } 
        }

        if arg1 == "admin" {
            if len(os.Args) > 2 {
                admin_option = os.Args[2]
                if admin_option == "init" {
                    if checkCreds("admin.json") {
                        fmt.Println("Cant run init, admin already exists")
                        os.Exit(1)
                    }
                    internal.Init()
                } else if admin_option == "add-user" {
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        // Register can only run after admin is authenticated
                        DEK, valid := internal.Auth(arg1)
                        if valid {
                            internal.Register(username, DEK)
                            fmt.Println("\nadding user:", username)
                        } else {
                            log.Fatal("Not authorized.")
                        }
                    } else {
                        fmt.Println("Specify user for add-user.")
                    }
                } else {
                    fmt.Println("Invalid argument.")
                }
            } else {
                fmt.Println("Not enough arguments, try again")
            }
        }
    }
}
