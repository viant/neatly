package neatly

import (
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
	OwnerSource *url.Resource
	OwnerName   string
	Name        string
	Group       string
	IsArray     bool
	Iterator    *TagIterator
	LineNumber  int
	Subpath     string
	PathMatch   string
	tagIdPrefix string
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
		t.Subpath, t.PathMatch = t.expandPathIfNeeded(toolbox.AsString(value))
	}
	if t.Subpath != "" {
		context.Subpath = t.Subpath
	}
	context.tagID = t.TagID()

	if includeMeta {
		t.setMeta(result, record)
	}
	context.tagObject = result
	return result
}

func (t *Tag) expandPathIfNeeded(subpath string) (string, string) {
	if !strings.HasSuffix(subpath, "*") {
		return subpath, ""
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
	storageService, err := storage.NewServiceForURL(parentURL, t.OwnerSource.Credentials)
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
						return path.Join(subPathParent, candidateName), string(candidateName[len(leafDirectory):])
					}
					return candidateName, string(candidateName[len(leafDirectory):])
				}
			}
		}
	}
	return subpath, ""
}

//setMeta sets Tag, optionally TagIndex and Subpath to the provided object
func (t *Tag) setMeta(object data.Map, record map[string]interface{}) {
	object["Tag"] = t.Name
	if t.HasActiveIterator() {
		object["TagIndex"] = t.Iterator.Index()
	}

	value, has := record["Subpath"]
	if has {
		t.SetSubPath(toolbox.AsString(value))
	}
	if t.Subpath != "" {
		object["Subpath"] = t.Subpath
	}

	if t.PathMatch != "" {
		object["PathMatch"] = t.PathMatch
	}
	if value, has := record["Group"]; has {
		t.Group = toolbox.AsString(value)
	}
	object["TagID"] = t.TagID()

}

func (t *Tag) Expand(text string) string {
	var aMap = data.NewMap()
	aMap.Put("pathMatch", t.PathMatch)
	aMap.Put("subPath", t.Subpath)
	if t.HasActiveIterator() {
		aMap.Put("index", t.Iterator.Index())
		aMap.Put("idx", toolbox.AsInt(t.Iterator.Index()))
	}
	return aMap.ExpandAsText(text)
}

//SetSubPath set subpath for the tag
func (t *Tag) SetSubPath(subpath string) {
	t.Subpath, t.PathMatch = t.expandPathIfNeeded(subpath)
}

//TagID returns tag ID
func (t *Tag) TagID() string {
	var index = ""
	if t.HasActiveIterator() {
		index = t.Iterator.Index()
	}

	var subPath = t.Subpath
	if subPath != "" {
		if strings.Contains(subPath, index) {
			index = ""
		}
	}
	if strings.Contains(t.Name, "$") {
		expandedName := t.Expand(t.Name)
		if strings.Contains(subPath, expandedName) {
			subPath = ""
		}
	}
	var tagIdPostfix = t.Group + index + subPath
	if tagIdPostfix != "" && t.tagIdPrefix != "" {
		tagIdPostfix = "_" + tagIdPostfix
	}

	value := t.Expand(t.tagIdPrefix + tagIdPostfix)
	var result = make([]byte, 0)
	for _, r := range value {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '_' {
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

	if rangeIndex := strings.LastIndex(result.Name, "{"); rangeIndex != -1 {
		result.Name = string(result.Name[:rangeIndex])
	}

	if ownerName != "" {
		ownerName = ownerName + "_"
	}
	result.tagIdPrefix = ownerName + result.Name
	return result
}
