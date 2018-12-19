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
