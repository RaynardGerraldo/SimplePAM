package cmd

import (
    "bytes"
    "encoding/json"
    "net/http"
    "io/ioutil"
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

type StatusResp struct {
    Error string `json:"error"`
}

func StatusCall(username string) error {
    values := map[string]string{
        "username": username,
    }
    jsondata, err := json.Marshal(values)

    if err != nil {
        return err
    }

    resp, err := http.Post("http://localhost:8080/status", "application/json", bytes.NewBuffer(jsondata))

    if err != nil {
        return fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return fmt.Errorf("access denied: %s", string(body))
    }

    var result StatusResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return fmt.Errorf("bad response: %v", result.Error)
    }
    return nil
}


func LoginCall(username string) (string, error){
    password, err := parser.Prompt(username)
    if err != nil {
        return "", err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
    }
    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", err
    }

    resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("access denied: %s", string(body))
    }

    var result LoginResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "", fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return "", fmt.Errorf("bad response: %v", result.Error)
    }
    return result.Token, nil
}


// todo
func RegisterCall(username string, key string) (string, error) {
    password,err := parser.Prompt(username)
    if err != nil {
        return "", err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
        "key": key,
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

func AdminCall() (string, error){
    fmt.Println("Your admin username is 'admin' by default")
    username := "admin"
    password, err := parser.Prompt(username)
    if err != nil {
        return "", err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", err
    }

    resp, err := http.Post("http://localhost:8080/initadmin", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("%s\n", string(body))
    }

    var result LoginResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "", fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return "", fmt.Errorf("bad response: %w", result.Error)
    }

    return result.Token, nil
}

func ServerCall(key string) (string, error) {
    var name string
    fmt.Println("\nTry it out with your localhost")
    fmt.Printf("Server username? ")
    fmt.Scan(&name)

    password,err := parser.Prompt("server " + name)
    if err != nil {
        return "", err
    }
    
    values := map[string]string{
        "username": name,
        "password": string(password),
        "key": key,
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", err
    }

    resp, err := http.Post("http://localhost:8080/initserver", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body,_ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return "", fmt.Errorf("%s\n", string(body))
    }
    
    var result RegResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "", fmt.Errorf("Cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return "", fmt.Errorf("bad response: %w", result.Error)
    }
    return result.Success, nil
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
                db,err := parser.OpenCon()
                if err != nil {
                    fmt.Fprintf(os.Stderr, "Failed to open connection to db: %v\n", err)
                    os.Exit(1)
                }
                err = StatusCall(username)
                if err != nil {
                    fmt.Fprintf(os.Stderr, "%v\n", err)
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
                    err := StatusCall(arg1)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "%v\n", err)
                        os.Exit(1)
                    }
                    key, err := AdminCall()
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to init admin: %v\n", err)
                        os.Exit(1)
                    }
                    success, err := ServerCall(key)
                    if err != nil {
                        fmt.Fprintf(os.Stderr, "Failed to init server: %v\n", err)
                        os.Exit(1)
                    }
                    fmt.Println(success)
                } else if admin_option == "add-user" {
                    err := StatusCall(arg1)
                    if err == nil {
                        fmt.Fprintf(os.Stderr, "Run init first.\n")
                        os.Exit(1)
                    }
                    if len(os.Args) > 3 {
                        username = os.Args[3]
                        err := StatusCall(username)
                        if err == nil {
                            fmt.Fprintf(os.Stderr, "User already exists\n")
                            os.Exit(1)
                        }
                        token, err := LoginCall(arg1)
                        if err != nil {
                            fmt.Fprintf(os.Stderr, "Cant login to admin: %v\n", err)
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
