package internal

import (
    "SimplePAM/models"
    "SimplePAM/internal"
    "log"
    "encoding/json"
    "os"
)

func toAdmin() {
    var admin models.Admin
    // generate dek
    // aes_gcm
    admin.Username = admin
    // admin.Password_Hash = ask for pass, hash, and put here
    admin.Enc_Master_Key = ""
    admin_ist := []models.Admin{admin}
    internal.Writer(admin_ist, "admin.json")
}
func toServer() {
    var server models.Server
    server.Server = "server-prod"
    server.Name = "rayray"
    server.IP = "localhost"
    server.Password = ""

    servers := []models.Server{server}
    internal.Writer(servers, "servers.json")
}

func Init(){
    toAdmin()
    toServer()
}
