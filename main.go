package main


import (
	"fmt"
	"flag"
	"time"
	"os"

	"github.com/nsf/termbox-go"
)


func main() {
	var clockspeed uint64
	var debug bool
	var rom string

	flag.Uint64Var(&clockspeed, "clockspeed", 500, "Clockspeed in Hz")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.StringVar(&rom, "rom", "", "ROM to run")
	flag.Parse()

	if rom == "" {
		fmt.Println("Please supply a ROM")
		os.Exit(1)
	}

	sys := newSystem(clockspeed, debug)
	sys.loadFont()
	sys.loadROMFile(rom)
	sys.timers()

	err := termbox.Init()
	if err != nil {
		fmt.Printf("Error initializing termbox: %v\n", err)
	}
	termbox.HideCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for range time.Tick(time.Duration(1000 / sys.clockspeed) * time.Millisecond) {
		sys.readInstruction()
		sys.parseInstruction()
		//fmt.Printf("0x%X\n", sys.opcode)
		/*
		 *if err != nil {
		 *    fmt.Println(err)
		 *    break
		 *}
		 */

		if sys.halt {
			break
		}
	}

	for i, ch := range "Press esc to quit" {
		termbox.SetCell(i, 0, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break
		}
	}

	termbox.Close()
}
