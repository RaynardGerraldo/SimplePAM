package cmd

import (
    "SimplePAM/service"
    "SimplePAM/internal"
    "os"
    "fmt"
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
                    os.Exit(1)
                }
                err := internal.StatusCall(username)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "%v\n", err)
                    os.Exit(1)
                }

                token, err := internal.LoginCall(username)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "SSH Failed: %v\n", err)
                    os.Exit(1)
                }
                allowed_servers, servers_list, err := internal.AllowedListCall(username)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Failed to get allowed servers and servers list: %v\n", err)
                    os.Exit(1)
                }
                service.SSH(token, username, allowed_servers, servers_list)
           } else {
                fmt.Println("Not enough arguments, try again.")
           } 
        }

        if arg1 == "admin" {
            if len(os.Args) > 2 {
                admin_option = os.Args[2]
                if admin_option == "init" {
                    err := internal.StatusCall(arg1)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "%v\n", err)
                        os.Exit(1)
                    }
                    key, err := internal.AdminCall()
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to init admin: %v\n", err)
                        os.Exit(1)
                    }
                    success, err := internal.ServerCall(key)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to init server: %v\n", err)
                        os.Exit(1)
                    }
                    fmt.Println(success)
                } else if admin_option == "add-user" {
                    err := internal.StatusCall(arg1)
                    if err == nil {
                        fmt.Fprintf(os.Stderr, "Run init first.\n")
                        os.Exit(1)
                    }
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        err := internal.StatusCall(username)
                        if err == nil {
                            fmt.Fprintf(os.Stderr, "User already exists\n")
                            os.Exit(1)
                        }
                        token, err := internal.LoginCall(arg1)
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Cant login to admin: %v\n", err)
                            os.Exit(1)
                        }
                        success, err := internal.RegisterCall(username, token)
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Failed to register: %v", err)
                            os.Exit(1)
                        }
                        fmt.Println(success)
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
