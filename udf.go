package neatly

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"

	"github.com/klauspost/pgzip"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/url"
	"gopkg.in/russross/blackfriday.v2"
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
	return toolbox.ToMap(source)
}

//AsInt converts source into int
func AsInt(source interface{}, state data.Map) (interface{}, error) {
	return toolbox.ToInt(source)
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
	if timeText == "now" {
		var now = time.Now()
		timeValue = &now
	} else {
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
	result := blackfriday.Run([]byte(input))
	return string(result), nil
}

//Cat returns content of supplied file name
func Cat(source interface{}, state data.Map) (interface{}, error) {
	candidate := url.NewResource(toolbox.AsString(source))
	filename := candidate.ParsedURL.Path
	if !toolbox.FileExists(filename) {
		var parentDirectory = ""
		if state.Has(OwnerURL) {
			parentDirectory, _ = GetOwnerDirectory(state)
		}
		filename = path.Join(parentDirectory, toolbox.AsString(source))
	}

	if !toolbox.FileExists(filename) {
		var resource = url.NewResource(state.GetString(OwnerURL))
		parentURL, _ := toolbox.URLSplit(resource.URL)
		var URL = toolbox.URLPathJoin(parentURL, toolbox.AsString(source))
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
