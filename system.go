package main


import (
	"errors"
	"fmt"
	"io/ioutil"
	"math/rand"

	"github.com/nsf/termbox-go"
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
	memory []byte

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
	sys.memory = make([]byte, 4096)
	sys.registers = make([]byte, 16)
	sys.stack = newStack(16)
	sys.programCounter = 0x200

	return sys
}

func (sys *System) incrementPC(skip bool) {
	if !skip {
		sys.programCounter += 2
	} else {
		sys.programCounter += 4
	}
}

func (sys *System) loadROMFile(path string) error {
	data, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	sys.loadROM(data)

	return nil
}

func (sys *System) loadROM(data []byte) {
	for i, b := range data {
		sys.memory[0x200 + i] = b
	}
}

func (sys *System) parseInstruction() error {
	op := sys.opcode

	switch op & 0xF000 {
	case 0x0000:
		switch op {
		// CLS - Clear display
		case 0x00E0:
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			err := termbox.Flush()
			if err != nil {
				return err
			}

			sys.incrementPC(false)
			break

		// RET - return from subroutine
		case 0x00EE:

		// SYS - jump to machine code routine at address
		default:

		}

		break

	// JMP - jump to address
	case 0x1000:
		sys.programCounter = op & 0x0FFF

		sys.incrementPC(false)
		break

	// CALL - call subroutine
	case 0x2000:
		// FIXME handle error
		sys.stack.push(sys.programCounter)
		sys.programCounter = op & 0x0FFF

		sys.incrementPC(false)
		break

	// SE - Skip next instruction if Vx == val
	case 0x3000:
		registerIndex := (op & 0x0F00) >> 0x2
		lastHalf := byte(op & 0x00FF)
		if sys.registers[registerIndex] == lastHalf {
			sys.incrementPC(true)
		} else {
			sys.incrementPC(false)
		}
		break

	// SNE - skip next instruction if Vx != val
	case 0x4000:
		registerIndex := (op & 0x0F00) >> 0x2
		lastHalf := byte(op & 0x00FF)
		if sys.registers[registerIndex] == lastHalf {
			sys.incrementPC(false)
		} else {
			sys.incrementPC(true)
		}
		break

	// SE - skip if Vx == Vy
	case 0x5000:
		registerA := (op & 0x0F00) >> 0x2
		registerB := (op & 0x00F0) >> 0x1

		if sys.registers[registerA] == sys.registers[registerB] {
			sys.incrementPC(true)
		} else {
			sys.incrementPC(false)
		}
		break

	// LD - sets register
	case 0x6000:
		registerIndex := (op & 0x0F00) >> 0x2
		val := byte(op & 0x00FF)
		sys.registers[registerIndex] = val

		sys.incrementPC(false)
		break

	// ADD - Vx = Vx + val
	case 0x7000:
		registerIndex := (op & 0x0F00) >> 0x2
		val := byte(op & 0x00FF)
		sys.registers[registerIndex] += val

		sys.incrementPC(false)
		break

	// Operation between two registers
	case 0x8000:
		registerA := (op & 0x0F00) >> 0x2
		registerB := (op & 0x00F0) >> 0x1

		switch op & 0x000F {
			// OR
			case 0x1:
				sys.registers[registerA] |= sys.registers[registerB]

				sys.incrementPC(false)
				break

			// AND
			case 0x2:
				sys.registers[registerA] &= sys.registers[registerB]

				sys.incrementPC(false)
				break

			// XOR
			case 0x3:
				sys.registers[registerA] ^= sys.registers[registerB]

				sys.incrementPC(false)
				break

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

				sys.incrementPC(false)
				break

			// SUB
			case 0x5:
				if (sys.registers[registerA] > sys.registers[registerB]) {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] -= sys.registers[registerB]

				sys.incrementPC(false)
				break

			// SHR
			case 0x6:
				if (sys.registers[registerA] & 0x1) == 1 {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] /= 2

				sys.incrementPC(false)
				break

			// SUBN
			case 0x7:
				if (sys.registers[registerB] > sys.registers[registerA]) {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] = sys.registers[registerB] - sys.registers[registerA]

				sys.incrementPC(false)
				break

			// SHL
			case 0xE:
				if (sys.registers[registerA] & 0x1) == 1 {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[registerA] *= 2

				sys.incrementPC(false)
				break
		}
		break

	// SNE
	case 0x9000:
		registerA := (op & 0x0F00) >> 0x2
		registerB := (op & 0x00F0) >> 0x1

		if sys.registers[registerA] != sys.registers[registerB] {
			sys.incrementPC(true)
		} else {
			sys.incrementPC(false)
		}
		break

	// LD
	case 0xA000:
		sys.iregister = op & 0x0FFF

		sys.incrementPC(false)
		break

	// JMP
	case 0xB000:
		sys.programCounter = (op & 0x0FFF) + uint16(sys.registers[0x0])
		break

	// RND
	case 0xC000:
		registerIndex := (op & 0x0F00) >> 0x2
		val := byte(op & 0x00FF)
		sys.registers[registerIndex] = val & byte(rand.Intn(256))

		sys.incrementPC(false)
		break

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
			registerIndex := (op & 0x0F00) >> 0x2
			sys.registers[registerIndex] = sys.delayTimer

			sys.incrementPC(false)
			break

		// LD - load from input
		case 0x000A:
			// TODO

		// LD - Set delay timer
		case 0x0015:
			registerIndex := (op & 0x0F00) >> 0x2
			sys.delayTimer = sys.registers[registerIndex]

			sys.incrementPC(false)
			break

		// LD - Set sound timer
		case 0x0018:
			registerIndex := (op & 0x0F00) >> 0x2
			sys.soundTimer = sys.registers[registerIndex]

			sys.incrementPC(false)
			break

		// ADD - I and Vx
		case 0x001E:
			registerIndex := (op & 0x0F00) >> 0x2
			sys.iregister += registerIndex

			sys.incrementPC(false)
			break

		// LD - Set I to the value of the location of the sprite
		case 0x0029:
			// TODO

		// LD - Store BCD representation in to I, I+1, I+2
		case 0x0033:
			registerIndex := (op & 0x0F00) >> 0x2
			val := sys.registers[registerIndex]

			hundred := (val / 100) * 100
			val -= hundred
			ten := (val / 10) * 10
			val -= ten
			one := val

			sys.memory[sys.iregister] = hundred
			sys.memory[sys.iregister + 1] = ten
			sys.memory[sys.iregister + 2] = one

			sys.incrementPC(false)
			break

		// LD - store registers in memory
		case 0x0055:
			limit := (op & 0x0F00) >> 0x2

			for i := uint16(0); i <= limit; i++ {
				sys.memory[sys.iregister + i] = sys.registers[i]
			}

			sys.incrementPC(false)
			break

		// LD - load register from memory
		case 0x0065:
			limit := (op & 0x0F00) >> 0x2

			for i := uint16(0); i <= limit; i++ {
				sys.registers[i] = sys.memory[sys.iregister + i]
			}

			sys.incrementPC(false)
			break
		}
		break

	default:
		return errors.New(fmt.Sprintf("Invalid operation: %x", op))
	}

	return nil
}
