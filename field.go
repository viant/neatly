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
	HasArrayComponent bool   //flag indicating if this or child fieds have an array component
	IsRoot            bool   //flag indicating if this field is Root
	IsVirtual         bool   //flag indicating if this field belong to virtual object
	IsIndex           bool   //flag indicating if this filed is actual array index, as opposed to sub field name
	Leaf              *Field //leaf field
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
			action = func(object data.Map, indexes ...int) {
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
			if f.Child.IsIndex {
				collection := target.GetCollection(f.Field)
				var index = toolbox.AsInt(f.Child.Field)
				for i := len(*collection); i < index+1; i++ {
					*collection = append(*collection, nil)
				}
				(*collection)[index] = value
				return

			}
			f.Child.Set(value, data, indexes...)
		}
	}

	if f.IsArray {
		if collectionPointer, ok := value.(*data.Collection); ok {
			target.Put(f.Field, collectionPointer)
			return
		}

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

//ArrayPath returns a field array path
func (f *Field) ArrayPath() string {
	if !f.HasArrayComponent {
		return ""
	}
	var result = make([]string, 0)
	field := f
	for {
		result = append(result, field.Field)
		if field.IsArray || !field.HasSubPath {
			break
		}

		field = field.Child
	}
	return strings.Join(result, ".")
}

func (f *Field) GetArraySize(value data.Map) int {
	if !f.HasArrayComponent {
		return 0
	}
	field := f
	for {
		subValue := value.Get(field.Field)
		if subValue == nil {
			return 0
		}
		if field.IsArray {
			return len(toolbox.AsSlice(subValue))
		}
		if !field.HasSubPath {
			break
		}

		value = data.Map(toolbox.AsMap(subValue))
		field = field.Child
	}
	return 0
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

	isArray := strings.HasPrefix(parsedExpression, "[]")
	if isArray {
		parsedExpression = string(parsedExpression[2:])
	}
	runes := []rune(parsedExpression)
	if unicode.IsLower(runes[0]) {
		isVirtual = true
	}

	var result = &Field{
		expression:        expression,
		HasArrayComponent: isArray || strings.Contains(parsedExpression, "[]"),
		IsArray:           isArray,
		HasSubPath:        strings.Contains(parsedExpression, "."),
		Field:             parsedExpression,
		IsRoot:            isRoot,
		IsVirtual:         isVirtual,
	}

	if result.HasSubPath {
		dotPosition := strings.Index(parsedExpression, ".")
		result.Field = string(result.Field[:dotPosition])
		result.Child = NewField(string(parsedExpression[dotPosition+1:]))
		if result.IsArray {
			_, err := toolbox.ToInt(result.Child.Field)
			if err == nil {
				result.Child.IsIndex = true
			}
		}
		result.Leaf = result.Child.Leaf

	} else {
		result.Leaf = result
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
