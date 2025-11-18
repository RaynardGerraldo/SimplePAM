package cmd

import (
    "fmt"
    "SimplePAM/internal"
    "os"
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
                _, _, err := internal.Auth(username)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Error during auth: %v\n", err)
                    os.Exit(1)
                }
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
                    err := internal.Init()
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to init admin: %v\n", err)
                        os.Exit(1)
                    }
                } else if admin_option == "add-user" {
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        // Register can only run after admin is authenticated
                        DEK, valid, err := internal.Auth(arg1)

                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Error during auth: %v\n", err)
                            os.Exit(1)
                        }
                        if valid {
                            err := internal.Register(username, DEK, checkCreds("users.json"))
                            if err != nil {
                                fmt.Fprintf(os.Stderr, "Error during register: %v\n", err)
                                os.Exit(1)
                            }
                            fmt.Println("\nadding user:", username)
                        } else {
                            fmt.Println("Not authorized")
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
