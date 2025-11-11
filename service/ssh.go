package service

import (
    "fmt"
    "os"
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
}

func initialModel() TUI {
    return TUI{
        Choices: []string{"server-prod", "server-test", "server-misc"},
        Selected: make(map[int]struct{}),
    }
}

func (t TUI) Init() tea.Cmd {
    return nil
}

func (t TUI) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
            cmd := exec.Command("sshpass", "-p", "yourpasswordhere", "ssh", "localhost")
            return t, tea.ExecProcess(cmd, func(err error) tea.Msg {
                return sshFinishedMsg{err: err}
            })
        }

    case sshFinishedMsg:
        t.Cursor = 0
        // redefined new Selected, doesnt use previous one, forces checks to be removed
        t.Selected = make(map[int]struct{})
        return t, tea.ClearScreen
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

    s += "\nPress q to quit.\n"
    return s
}

func SSH(auth bool) {
    if auth {
        p := tea.NewProgram(initialModel())
        if _, err := p.Run(); err != nil {
            fmt.Printf("Alas, there's been an error: %v", err)
            os.Exit(1)
        }
    } else {
        fmt.Println("\nYou are not logged in. Try again.")
    }
}
