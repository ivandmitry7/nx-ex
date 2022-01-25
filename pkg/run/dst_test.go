package run

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io/ioutil"
	"os"
	"path"
	"testing"
)

func TestFileWriter_OpenClose(t *testing.T) {
	dst := NewFileWriter()
	err := dst.Open("QQW://Invalid+DirName", "fileName.wav")
	require.Error(t, err)
	assert.Equal(t, err.Error()[:26], string("failed to create directory")[:26])

	err = dst.Open(os.TempDir(), "dmd-test.wav")
	defer func() {
		dst.Close()
		_ = os.Remove(path.Join(os.TempDir(), "dmd-test.wav.log"))
	}()
	assert.NoError(t, err)
}

func TestFileWriter_SaveResults(t *testing.T) {
	dst := NewFileWriter()
	err := dst.Open(os.TempDir(), "dmd-test.wav")
	defer func() {
		dst.Close()
		_ = os.Remove(path.Join(os.TempDir(), "dmd-test.wav.log"))
		_ = os.Remove(path.Join(os.TempDir(), "dmd-test.wav.ale"))
		_ = os.Remove(path.Join(os.TempDir(), "dmd-test.wav.txt"))
		_ = os.Remove(path.Join(os.TempDir(), "dmd-test.wav.bin"))
	}()
	require.NoError(t, err)

	r := Results{
		HasData: true,
		Freq:    1800,
		Text:    "This is text result",
		Bits:    "0000010110010100",
		Ale:     "This is ALE result",
		Log:     "This is log result",
		Phase:   nil,
		Params:  nil,
	}

	err = dst.SaveResults(r)
	require.NoError(t, err)

	b, err := ioutil.ReadFile(path.Join(os.TempDir(), "dmd-test.wav.log"))
	require.NoError(t, err)
	assert.Equal(t, string(b), "This is log result\n")

	b, err = ioutil.ReadFile(path.Join(os.TempDir(), "dmd-test.wav.txt"))
	require.NoError(t, err)
	assert.Equal(t, string(b), "This is text result")

	b, err = ioutil.ReadFile(path.Join(os.TempDir(), "dmd-test.wav.ale"))
	require.NoError(t, err)
	assert.Equal(t, string(b), "This is ALE result")

	b, err = ioutil.ReadFile(path.Join(os.TempDir(), "dmd-test.wav.bin"))
	require.NoError(t, err)
	assert.Equal(t, b, []byte{0x05, 0x94})
}
