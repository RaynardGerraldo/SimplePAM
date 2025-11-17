package parser

import (
    "SimplePAM/models"
    "encoding/json"
    "os"
    "io/ioutil"
)

func Unmarshal(filename string) (any, error) {
    var unmarshalled any
    jsonfile, err := os.Open(filename)
    if err != nil {
        return nil, err
    }
    defer jsonfile.Close()

    bytes, err := ioutil.ReadAll(jsonfile)
    if err != nil {
        return nil, err
    }
    
    if filename == "users.json" || filename == "admin.json" {
        var users []models.User
        err := json.Unmarshal(bytes, &users)
        if err != nil {
            return nil, err
        }
        unmarshalled = users
    } else if filename == "servers.json" {
        var servers []models.Server
        err := json.Unmarshal(bytes, &servers)
        if err != nil {
            return nil, err
        }
        unmarshalled = servers
    }

    return unmarshalled, err
}
