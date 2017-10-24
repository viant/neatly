package neatly_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/neatly"
	"testing"
)

func Test_Tag(t *testing.T) {
	var tag = neatly.NewTag("[]Test{1 .. 003}", 1)
	assert.True(t, tag.IsArray)
	assert.Equal(t, "Test", tag.Name)
	assert.Equal(t, 1, tag.Iterator.Min)
	assert.Equal(t, 3, tag.Iterator.Max)
	assert.Equal(t, "%03d", tag.Iterator.Template)
	assert.Equal(t, 1, tag.LineNumber)
}
