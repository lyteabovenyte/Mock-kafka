package log

import (
	"bufio"
	"encoding/binary"
	"os"
	"sync"
)

var (
	enc = binary.BigEndian
)

const (
	lenWidth = 8 // a placeholder to determine how many bytes is this specific log.
)

// store struct is a simple wrapper around a file
// with two simple APIs to append and read bytes
// to and from the file.
type store struct {
	*os.File
	mu   sync.Mutex
	buf  *bufio.Writer
	size uint64
}

func newStore(f *os.File) (*store, error) {
	fi, err := os.Stat(f.Name())
	if err != nil {
		return nil, err
	}

	size := uint64(fi.Size())
	return &store{
		File: f,
		buf:  bufio.NewWriter(f),
		size: size,
	}, nil
}

// Append persists the given bytes to the store.
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	pos = s.size
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, nil
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, nil
	}
	w += lenWidth
	s.size = uint64(w)
	return uint64(w), pos, nil // returns how many bytes we just wrote, the position, and nil error
}

// Read returns the record stored at the given position.
func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	size := make([]byte, lenWidth)
	_, err := s.File.ReadAt(size, int64(pos))
	if err != nil {
		return nil, err
	}
	b := make([]byte, enc.Uint64(size))
	_, err = s.File.ReadAt(b, int64(pos+lenWidth)) // the actual log is at pos+lenWidth position.
	if err != nil {
		return nil, err
	}
	return b, nil
}

// ReadAt implements io.ReaderAt on the store type.
// ReadAt reads len(p) bytes into p beginning at the off offset.
func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if err := s.buf.Flush(); err != nil {
		return 0, err
	}
	return s.File.ReadAt(p, off)
}

// Close persist any buffered data before closing the file.
func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	err := s.buf.Flush()
	if err != nil {
		return err
	}
	return s.File.Close()
}
