package server

import (
	"github.com/springstar/robot/core"
	"strings"
	"fmt"
)

type Instruction struct {
	cmd string
	params []interface{}
}

func newInstruction() *Instruction {
	return &Instruction{

	}
}

func loadInstructions() {
	lines := core.ReadLines("../server/orders.txt")
	for _, line := range lines {
		command := strings.Split(line, " ")
		instruction := newInstruction()
		instruction.cmd = command[0]
		for i := 1; i < len(command); i++ {
			instruction.params = append(instruction.params, command[i])
		}

		fmt.Println(instruction)

	}	
}