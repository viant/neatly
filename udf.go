package neatly

import (
	"bytes"
	"compress/gzip"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/gomarkdown/markdown"
	"github.com/klauspost/pgzip"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/data/udf"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/url"
	"io"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

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
	content, err := LoadBinary(source, state)
	if err != nil {
		return nil, err
	}
	return toolbox.AsString(content), err
}

//LoadBinary returns []byte content of supplied file name
func LoadBinary(source interface{}, state data.Map) (interface{}, error) {
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
	return content, nil
}

//AssetsToMap loads assets into map[string]string, it takes url, with optional list of extension as filter
func AssetsToMap(source interface{}, state data.Map) (interface{}, error) {
	if source == nil {
		return nil, nil
	}
	var result = make(map[string]string)
	updator := func(key string, data []byte) {
		result[key] = string(data)
	}
	return assetToMap(source, state, updator, result)
}

//BinaryAssetsToMap loads binary assets into map[string]string, it takes url, with optional list of extension as filter
func BinaryAssetsToMap(source interface{}, state data.Map) (interface{}, error) {
	if source == nil {
		return nil, nil
	}
	var result = make(map[string][]byte)
	updator := func(key string, data []byte) {
		result[key] = data
	}
	return assetToMap(source, state, updator, result)
}

func assetToMap(source interface{}, state data.Map, updator func(key string, data []byte), result interface{}) (interface{}, error) {
	URL, ok := source.(string) //URL param case
	if ok {
		return result, loadAssetToMap(url.NewResource(URL), updator)
	}
	//url.Resource param case
	resource := &url.Resource{}
	if toolbox.IsStruct(source) || toolbox.IsMap(source) {
		if err := toolbox.DefaultConverter.AssignConverted(&resource, source); err == nil {
			return result, loadAssetToMap(resource, updator)
		}
	}
	if toolbox.IsSlice(source) { //URL, credentials params case
		params := toolbox.AsSlice(source)
		return result, loadAssetToMap(url.NewResource(params...), updator)
	}
	return nil, fmt.Errorf("unsupported source %T", source)
}

func loadAssetToMap(resource *url.Resource, updator func(key string, data []byte)) error {
	storageService, err := storage.NewServiceForURL(resource.URL, resource.Credentials)
	if err != nil {
		return err
	}
	objects, err := storageService.List(resource.URL)
	if err != nil {
		return err
	}
	for _, object := range objects {
		if object.IsFolder() {
			continue
		}
		reader, err := storageService.Download(object)
		if err != nil {
			return err
		}
		defer reader.Close()
		content, err := ioutil.ReadAll(reader)
		if err != nil {
			return err
		}
		info := object.FileInfo()
		updator(info.Name(), content)
	}
	return err
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

//AddStandardUdf register building udf to the context
func AddStandardUdf(aMap data.Map) {
	udf.Register(aMap)
	aMap.Put("IsJSON", IsJSON)
	aMap.Put("WorkingDirectory", WorkingDirectory)
	aMap.Put("Pwd", WorkingDirectory)
	aMap.Put("HasResource", HasResource)
	aMap.Put("Md5", Md5)
	aMap.Put("LoadNeatly", LoadNeatly)
	aMap.Put("Zip", Zip)
	aMap.Put("Unzip", Unzip)
	aMap.Put("UnzipText", UnzipText)
	aMap.Put("Markdown", Markdown)
	aMap.Put("Cat", Cat)
	aMap.Put("LoadBinary", LoadBinary)
	aMap.Put("AssetsToMap", AssetsToMap)
	aMap.Put("BinaryAssetsToMap", BinaryAssetsToMap)
}
