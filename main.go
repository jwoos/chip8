package main


import (
	"fmt"
	"time"

	"github.com/nsf/termbox-go"
)


func main() {
	sys := newSystem()
	sys.loadFont()
	sys.loadROMFile("test-bin.ch8")

	err := termbox.Init()
	if err != nil {
		fmt.Printf("Error initializing termbox: %v\n", err)
	}
	termbox.HideCursor()
	termbox.Clear(termbox.ColorDefault, termbox.ColorDefault)

	for range time.Tick(time.Duration(1000 / 60) * time.Millisecond) {
		sys.readInstruction()
		err := sys.parseInstruction()
		if err != nil {
			fmt.Println(err)
			break
		}

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
