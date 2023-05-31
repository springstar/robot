package server

import (
	"testing"
	"github.com/stretchr/testify/assert"

)

func TestLoadInstruction(t *testing.T) {
	lines := []string {
		"move 127.23 45.25 54.34 45.66",
		"quest",
	}

	instructions := parseInstructions(lines)
	assert.Equal(t, 2, instructions.icount())
	assert.Equal(t, "move", instructions.fetch(0).cmd)
	assert.Equal(t, "quest", instructions.fetch(1).cmd)
	assert.Equal(t, (*Instruction)(nil), instructions.fetch(2))

	
}