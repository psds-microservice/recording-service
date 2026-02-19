package storage

import (
	"io"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
)

// SessionWriter writes stream chunks to a file per session.
type SessionWriter struct {
	path   string
	file   *os.File
	log    *zap.Logger
	mu     sync.Mutex
	closed bool
}

// NewSessionWriter creates a new writer for the given session file path.
// Caller must ensure the parent directory exists.
func NewSessionWriter(path string, log *zap.Logger) (*SessionWriter, error) {
	if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
		return nil, err
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	return &SessionWriter{path: path, file: f, log: log}, nil
}

// Write appends data to the recording file.
func (w *SessionWriter) Write(data []byte) (int, error) {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return 0, io.ErrClosedPipe
	}
	return w.file.Write(data)
}

// Close flushes and closes the file.
func (w *SessionWriter) Close() error {
	w.mu.Lock()
	defer w.mu.Unlock()
	if w.closed {
		return nil
	}
	w.closed = true
	return w.file.Close()
}
