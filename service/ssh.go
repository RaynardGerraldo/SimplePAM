package service

import (
    "SimplePAM/models"
    "SimplePAM/crypto"
    "encoding/base64"
    "fmt"
    "golang.org/x/crypto/ssh"
    "os"
    "io"
    //"golang.org/x/term"
    tea "github.com/charmbracelet/bubbletea"
)

type sshFinishedMsg struct{
    err error
}

type TUI struct {
    Choices  []string
    Cursor   int
    Selected map[int]struct{}
    Server_List []models.Server
    Allowed  []string
    ErrorMessage string
    Key []byte
    // use model in memory, gives nil value as well
    Target *models.Server
}

func InternalSSH(reader io.Reader, writer io.Writer, username string, password string, ip string, port uint16) error {
    config := &ssh.ClientConfig {
        User: username,
        Auth: []ssh.AuthMethod {
            ssh.Password(password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    address := fmt.Sprintf("%s:%d", ip, port)   
    client, err := ssh.Dial("tcp", address, config)
    if err != nil {
        return fmt.Errorf("failed to dial: %w", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return fmt.Errorf("Failed to create session: %w", err)
    }
    defer session.Close()

    session.Stdin = reader
    session.Stdout = writer
    session.Stderr = writer

    modes := ssh.TerminalModes {
        ssh.ECHO: 1,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }

    if err := session.RequestPty("xterm-256color", 40, 80, modes); err != nil {
        return fmt.Errorf("request pty failed: %w", err)
    }

    if err := session.Shell(); err != nil {
        return fmt.Errorf("failed to start shell: %w", err)
    }
    
    return session.Wait()
}

func initialModel(username string, key []byte, server_list []models.Server, allowed []string) (TUI,error) {
    return TUI{
        Choices: allowed,
        Selected: make(map[int]struct{}),
        Server_List: server_list,
        Allowed: allowed,
        Key: key,
    }, nil
}

func (t TUI) Init() tea.Cmd {
    return nil
}

func (t TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
    var matched bool
    switch msg := msg.(type) {

    case tea.KeyMsg:
        switch msg.String() {
        case "ctrl+c", "q":
            return t, tea.Quit

        case "up", "k":
            if t.Cursor > 0 {
                t.Cursor--
            }

        case "down", "j":
            if t.Cursor < len(t.Choices)-1 {
                t.Cursor++
            }

        case "enter", " ":
            _, ok := t.Selected[t.Cursor]
            if ok {
                delete(t.Selected, t.Cursor)
            } else {
                t.Selected[t.Cursor] = struct{}{}
            }
            server_name := t.Choices[t.Cursor]
            for _,s := range t.Allowed {
                if s == server_name {
                    matched = true
                    t.ErrorMessage = ""
                    for _, sl := range t.Server_List {
                        if sl.Server == server_name {
                            t.Target = &sl
                            return t, tea.Quit
                        }
                    }
                }
            }
 
            if !matched {
                t.Cursor = 0
                t.Selected = make(map[int]struct{})
                t.ErrorMessage = "You need to request access for this server."
                return t, nil
            }
        }

    case sshFinishedMsg:
        t.Cursor = 0
        // redefined new Selected, doesnt use previous one, forces checks to be removed
        t.Selected = make(map[int]struct{})
        return t, nil
    }
    return t, nil
}

func (t TUI) View() string {
    s := "What server to login to?\n\n"
    for i, choice := range t.Choices {
        cursor := " " // no Cursor
        if t.Cursor == i {
            cursor = ">" // Cursor
        }
        checked := " " // not Selected

        if _, ok := t.Selected[i]; ok {
            checked = "x"  // Selected
        }
        s += fmt.Sprintf("%s [%s] %s\n", cursor, checked, choice)
    }

    if t.ErrorMessage != "" {
        s += "\nError: " + t.ErrorMessage
    }
    s += "\nPress q to quit.\n"
    return s
}

func SSH(key string, username string, allowed []string, servers_list []models.Server) error {
    if len(key) == 0 {
        return fmt.Errorf("\nYou are not logged in. Try again.")
    }

    decodedKey, err := base64.StdEncoding.DecodeString(key)
    if err != nil {
        return fmt.Errorf("Couldnt decode base64")
    }

    // loop, load up TUI, wait for either "q" or server selection, then quit or ssh in. if ssh in loop back to TUI after.
    for {
        model, err := initialModel(username, decodedKey, servers_list, allowed)
        if err != nil {
            return fmt.Errorf("init failed: %w", err)
        }
        p := tea.NewProgram(model)

        t, err := p.Run()
        if err != nil {
            return fmt.Errorf("TUI failed: %w", err)
        }

        final_t, ok := t.(TUI)
        if !ok {
            return fmt.Errorf("internal model error")
        }

        if final_t.Target == nil {
            return nil
        }

        target := *final_t.Target
        password, err := crypto.Decrypt(target.Password, decodedKey)
        if err != nil {
            fmt.Errorf("Cannot decrypt password: %w", err)
            continue
        }

        err = InternalSSH(os.Stdin, os.Stdout, target.Name, string(password), target.IP, target.Port)

        if err != nil {
            fmt.Errorf("SSH connection error: %w", err)
        }
    }
   
    return nil
}
