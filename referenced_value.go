package neatly

import (
	"fmt"
	"github.com/viant/toolbox/data"
	"strings"
)

type referenceValues map[string]*referenceValue

func (v *referenceValues) CheckUnused() error {
	var unused = make([]string, 0)
	for k, value := range *v {
		if !value.Used {
			unused = append(unused, k)
		}
	}
	if len(unused) == 0 {
		return nil
	}
	return fmt.Errorf("Unresolved references: '%v' ", strings.Join(unused, ","))
}

func (v *referenceValues) Add(tagName string, field *Field, object data.Map) error {
	var referencedValue = &referenceValue{
		Key:    tagName,
		Field:  field,
		Object: object,
	}
	referencedValue.Setter = func(value interface{}) {
		referencedValue.Used = true
		field.Set(value, referencedValue.Object)
	}
	(*v)[tagName] = referencedValue
	return nil
}

func (v *referenceValues) Apply(tagName string, value interface{}) error {
	referencedValue, ok := (*v)[tagName]
	if !ok {

		var referencesSoFar = make([]string, 0)
		for k := range *v {
			referencesSoFar = append(referencesSoFar, k)
		}

		return fmt.Errorf("Missing referenceValue %v in the previous rows, available[%v]", tagName, strings.Join(referencesSoFar, ","))
	}
	referencedValue.Setter(value)
	return nil
}

func newReferenceValues() referenceValues {
	var result referenceValues = make(map[string]*referenceValue)
	return result
}

//referenceValue represent reference value
type referenceValue struct {
	Setter func(value interface{}) //setter handler
	Key    string                  //reference key
	Field  *Field                  //field
	Object data.Map                //target object
	Used   bool                    //flag indicating if reference was used, if not then error
}
