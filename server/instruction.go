package server

import (
	_ "fmt"
	"github.com/springstar/robot/core"
	"strings"
	
)

type Instruction struct {
	pc int
	cmd string
	params []string
}

type InstructionList []*Instruction

func newInstruction() *Instruction {
	return &Instruction{

	}
}

func cloneInstructions(il *InstructionList) *InstructionList {
	insList := newInstructionList()
	for _, i := range *il {
		ins := newInstruction()
		ins.cmd = i.cmd
		ins.params = i.params
		insList.addInstrution(ins)
	}

	return insList

}

func loadInstructions(f string) *InstructionList {
	lines := core.ReadLines(f)
	il := parseInstructions(lines)
	return il
}

func parseInstructions(lines []string) *InstructionList {
	il := newInstructionList()
	for _, line := range lines {
		command := strings.Split(line, " ")
		instruction := newInstruction()
		instruction.cmd = command[0]
		for i := 1; i < len(command); i++ {
			instruction.params = append(instruction.params, command[i])
		}

		il.addInstrution(instruction)
	}	

	return il

}

func newInstructionList() *InstructionList{
	return &InstructionList{}
}

func (il *InstructionList) addInstrution(i *Instruction) {
	*il = append(*il, i)	
}

func (il *InstructionList) icount() int {
	return len(*il)
}

func (il *InstructionList) fetch(i int) *Instruction {
	if i >= il.icount() {
		return nil
	}

	return (*il)[i]
}

func (il *InstructionList) next(i int) (int, *Instruction) {
	i = i + 1
	if i >= il.icount() {
		i = 0
	}

	return i, (*il)[i]
}

