package tui

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	"github.com/charmbracelet/bubbles/progress"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/pkg/sftp"
)

var (
	docStyle           = lipgloss.NewStyle().Margin(2, 2)
	statusMessageStyle = lipgloss.NewStyle().
				Foreground(lipgloss.AdaptiveColor{Light: "#04B575", Dark: "#04B575"}).
				Render

	fileItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#ffffff", Dark: "#ffffff"}).
			Render
	dirItemStyle = lipgloss.NewStyle().
			Foreground(lipgloss.AdaptiveColor{Light: "#64CDEF", Dark: "#64CDEF"}).
			Render
)

// Struct that keeps the progress bar percantage
type barPercentage float64

// Holds the state of the tui
type Model struct {
	List       list.Model   // the list of items
	SftpClient *sftp.Client // the sftp client
	currentDir string       // current directory
	progress   progress.Model
}

func (m Model) Init() tea.Cmd {
	return nil
}

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	var cmds []tea.Cmd

	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit
		case "backspace":
			cmds = moveDir(&m, "..", cmds)
			return m, tea.Batch(cmds...)
		case "enter":
			var cmd tea.Cmd
			selectedItem := m.List.SelectedItem().(*item).rawValue

			selectedItemName := selectedItem.Name()
			if selectedItem.IsDir() {
				cmds = moveDir(&m, selectedItemName, cmds)
			} else {
				cmd = m.List.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Downloading %s", selectedItemName)))
				cmds = append(cmds, cmd)
				cmds = append(cmds, m.downloadFile(selectedItem))
			}

			return m, tea.Batch(cmds...)
		}

	case *barPercentage:
		cmd := m.progress.SetPercent(float64(*msg) / 100.0)
		//fmt.Println(int(*msg))
		if int(*msg) != 100 {
			return m, tea.Batch(cmd, tea.Tick(time.Millisecond*500, func(t time.Time) tea.Msg {
				return msg
			}))
		} else {
		}

	case progress.FrameMsg:
		progressModel, cmd := m.progress.Update(msg)
		m.progress = progressModel.(progress.Model)
		return m, cmd

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.List.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.List, cmd = m.List.Update(msg)
	return m, cmd
}

func moveDir(m *Model, selectedItemName string, cmds []tea.Cmd) []tea.Cmd {
	currentWd, err := m.SftpClient.RealPath(m.SftpClient.Join(m.currentDir, selectedItemName))
	handleError(err)
	m.currentDir = currentWd

	cmd := m.List.SetItems(CreateItemListModel(currentWd, m.SftpClient))
	cmds = append(cmds, cmd)
	cmd = m.List.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Entered %s", selectedItemName)))
	cmds = append(cmds, cmd)
	return cmds
}

// Donwload a file based on the path provided
func (m *Model) downloadFile(fileItem fs.FileInfo) tea.Cmd {
	return func() tea.Msg {
		var srcFile io.Reader
		srcFile, err := m.SftpClient.Open(m.SftpClient.Join(m.currentDir, fileItem.Name()))
		handleError(err)
		// Instrument with our counter.
		barPercentage := barPercentage(0)
		counter := &writeProgressCounter{
			TotalFileSize: fileItem.Size(),
			percentage:    &barPercentage,
		}
		srcFile = io.TeeReader(srcFile, counter)

		destFile, err := os.Create(filepath.Join(".", fileItem.Name()))
		handleError(err)
		go func() {
			defer destFile.Close()
			_, err = io.Copy(destFile, srcFile)
			handleError(err)
		}()
		return &barPercentage
	}
}

func (m Model) View() string {
	f, err := tea.LogToFile("debug.log", "debug")
	handleError(err)
	f.WriteString(fmt.Sprintf("Percentuale: %f", m.progress.Percent()))
	if m.progress.Percent() != 0 && m.progress.Percent() != 1 {
		return docStyle.Render(
			lipgloss.JoinHorizontal(
				lipgloss.Center,
				m.progress.View(),
				// lipgloss.NewStyle().Render("Banana"),
			),
		)
	} else {
		return docStyle.Render(m.List.View())
	}
}

// Create the list of item by fetching the server
func CreateItemListModel(dirPath string, sftpClient *sftp.Client) []list.Item {
	fileList, err := sftpClient.ReadDir(dirPath)
	handleError(err)

	previousDir := PreviousDir{}
	// Insert the .. dir
	items := []list.Item{
		&item{
			rawValue: &previousDir,
		},
	}

	for _, file := range fileList {
		items = append(items, &item{rawValue: file})
	}
	return items
}
