package internal

import (
    "SimplePAM/models"
    "encoding/json"
    "SimplePAM/parser"
    "net/http"
    "bytes"
    "io/ioutil"
    "fmt"
)


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

type AllowedListResp struct {
    List []models.Server `json:"list"`
    Allowed []string `json:"allowed"`
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

func AllowedListCall(username string) ([]string, []models.Server, error){
    var servers_list []models.Server
    var allowed_servers []string
    values := map[string]string{
        "username": username,
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return nil, nil, err
    }

    resp, err := http.Post("http://localhost:8080/allowedservers", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body,_ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return nil, nil, fmt.Errorf("%s\n", string(body))
    }
    
    var result AllowedListResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, nil, fmt.Errorf("Cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return nil, nil, fmt.Errorf("bad response: %w", result.Error)
    }
    
    allowed_servers = result.Allowed
    
    resp, err = http.Get("http://localhost:8080/serverslist")
    if err != nil {
        return nil, nil, fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body,_ = ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return nil, nil, fmt.Errorf("%s\n", string(body))
    }
    
    err = json.Unmarshal(body, &result)
    if err != nil {
        return nil, nil, fmt.Errorf("Cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return nil,nil,fmt.Errorf("bad response: %w", result.Error)
    }
    
    servers_list = result.List

    return allowed_servers, servers_list, nil
}

