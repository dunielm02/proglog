package log

import (
	"io"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestIndex(t *testing.T) {
	f, err := os.CreateTemp("", "index_test")
	require.NoError(t, err)
	defer os.Remove(f.Name())

	c := Config{}
	c.Segment.MaxIndexBytes = 37

	index, err := newIndex(f, c)
	require.NoError(t, err)

	_, _, err = index.Read(-1)
	require.Error(t, err)
	require.Equal(t, f.Name(), index.Name())

	entries := []struct {
		Off uint32
		Pos uint64
	}{
		{Off: 0, Pos: 0},
		{Off: 1, Pos: 10},
		{Off: 2, Pos: 100},
	}

	for i, v := range entries {
		err := index.Write(v.Off, v.Pos)
		require.NoError(t, err)
		off, pos, err := index.Read(int64(i))
		require.NoError(t, err)
		require.Equal(t, pos, v.Pos)
		require.Equal(t, off, v.Off)
	}

	_, _, err = index.Read(int64(len(entries)))
	require.Equal(t, io.EOF, err)
	err = index.Close()
	require.NoError(t, err)

	f, _ = os.OpenFile(f.Name(), os.O_RDWR, 0600)
	index, err = newIndex(f, c)
	require.NoError(t, err)
	off, pos, err := index.Read(-1)
	require.NoError(t, err)
	require.Equal(t, uint32(2), off)
	require.Equal(t, entries[2].Pos, pos)
}
