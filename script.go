package main

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

type model struct {
	cursor   int
	choices  []string
	selected string
	quit     bool
}

const (
	menuTitle = "Choisissez un routeur"
	quitLabel = "Quitter"
	columns   = 3
	fileName  = "pe.txt"
)

type message string

const (
	messageQuit    message = "quit"
	messageUp      message = "up"
	messageDown    message = "down"
	messageInvalid message = "invalid"
	messageEnter   message = "enter"
)

func (m *model) Init() tea.Cmd {
	return nil
}

// Ajoutez cette fonction pour rediriger la sortie standard vers un fichier
func redirectStdoutToFile(fileName string) (*os.File, *os.File, error) {
	file, err := os.Create(fileName)
	if err != nil {
		return nil, nil, err
	}

	originalStdout := os.Stdout
	os.Stdout = file

	return file, originalStdout, nil
}

// Ajoutez cette fonction pour restaurer la sortie standard
func restoreStdout(file *os.File, originalStdout *os.File) {
	file.Close()
	os.Stdout = originalStdout
}

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "q":
			return m, tea.Quit
		case "up", "k":
			return m, m.update(messageUp)
		case "down", "j":
			return m, m.update(messageDown)
		case "left", "h":
			return m, m.update(messageInvalid)
		case "right", "l":
			return m, m.update(messageInvalid)
		case "enter":
			return m, m.update(messageEnter)
		}
	}

	return m, nil
}

func (m *model) View() string {
	var b strings.Builder

	clearScreen()

	fmt.Fprint(&b, menuTitle+"\n\n")

	printTable(&b, m.choices, columns, m.cursor)

	return b.String()
}

func printTable(b *strings.Builder, items []string, columns int, cursor int) {
	colWidth := 20

	for i := 0; i < len(items); i += columns {
		for j := 0; j < columns && i+j < len(items); j++ {
			cursorChar := " "
			if i+j == cursor {
				cursorChar = ">"
			}

			fmt.Fprintf(b, "%s %-*s", cursorChar, colWidth, items[i+j])
		}
		fmt.Fprintln(b)
	}
}

func clearScreen() {
	cmd := exec.Command("clear") // pour Linux et MacOS
	if _, err := cmd.Output(); err != nil {
		cmd = exec.Command("cmd", "/c", "cls") // pour Windows
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
}

// Modifiez votre fonction main pour utiliser ces fonctions
func main() {
	choices, err := readValuesFromFile(fileName)
	if err != nil {
		fmt.Printf("Erreur lors de la lecture du fichier : %v\n", err)
		os.Exit(1)
	}

	// Rediriger la sortie standard vers un fichier
	file, originalStdout, err := redirectStdoutToFile("output.txt")
	if err != nil {
		fmt.Printf("Erreur lors de la redirection de la sortie standard : %v\n", err)
		os.Exit(1)
	}
	defer restoreStdout(file, originalStdout)

	m := &model{
		choices: choices,
	}
	if err := tea.NewProgram(m).Start(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}

	// Imprimez la valeur sélectionnée à la sortie standard
	fmt.Println(m.selected)
}

func (m *model) update(msg message) tea.Cmd {
	switch msg {
	case messageUp:
		m.cursor = (m.cursor - 1 + len(m.choices)) % len(m.choices)
	case messageDown:
		m.cursor = (m.cursor + 1) % len(m.choices)
	case messageEnter:
		if m.choices[m.cursor] == quitLabel {
			return tea.Quit
		}
		m.selected = m.choices[m.cursor]
		m.quit = true
		return tea.Quit
	}
	return nil
}

func readValuesFromFile(fileName string) ([]string, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var values []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		values = append(values, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return values, nil
}
