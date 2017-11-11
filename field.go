package neatly

import (
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"strings"
	"unicode"
)

//Field represent field of object
type Field struct {
	expression        string //expression Used to create this field
	Field             string //actual expression field
	Child             *Field //child expression if expression contains .
	IsArray           bool   //is this filed an array type
	HasSubPath        bool   //flag indicating if this field has sub fields
	HasArrayComponent bool   //flag indicating if this or child fileds have an array component
	IsRoot            bool   //flag indicating if this field is Root
	IsVirtual         bool   //flag indicating if this field belong to virtual object
}

//Set sets value into target map, if indexes are provided value will be pushed into a slice
func (f *Field) Set(value interface{}, target data.Map, indexes ...int) {
	var index = 0
	if !target.Has(f.Field) {
		if f.IsArray {
			target.Put(f.Field, data.NewCollection())
		} else if f.HasSubPath {
			target.Put(f.Field, data.NewMap())
		}
	}
	var aMap data.Map
	var action func(data data.Map, indexes ...int)
	if !f.HasSubPath {
		if f.IsArray {
			action = func(data data.Map, indexes ...int) {
				collection := target.GetCollection(f.Field)
				(*collection)[index] = value
			}
		} else {
			action = func(data data.Map, indexes ...int) {
				var isValueSet = false
				if data.Has(f.Field) {
					existingValue := data.Get(f.Field)
					if toolbox.IsMap(existingValue) && toolbox.IsMap(value) { //is existing Value is a aMap.Map and existing Value is aMap.Map
						// then add keys to existing aMap.Map
						existingMap := data.GetMap(f.Field)
						existingMap.Apply(toolbox.AsMap(value))
						isValueSet = true
					} else if toolbox.IsSlice(existingValue) { //is existing value is a slice append elements
						existingSlice := data.GetCollection(f.Field)
						if toolbox.IsSlice(value) {
							for _, item := range toolbox.AsSlice(value) {
								existingSlice.Push(item)
							}
						} else {
							existingSlice.Push(value)
						}
						data.Put(f.Field, existingSlice)
						isValueSet = true
					}
				}
				if !isValueSet {
					data.Put(f.Field, value)
				}
			}
		}

	} else {
		action = func(data data.Map, indexes ...int) {
			f.Child.Set(value, data, indexes...)
		}
	}

	if f.IsArray {
		index, indexes = shiftIndex(indexes...)
		collection := target.GetCollection(f.Field)
		collection.PadWithMap(index + 1)
		aMap, _ = (*collection)[index].(data.Map)

	} else if f.HasSubPath {
		aMap = target.GetMap(f.Field)
	} else {
		aMap = target
	}
	action(aMap, indexes...)
}

//NewField return a new Field for provided expression.
func NewField(expression string) *Field {
	var parsedExpression = expression
	isRoot := strings.HasPrefix(parsedExpression, "/")
	if isRoot {
		parsedExpression = string(parsedExpression[1:])
	}
	isVirtual := strings.HasPrefix(parsedExpression, ":")
	if isVirtual {
		parsedExpression = string(parsedExpression[1:])
	}
	runes := []rune(expression)
	if unicode.IsLower(runes[0]) {
		isVirtual = true
	}
	var result = &Field{
		expression:        expression,
		HasArrayComponent: strings.Contains(parsedExpression, "[]"),
		IsArray:           strings.HasPrefix(parsedExpression, "[]"),
		HasSubPath:        strings.Contains(parsedExpression, "."),
		Field:             parsedExpression,
		IsRoot:            isRoot,
		IsVirtual:         isVirtual,
	}

	if result.HasSubPath {
		dotPosition := strings.Index(parsedExpression, ".")
		result.Field = string(result.Field[:dotPosition])
		result.Child = NewField(string(parsedExpression[dotPosition+1:]))
	}
	if result.IsArray {
		result.Field = string(result.Field[2:])
	}
	return result
}

func shiftIndex(indexes ...int) (int, []int) {
	var index int
	if len(indexes) > 0 {
		index = indexes[0]
		indexes = indexes[1:]
	}
	return index, indexes
}
