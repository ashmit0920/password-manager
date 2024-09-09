package main

import (
	"fmt"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   int
	choices  []string
	selected string
}

// Initial model with choices for the menu
func initialModel() model {
	return model{
		choices: []string{"View Stored Passwords", "Add New Password", "Quit"},
	}
}

func (m model) Init() tea.Cmd {
	// No initial command
	return nil
}

// Update the model based on the user input (navigation)
func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c", "q":
			return m, tea.Quit
		case "up":
			if m.cursor > 0 {
				m.cursor--
			}
		case "down":
			if m.cursor < len(m.choices)-1 {
				m.cursor++
			}
		case "enter":
			m.selected = m.choices[m.cursor]
			switch m.selected {
			case "View Stored Passwords":
				fmt.Println("\nLoading stored passwords...")
				passwords, err := loadPasswords("passwords.enc")
				if err != nil {
					fmt.Println("No passwords found!")
				} else {
					fmt.Println("Stored Passwords:")
					for _, entry := range passwords {
						fmt.Printf("Website: %s, Username: %s, Password: %s\n", entry.Website, entry.Username, entry.Password)
					}
				}
			case "Add New Password":
				var entry PasswordEntry
				fmt.Println("\nEnter Website: ")
				fmt.Scanln(&entry.Website)
				fmt.Println("Enter Username: ")
				fmt.Scanln(&entry.Username)
				fmt.Println("Enter Password: ")
				fmt.Scanln(&entry.Password)
				savePassword(entry, "passwords.enc")
				fmt.Println("Password added successfully!")
			case "Quit":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

// The view to render the TUI menu
func (m model) View() string {
	s := "Welcome to Go Password Manager\n\n"
	s += "Choose an option:\n\n"

	// Render the menu choices
	for i, choice := range m.choices {
		if m.cursor == i {
			s += fmt.Sprintf("> %s\n", choice) // Highlight the selected choice
		} else {
			s += fmt.Sprintf("  %s\n", choice)
		}
	}

	s += "\nPress q to quit.\n"
	return s
}
