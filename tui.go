package main

import (
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"time"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/knipferrc/teacup/icons"
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

type item struct {
	rawValue fs.FileInfo
}

func (i item) Title() string {
	if i.rawValue.Name() == ".." {
		return ".."
	}

	var title string
	if i.rawValue.IsDir() {
		title = dirItemStyle(i.rawValue.Name())
	} else {
		title = fileItemStyle(i.rawValue.Name())
	}
	return getFileIcon(i.rawValue) + " " + title
}
func (i item) Description() string {
	if i.rawValue.Name() == ".." {
		return ""
	}
	return getFileDescription(i.rawValue)
}
func (i item) FilterValue() string { return i.rawValue.Name() }

type model struct {
	list       list.Model
	sftpClient *sftp.Client
	currentDir string
}

func (m model) Init() tea.Cmd {
	return nil
}

func (m model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
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
			selectedItem := m.list.SelectedItem().(*item).rawValue

			//if it's nil then it is a ".." dir
			if selectedItem == nil {
				cmds = moveDir(&m, "..", cmds)
			} else {
				selectedItemName := selectedItem.Name()
				if selectedItem.IsDir() {
					cmds = moveDir(&m, selectedItemName, cmds)
				} else {
					cmd = m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Downloading %s", selectedItemName)))
					cmds = append(cmds, cmd)
					cmds = append(cmds, m.list.ToggleSpinner())
					err := m.downloadFile(m.currentDir, selectedItemName)
					handleError(err)
				}
			}

			cmds = append(cmds, cmd)

			return m, tea.Batch(cmds...)
		}

	case tea.WindowSizeMsg:
		h, v := docStyle.GetFrameSize()
		m.list.SetSize(msg.Width-h, msg.Height-v)

	}

	var cmd tea.Cmd
	m.list, cmd = m.list.Update(msg)
	return m, cmd
}

func moveDir(m *model, selectedItemName string, cmds []tea.Cmd) []tea.Cmd {
	currentWd, err := m.sftpClient.RealPath(m.sftpClient.Join(m.currentDir, selectedItemName))
	handleError(err)
	m.currentDir = currentWd

	cmd := m.list.SetItems(createItemListModel(currentWd, sftpClient))
	cmds = append(cmds, cmd)
	cmd = m.list.NewStatusMessage(statusMessageStyle(fmt.Sprintf("Entered %s", selectedItemName)))
	cmds = append(cmds, cmd)
	return cmds
}

func (m model) downloadFile(filePath, fileName string) error {
	srcFile, err := m.sftpClient.Open(m.sftpClient.Join(filePath, fileName))
	handleError(err)
	defer srcFile.Close()
	destFile, err := os.Create(filepath.Join(".", fileName))
	defer destFile.Close()
	handleError(err)
	_, err = io.Copy(destFile, srcFile)
	return err
}

func (m model) View() string {
	return docStyle.Render(
		lipgloss.JoinHorizontal(
			lipgloss.Top,
			m.list.View(),
		),
	)
}

func createItemListModel(dirPath string, sftpClient *sftp.Client) []list.Item {
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

func getFileDescription(value fs.FileInfo) string {
	status := fmt.Sprintf("%s %s %s",
		value.ModTime().Format("2006-01-02 15:04:05"),
		value.Mode().String(),
		ConvertBytesToSizeString(value.Size()))
	return status
}

func getFileIcon(value fs.FileInfo) string {
	icon, _ := icons.GetIcon(
		value.Name(),
		filepath.Ext(value.Name()),
		icons.GetIndicator(value.Mode()),
	)
	return icon
}

// ConvertBytesToSizeString converts a byte count to a human readable string.
func ConvertBytesToSizeString(size int64) string {
	const (
		thousand    = 1000
		ten         = 10
		fivePercent = 0.0499
	)

	if size < thousand {
		return fmt.Sprintf("%dB", size)
	}

	suffix := []string{
		"K", // kilo
		"M", // mega
		"G", // giga
		"T", // tera
		"P", // peta
		"E", // exa
		"Z", // zeta
		"Y", // yotta
	}

	curr := float64(size) / thousand
	for _, s := range suffix {
		if curr < ten {
			return fmt.Sprintf("%.1f%s", curr-fivePercent, s)
		} else if curr < thousand {
			return fmt.Sprintf("%d%s", int(curr), s)
		}
		curr /= thousand
	}

	return ""
}

type PreviousDir struct{}

func (p *PreviousDir) IsDir() bool        { return true }
func (p *PreviousDir) Name() string       { return ".." }
func (p *PreviousDir) Size() int64        { return 0 }
func (p *PreviousDir) Mode() fs.FileMode  { return os.FileMode(0) }
func (p *PreviousDir) ModTime() time.Time { return time.Time{} }
func (p *PreviousDir) Sys() any           { return nil }
