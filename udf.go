package neatly

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/gomarkdown/markdown"
	"github.com/klauspost/pgzip"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/url"
)

//Keys returns keys of the supplied map
func Keys(source interface{}, state data.Map) (interface{}, error) {
	aMap, err := AsMap(source, state)
	if err != nil {
		return nil, err
	}
	var result = make([]interface{}, 0)
	err = toolbox.ProcessMap(aMap, func(key, value interface{}) bool {
		result = append(result, key)
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

//Values returns values of the supplied map
func Values(source interface{}, state data.Map) (interface{}, error) {
	aMap, err := AsMap(source, state)
	if err != nil {
		return nil, err
	}
	var result = make([]interface{}, 0)
	err = toolbox.ProcessMap(aMap, func(key, value interface{}) bool {
		result = append(result, value)
		return true
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

//IndexOf returns index of the matched slice elements or -1
func IndexOf(source interface{}, state data.Map) (interface{}, error) {
	if toolbox.IsSlice(source) {
		return nil, fmt.Errorf("expected arguments but had: %T", source)
	}
	args := toolbox.AsSlice(source)
	if len(args) != 2 {
		return nil, fmt.Errorf("expected 2 arguments but had: %v", len(args))
	}

	collection, err := AsCollection(args[0], state)
	if err != nil {
		return nil, err
	}
	for i, candidate := range toolbox.AsSlice(collection) {
		if candidate == args[1] || toolbox.AsString(candidate) == toolbox.AsString(args[1]) {
			return i, nil
		}
	}
	return -1, nil
}

//AsMap converts source into map
func AsMap(source interface{}, state data.Map) (interface{}, error) {

	if source == nil || toolbox.IsMap(source) {
		return source, nil
	}
	source = convertToTextIfNeeded(source)
	if toolbox.IsString(source) {
		aMap := make(map[string]interface{})
		err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(toolbox.AsString(source))).Decode(&aMap)
		if err != nil {
			return nil, err
		}
		return aMap, nil
	}
	return toolbox.ToMap(source)
}

//AsCollection converts source into a slice
func AsCollection(source interface{}, state data.Map) (interface{}, error) {
	source = convertToTextIfNeeded(source)
	if source == nil || toolbox.IsMap(source) {
		return source, nil
	}
	source = convertToTextIfNeeded(source)
	if toolbox.IsString(source) {
		aSlice := []interface{}{}
		err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(toolbox.AsString(source))).Decode(&aSlice)
		if err != nil {
			if e := yaml.NewDecoder(strings.NewReader(toolbox.AsString(source))).Decode(&aSlice);e!= nil {
				return nil, err
			}
		}
		return aSlice, nil
	}
	return toolbox.AsSlice(source), nil
}

//AsData converts source into map or slice
func AsData(source interface{}, state data.Map) (interface{}, error) {
	source = convertToTextIfNeeded(source)
	if source == nil || toolbox.IsMap(source) || toolbox.IsSlice(source) {
		return source, nil
	}

	if toolbox.IsString(source) {
		var result interface{}
		err := toolbox.NewJSONDecoderFactory().Create(strings.NewReader(toolbox.AsString(source))).Decode(&result)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	return source, nil
}

func convertToTextIfNeeded(data interface{}) interface{} {
	if data == nil {
		return data
	}
	if bs, ok := data.([]byte); ok {
		return string(bs)
	}
	return data
}

//AsInt converts source into int
func AsInt(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.ToInt(source)
}

//AsInt converts source into int
func AsString(source interface{}, state data.Map) (interface{}, error) {
	if toolbox.IsSlice(source) || toolbox.IsMap(source) || toolbox.IsStruct(source) {
		text, err := toolbox.AsJSONText(source)
		if err == nil {
			return text, nil
		}
	}
	return toolbox.AsString(source), nil

}

//Add increment supplied state key with delta ['key', -2]
func Increment(args interface{}, state data.Map) (interface{}, error) {
	if toolbox.IsSlice(args) {
		return nil, fmt.Errorf("args were not slice: %T", args)
	}
	aSlice := toolbox.AsSlice(args)
	if len(aSlice) != 2 {
		return nil, fmt.Errorf("expeted 2 arguments but had: %v", len(aSlice))

	}
	var delta = toolbox.AsInt(aSlice[1])
	var exrp = toolbox.AsString(aSlice[0])
	value, has := state.GetValue(exrp)
	if !has {
		state.SetValue(exrp, delta)
	} else {
		state.SetValue(exrp, delta+toolbox.AsInt(value))
	}
	value, _ = state.GetValue(exrp)
	return value, nil
}

//AsFloat converts source into float64
func AsFloat(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.AsFloat(source), nil
}

//Length returns length of slice or string
func Length(source interface{}, state data.Map) (interface{}, error) {

	if toolbox.IsSlice(source) {
		return len(toolbox.AsSlice(source)), nil
	}
	if toolbox.IsMap(source) {
		return len(toolbox.AsMap(source)), nil
	}
	if text, ok := source.(string); ok {
		return len(text), nil
	}
	return 0, nil
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

//GetOwnerDirectory returns owner neatly document directory
func GetOwnerDirectory(state data.Map) (string, error) {
	if !state.Has(OwnerURL) {
		return "", fmt.Errorf("OwnerURL was empty")
	}
	var resource = url.NewResource(state.GetString(OwnerURL))
	return resource.DirectoryPath(), nil
}

//HasResource check if patg/url to external resource exists
func HasResource(source interface{}, state data.Map) (interface{}, error) {
	filename := toolbox.AsString(source)

	if !strings.HasPrefix(filename, "/") {
		var parentDirectory = ""
		if state.Has(OwnerURL) {
			parentDirectory, _ = GetOwnerDirectory(state)
		}
		candidate := path.Join(parentDirectory, toolbox.AsString(source))
		if toolbox.FileExists(candidate) {
			return true, nil
		}
	}
	var result = url.NewResource(filename).ParsedURL.Path
	return toolbox.FileExists(result), nil
}

//LoadNeatly loads neatly document as data structure, source represents path to nearly document
func LoadNeatly(source interface{}, state data.Map) (interface{}, error) {
	var filename = toolbox.AsString(source)
	var parentDirectory = ""
	if !strings.HasPrefix(filename, "/") {
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
		return nil, fmt.Errorf("failed to get neatly loader %T", state.Get(NeatlyDao))
	}
	var aMap = make(map[string]interface{})
	newState := data.NewMap()
	newState.Put(OwnerURL, state.Get(OwnerURL))
	newState.Put(NeatlyDao, state.Get(NeatlyDao))

	for k, v := range state {
		if toolbox.IsFunc(v) {
			newState.Put(k, v)
		}
	}
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
	if subPath == "" {
		return currentDirectory, nil
	}
	return path.Join(currentDirectory, subPath), nil
}

//FormatTime return formatted time, it takes an array of two arguments, the first id time, or now followed by java style time format.
func FormatTime(source interface{}, state data.Map) (interface{}, error) {
	if !toolbox.IsSlice(source) {
		return nil, fmt.Errorf("unable to run FormatTime: expected %T, but had: %T", []interface{}{}, source)
	}
	aSlice := toolbox.AsSlice(source)
	if len(aSlice) < 2 {
		return nil, fmt.Errorf("unable to run FormatTime, expected 2 parameters, but had: %v", len(aSlice))
	}
	var err error
	var timeText = toolbox.AsString(aSlice[0])
	var timeFormat = toolbox.AsString(aSlice[1])
	var timeLayout = toolbox.DateFormatToLayout(timeFormat)
	var timeValue *time.Time
	timeValue, err = toolbox.TimeAt(timeText)
	if err != nil {
		timeValue, err = toolbox.ToTime(aSlice[0], timeLayout)
	}
	if err != nil {
		return nil, err
	}
	if len(aSlice) > 2 {
		timeLocation, err := time.LoadLocation(toolbox.AsString(aSlice[2]))
		if err != nil {
			return nil, err
		}
		timeInLocation := timeValue.In(timeLocation)
		timeValue = &timeInLocation
	}
	return timeValue.Format(timeLayout), nil
}

//Unzip uncompress supplied []byte or error
func Unzip(source interface{}, state data.Map) (interface{}, error) {
	payload, ok := source.([]byte)
	if !ok {
		return nil, fmt.Errorf("invalid Unzip input, expected %T, but had %T", []byte{}, source)
	}
	reader, err := pgzip.NewReader(bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("failed to create gzip reader %v", err)
	}
	payload, err = ioutil.ReadAll(reader)
	if err != nil {
		return nil, fmt.Errorf("failed to read gzip reader %v", err)
	}
	return payload, err
}

//UnzipText uncompress supplied []byte into text or error
func UnzipText(source interface{}, state data.Map) (interface{}, error) {
	payload, err := Unzip(source, state)
	if err != nil {
		return nil, err
	}
	return toolbox.AsString(payload), nil
}

//Zip compresses supplied []byte or test or error
func Zip(source interface{}, state data.Map) (interface{}, error) {
	payload, ok := source.([]byte)
	if !ok {
		if text, ok := source.(string); ok {
			payload = []byte(text)
		} else {
			return nil, fmt.Errorf("invalid Zip input, expected %T, but had %T", []byte{}, source)
		}
	}
	buffer := new(bytes.Buffer)
	writer, err := pgzip.NewWriterLevel(buffer, gzip.BestSpeed)
	if err != nil {
		return nil, err
	}
	_, err = writer.Write(payload)
	if err != nil {
		return nil, fmt.Errorf("error in Zip, failed to write %v", err)
	}
	_ = writer.Flush()
	err = writer.Close()
	return buffer.Bytes(), err
}

//Markdown returns html fot supplied markdown
func Markdown(source interface{}, state data.Map) (interface{}, error) {
	var input = toolbox.AsString(source)
	response, err := Cat(input, state)
	if err == nil && response != nil {
		input = toolbox.AsString(response)
	}
	result := markdown.ToHTML([]byte(input), nil, nil)
	return string(result), nil
}

//Cat returns content of supplied file name
func Cat(source interface{}, state data.Map) (interface{}, error) {
	filename := toolbox.AsString(source)
	candidate := url.NewResource(filename)
	if candidate != nil || candidate.ParsedURL != nil {
		filename = candidate.ParsedURL.Path
	}
	if !toolbox.FileExists(filename) {
		var parentDirectory = ""
		if state.Has(OwnerURL) {
			parentDirectory, _ = GetOwnerDirectory(state)
		}
		filename = path.Join(parentDirectory, toolbox.AsString(source))
	}
	if !toolbox.FileExists(filename) {
		filename := toolbox.AsString(source)
		var resource = url.NewResource(state.GetString(OwnerURL))
		parentURL, _ := toolbox.URLSplit(resource.URL)
		var URL = toolbox.URLPathJoin(parentURL, filename)
		service, err := storage.NewServiceForURL(URL, "")
		if err == nil {
			if exists, _ := service.Exists(URL); exists {
				resource = url.NewResource(URL)
				if text, err := resource.DownloadText(); err == nil {
					return text, nil
				}
			}
		}
		return nil, fmt.Errorf("no such file or directory %v", filename)
	}
	file, err := toolbox.OpenFile(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	content, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}
	return string(content), nil
}

// Validate if JSON file is a well-formed JSON
// Returns true if file content is valid JSON
func IsJSON(fileName interface{}, state data.Map) (interface{}, error) {
	content, err := Cat(fileName, state)
	if err != nil {
		return false, err
	}
	var m json.RawMessage
	if err := json.Unmarshal([]byte(toolbox.AsString(content)), &m); err != nil {
		return false, err
	}
	return true, nil
}

// Join joins slice by separator
func Join(args interface{}, state data.Map) (interface{}, error) {
	if !toolbox.IsSlice(args) {
		return nil, fmt.Errorf("expected 2 arguments but had: %T", args)
	}
	arguments := toolbox.AsSlice(args)
	if len(arguments) != 2 {
		return nil, fmt.Errorf("expected 2 arguments but had: %v", len(arguments))
	}

	if !toolbox.IsSlice(arguments[0]) {
		return nil, fmt.Errorf("expected 1st arguments as slice but had: %T", arguments[0])
	}
	var result = make([]string, 0)
	toolbox.CopySliceElements(arguments[0], &result)
	return strings.Join(result, toolbox.AsString(arguments[1])), nil
}

//AddStandardUdf register building udf to the context
func AddStandardUdf(state data.Map) {
	state.Put("AsMap", AsMap)
	state.Put("AsData", AsData)
	state.Put("AsCollection", AsCollection)
	state.Put("AsInt", AsInt)
	state.Put("AsString", AsString)
	state.Put("AsFloat", AsFloat)
	state.Put("AsBool", AsBool)

	state.Put("WorkingDirectory", WorkingDirectory)
	state.Put("Pwd", WorkingDirectory)
	state.Put("HasResource", HasResource)
	state.Put("Md5", Md5)
	state.Put("Length", Length)
	state.Put("Len", Length)
	state.Put("LoadNeatly", LoadNeatly)
	state.Put("FormatTime", FormatTime)
	state.Put("Zip", Zip)
	state.Put("Unzip", Unzip)
	state.Put("UnzipText", UnzipText)
	state.Put("Markdown", Markdown)
	state.Put("Cat", Cat)
	state.Put("IsJSON", IsJSON)
	state.Put("Join", Join)
	state.Put("Keys", Keys)
	state.Put("Values", Values)
}
