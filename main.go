package main


import (
	"fmt"
)


func main() {
	sys := newSystem()
	sys.loadROMFile("IBM Logo.ch8")
	fmt.Println(sys.memory[0x200:])
}
