package neatly_test

import (
	"github.com/viant/assertly"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/viant/neatly"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/url"
)

func Test_Md5(t *testing.T) {

	var md5, err = neatly.Md5("554257_popularmechanics.com", nil)
	assert.Nil(t, err)
	assert.EqualValues(t, "ed045d398e8e1924486afa44acbb6b82", md5)

	aMap := data.NewMap()
	aMap.Put("md5", neatly.Md5)

	var text = "11$md5(554257_popularmechanics.com)22"
	expanded := aMap.ExpandAsText(text)
	assert.EqualValues(t, "11ed045d398e8e1924486afa44acbb6b8222", expanded)

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

	state.Put(neatly.OwnerURL, url.NewResource("test/usecase7/001/customer.json").URL)
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

func Test_Keys(t *testing.T) {

	{
		var keys, err = neatly.Keys(map[string]interface{}{}, nil)
		assert.Nil(t, err)
		assert.NotNil(t, keys)
	}
	{
		var keys, err = neatly.Keys("{\"abc\":1}", nil)
		assert.Nil(t, err)
		assert.EqualValues(t, []interface{}{"abc"}, keys)
	}
}

func Test_Values(t *testing.T) {

	{
		var keys, err = neatly.Values(map[string]interface{}{}, nil)
		assert.Nil(t, err)
		assert.NotNil(t, keys)
	}
	{
		var keys, err = neatly.Values("{\"abc\":1}", nil)
		assert.Nil(t, err)
		assert.EqualValues(t, []interface{}{1.0}, keys)
	}
}

func Test_AsCollection(t *testing.T) {

	{
		var aSlice, err = neatly.AsCollection([]interface{}{1}, nil)
		assert.Nil(t, err)
		assert.NotNil(t, aSlice)
	}
	{
		var aSlice, err = neatly.AsCollection("[1,2]", nil)
		assert.Nil(t, err)
		assert.NotNil(t, aSlice)
		assertly.AssertValues(t, []interface{}{1, 2}, aSlice)
	}

	{
		_, err := neatly.AsCollection("[\"a,2]", nil)
		assert.NotNil(t, err)
	}
}

func Test_AsData(t *testing.T) {
	{
		var aSlice, err = neatly.AsData("[1,2]", nil)
		assert.Nil(t, err)
		assert.NotNil(t, aSlice)
		assertly.AssertValues(t, []interface{}{1, 2}, aSlice)
	}
	{
		var aMap, err = neatly.AsData("{\"abc\":1}", nil)
		assert.Nil(t, err)
		assert.NotNil(t, aMap)
		assertly.AssertValues(t, "{\"abc\":1}", aMap)
	}

}

func Test_Join(t *testing.T) {
	{
		var joined, err = neatly.Join([]interface{}{
			[]interface{}{1, 2, 3},
			",",
		}, nil)
		assert.Nil(t, err)
		assert.NotNil(t, joined)
		assert.EqualValues(t, "1,2,3", joined)
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

func Test_Length(t *testing.T) {
	{
		value, err := neatly.Length(4.3, nil)
		assert.Nil(t, err)
		assert.Equal(t, 0, value)
	}
	{
		value, err := neatly.Length("abcd", nil)
		assert.Nil(t, err)
		assert.Equal(t, 4, value)
	}
	{
		value, err := neatly.Length(map[int]int{
			2: 3,
			1: 1,
			6: 3,
		}, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, value)
	}
	{
		value, err := neatly.Length([]int{1, 2, 3}, nil)
		assert.Nil(t, err)
		assert.Equal(t, 3, value)
	}
}

func Test_FormatTime(t *testing.T) {

	{
		value, err := neatly.FormatTime([]interface{}{"now", "yyyy"}, nil)
		assert.Nil(t, err)
		now := time.Now()
		assert.Equal(t, now.Year(), toolbox.AsInt(value))
	}
	{
		value, err := neatly.FormatTime([]interface{}{"2015-02-11", "yyyy"}, nil)
		assert.Nil(t, err)
		assert.Equal(t, 2015, toolbox.AsInt(value))
	}
	{
		_, err := neatly.FormatTime([]interface{}{"2015-02-11"}, nil)
		assert.NotNil(t, err)
	}
	{
		_, err := neatly.FormatTime([]interface{}{"201/02/11 2", "y-d"}, nil)
		assert.NotNil(t, err)
	}
	{
		_, err := neatly.FormatTime("a", nil)
		assert.NotNil(t, err)
	}

	{
		value, err := neatly.FormatTime([]interface{}{"now", "yyyy", "UTC"}, nil)
		assert.Nil(t, err)
		now := time.Now()
		assert.Equal(t, now.Year(), toolbox.AsInt(value))
	}

}

func Test_Zip_Unzip(t *testing.T) {
	{
		compressed, err := neatly.Zip("abc", nil)
		assert.Nil(t, err)

		{
			origin, err := neatly.Unzip(compressed, nil)
			assert.Nil(t, err)
			assert.Equal(t, "abc", toolbox.AsString(origin))
		}
		{
			origin, err := neatly.UnzipText(compressed, nil)
			assert.Nil(t, err)
			assert.Equal(t, "abc", origin)
		}
	}

	{
		compressed, err := neatly.Zip([]byte("abc"), nil)
		assert.Nil(t, err)
		origin, err := neatly.Unzip(compressed, nil)
		assert.Nil(t, err)
		assert.Equal(t, "abc", toolbox.AsString(origin))
	}

	{ //Error case
		_, err := neatly.Zip(1, nil)
		assert.NotNil(t, err)
		_, err = neatly.Unzip(1, nil)
		assert.NotNil(t, err)
		_, err = neatly.Unzip([]byte{}, nil)
		assert.NotNil(t, err)

	}

}

func Test_Markdown(t *testing.T) {
	{
		html, err := neatly.Markdown("*Hello*", nil)
		assert.Nil(t, err)
		assert.EqualValues(t, "<p><em>Hello</em></p>\n", html)
	}

}

func Test_Cat(t *testing.T) {
	{
		content, err := neatly.Cat("udf.go", nil)
		assert.Nil(t, err)
		assert.True(t, content != "")
	}
	{
		_, err := neatly.Cat("uaaadf.go", nil)
		assert.NotNil(t, err)
	}
	{
		content, err := neatly.Cat("test/../udf.go", nil)
		assert.Nil(t, err)
		assert.True(t, content != "")
	}
}

func Test_IsSON(t *testing.T) {
	//Table driven tests
	useCases := []struct {
		name     string
		input    string
		expected bool
	}{
		{"Case to check if invalid JSON is validated", "test/invalid_file.json", false},
		{"Case to check if valid JSON is validated", "test/valid_json_file.json", true},
	}

	for _, useCase := range useCases {
		t.Run(useCase.name, func(t *testing.T) {
			isJson, _ := neatly.IsJSON(useCase.input, nil)
			assert.Equal(t, useCase.expected, toolbox.AsBoolean(isJson))

		})
	}
}
