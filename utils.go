package main

import (
	"strconv"
	"fmt"
)


func bits(num interface{}) ([]bool, error) {
	var size int
	var concreteInt int64
	var concreteUint uint64
	var flag bool

	switch num.(type) {
	case int:
		size = strconv.IntSize
		concreteInt = int64(num.(int))
		flag = true
	case uint:
		size = strconv.IntSize
		concreteUint = uint64(num.(uint))
		flag = false

	case int8:
		size = 8
		concreteInt = int64(num.(int8))
		flag = true
	case uint8:
		size = 8
		concreteUint = uint64(num.(uint8))
		flag = false

	case int16:
		size = 16
		concreteInt = int64(num.(int16))
		flag = true
	case uint16:
		size = 16
		concreteUint = uint64(num.(uint16))
		flag = false

	case int32:
		size = 32
		concreteInt = int64(num.(int32))
		flag = true
	case uint32:
		size = 32
		concreteUint = uint64(num.(uint32))
		flag = false

	case int64:
		size = 64
		concreteInt = num.(int64)
		flag = true
	case uint64:
		size = 64
		concreteUint = num.(uint64)
		flag = false

	default:
		return nil, fmt.Errorf("Invalid type")
	}

	data := make([]bool, size)

	if flag {
		for i := 0; i < size; i++ {
			data[i] = ((concreteInt >> uint(size - 1 - i)) & 0x1) == 0x1
		}
	} else {
		for i := 0; i < size; i++ {
			data[i] = ((concreteUint >> uint(size - 1 - i)) & 0x1) == 0x1
		}
	}

	return data, nil
}
