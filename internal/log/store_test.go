package log

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

var (
	write = []byte("hello world")
	width = uint64(len(write) + lenWidth)
)

func TestStoreAppendRead(t *testing.T) {
	f, err := ioutil.TempFile("", "store_append_read_test")
	require.NoError(t, err)

	defer os.Remove(f.Name())

	s, err := newStore(f)

	require.NoError(t, err)

	testAppend(t, s)
	testRead(t, s)
	testReadAt(t, s)

	s, err = newStore(f)

	require.NoError(t, err)
	testRead(t, s)
}

func testAppend(t *testing.T, s *store) {
	/* Given a store instance
	   When I append a sequence s made of of n bytes
	   Then I expect its size to be n + lenWidth=8

	   lenWidth is the number of bytes used to represent
	   the length of the actual data being stored in the store
	(*/
	t.Helper()
	for i := uint64(1); i <= 4; i++ {
		n, pos, err := s.Append(write)
		require.NoError(t, err)
		require.Equal(t, pos+n, width*i)
	}

}

func testRead(t *testing.T, s *store) {
	/* Given a store instance
	   When I append a sequence of bytes to it
	   And then I read from it
	   Then I read back the same sequence of bytes that I
	   initially appended to it
	*/
	t.Helper()
	var pos uint64
	for i := uint64(1); i <= 4; i++ {
		read, err := s.Read(pos)
		require.NoError(t, err)
		require.Equal(t, read, write)
		pos += width
	}
}

func testReadAt(t *testing.T, s *store) {
	/* Given a prepopulated store instance
	   When I iteratively read it using the ReadAt function
	   Then I correctly read the contents I wrote to it
	*/
	t.Helper()
	for i, off := uint64(1), int64(0); i <= 4; i++ {
		b := make([]byte, lenWidth)
		n, err := s.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, lenWidth, n)
		off += int64(n)

		size := enc.Uint64(b)
		b = make([]byte, size)
		n, err = s.ReadAt(b, off)
		require.NoError(t, err)
		require.Equal(t, write, b)
		require.Equal(t, int(size), n)
		off += int64(n)
	}
}

func TestStoreClose(t *testing.T) {
	/* Given I write the `write` variable to the store
	   When I read the size of the store before and after I close the store
	   Then the size after is bigger than that before (because of buffer flushing)
	*/
	f, err := ioutil.TempFile("", "store_close_test")
	require.NoError(t, err)

	defer os.Remove(f.Name())
	s, err := newStore(f)
	require.NoError(t, err)

	_, _, err = s.Append(write)
	require.NoError(t, err)

	f, beforeSize, err := openFile(f.Name())
	require.NoError(t, err)

	err = s.Close()
	require.NoError(t, err)

	f, afterSize, err := openFile(f.Name())
	require.NoError(t, err)
	require.True(t, beforeSize < afterSize)
}

func openFile(name string) (*os.File, int64, error) {
	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0644)

	if err != nil {
		return nil, 0, err
	}

	fi, err := f.Stat()
	if err != nil {
		return nil, 0, err
	}
	return f, fi.Size(), err
}
