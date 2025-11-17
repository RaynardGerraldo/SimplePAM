package service

import (
    "SimplePAM/models"
    "SimplePAM/crypto"
    "SimplePAM/parser"
    "log"
    "fmt"
    "os/exec"
    tea "github.com/charmbracelet/bubbletea"
)

type sshFinishedMsg struct{
    err error
}

type TUI struct {
    Choices  []string
    Cursor   int
    Selected map[int]struct{}
    Servers  []string
    ErrorMessage string
    Key []byte
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

    server, ok := raw.([]models.Server)
    if !ok {
        return nil, fmt.Errorf("Invalid server format")
    }

    return server, nil
}

func initialModel(username string, key []byte) (TUI,error) {
    servers, err := allowed(username)
    if err != nil {
        return TUI{}, err
    }
    return TUI{
        Choices: []string{"server-prod", "server-test", "server-misc"},
        Selected: make(map[int]struct{}),
        Servers: servers,
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
            for _,s := range t.Servers {
                if s == server_name {
                    matched = true
                    t.ErrorMessage = ""
                    servers_list, err := parseServers()
                    if err != nil {
                        t.ErrorMessage = fmt.Sprintf("Couldnt parse servers: %v", err)
                        return t, nil
                    }
                    for _, sl := range servers_list {
                        if sl.Server == server_name {
                            login := sl.Name + "@" + sl.IP
                            password,err := crypto.Decrypt(sl.Password, t.Key)
                            if err != nil {
                                log.Fatal("Cannot decrypt password: %v", err)
                            }
                            cmd := exec.Command("sshpass", "-p", string(password), "ssh", login)
                            return t, tea.ExecProcess(cmd, func(err error) tea.Msg {
                                return sshFinishedMsg{err: err}
                            })
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

func SSH(key []byte, username string) error{
    if len(key) > 0 {
        model, err := initialModel(username, key)
        if err != nil {
            return fmt.Errorf("init failed: %w", err)
        }
        p := tea.NewProgram(model)
        if _, err := p.Run(); err != nil {
            return fmt.Errorf("TUI failed: %w", err)
        }
    }
    return fmt.Errorf("\nYou are not logged in. Try again.")
}
