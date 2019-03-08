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
	var disassemble bool
	var debug bool
	var rom string
	var keyTimeOut uint

	flag.Uint64Var(&clockspeed, "clockspeed", 500, "Clockspeed in Hz")
	flag.BoolVar(&debug, "debug", false, "Debug mode")
	flag.BoolVar(&disassemble, "disassemble", false, "Disassemble ROM")
	flag.StringVar(&rom, "rom", "", "ROM to run")
	flag.UintVar(&keyTimeOut, "keytimeout", 100, "Key presses are held this amount of milliseconds")
	flag.Parse()

	if rom == "" {
		fmt.Println("Please supply a ROM")
		os.Exit(1)
	}

	sys := newSystem(clockspeed, keyTimeOut, debug)
	sys.loadFont()
	sys.loadROMFile(rom)

	if disassemble {
		sys.disassemble()
		return
	}

	err := termbox.Init()
	if err != nil {
		fmt.Printf("Error initializing termbox: %v\n", err)
	}
	termbox.HideCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	sys.keyEvents()
	sys.timers()

clockLoop:
	for range time.Tick(time.Duration(1000 / sys.clockspeed) * time.Millisecond) {
		select {
		case <-sys.halt:
			break clockLoop
		default:
			sys.readInstruction()
			err := sys.parseInstruction()
			if err != nil {
				fmt.Printf("Error parsing instruction: %v", err)
				break clockLoop
			}
		}
	}

	for i, ch := range "Press any key to quit" {
		termbox.SetCell(i, 0, ch, termbox.ColorDefault, termbox.ColorDefault)
	}
	termbox.Flush()

	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey {
			break
		}
	}

	termbox.Close()
}
