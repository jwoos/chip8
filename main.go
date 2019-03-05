package main


import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)


func main() {
	sys := newSystem(500, false)
	sys.loadFont()
	sys.loadROMFile("IBM Logo.ch8")
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
