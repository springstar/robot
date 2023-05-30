package server

import (
	"testing"
)

func TestLoadInstruction(t *testing.T) {
	instructions := loadInstructions()
	t.Log(instructions)
}