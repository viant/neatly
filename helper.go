package neatly

import (
	"fmt"
	"github.com/viant/toolbox"
	"strings"
)

func unescapeSpecialCharacters(input string) (string, bool) {
	result := strings.TrimSpace(input)
	unescaped := false;
	if len(input) < 2 {
		return input, unescaped
	}
	firstSequence := string(result[:2])
	switch firstSequence {
	case "@@":
		fallthrough
	case "$$":
		fallthrough
	case "##":
		fallthrough
	case "[[":
		fallthrough
	case "{{":
		result = string(result[1:])
		unescaped = true
	}

	lastSequence := string(result[len(result)-2:])
	switch lastSequence {
	case "]]":
		fallthrough
	case "}}":
		result = string(result[:len(result)-1])
		unescaped = true
	}
	return result, unescaped
}

func asDataStructure(value string) (interface{}, error) {
	if len(value) == 0 {
		return nil, nil
	}
	if value, unescaped := unescapeSpecialCharacters(value); unescaped {
		return value, nil
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


func getAssetURIs(value string) []string {
	var separator = " "
	if strings.Contains(value, "[") || strings.Contains(value, "{", ) {
		separator = "|"
	}
	if separator != "|" {
		value = strings.Replace(value, "|", separator, len(value))
	}
	var result = make([]string, 0)
	for _, item := range strings.Split(value, separator) {
		item = strings.TrimSpace(item)
		if item == "" {
			continue
		}
		result = append(result, item)
	}
	return result
}
