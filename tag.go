package neatly

import (
	"fmt"
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
	"github.com/viant/toolbox/storage"
	"github.com/viant/toolbox/url"
	"path"
	"strings"
	"unicode"
)

//Tag represents a nearly tag
type Tag struct {
	OwnerSource   *url.Resource
	OwnerName     string
	Name          string
	Group         string
	IsArray       bool
	Iterator      *TagIterator
	LineNumber    int
	Subpath       string
	TagIDTemplate string
}

//HasActiveIterator returns true if tag has active iterator
func (t *Tag) HasActiveIterator() bool {
	if t == nil {
		return false
	}
	return t.Iterator != nil && t.Iterator.Has()
}

func (t *Tag) setTagObject(context *tagContext, record map[string]interface{}, includeMeta bool) data.Map {
	var result data.Map
	if t.IsArray {
		result = data.NewMap()
		context.objectContainer.GetCollection(t.Name).Push(result)
	} else {
		result = context.objectContainer.GetMap(t.Name)
	}

	context.tagName = t.Name
	if t.HasActiveIterator() {
		context.tagIndex = t.Iterator.Index()
	}
	value, has := record["Subpath"]
	if has {
		t.Subpath = t.expandPathIfNeeded(toolbox.AsString(value))
	}
	if t.Subpath != "" {
		context.Subpath = t.Subpath
	}
	context.tagID = t.TagID()

	/*

		tagName  string
		tagIndex string
		Subpath string
		tagID string



	*/

	if includeMeta {
		t.setMeta(result, record)
	}
	context.tagObject = result
	return result
}

func (t *Tag) expandPathIfNeeded(subpath string) string {
	if !strings.HasSuffix(subpath, "*") {
		return subpath
	}
	parentURL, _ := toolbox.URLSplit(t.OwnerSource.URL)
	var leafDirectory = ""
	var subPathParent = ""
	subPathElements := strings.Split(subpath, "/")
	for _, candidate := range subPathElements {
		if strings.Contains(candidate, "*") {
			leafDirectory = strings.Replace(candidate, "*", "", 1)
			break
		}
		subPathParent = path.Join(subPathParent, candidate)
		parentURL = toolbox.URLPathJoin(parentURL, candidate)
	}
	storageService, err := storage.NewServiceForURL(parentURL, t.OwnerSource.Credential)
	if err == nil {
		candidates, err := storageService.List(parentURL)
		if err == nil {
			for _, candidate := range candidates {
				if candidate.URL() == parentURL {
					continue
				}
				_, candidateName := toolbox.URLSplit(candidate.URL())
				if strings.HasPrefix(candidateName, leafDirectory) {
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
	object["TagID"] = t.TagID()

	if value, has := record["Group"];has {
		t.Group = toolbox.AsString(value)
	}

}

//TagID returns tag ID
func (t *Tag) TagID() string {
	var index = ""
	if t.HasActiveIterator() {
		index = t.Iterator.Index()
	}
	value := fmt.Sprintf(t.TagIDTemplate, t.Group, index, t.Subpath)
	var result = make([]byte, 0)
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
		OwnerName:   ownerName,
		OwnerSource: ownerSource,
		Name:        key,
		LineNumber:  lineNumber,
	}
	key = decodeIteratorIfPresent(key, result)
	if len(key) > 2 && string(key[0:2]) == "[]" {
		result.Name = string(key[2:])
		result.IsArray = true
	}
	result.TagIDTemplate = ownerName + result.Name + "%v%v%v"
	return result
}
