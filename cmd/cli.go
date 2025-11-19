package cmd

import (
    "fmt"
    "SimplePAM/internal"
    "SimplePAM/parser"
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
                db,err := parser.OpenCon()
                if err != nil {
                    fmt.Println("Failed to open connection to db: %w", err)
                    os.Exit(1)
                }
                _, err = parser.ReadUsernameDB(db, username)
                if err != nil {
                    fmt.Println("No users exist, run add-user.")
                    os.Exit(1)
                }
                _, _, err = internal.Auth(db, username)
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
                    db,err := parser.OpenCon()
                    if err != nil {
                        fmt.Println("Failed to open connection to db: %w", err)
                        os.Exit(1)
                    }
                    _, err = parser.ReadUsernameDB(db, arg1)
                    if err == nil {
                        fmt.Println("Cant run init, admin already exists")
                        os.Exit(1)
                    }
                    err = internal.Init(db)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to init admin: %v\n", err)
                        os.Exit(1)
                    }
                } else if admin_option == "add-user" {
                    if !checkCreds("pam.db") {
                        fmt.Fprintf(os.Stderr, "Run init first.\n")
                        os.Exit(1)
                    }
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        // Register can only run after admin is authenticated
                        db,err := parser.OpenCon()
                        DEK, valid, err := internal.Auth(db, arg1)
                        if err != nil {
                            fmt.Println("Failed to open connection to db: %w", err)
                            os.Exit(1)
                        }
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Error during auth: %v\n", err)
                            os.Exit(1)
                        }
                        if valid {
                            err := internal.Register(db, username, DEK)
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
