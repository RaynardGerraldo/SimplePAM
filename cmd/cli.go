package cmd

import (
    "bytes"
    "encoding/json"
    "net/http"
    "io/ioutil"
    "SimplePAM/internal"
    "SimplePAM/parser"
    "SimplePAM/service"
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

type LoginResp struct {
    Token string `json:"token"`
    Error string `json:"error"`
}

type RegResp struct {
    Success string `json:"success"`
    Error string `json:"error"`
}

func LoginCall(username string) ([]byte, error){
    password, err := parser.Prompt(username)
    if err != nil {
        return nil, err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
    }
    jsondata, err := json.Marshal(values)

    if err != nil {
        return nil, err
    }

    resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return nil, fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return nil, fmt.Errorf("access denied: %s", string(body))
    }

    var result LoginResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return nil, fmt.Errorf("bad response: %v", result.Error)
    }
    return []byte(result.Token), nil
}


// todo
func RegisterCall(username string, key []byte) (string, error) {
    password,err := parser.Prompt(username)
    if err != nil {
        return "", err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
        "key": string(key),
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", err
    }

    resp, err := http.Post("http://localhost:8080/register", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("%s\n", string(body))
    }

    var result RegResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "", fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return "", fmt.Errorf("bad response: %v", result.Error)
    }

    return result.Success, nil
}

//func InitCall()

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
                // replace here with api call?
                db,err := parser.OpenCon()
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Failed to open connection to db: %v\n", err)
                    os.Exit(1)
                }
                _, err = parser.ReadUsernameDB(db, username)
                if err != nil {
                    fmt.Println("User doesnt exist.")
                    os.Exit(1)
                }
                token, err := LoginCall(username)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "SSH Failed: %v\n", err)
                    os.Exit(1)
                }
                service.SSH(db, token, username)
           } else {
                fmt.Println("Not enough arguments, try again.")
           } 
        }

        if arg1 == "admin" {
            if len(os.Args) > 2 {
                admin_option = os.Args[2]
                if admin_option == "init" {
                    // replace here with api call?
                    db,err := parser.OpenCon()
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to open connection to db: %v\n", err)
                        os.Exit(1)
                    }
                    _, err = parser.ReadUsernameDB(db, arg1)
                    if err == nil {
                        fmt.Println("Cant run init, admin already exists")
                        os.Exit(1)
                    }

                    // replace with InitCall
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
                        token, err := LoginCall(arg1)
                        if err != nil {
                            fmt.Println("whoopsie: %s", err)
                            os.Exit(1)
                        }
                        success, err := RegisterCall(username, token)
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
