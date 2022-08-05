package tui

import "fmt"

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	BytesWritten  int64 // Total # of bytes written
	TotalFileSize int64
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.BytesWritten += int64(n)
	percentage := 100 * wc.BytesWritten / wc.TotalFileSize
	fmt.Printf("Read %d bytes for a total of %d\n", percentage, 100)
	return n, nil
}
