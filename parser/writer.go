package parser

import (
    "encoding/json"
    "log"
    "os"
)

func Writer[T any](target []T, filename string) {
    toJson, err := json.MarshalIndent(target, "", " ")
    if err != nil {
        log.Fatal("Couldnt parse to JSON", err)
    }

    err = os.WriteFile(filename, toJson, 0644)
    if err != nil{
        log.Fatal("Couldnt write to json", err)
    }
}
