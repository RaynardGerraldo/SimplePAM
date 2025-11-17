package parser

import (
    "SimplePAM/models"
    "encoding/json"
    "os"
    "io/ioutil"
    "log"
)

func Unmarshal(filename string) any {
    var unmarshalled any
    jsonfile, err := os.Open(filename)
    if err != nil {
        log.Fatal(err)
    }
    defer jsonfile.Close()

    bytes, err := ioutil.ReadAll(jsonfile)
    if err != nil {
        log.Fatal(err)
    }
    
    if filename == "users.json" || filename == "admin.json" {
        var users []models.User
        json.Unmarshal(bytes, &users)
        unmarshalled = users
    } else if filename == "servers.json" {
        var servers []models.Server
        json.Unmarshal(bytes, &servers)
        unmarshalled = servers
    }

    return unmarshalled
}
