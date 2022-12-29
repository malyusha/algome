package structproviders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CollectInfoFromSourceDirectory(t *testing.T) {
	expectOut := map[string][]string{
		"testdata/dir/a": {"one.file", "two.file"},
		"testdata/dir/b": {"one.file"},
	}
	gotOut, err := collectFilesOfDirectory("testdata/dir")
	assert.NoError(t, err)
	assert.Equal(t, expectOut, gotOut)
}
