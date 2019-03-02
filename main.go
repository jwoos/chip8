package main


import (
	"fmt"
)


func main() {
	sys := newSystem()
	sys.loadROMFile("IBM Logo.ch8")
	for _, x := range sys.memory[0x200:] {
		fmt.Printf("%X ", x)
	}
}
