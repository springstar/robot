package server

import (
	"github.com/springstar/robot/core"
	"strings"
	
)

type Instruction struct {
	cmd string
	params []interface{}
}

func newInstruction() *Instruction {
	return &Instruction{

	}
}

func loadInstructions(f string) (instructions []*Instruction) {
	lines := core.ReadLines(f)
	instructions = parseInstructions(lines)
	return
}

func parseInstructions(lines []string) (instructions []*Instruction) {
	for _, line := range lines {
		command := strings.Split(line, " ")
		instruction := newInstruction()
		instruction.cmd = command[0]
		for i := 1; i < len(command); i++ {
			instruction.params = append(instruction.params, command[i])
		}

		instructions = append(instructions, instruction)
	}	

	return

}