package generator

import (
	"io/ioutil"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFileWriter_Write(t *testing.T) {
	t.Run("creates new file", func(t *testing.T) {
		w := &fileWriter{
			filepath: "testdata/readme/readme.md",
		}
		defer cleanupFile("testdata/readme")

		contents := []byte("test")
		n, err := w.Write(contents)
		assert.NoError(t, err)

		// check if file created
		finfo, err := os.Stat("testdata/readme/readme.md")
		assert.NoError(t, err)
		assert.Equal(t, int64(n), finfo.Size())
	})

	t.Run("appends to an existing file", func(t *testing.T) {
		w := &fileWriter{
			filepath: "testdata/exists/readme.md",
		}
		initialFinfo, err := ioutil.ReadFile("testdata/exists/readme.md")
		if err != nil {
			t.Fatalf("unexpected error reading file stat: %s", err.Error())
		}

		defer func() {
			f, err := os.OpenFile("testdata/exists/readme.md", os.O_RDWR|os.O_TRUNC, 0666)
			if err != nil {
				panic(err)
			}
			if _, err := f.Write(initialFinfo); err != nil {
				panic(err)
			}
			f.Close()
		}()

		contents := []byte("\ntest")
		n, err := w.Write(contents)
		assert.NoError(t, err)

		// check if file contents appended, not replaced
		updatedFinfo, err := os.Stat("testdata/exists/readme.md")
		assert.NoError(t, err)
		assert.Equal(t, int64(len(initialFinfo))+int64(n), updatedFinfo.Size())

	})
}

func cleanupFile(path string) {
	if err := os.RemoveAll(path); err != nil {
		panic(err)
	}
}
