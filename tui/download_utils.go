package tui

// WriteCounter counts the number of bytes written to it.
type WriteCounter struct {
	BytesWritten  int64 // Total # of bytes written
	TotalFileSize int64
	percentage    *barPercentage
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *WriteCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.BytesWritten += int64(n)
	*wc.percentage = barPercentage(100 * wc.BytesWritten / wc.TotalFileSize)

	return n, nil
}
