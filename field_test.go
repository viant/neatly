package neatly_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/viant/neatly"
	"github.com/viant/toolbox/data"
	"testing"
)

func TestFieldExpression_Set(t *testing.T) {

	{
		var object = data.NewMap()
		field1 := neatly.NewField("Field1")
		field1.Set(123, object)
		assert.Equal(t, 123, object.GetInt("Field1"))

	}

	{
		var object = data.NewMap()
		field1 := neatly.NewField("Req.[]Array.H")
		assert.True(t, field1.HasArrayComponent)
		assert.False(t, field1.IsArray)
		assert.True(t, field1.Child.IsArray)
		assert.Equal(t, "H", field1.Child.Child.Field)

		field1.Set("v1H", object)
		field1.Set("v2H", object, 1)

		field2 := neatly.NewField("Req.[]Array.A")
		field2.Set("v1A", object)
		field2.Set("v2A", object, 1)

		field3 := neatly.NewField("Req.Field")

		field3.Set("v", object)

		assert.True(t, object.Has("Req"))
		var reqObject = object.GetMap("Req")
		assert.NotNil(t, reqObject)

		assert.Equal(t, "v", reqObject.GetString("Field"))
		assert.True(t, reqObject.Has("Array"))
		assert.True(t, reqObject.Has("Field"))

		array := reqObject.GetCollection("Array")
		assert.NotNil(t, array)
		assert.Equal(t, 2, len(*array))

		err := array.RangeMap(func(item data.Map, index int) (bool, error) {
			switch index {

			case 0:
				assert.Equal(t, "v1H", item.GetString("H"))
				assert.Equal(t, "v1A", item.GetString("A"))

			case 1:
				assert.Equal(t, "v2H", item.GetString("H"))
				assert.Equal(t, "v2A", item.GetString("A"))

			}
			return true, nil
		})
		assert.Nil(t, err)

	}

	{
		var object = data.NewMap()
		field1 := neatly.NewField("/Field1")
		field1.Set(123, object)
		assert.Equal(t, 123, object.GetInt("Field1"))
		assert.True(t, field1.IsRoot)
	}

	{
		var object = data.NewMap()
		field1 := neatly.NewField(":Field1")
		field1.Set(123, object)
		assert.Equal(t, 123, object.GetInt("Field1"))
		assert.True(t, field1.IsVirtual)
	}

}
