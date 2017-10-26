package neatly_test

import (
	"testing"
	"github.com/stretchr/testify/assert"
	"github.com/viant/neatly"
)

func Test_Md5(t *testing.T) {

	//!Md5(www.viewability.com) 8c505168697a000f0946c04e09f2d524
	var md5, err = neatly.Md5("www.viewability.com", nil)
	assert.Nil(t, err)
	assert.EqualValues(t, "123c274fb9a25ddbc77c1634f1e55525", md5)
}