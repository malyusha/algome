package generator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_buildLangsMap(t *testing.T) {
	defs := []string{"Go|go", "Test|t|test|another"}
	expectOut := map[string]string{
		"go":      "Go",
		"t":       "Test",
		"test":    "Test",
		"another": "Test",
	}
	assert.Equal(t, expectOut, buildLangsMap(defs))
}
