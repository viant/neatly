package neatly

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/toolbox"
	"testing"
)

func Test_asDataStructure(t *testing.T) {

	{
		input := `[1,2,3]`
		output, err := asDataStructure(input)
		assert.Nil(t, err)
		assert.EqualValues(t, []interface{}{float64(1), float64(2), float64(3)}, output)
	}
	{
		input := `{"a":1, "b":2}`
		output, err := asDataStructure(input)
		if assert.Nil(t, err) {
			outputMap := toolbox.AsMap(output)
			assert.EqualValues(t, map[string]interface{}{
				"a": float64(1),
				"b": float64(2),
			}, outputMap)
		}
	}
	{
		input := `{"a":1, "b":2}
{"a2":2, "b3":21}
{"a3":3, "b4:22}
`
		output, err := asDataStructure(input)
		if assert.Nil(t, err) {
			outputMap := toolbox.AsString(output)
			assert.EqualValues(t, input, outputMap)
		}
	}
}
