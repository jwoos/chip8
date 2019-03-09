# CHIP8 emulator and disassembler
A CHIP8 emulator and disassembler. Everything is implemented except for sound!

## Building
You should have a Go version that supports modules (> 1.11). Just clone and run `go build`.

## Running
Run a ROM by doing:
```
$ ./go_chip8 --rom <PATH_TO_ROM>
```

You can also disassemble it instead of running by doing:
```
$ ./go_chip8 --rom <PATH_TO_ROM> --disassemble
```

### Clock speed
You can determine the clock speed of the emulator (the default is 500 Hz), up to 1000 Hz by using `--clockspeed`. This does not affect the timers as they will decrement at a steady rate of 60 Hz, thanks to being in their own goroutine.

### Key timeout
Due to running in a terminal, it's impossible to detect whether a key is being held down. That's what the key timeout is for. It will leave a key "pressed" for that number of milliseconds. One thing to note is that, instructions that read input will reset key presses.

### Keyboard
The keyboard mapping is as follows:
```
1 | 2 | 3 | C               1 | 2 | 3 | 4
4 | 5 | 6 | D               Q | W | E | R
7 | 8 | 9 | E               A | S | D | F
A | 0 | B | F               Z | X | C | V
```

## References
- http://devernay.free.fr/hacks/chip8/C8TECH10.HTM
- https://en.wikipedia.org/wiki/CHIP-8
