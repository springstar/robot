package core

import (
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestVec2(t *testing.T) {
	v1 := NewVec2(34.5, 43.7)
	v2 := NewVec2(70, 25)
	v3 := v1.Subed(v2.T)
	assert.Equal(t, float32(-35.5), v3.Get(0, 0))
	assert.Equal(t, float32(18.7), v3.Get(0, 1))

}

