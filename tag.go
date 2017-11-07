package neatly

import (
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/url"
	"fmt"
	"strings"
	"github.com/viant/toolbox/storage"
	"path"
	"unicode"
)

//Tag represents a nearly tag
type Tag struct {
	OwnerSource   *url.Resource
	OwnerName     string
	Name          string
	IsArray       bool
	Iterator      *TagIterator
	LineNumber    int
	Subpath       string
	TagIdTemplate string
}

//HasActiveIterator returns true if tag has active iterator
func (t *Tag) HasActiveIterator() bool {
	if t == nil {
		return false
	}
	return t.Iterator != nil && t.Iterator.Has()
}

func (t *Tag) setTagObject(context *tagContext, record map[string]interface{}) data.Map {
	var result data.Map
	if t.IsArray {
		result = data.NewMap()
		context.objectContainer.GetCollection(t.Name).Push(result)
	} else {
		result = context.objectContainer.GetMap(t.Name)
	}
	t.setMeta(result, record)
	context.tagObject = result
	return result
}

func (t *Tag) expandPathIfNeeded(subpath string) string {
	if ! strings.HasSuffix(subpath, "*") {
		return subpath
	}
	subPathElements := strings.Split(subpath, "/")
	var subPathParent = ""
	subPathPrefix := strings.Replace(subpath, "*", "", 1)

	parentURL, _ := toolbox.URLSplit(t.OwnerSource.URL)
	if len(subPathPrefix) > 1 {
		subPathPrefix = strings.Replace(subPathElements[len(subPathElements)-1], "*", "", 1)
		subPathParent = path.Join(subPathElements[:len(subPathElements)-1]...)
		parentURL = toolbox.URLPathJoin(parentURL, subPathParent)
	}
	storageService, err := storage.NewServiceForURL(parentURL, t.OwnerSource.Credential)
	if err == nil {
		candidates, err := storageService.List(parentURL)
		if err == nil {
			for _, candidate := range candidates {
				_, candidateName := toolbox.URLSplit(candidate.URL())
				if strings.HasPrefix(candidateName, subPathPrefix) {
					if subPathParent != "" {
						return path.Join(subPathParent, candidateName)
					}
					return candidateName
				}
			}
		}
	}
	return subpath
}


//setMeta sets Tag, optionally TagIndex and Subpath to the provided object
func (t *Tag) setMeta(object data.Map, record map[string]interface{}) {
	object["Tag"] = t.Name
	if t.HasActiveIterator() {
		object["TagIndex"] = t.Iterator.Index()
	}
	value, has := record["Subpath"]
	if has {
		t.Subpath = t.expandPathIfNeeded(toolbox.AsString(value))
	}
	if t.Subpath != "" {
		object["Subpath"] = t.Subpath
	}
	object["TagId"] = t.TagId()
}

func (t *Tag) TagId() string {
	var index = ""
	if t.HasActiveIterator() {
		index  = t.Iterator.Index()
	}
	value:=fmt.Sprintf(t.TagIdTemplate, index, t.Subpath)
	var result  = make([]byte, 0)
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) {
			result = append(result, byte(r))
		}
	}
	return string(result)
}


//NewTag creates a new neatly tag
func NewTag(ownerName string, ownerSource *url.Resource, key string, lineNumber int) *Tag {
	var result = &Tag{
		OwnerName:     ownerName,
		OwnerSource:   ownerSource,
		Name:          key,
		LineNumber:    lineNumber,
	}
	key = decodeIteratorIfPresent(key, result)
	if string(key[0:2]) == "[]" {
		result.Name = string(key[2:])
		result.IsArray = true
	}
	result.TagIdTemplate = ownerName + result.Name + "%v%v"
	return result
}
