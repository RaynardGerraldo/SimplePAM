package parser

import (
    "fmt"
    "os"
    "os/signal"
    "syscall"
    "golang.org/x/term"
    "golang.org/x/crypto/ssh/terminal"
    "log"
)

func Prompt() []byte {
    // old term state
    before, err := term.GetState(int(syscall.Stdin))
    if err != nil {
        panic(err)
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

    fmt.Print("Enter password: ")
    password, err := terminal.ReadPassword(int(syscall.Stdin))
    if err != nil {
        term.Restore(int(syscall.Stdin), before)
        log.Fatal(err)
    }
    return password
}
