package main


import (
	"fmt"
)


func (sys *System) disassemble() {
	for i := 0x200; i < len(sys.memory); i += 2 {
		opcode := (uint16(sys.memory[i]) << 8) | uint16(sys.memory[i + 1])
		if opcode == 0x0A00 || opcode == 0x0000 {
			break
		}

		fmt.Printf("0x%04X: 0x%04X %s\n", sys.memory[i], opcode, describeOp(opcode))
	}
}


func describeOp(op uint16) string {
	switch op & 0xF000 {
	case 0x0000:
		switch op {
		// CLS - Clear display
		case 0x00E0:
			return "[CLS] - Clear display"

		// RET - return from subroutine
		case 0x00EE:
			return "[RET] - Return from subroutine"

		// SYS - jump to machine code routine at address

		default:
			return "[N/A] - Instruction not found"
		}

	// JMP - jump to address
	case 0x1000:
		return fmt.Sprintf("[JMP] - Jump to 0x%X", op & 0x0FFF)

	// CALL - call subroutine
	case 0x2000:
		return fmt.Sprintf("[CALL] - Call subroutine at 0x%X", op & 0x0FFF)

	// SE - Skip next instruction if Vx == val
	case 0x3000:
		return fmt.Sprintf("[SE] - Skip next instruction if registers[0x%X] == 0x%X", (op & 0x0F00) >> 8, op & 0x00FF)

	// SNE - skip next instruction if Vx != val
	case 0x4000:
		return fmt.Sprintf("[SNE] - Skip next instruction if registers[0x%X] != 0x%X", (op & 0x0F00) >> 8, op & 0x00FF)

	// SE - skip if Vx == Vy
	case 0x5000:
		return fmt.Sprintf("[SE] - Skip next instruction if registers[0x%X] == registers[0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

	// LD - sets register
	case 0x6000:
		return fmt.Sprintf("[LD] - Set registers[0x%X] to 0x%X", (op & 0x0F00) >> 8, op & 0x00FF)

	// ADD - Vx = Vx + val
	case 0x7000:
		return fmt.Sprintf("[ADD] - registers[0x%X] += 0x%X", (op & 0x0F00) >> 8, op & 0x00FF)

	// Operation between two registers
	case 0x8000:
		switch op & 0x000F {
			// OR
			case 0x1:
				return fmt.Sprintf("[OR] - registers[0x%X] |= registers[0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

			// AND
			case 0x2:
				return fmt.Sprintf("[AND] - registers[0x%X] &= registers[0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

			// XOR
			case 0x3:
				return fmt.Sprintf("[XOR] - registers[0x%X] ^= registers[0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

			// ADD
			case 0x4:
				return fmt.Sprintf("[ADD] - registers[0x%X] += registers[0x%X] and set registers[0xF] for overflow", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

			// SUB
			case 0x5:
				return fmt.Sprintf("[SUB] - registers[0x%X] -= registers[0x%X] and set registers[0xF] for borrow", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

			// SHR
			case 0x6:
				return fmt.Sprintf("[SHR] - registers[0x%X] /= 2 and set registers[0xF] if odd", (op & 0x0F00) >> 8)

			// SUBN
			case 0x7:
				return fmt.Sprintf("[SUBN] - registers[0x%X] = registers[0x%X] - registers[0x%X] and set registers[0xF] for borrow", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4, (op & 0x0F00) >> 8)

			// SHL
			case 0xE:
				return fmt.Sprintf("[SHL] - registers[0x%X] *= 2 and set registers[0xF] if odd", (op & 0x0F00) >> 8)

			default:
				return "[N/A] - Instruction not found"
		}
		break

	// SNE
	case 0x9000:
		return fmt.Sprintf("[SNE] - Skip next instruction if registers[0x%X] != registers[0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)

	// LD
	case 0xA000:
		return fmt.Sprintf("[LD] - set I to 0x%X", op & 0x0FFF)

	// JMP
	case 0xB000:
		return fmt.Sprintf("[JMP] - Jump to 0x%X + registers[0x0]", op & 0x0FFF)

	// RND
	case 0xC000:
		return fmt.Sprintf("[RND] - registers[0x%X] = rnd & 0x%X", (op & 0x0F00) >> 8, op & 0x00FF)

	// DRW
	case 0xD000:
		return fmt.Sprintf("[DRW] - draws sprite from I to I + %d - 1 starting at (%d, %d)", op & 0x000F, (op & 0x0F00) >> 8, (op & 0x0F00) >> 4)
		break

	case 0xE000:

		switch op & 0x00FF {
		// SKP
		case 0x009E:
			return fmt.Sprintf("[SKP] - Skip next instruction if key pressed == registers[0x%X]", op & 0x0F00)

		// SKNP
		case 0x00A1:
			return fmt.Sprintf("[SKNP] - Skip next instruction if key pressed != registers[0x%X]", op & 0x0F00)

		default:
			return "[N/A] - Instruction not found"
		}

	case 0xF000:
		switch op & 0x00FF {
		// LD - Load delay timer value into vx
		case 0x0007:
			return fmt.Sprintf("[LD] - registers[0x%X] = delay timer", (op & 0x0F00) >> 8)

		// LD - load from input
		case 0x000A:
			return fmt.Sprintf("[LD] - registers[0x%X] = input", (op & 0x0F00) >> 8)

		// LD - Set delay timer
		case 0x0015:
			return fmt.Sprintf("[LD] - delay timer = registers[0x%X]", (op & 0x0F00) >> 8)

		// LD - Set sound timer
		case 0x0018:
			return fmt.Sprintf("[LD] - sound timer = registers[0x%X]", (op & 0x0F00) >> 8)

		// ADD - I and Vx
		case 0x001E:
			return fmt.Sprintf("[ADD] - I += registers[0x%X]", (op & 0x0F00) >> 8)

		// LD - Set I to the value of the location of the sprite
		case 0x0029:
			return fmt.Sprintf("[LD] - I = location of sprite at registers[0x%X]", (op & 0x0F00) >> 8)

		// LD - Store BCD representation in to I, I+1, I+2
		case 0x0033:
			return fmt.Sprintf("[LD] - Store BCD representation of registers[0x%X] into I to I + 2", (op & 0x0F00) >> 8)

		// LD - store registers in memory
		case 0x0055:
			return fmt.Sprintf("[LD] - Store registers[0x0] to registers[0x%X] into memory[I] to memory[I + 0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 8)

		// LD - load register from memory
		case 0x0065:
			return fmt.Sprintf("[LD] - Load registers[0x0] to registers[0x%X] from memory[I] to memory[I + 0x%X]", (op & 0x0F00) >> 8, (op & 0x0F00) >> 8)

		default:
			return "[N/A] - Instruction not found"
		}
		break

	default:
		return "[N/A] - Instruction not found"
	}

	return "[N/A] - Instruction not found"
}
