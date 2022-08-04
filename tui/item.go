package tui

import "io/fs"

// Rapresents an a file as an item of the list of the tui client 
type item struct {
	rawValue fs.FileInfo // File properties
}

// Get the stiled title for the file item
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

// Get fancy description for the file item
func (i item) Description() string {
	if i.rawValue.Name() == ".." {
		return ""
	}
	return getFileDescription(i.rawValue)
}

// The value to filter when searching
func (i item) FilterValue() string { return i.rawValue.Name() }
