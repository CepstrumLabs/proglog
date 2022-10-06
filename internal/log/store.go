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
	lenWidth = 8 // the size of the store is uint64
) // which can be represented by 8 bytes
// lenWidth is the number of bytes used to represent
//   the length of the actual data being stored in the store

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
		size: size,
		buf:  bufio.NewWriter(f),
	}, nil

}

// Append adds the sequence of bytes p to the store instance
// it first writes the length of the sequence to add as binary data,
// then it writes the contents of p and keeps track of the bytes written
// the total number of bytes written is the number of bytes used to write
// the length of the data to be written and then the number of bytes that represents
// the data itself.
// It returns the number of bytes written n, the position that they were inserted into
// and an error or nil if no error occurs
func (s *store) Append(p []byte) (n uint64, pos uint64, err error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	pos = s.size
	if err := binary.Write(s.buf, enc, uint64(len(p))); err != nil {
		return 0, 0, err
	}
	w, err := s.buf.Write(p)
	if err != nil {
		return 0, 0, err
	}
	w += lenWidth
	s.size += uint64(w)
	return uint64(w), pos, nil
}

func (s *store) Read(pos uint64) ([]byte, error) {
	s.mu.Lock()
	defer s.mu.Unlock()
	// Flush the buffer in case we're about to read data
	// that is still on the buffer and not on the disk
	if err := s.buf.Flush(); err != nil {
		return nil, err
	}

	size := make([]byte, lenWidth)

	// Read the size of the thing we are about to
	// read so we know how long the record is
	if _, err := s.File.ReadAt(size, int64(pos)); err != nil {
		return nil, err
	}

	// create a list of bytes of size enc.Uint64(size)
	b := make([]byte, enc.Uint64(size))

	// Read len(b) bytes at position pos + lenWidth
	if _, err := s.File.ReadAt(b, int64(pos+lenWidth)); err != nil {
		return nil, err
	}

	return b, nil
}

func (s *store) ReadAt(p []byte, off int64) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return 0, err
	}

	return s.File.ReadAt(p, off)
}

func (s *store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if err := s.buf.Flush(); err != nil {
		return err
	}

	return s.File.Close()
}
