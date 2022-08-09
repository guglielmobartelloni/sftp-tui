package tui

// writeProgressCounter counts the number of bytes written to it.
type writeProgressCounter struct {
	BytesWritten  int64 // Total # of bytes written
	TotalFileSize int64 // Total file size
	percentage    *barPercentage // Write percentage calculated
}

// Write implements the io.Writer interface.
//
// Always completes and never returns an error.
func (wc *writeProgressCounter) Write(p []byte) (int, error) {
	n := len(p)
	wc.BytesWritten += int64(n)
	*wc.percentage = barPercentage(100 * wc.BytesWritten / wc.TotalFileSize)

	return n, nil
}
