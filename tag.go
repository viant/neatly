package neatly

import (
	"github.com/viant/toolbox"
	"github.com/viant/toolbox/data"
)

//Tag represents a nearly tag
type Tag struct {
	Name       string
	IsArray    bool
	Iterator   *TagIterator
	LineNumber int
	Subpath    string
}

//HasActiveIterator returns true if tag has active iterator
func (t *Tag) HasActiveIterator() bool {
	if t == nil {
		return false
	}
	return t.Iterator != nil && t.Iterator.Has()
}

func (t *Tag) getObject(objectContainer data.Map, record map[string]interface{}) data.Map {
	var result data.Map
	if t.IsArray {
		result = data.NewMap()
		objectContainer.GetCollection(t.Name).Push(result)
	} else {
		result = objectContainer.GetMap(t.Name)
	}
	t.setMeta(result, record)
	return result
}

//setMeta sets Tag, optionally TagIndex and Subpath to the provided object
func (t *Tag) setMeta(object data.Map, record map[string]interface{}) {
	object["Tag"] = t.Name
	if t.HasActiveIterator() {
		object["TagIndex"] = t.Iterator.Index()
	}
	subpath, has := record["Subpath"]
	if has {
		t.Subpath = toolbox.AsString(subpath)
	}
	if t.Subpath != "" {
		object["Subpath"] = t.Subpath
	}
}

//NewTag creates a new neatly tag
func NewTag(key string, lineNumber int) *Tag {
	var result = &Tag{
		Name:       key,
		LineNumber: lineNumber,
	}
	key = decodeIteratorIfPresent(key, result)
	if string(key[0:2]) == "[]" {
		result.Name = string(key[2:])
		result.IsArray = true
	}

	return result
}
