package main


import (
	"fmt"
	"math/rand"

	"github.com/nsf/termbox-go"
)


type Instruction struct {
	opcode uint16
	x uint8
	y uint8
	address uint16
	n uint8
	k uint8
}


func (sys *System) parseInstruction() error {
	op := sys.opcode

	x := (op & 0x0F00) >> 8
	y := (op & 0x00F0) >> 4
	n := op & 0x000F
	nnn := op & 0x0FFF
	kk := byte(op & 0x00FF)

	switch op & 0xF000 {
	case 0x0000:
		switch op {
		// CLS - Clear display
		case 0x00E0:
			sys.clearDisplay()
			termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)
			err := termbox.Flush()
			if err != nil {
				return err
			}

			sys.incrementPC(false)
			break

		// RET - return from subroutine
		case 0x00EE:
			addr, err := sys.stack.pop()
			if err != nil {
				return err
			}
			sys.programCounter = addr
			break

		// exit
		case 0x0A00:
			fallthrough
		case 0x0000:
			sys.halt <- true

			sys.incrementPC(false)
			break

		// SYS - jump to machine code routine at address
		// will not implement

		default:
			sys.incrementPC(false)
			return fmt.Errorf("Invalid operation 0x%X", op)
		}

		break

	// JMP - jump to address
	case 0x1000:
		sys.programCounter = nnn
		break

	// CALL - call subroutine
	case 0x2000:
		// FIXME handle error
		sys.stack.push(sys.programCounter)
		sys.programCounter = nnn

		sys.incrementPC(false)
		break

	// SE - Skip next instruction if Vx == val
	case 0x3000:
		if sys.registers[x] == kk {
			sys.incrementPC(true)
		} else {
			sys.incrementPC(false)
		}
		break

	// SNE - skip next instruction if Vx != val
	case 0x4000:
		if sys.registers[x] == kk {
			sys.incrementPC(false)
		} else {
			sys.incrementPC(true)
		}
		break

	// SE - skip if Vx == Vy
	case 0x5000:
		if sys.registers[x] == sys.registers[y] {
			sys.incrementPC(true)
		} else {
			sys.incrementPC(false)
		}
		break

	// LD - sets register
	case 0x6000:
		sys.registers[x] = kk

		sys.incrementPC(false)
		break

	// ADD - Vx = Vx + val
	case 0x7000:
		sys.registers[x] += kk

		sys.incrementPC(false)
		break

	// Operation between two registers
	case 0x8000:
		switch op & 0x000F {
			// OR
			case 0x1:
				sys.registers[x] |= sys.registers[y]

				sys.incrementPC(false)
				break

			// AND
			case 0x2:
				sys.registers[x] &= sys.registers[y]

				sys.incrementPC(false)
				break

			// XOR
			case 0x3:
				sys.registers[x] ^= sys.registers[y]

				sys.incrementPC(false)
				break

			// ADD
			case 0x4:
				sum := sys.registers[x] + sys.registers[y]
				if (sum > sys.registers[x]) == (sys.registers[y] > 0) {
					sys.registers[0xF] = 0
				} else {
					sys.registers[0xF] = 1
				}
				sys.registers[x] = sum

				sys.incrementPC(false)
				break

			// SUB
			case 0x5:
				if (sys.registers[x] > sys.registers[y]) {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[x] -= sys.registers[y]

				sys.incrementPC(false)
				break

			// SHR
			case 0x6:
				if (sys.registers[x] & 0x1) == 1 {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[x] >>= 1

				sys.incrementPC(false)
				break

			// SUBN
			case 0x7:
				if (sys.registers[y] > sys.registers[x]) {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[x] = sys.registers[y] - sys.registers[x]

				sys.incrementPC(false)
				break

			// SHL
			case 0xE:
				if (sys.registers[x] & 0x1) == 1 {
					sys.registers[0xF] = 1
				} else {
					sys.registers[0xF] = 0
				}

				sys.registers[x] <<= 1

				sys.incrementPC(false)
				break

			default:
				return fmt.Errorf("Invalid operation 0x%04X", op)
		}
		break

	// SNE
	case 0x9000:
		if sys.registers[x] != sys.registers[y] {
			sys.incrementPC(true)
		} else {
			sys.incrementPC(false)
		}
		break

	// LD
	case 0xA000:
		sys.iregister = nnn

		sys.incrementPC(false)
		break

	// JMP
	case 0xB000:
		sys.programCounter = nnn + uint16(sys.registers[0x0])
		break

	// RND
	case 0xC000:
		sys.registers[x] = kk & byte(rand.Intn(256))

		sys.incrementPC(false)
		break

	// DRW
	case 0xD000:
		sys.registers[0xF] = 0

		for i := uint16(0); i < n; i++ {
			toDraw := sys.memory[sys.iregister + i]
			toDrawBits, err := bits(toDraw)
			if err != nil {
				return err
			}

			cells := termbox.CellBuffer()
			width, _ := termbox.Size()

			for j := uint16(0); j < uint16(len(toDrawBits)); j++ {
				prev := sys.display[y + i][x + j]
				sys.display[y + i][x + j] = sys.display[y + i][x + j] != toDrawBits[j]

				if (prev == true) && (sys.display[y + i][x + j] != true) {
					sys.registers[0xF] = 1
				}

				if sys.display[y + i][x + j] {
					cells[(uint16(width) * (y + i)) + (x + j)].Ch = '█'
				} else {
					cells[(uint16(width) * (y + i)) + (x + j)].Ch = ' '
				}
			}
		}

		termbox.Flush()
		sys.incrementPC(false)
		break

	case 0xE000:

		switch op & 0x00FF {
		// SKP
		case 0x009E:
			// TODO
			panic("Instruction not implemented 0xEx9E")

		// SKNP
		case 0x00A1:
			// TODO
			panic("Instruction not implemented 0xExA1")

		default:
			return fmt.Errorf("Invalid operation 0x%04X", op)
		}

	case 0xF000:
		switch op & 0x00FF {
		// LD - Load delay timer value into vx
		case 0x0007:
			sys.registers[x] = sys.delayTimer

			sys.incrementPC(false)
			break

		// LD - load from input
		case 0x000A:
			ev := termbox.PollEvent()
			for {
				if ev.Type == termbox.EventKey {
					val, ok := INPUT_MAP[ev.Ch]

					if ok {
						sys.registers[x] = val
						break
					}
				}
			}

			sys.incrementPC(false)
			break

		// LD - Set delay timer
		case 0x0015:
			sys.delayTimer = sys.registers[x]

			sys.incrementPC(false)
			break

		// LD - Set sound timer
		case 0x0018:
			sys.soundTimer = sys.registers[x]

			sys.incrementPC(false)
			break

		// ADD - I and Vx
		case 0x001E:
			sys.iregister += x

			sys.incrementPC(false)
			break

		// LD - Set I to the value of the location of the sprite
		case 0x0029:
			sys.iregister = uint16(sys.registers[x] * 5)

			sys.incrementPC(false)
			break

		// LD - Store BCD representation in to I, I+1, I+2
		case 0x0033:
			val := sys.registers[x]

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
			for i := uint16(0); i <= x; i++ {
				sys.memory[sys.iregister + i] = sys.registers[i]
			}

			sys.incrementPC(false)
			break

		// LD - load register from memory
		case 0x0065:
			for i := uint16(0); i <= x; i++ {
				sys.registers[i] = sys.memory[sys.iregister + i]
			}

			sys.incrementPC(false)
			break

		default:
			return fmt.Errorf("Invalid operation 0x%04X", op)
		}
		break

	default:
		return fmt.Errorf("Invalid operation 0x%04X", op)
	}

	return nil
}
