package compile

func Optimize(original []Instruction) []Instruction {
	var instructions []Instruction
	if original == nil || len(original) == 0 {
		return instructions
	}

	var cur *Instruction
	cur = &original[0]
	for _, instruction := range original[1:] {
		if cur == nil {
			cur = &instruction
			continue
		}
		switch {
		case cur.Opcode == OpPrint && cur.Opcode == instruction.Opcode:
			cur.Arg = cur.Arg.(string) + instruction.Arg.(string)
		default:
			instructions = append(instructions, *cur)
			cur = &instruction
		}
	}

	if cur != nil {
		instructions = append(instructions, *cur)
	}

	return instructions
}
