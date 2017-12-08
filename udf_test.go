package neatly_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/neatly"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/url"
	"strings"
	"testing"
)

func Test_Md5(t *testing.T) {

	//!Md5(www.viewability.com) 8c505168697a000f0946c04e09f2d524
	var md5, err = neatly.Md5("www.viewability.com", nil)
	assert.Nil(t, err)
	assert.EqualValues(t, "123c274fb9a25ddbc77c1634f1e55525", md5)
}

func Test_WorkingDirectory(t *testing.T) {

	var path, err = neatly.WorkingDirectory("../../abc.txt", nil)
	assert.Nil(t, err)
	assert.True(t, strings.HasSuffix(toolbox.AsString(path), "/abc.txt"))
	assert.True(t, !strings.Contains(toolbox.AsString(path), ".."))

}

func Test_HasResource(t *testing.T) {
	var state = data.NewMap()

	{
		has, err := neatly.HasResource("abc", state)
		assert.Nil(t, err)
		assert.False(t, toolbox.AsBoolean(has))
	}
	{
		has, err := neatly.HasResource("udf.go", state)
		assert.Nil(t, err)
		assert.True(t, toolbox.AsBoolean(has))
	}

	state.Put(neatly.OwnerURL, url.NewResource("test/use_case1.csv").URL)
	{
		has, err := neatly.HasResource("abc", state)
		assert.Nil(t, err)
		assert.False(t, toolbox.AsBoolean(has))
	}

	state.Put(neatly.OwnerURL, url.NewResource("test/usecase7/001/a.csv").URL)
	{
		has, err := neatly.HasResource("use_case.txt", state)
		assert.Nil(t, err)
		assert.True(t, toolbox.AsBoolean(has))
	}
}

func Test_AsMap(t *testing.T) {

	{
		var aMap, err = neatly.AsMap(map[string]interface{}{}, nil)
		assert.Nil(t, err)
		assert.NotNil(t, aMap)
	}
	{
		var aMap, err = neatly.AsMap("{\"abc\":1}", nil)
		assert.Nil(t, err)
		assert.NotNil(t, aMap)
	}

	{
		_, err := neatly.AsMap("{\"abc\":1, \"a}", nil)
		assert.NotNil(t, err)
	}
}

func Test_AsBool(t *testing.T) {
	ok, err := neatly.AsBool("true", nil)
	assert.Nil(t, err)
	assert.Equal(t, true, ok)
}

func Test_AsFloat(t *testing.T) {
	value, err := neatly.AsFloat(0.3, nil)
	assert.Nil(t, err)
	assert.Equal(t, 0.3, value)
}

func Test_AsInt(t *testing.T) {
	value, err := neatly.AsInt(4.3, nil)
	assert.Nil(t, err)
	assert.Equal(t, 4, value)
}
