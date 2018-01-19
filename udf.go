package neatly

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"fmt"
	"github.com/klauspost/pgzip"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/url"
	"gopkg.in/russross/blackfriday.v2"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
	"time"
	"github.com/viant/toolbox/storage"
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
	return path.Join(currentDirectory, subPath), nil
}

//FormatTime return formatted time, it takes an array of two arguments, the first id time, or now followed by java style time format.
func FormatTime(source interface{}, state data.Map) (interface{}, error) {
	if !toolbox.IsSlice(source) {
		return nil, fmt.Errorf("unable to run FormatTime: expected %T, but had: %T", []interface{}{}, source)
	}
	aSlice := toolbox.AsSlice(source)
	if len(aSlice) != 2 {
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
	var parentDirectory = ""
	if state.Has(OwnerURL) {
		parentDirectory, _ = GetOwnerDirectory(state)
	}
	filename := path.Join(parentDirectory, toolbox.AsString(source))
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
	file, err := os.Open(filename)
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
