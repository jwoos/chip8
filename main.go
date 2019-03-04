package main


import (
	"fmt"

	"github.com/nsf/termbox-go"
)


func main() {
	sys := newSystem()
	sys.loadFont()
	sys.loadROMFile("test-bin.ch8")
	for _, x := range sys.memory[0x200:] {
		fmt.Printf("%X ", x)
	}

	err := termbox.Init()
	if err != nil {
		fmt.Printf("Error initializing termbox: %v\n", err)
	}
	termbox.HideCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	/*
	 *for i, ch := range "Press esc to quit" {
	 *    termbox.SetCell(i, 0, ch, termbox.ColorDefault, termbox.ColorDefault)
	 *}
	 */

	for i := 0; i < 2; i++ {
		sys.readInstruction()
		sys.parseInstruction()
	}
	for {
		ev := termbox.PollEvent()
		if ev.Type == termbox.EventKey && ev.Key == termbox.KeyEsc {
			break
		}
	}

	termbox.Close()
}
