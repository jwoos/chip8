package main


import (
	"math/rand"
)


/* Memory map
 * +---------------+= 0xFFF (4095) End of Chip-8 RAM
 * |               |
 * |               |
 * |               |
 * |               |
 * |               |
 * | 0x200 to 0xFFF|
 * |     Chip-8    |
 * | Program / Data|
 * |     Space     |
 * |               |
 * |               |
 * |               |
 * +- - - - - - - -+= 0x600 (1536) Start of ETI 660 Chip-8 programs
 * |               |
 * |               |
 * |               |
 * +---------------+= 0x200 (512) Start of most Chip-8 programs
 * | 0x000 to 0x1FF|
 * | Reserved for  |
 * |  interpreter  |
 * +---------------+= 0x000 (0) Start of Chip-8 RAM
*/


type System struct {
	// 4096 bytes
	memory []uint16

	// 16 byte registers
	registers []byte
	iregister uint16

	delayTimer byte
	soundTimer byte

	programCounter uint16

	stack *Stack

	opcode uint16
}


func newSystem() *System {
	sys := new(System)
	sys.memory = make([]byte, 2048)
	sys.registers = make([]byte, 16)
	sys.stack = newStack(16)
	sys.programCounter = 0x200

	return mem
}


func (sys *System) parseInstruction() {
	op := sys.opcode

	switch op & 0xF000 {
	case 0x0000:
		switch op {
		// CLS - Clear display
		case 0x00E0:

		// RET - return from subroutine
		case 0x00EE:

		// SYS - jump to machine code routine at address
		default:

		}

		break

	// JMP - jump to address
	case 0x1000:
		sys.pc = op & 0x0FFF
		break

	// CALL - call subroutine
	case 0x2000:
		// FIXME handle error
		sys.stack.push(sys.programCounter)
		sys.pc = op & 0x0FFF
		break

	// SE - Skip next instruction if Vx == val
	case 0x3000:
		registerIndex := (op & 0x0F00) >> 2
		lastHalf := op & 0x00FF
		if sys.registers[registerIndex] == lastHalf {
			sys.programCounter += 2
		} else {
			sys.programCounter++
		}
	}

	// SNE - skip next instruction if Vx != val
	case 0x4000:
		registerIndex := (op & 0x0F00) >> 2
		lastHalf := op & 0x00FF
		if sys.registers[registerIndex] == lastHalf {
			sys.programCounter++
		} else {
			sys.programCounter += 2
		}

	// SE - skip if Vx == Vy
	case 0x5000:
		registerA := (op & 0x0F00) >> 2
		registerB := (op & 0x00F0) >> 1

		if sys.registers[registerA] == sys.registers[registerB] {
			sys.programCounter += 2
		} else {
			sys.programCounter++
		}

	// LD - sets register
	case 0x6000:
		registerIndex := (op & 0x0F00) >> 2
		val := op & 0x00FF
		sys.registers[registerIndex] = val

	// ADD - Vx = Vx + val
	case 0x7000:
		registerIndex := (op & 0x0F00) >> 2
		val := op & 0x00FF
		sys.registers[registerIndex] += val

	// Operation between two registers
	case 0x8000:
		registerA := (op & 0x0F00) >> 2
		registerB := (op & 0x00F0) >> 1

		switch op & 0x000F {
			// OR
			case 0x1:
				sys.registers[registerA] |= sys.registers[registerB]

			// AND
			case 0x2:
				sys.registers[registerA] &= sys.registers[registerB]

			// XOR
			case 0x3:
				sys.registers[registerA] ^= sys.registers[registerB]

			// ADD
			case 0x4:
				sum := sys.registers[registerA] + sys.registers[registerB]

				if (sum > sys.registers[registerA]) == (sys.registers[registerB] > 0) {
					sys.registers[registerA] = sum
					sys.registers[0xF] = 0
				} else {
					sys.registers[registerA] = sum
					sys.registers[0xF] = 1
				}

			// SUB
			case 0x5:
				if (sys.registers[registerA] > sys.registers[registerB]) {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] -= sys.registers[registerB]

			// SHR
			case 0x6:
				if (sys.registers[registerA] & 0x1) == 1 {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] /= 2

			// SUBN
			case 0x7:
				if (sys.registers[registerB] > sys.registers[registerA]) {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] = sys.registers[registerB] - sys.registers[registerA]

			// SHL
			case 0xE:
				if (sys.registers[registerA] & 0x1) == 1 {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] *= 2
		}

	// SNE
	case 0x9000:
		registerA := (op & 0x0F00) >> 2
		registerB := (op & 0x00F0) >> 1

		if sys.registers[registerA] != sys.registers[registerB] {
			sys.programCounter += 2
		} else {
			sys.programCounter += 1
		}

	// LD
	case 0xA000:
		sys.iregister = op & 0x0FFF
		sys.programCounter++

	// JMP
	case 0xB000:
		sys.programCounter = (op & 0x0FFF) + sys.registers[0x0]

	// RND
	case 0xC000:
		registerIndex := (op & 0x0F00) >> 2
		val := op & 0x00FF
		sys.registers[registerIndex] = val & byte(rand.Intn(256))
		sys.programCounter++

	// DRW
	case 0xD000:
		// TODO

	case 0xE000:

		switch op & 0x00FF {
		// SKP
		case 0x009E:
			// TODO

		// SKNP
		case 0x00A1:
			// TODO

		}

	case 0xF000:
		switch op & 0x00FF {
		// LD - Load delay timer value into vx
		case 0x0007:
			registerIndex := (op & 0x0F00) >> 2
			sys.registers[registerIndex] = sys.delayTimer
			sys.programCounter++

		// LD - load from input
		case 0x000A:
			// TODO

		// LD - Set delay timer
		case 0x0015:
			registerIndex := (op & 0x0F00) >> 2
			sys.delayTimer = sys.registers[registerIndex]
			sys.programCounter++

		// LD - Set sound timer
		case 0x0018:
			registerIndex := (op & 0x0F00) >> 2
			sys.soundTimer = sys.registers[registerIndex]
			sys.programCounter++

		// ADD - I and Vx
		case 0x001E:
			registerIndex := (op & 0x0F00) >> 2
			sys.iregister += registerIndex
			sys.programCounter++

		// LD - Set I to the value of the location of the sprite
		case 0x0029:
			// TODO

		// LD - Store BCD representation in to I, I+1, I+2
		case 0x0033:
			registerIndex := (op & 0x0F00) >> 2
			val := sys.registers[registerIndex]

			hundred := (val / 100) * 100
			val -= hundred
			ten := (val / 10) * 10
			val -= ten
			one := val

			sys.memory[sys.iregister] = hundred
			sys.memory[sys.iregister + 1] = ten
			sys.memory[sys.iregister + 2] = one

			sys.programCounter++

		// LD - store registers in memory
		case 0x0055:
			limit := (op & 0x0F00) >> 2

			for i := 0; i <= limit; i++ {
				sys.memory[sys.iregister + i] = sys.registers[i]
			}

			sys.programCounter++

		// LD - load register from memory
		case 0x0065:
			limit := (op & 0x0F00) >> 2

			for i := 0; i <= limit; i++ {
				sys.registers[i] = sys.memory[sys.iregister + i]
			}

			sys.programCounter++
		}

	default:
}
