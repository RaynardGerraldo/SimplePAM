package parser

import (
    "encoding/json"
    "os"
    "fmt"
)

func Writer[T any](target []T, filename string) error {
    // append to existing
    var existing []T
    exist,err := os.ReadFile(filename)
    if err == nil {
        // unmarshal existing data to existing var
        err := json.Unmarshal(exist, &existing)
        if err != nil {
            return fmt.Errorf("Error during unmarshall: %w", err)
        }
    }

    // append to existing, whether existing is empty or not
    existing = append(existing, target...)

    toJson, err := json.MarshalIndent(existing, "", " ")
    if err != nil {
        return fmt.Errorf("Couldnt parse to JSON: %w", err)
    }

    err = os.WriteFile(filename, toJson, 0644)
    if err != nil {
        return fmt.Errorf("Couldnt write to JSON: %w", err)
    }
    return nil
}
