package main

import (
	"errors"
	"fmt"
)

type Stack struct {
	memory []uint16
	index uint
	capacity uint
}


func newStack(size uint) *Stack {
	stack := new(Stack)
	stack.memory = make([]uint16, size)
	stack.capacity = size
	return stack
}


func (stack *Stack) push(item uint16) error {
	if stack.index == stack.capacity {
		return errors.New("Stack is full")
	}

	stack.memory[stack.index] = item
	stack.index++

	return nil
}


func (stack *Stack) pop() (uint16, error) {
	if stack.index == 0 {
		return 0, errors.New("Stack is empty")
	}

	stack.index--
	retrieved := stack.memory[stack.index]
	stack.memory[stack.index] = 0

	return retrieved, nil
}


func (stack *Stack) peek() (uint16, error) {
	if stack.index == 0 {
		return 0, errors.New("Stack is empty")
	}

	return stack.memory[stack.index - 1], nil
}


func (stack *Stack) size() uint {
	return uint(len(stack.memory))
}

func (stack *Stack) String() string {
	return fmt.Sprintf("%04X", stack.memory)
}
