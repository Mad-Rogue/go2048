package main

import (
	"fmt"
	"math"
	"math/rand"
)

type Field struct {
	size      int
	field     []byte
	processed []bool
}

type Point struct {
	x, y  int
	value byte
}

var strightDirect = func(value, _ int) int {
	return value
}

var invertDirect = func(value, max int) int {
	return max - value - 1
}

var xPicker = func(value, _ int) int {
	return value
}
var yPicker = func(_, value int) int {
	return value
}

func realCoordsFnc(direct int) func(int, int, int) (int, int) {
	xPickerFnc := xPicker
	yPickerFnc := yPicker

	switch direct {
	case TO_UP:
		xPickerFnc = yPicker
		yPickerFnc = xPicker
		break
	case TO_DOWN:
		xPickerFnc = yPicker
		yPickerFnc = xPicker
		break
	}

	xDirectFnc := strightDirect

	switch direct {
	case TO_DOWN:
		xDirectFnc = invertDirect
		break
	case TO_RIGHT:
		xDirectFnc = invertDirect
		break
	}

	return func(x, y, max int) (int, int) {
		fx := xDirectFnc(x, max)
		return xPickerFnc(fx, y), yPickerFnc(fx, y)
	}
}

func removeAt(s []int, i int) []int {
	s[i] = s[len(s)-1]
	return s[:len(s)-1]
}

func (field *Field) clone() Field {
	result := Field{0, nil, nil}
	result.Init(field.size)
	for i, value := range field.field {
		result.field[i] = value
	}
	return result
}

func (field *Field) get(x, y int) byte {
	return field.field[x+y*field.size]
}

func (field *Field) set(x, y int, value byte) {
	field.field[x+y*field.size] = value
}

func (field *Field) isProcessed(x, y int) bool {
	return field.processed[x+y*field.size]
}

func (field *Field) setProcessed(x, y int) {
	field.processed[x+y*field.size] = true
}

func (field *Field) prepare() {
	for i := 0; i < len(field.processed); i++ {
		field.processed[i] = false
	}
}

func (field *Field) Draw() {
	ClearScreen()
	var text string
	for y := 0; y < field.size; y++ {
		for x := 0; x < field.size; x++ {
			value := field.get(x, y)

			if value > 0 {
				text = fmt.Sprint(int(math.Pow(2, float64(value))))
			} else {
				text = "."
			}

			fmt.Print("[", Center(text, 5), "]")
		}
		fmt.Print("\n")
	}
}

func (field *Field) Init(size int) {
	field.size = size
	field.field = make([]byte, size*size)
	field.processed = make([]bool, size*size)
}

func (field *Field) getMaxValue() byte {
	var result byte = 0
	for _, value := range field.field {
		result = Max(value, result)
	}
	return result
}

func (field *Field) AddRandomValues(count int) (int, bool) {
	result := 0

	emptyCells := []int{}

	for i, cell := range field.field {
		if cell == 0 {
			emptyCells = append(emptyCells, i)
		}
	}

	for i := 0; i < count; i++ {
		if len(emptyCells) == 0 {
			break
		}
		index := rand.Intn(len(emptyCells))
		cellIndex := emptyCells[index]
		emptyCells = removeAt(emptyCells, index)

		maxValue := int(Min(MaxRandomValue, field.getMaxValue()))
		if maxValue == 0 {
			field.field[cellIndex] = 1
		} else {
			field.field[cellIndex] = byte(rand.Intn(maxValue)) + 1
		}
		result++
	}
	return result, len(emptyCells) > 0
}

func (field *Field) SlideTo(direct int) (result bool) {
	result = false
	coordsFnc := realCoordsFnc(direct)
	field.prepare()

	for y := 0; y < field.size; y++ {
		for x := 1; x < field.size; x++ {
			tx, ty := coordsFnc(x, y, field.size)
			value := field.get(tx, ty)
			point := Point{tx, ty, value}

			if point.value > 0 {
				target := Point{-1, -1, 0}
				for x2 := x - 1; x2 >= 0; x2-- {
					tx, ty := coordsFnc(x2, y, field.size)
					value := field.get(tx, ty)
					rpoint := Point{tx, ty, value}

					if value == 0 {
						target = rpoint
					} else if point.value == value && !field.isProcessed(tx, ty) {
						target = rpoint
						break
					} else {
						break
					}
				}
				if target.x != -1 && target.y != -1 {
					result = true
					if target.value == 0 {
						field.set(target.x, target.y, value)
					} else if point.value == target.value {
						field.set(target.x, target.y, value+1)
						field.setProcessed(target.x, target.y)
					} else {
						panic("oooops! it's impssible!")
					}

					field.set(point.x, point.y, 0)
				}
			}
		}
	}
	return result
}

func HasAvailableSteps(field Field) bool {
	result := false
	directs := []int{TO_RIGHT, TO_LEFT, TO_UP, TO_DOWN}
	for direct := range directs {
		copy := field.clone()
		if copy.SlideTo(direct) {
			result = true
			break
		}
	}
	return result
}
