package tui

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/charmbracelet/bubbles/list"
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

type Model struct {
	List       list.Model   // the list of items
	SftpClient *sftp.Client // the sftp client
	currentDir string       // current directory
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
				cmds = append(cmds, m.List.ToggleSpinner())
				err := m.downloadFile(m.currentDir, selectedItem)
				handleError(err)
			}

			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)
		}

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
func (m *Model) downloadFile(filePath string, fileItem fs.FileInfo) error {
	var srcFile io.Reader
	srcFile, err := m.SftpClient.Open(m.SftpClient.Join(filePath, fileItem.Name()))
	handleError(err)
	// Instrument with our counter.
	counter := &WriteCounter{
		TotalFileSize: fileItem.Size(),
	}
	srcFile = io.TeeReader(srcFile, counter)

	destFile, err := os.Create(filepath.Join(".", fileItem.Name()))
	defer destFile.Close()
	handleError(err)
	_, err = io.Copy(destFile, srcFile)
	return err
}

func (m Model) View() string {
	return docStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.List.View(),
		),
	)
}

// Create the list of item by fetching the server
func CreateItemListModel(dirPath string, sftpClient *sftp.Client) []list.Item {
	fileList, err := sftpClient.ReadDir(dirPath)
	handleError(err)

	previousDir := PreviousDir{}
	// Insert the previous dir
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
