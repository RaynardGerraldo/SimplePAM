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
    Jwt   string `json:"jwt"`
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


func LoginCall(username string) (string, string, error){
    password, err := parser.Prompt(username)
    if err != nil {
        return "", "", err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
    }
    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", "", err
    }

    resp, err := http.Post("http://localhost:8080/login", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", "", fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return "", "", fmt.Errorf("access denied: %s", string(body))
    }

    var result LoginResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "", "", fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return "", "", fmt.Errorf("bad response: %v", result.Error)
    }
    return result.Token, result.Jwt, nil
}

func RegisterCall(username string, key string, jwt string) (string, error) {
    var serverName string
    password,err := parser.Prompt(username)
    if err != nil {
        return "", err
    }

    fmt.Printf("Server name to assign? ")
    fmt.Scan(&serverName)

    values := map[string]string{
        "username": username,
        "password": string(password),
        "key": key,
        "servername": serverName,
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", "http://localhost:8080/register", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + jwt)

    client := &http.Client{}
    resp, err := client.Do(req)
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

func AdminCall() (string, string, error){
    fmt.Println("Your admin username is 'admin' by default")
    username := "admin"
    password, err := parser.Prompt(username)
    if err != nil {
        return "", "", err
    }

    values := map[string]string{
        "username": username,
        "password": string(password),
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", "", err
    }

    resp, err := http.Post("http://localhost:8080/initadmin", "application/json", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", "", fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    body, _ := ioutil.ReadAll(resp.Body)

    if resp.StatusCode != 200 {
        return "", "", fmt.Errorf("%s\n", string(body))
    }

    var result LoginResp
    err = json.Unmarshal(body, &result)
    if err != nil {
        return "", "", fmt.Errorf("cannot unmarshal: %w", err)
    }

    if result.Error != "" {
        return "", "", fmt.Errorf("bad response: %w", result.Error)
    }

    return result.Token, result.Jwt, nil
}

func ServerCall(key string, jwt string) (string, error) {
    var serverName string
    var name string
    var ip string
    var port uint16
    fmt.Printf("Server Name in PAM? (ex: server-prod) ")
    fmt.Scan(&serverName)

    fmt.Printf("Server IP? (default is localhost) ")
    fmt.Scan(&ip)

    fmt.Printf("Server Port? (default is 22) ")
    fmt.Scan(&port)

    fmt.Printf("Server username? ")
    fmt.Scan(&name)

    password,err := parser.Prompt("server " + name)
    if err != nil {
        return "", err
    }
    
    values := map[string]any{
        "servername": serverName,
        "username": name,
        "password": string(password),
        "key": key,
        "ip": ip,
        "port": port,
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return "", err
    }

    req, err := http.NewRequest("POST", "http://localhost:8080/initserver", bytes.NewBuffer(jsondata))
    if err != nil {
        return "", err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + jwt)

    client := &http.Client{}
    resp, err := client.Do(req)
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

func AllowedListCall(username string, jwt string) ([]string, []models.Server, error){
    var servers_list []models.Server
    var allowed_servers []string
    values := map[string]string{
        "username": username,
    }

    jsondata, err := json.Marshal(values)

    if err != nil {
        return nil, nil, err
    }

    req, err := http.NewRequest("POST", "http://localhost:8080/allowedservers", bytes.NewBuffer(jsondata))
    if err != nil {
        return nil, nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + jwt)

    client := &http.Client{}
    resp, err := client.Do(req)
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
   
    req, err = http.NewRequest("GET", "http://localhost:8080/serverslist", bytes.NewBuffer(jsondata))
    if err != nil {
        return nil, nil, err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + jwt)

    client = &http.Client{}
    resp, err = client.Do(req)
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

func AddtoUserCall(jwt string) error {
    var username string
    var servername string
    fmt.Printf("Username? ")
    fmt.Scan(&username)

    fmt.Printf("Server name to add to user? ")
    fmt.Scan(&servername)

    values := map[string]string{
        "username": username,
        "servername": servername,
    }
    jsondata, err := json.Marshal(values)

    if err != nil {
        return err
    }

    req, err := http.NewRequest("POST", "http://localhost:8080/addtouser", bytes.NewBuffer(jsondata))

    if err != nil {
        return err
    }

    req.Header.Set("Content-Type", "application/json")
    req.Header.Set("Authorization", "Bearer " + jwt)
    client := &http.Client{}
    resp, err := client.Do(req)
    
    if err != nil {
        return fmt.Errorf("failed to connect to PAM server: %w", err)
    }
    defer resp.Body.Close()

    return nil
}
