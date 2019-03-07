package main


import (
	"io/ioutil"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	MEMORY_SIZE = 4096
	REGISTER_COUNT = 16
	STACK_SIZE = 16
	PC_START = 0x200
	DISPLAY_WIDTH = 64
	DISPLAY_HEIGHT = 32
)

var INPUT_MAP = map[rune]byte{
	'1': 0x1,
	'2': 0x2,
	'3': 0x3,
	'4': 0xC,
	'q': 0x4,
	'w': 0x5,
	'e': 0x6,
	'r': 0xD,
	'a': 0x7,
	's': 0x8,
	'd': 0x9,
	'f': 0xE,
	'z': 0xA,
	'x': 0x0,
	'c': 0xB,
	'v': 0xF,
}


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

	display [][]bool

	halt chan bool

	// Hz
	clockspeed uint64

	debug bool
}


func newSystem(clockspeed uint64, debug bool) *System {
	sys := new(System)
	sys.memory = make([]byte, MEMORY_SIZE)
	sys.registers = make([]byte, REGISTER_COUNT)
	sys.stack = newStack(STACK_SIZE)
	sys.programCounter = PC_START
	sys.clockspeed = clockspeed

	sys.halt = make(chan bool, 1)

	sys.display = make([][]bool, DISPLAY_HEIGHT)
	for i := 0; i < len(sys.display); i++ {
		sys.display[i] = make([]bool, DISPLAY_WIDTH)
	}

	return sys
}

func (sys *System) incrementPC(skip bool) {
	if !skip {
		sys.programCounter += 2
	} else {
		sys.programCounter += 4
	}
}

// run in a goroutine
func (sys *System) timers() {
	go func() {
		for range time.Tick(time.Duration(1000 / 60) * time.Millisecond) {
			if sys.soundTimer > 0 {
				sys.soundTimer--
			}

			if sys.delayTimer > 0 {
				sys.delayTimer--
			}
		}
	}()
}

func (sys *System) keyEvents() {
	go func() {
		for {
			ev := termbox.PollEvent()
			if ev.Type == termbox.EventKey && ev.Key == termbox.KeyCtrlC {
				sys.halt <- true
				return
			}
		}
	}()
}

func (sys *System) write(val string) {
	for i, ch := range val {
		termbox.SetCell(i, 0, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()
}

func (sys *System) loadFont() error {
	fonts := []byte{
		0xF0, 0x90, 0x90, 0x90, 0xF0,
		0x20, 0x60, 0x20, 0x20, 0x70,
		0xF0, 0x10, 0xF0, 0x80, 0xF0,
		0xF0, 0x10, 0xF0, 0x10, 0xF0,
		0x90, 0x90, 0xF0, 0x10, 0x10,
		0xF0, 0x80, 0xF0, 0x10, 0xF0,
		0xF0, 0x80, 0xF0, 0x90, 0xF0,
		0xF0, 0x10, 0x20, 0x40, 0x40,
		0xF0, 0x90, 0xF0, 0x90, 0xF0,
		0xF0, 0x90, 0xF0, 0x10, 0xF0,
		0xF0, 0x90, 0xF0, 0x90, 0x90,
		0xE0, 0x90, 0xE0, 0x90, 0xE0,
		0xF0, 0x80, 0x80, 0x80, 0xF0,
		0xE0, 0x90, 0x90, 0x90, 0xE0,
		0xF0, 0x80, 0xF0, 0x80, 0xF0,
		0xF0, 0x80, 0xF0, 0x80, 0x80,
	}

	for i, x := range fonts {
		sys.memory[i] = x
	}

	return nil
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
		sys.memory[PC_START + i] = b
	}
}

func (sys *System) readInstruction() {
	sys.opcode = (uint16(sys.memory[sys.programCounter]) << 8) | uint16(sys.memory[sys.programCounter + 1])
}

func (sys *System) clearDisplay() {
	for i := 0; i < len(sys.display); i++ {
		for j := 0; j < len(sys.display[i]); j++ {
			sys.display[i][j] = false
		}
	}
}
