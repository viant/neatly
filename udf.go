package neatly

import (
	"crypto/md5"
	"fmt"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/url"
	"io"
	"os"
	"path"
	"strings"
)

//AsMap converts source into map
func AsMap(source interface{}, state data.Map) (interface{}, error) {
	if source == nil || toolbox.IsMap(source) {
		return source, nil
	}
	if toolbox.IsString(source) {
		aMap := make(map[string]interface{})
		err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(toolbox.AsString(source))).Decode(&aMap)
		if err != nil {
			return nil, err
		}
		return aMap, nil

	}
	return source, nil
}

//AsInt converts source into int
func AsInt(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.AsInt(source), nil
}

//AsFloat converts source into float64
func AsFloat(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.AsFloat(source), nil
}

//AsBool converts source into bool
func AsBool(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.AsBoolean(source), nil
}

//Md5 computes source md5
func Md5(source interface{}, state data.Map) (interface{}, error) {
	hash := md5.New()
	_, err := io.WriteString(hash, toolbox.AsString(source))
	if err != nil {
		return nil, err
	}
	var result = fmt.Sprintf("%x", hash.Sum(nil))
	return result, nil
}


func GetOwnerDirectory(state data.Map) (string, error) {
	if !state.Has(OwnerURL) {
		return "", fmt.Errorf("OwnerURL was empty")
	}
	var resource = url.NewResource(state.GetString(OwnerURL))
	return resource.DirectoryPath(), nil
}

//HasResource check if patg/url to external resource exists
func HasResource(source interface{}, state data.Map) (interface{}, error) {
	var parentDirectory = ""
	if state.Has(OwnerURL) {
		parentDirectory, _ = GetOwnerDirectory(state)
	}
	filename := path.Join(parentDirectory, toolbox.AsString(source))
	var result = toolbox.FileExists(filename)
	return result, nil
}

//LoadNeatly loads neatly document as data structure, source represents path to nearly document
func LoadNeatly(source interface{}, state data.Map) (interface{}, error) {
	var filename = toolbox.AsString(source)

	if !strings.HasPrefix(filename, "/") {
		var parentDirectory = ""
		if state.Has(OwnerURL) {
			parentDirectory, _ = GetOwnerDirectory(state)
		}
		filename = path.Join(parentDirectory, filename)
	}
	if !toolbox.FileExists(filename) {
		return nil, fmt.Errorf("File %v does not exists", filename)
	}
	var documentResource = url.NewResource(filename)
	var dao, ok = state.Get(NeatlyDao).(*Dao)
	if !ok {
		fmt.Errorf("failed to get neatly loader %T", state.Get(NeatlyDao))
	}


	var aMap = make(map[string]interface{})
	newState := data.NewMap()
	newState.Put(OwnerURL, state.Get(OwnerURL));
	newState.Put(NeatlyDao, state.Get(NeatlyDao));
	err := dao.Load(newState, documentResource, &aMap)
	return aMap, err
}

//WorkingDirectory return joined path with current directory, ../ is supported as subpath
func WorkingDirectory(source interface{}, state data.Map) (interface{}, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var subPath = toolbox.AsString(source)

	for strings.HasSuffix(subPath, "../") {
		currentDirectory, _ = path.Split(currentDirectory)
		if len(subPath) == 3 {
			subPath = ""
		} else {
			subPath = string(subPath[3:])
		}
	}
	return path.Join(currentDirectory, subPath), nil
}
