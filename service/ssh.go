package service

import (
    "SimplePAM/models"
    "SimplePAM/crypto"
    "SimplePAM/parser"
    "fmt"
    "golang.org/x/crypto/ssh"
    "os"
    "golang.org/x/term"
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

func internalSSH(username string, password string, ip string) error {
    config := &ssh.ClientConfig {
        User: username,
        Auth: []ssh.AuthMethod {
            ssh.Password(password),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }

    client, err := ssh.Dial("tcp", ip+":22", config)
    if err != nil {
        return fmt.Errorf("failed to dial: %w", err)
    }
    defer client.Close()

    session, err := client.NewSession()
    if err != nil {
        return fmt.Errorf("Failed to create session: ", err)
    }
    defer session.Close()

    // looks and interactive 
    fd := int(os.Stdin.Fd())
    state, err := term.MakeRaw(fd)
    if err != nil {
        return fmt.Errorf("failed to set raw mode: %w", err)
    }
    defer term.Restore(fd,state)

    w, h, err := term.GetSize(fd)
    if err != nil {
        return fmt.Errorf("failed to get term size: %w", err)
    }

    modes := ssh.TerminalModes {
        ssh.ECHO: 1,
        ssh.TTY_OP_ISPEED: 14400,
        ssh.TTY_OP_OSPEED: 14400,
    }

    if err := session.RequestPty("xterm-256color", h, w, modes); err != nil {
        return fmt.Errorf("request pty failed: %w", err)
    }

    session.Stdout = os.Stdout
    session.Stderr = os.Stderr
    session.Stdin = os.Stdin

    if err := session.Shell(); err != nil {
        return fmt.Errorf("failed to start shell: %w", err)
    }

    return session.Wait()
}

func allowed(username string) ([]string, error) {
    raw, err := parser.Unmarshal("users.json")
    if err != nil {
        return nil, err
    }

    users, ok := raw.([]models.User)
    if !ok {
        return nil, fmt.Errorf("Invalid user format")
    }

    for _,u := range users {
        if u.Username == username {
            return u.Servers,nil
        }
    }

    return nil, fmt.Errorf("User not found")
}

func parseServers() ([]models.Server, error) {
    raw, err := parser.Unmarshal("servers.json")
    if err != nil {
        return nil, err
    }

    servers, ok := raw.([]models.Server)
    if !ok || len(servers) == 0 {
        return nil, fmt.Errorf("Invalid servers.json format")
    }
    return servers, nil
}

func initialModel(username string, key []byte, server_list []models.Server) (TUI,error) {
    allowed_servers, err := allowed(username)
    if err != nil {
        return TUI{}, err
    }
    return TUI{
        Choices: []string{"server-prod", "server-test", "server-misc"},
        Selected: make(map[int]struct{}),
        Server_List: server_list,
        Allowed: allowed_servers,
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

func SSH(key []byte, username string) error {
    if len(key) == 0 {
        return fmt.Errorf("\nYou are not logged in. Try again.")
    }
    servers_list, err := parseServers()
    if err != nil {
        return err
    }

    // loop, load up TUI, wait for either "q" or server selection, then quit or ssh in. if ssh in loop back to TUI after.
    for {
        model, err := initialModel(username, key, servers_list)
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
        password, err := crypto.Decrypt(target.Password, key)
        if err != nil {
            fmt.Printf("Cannot decrypt password: %v", err)
            continue
        }

        err = internalSSH(target.Name, string(password), target.IP)

        if err != nil {
            fmt.Printf("SSH connection error: %v", err)
        }
    }
   
    return nil
}
