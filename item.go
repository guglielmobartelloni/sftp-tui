package main

import "io/fs"

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
