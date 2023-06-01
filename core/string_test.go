package core

import (
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestStr2Float64(t *testing.T) {
	f := Str2Float64("85.237")
	assert.Equal(t, 85.237, f)

	f = Str2Float64(" 85.237 ")
	assert.Equal(t, 85.237, f)

}