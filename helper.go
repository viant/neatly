package neatly

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
)

func asDataStructure(value string) (interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}
	if strings.HasPrefix(value, "{{") || strings.HasSuffix(value, "}}") {
		return string(value[1 : len(value)-1]), nil
	}

	if strings.HasPrefix(value, "[[") || strings.HasSuffix(value, "]]") {
		return string(value[1 : len(value)-1]), nil
	}

	if strings.HasPrefix(value, "{") {
		if toolbox.IsNewLineDelimitedJSON(value) {
			return value, nil
		}
		jsonFactory := toolbox.NewJSONDecoderFactory()
		var jsonObject = make(map[string]interface{})
		err := jsonFactory.Create(strings.NewReader(value)).Decode(&jsonObject)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %v %T, %v", value, value, err)
		}
		return jsonObject, nil
	} else if strings.HasPrefix(value, "[") {
		var jsonArray = make([]interface{}, 0)
		err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(value)).Decode(&jsonArray)
		if err != nil {
			return nil, fmt.Errorf("failed to decode: %v %v", value, err)
		}
		return jsonArray, nil
	}
	return value, nil
}
