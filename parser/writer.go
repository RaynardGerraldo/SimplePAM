package parser

import (
    "encoding/json"
    "log"
    "os"
)

func Writer[T any](target []T, filename string) {
    // append to existing
    var existing []T
    exist,err := os.ReadFile(filename)
    if err == nil {
        // unmarshal existing data to existing var
        unmarshal := json.Unmarshal(exist, &existing)
        if unmarshal != nil {
            log.Fatal(unmarshal)
        }
    }

    // append to existing, whether existing is empty or not
    existing = append(existing, target...)

    toJson, err := json.MarshalIndent(existing, "", " ")
    if err != nil {
        log.Fatal("Couldnt parse to JSON", err)
    }

    err = os.WriteFile(filename, toJson, 0644)
    if err != nil{
        log.Fatal("Couldnt write to json", err)
    }
}
