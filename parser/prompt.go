package parser

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "golang.org/x/term"
    "golang.org/x/crypto/ssh/terminal"
)

func Prompt(username string) ([]byte,error) {
    // old term state
    before, err := term.GetState(int(syscall.Stdin))
    if err != nil {
        return nil, fmt.Errorf("stdin is not a terminal: %w", err)
    }

    // monitor for sigint or sigterm (ctrl c)
    c := make(chan os.Signal, 1)
    signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)

    // restore term to old state on ctrl c
    go func() {
        <-c
        term.Restore(int(syscall.Stdin), before)
        os.Exit(1)
    }()

    fmt.Printf("Enter password for %s: ", username)
    password, err := terminal.ReadPassword(int(syscall.Stdin))
    fmt.Println()
    if err != nil {
        term.Restore(int(syscall.Stdin), before)
        return nil, fmt.Errorf("failed to read password: %w", err)
    }
    return password, nil
}
