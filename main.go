package main

import (
	"fmt"
	"log"
)

var field = Field{FieldSize, nil, nil}
var finished = false

func reset() {
	finished = false
	field.Init(FieldSize)
	field.AddRandomValues(2)
}

func main() {
	fmt.Println("Init...")
	reset()
	for {
		field.Draw()
		fmt.Println("\nUse 'wasd' keys for control. Press 'q' for exit. Press 'r' for reset.")
		if finished {
			fmt.Println("\nGame Over!")
		}
		key, err := GetKey()
		if err != nil {
			log.Print(err)
			return
		}
		switch key {
		case 'r', 'R':
			reset()
			break
		case 'q', 'Q':
			return
		}
		if !finished {
			addvalues := false
			switch key {
			case 'd', 'D':
				addvalues = field.SlideTo(TO_RIGHT)
				break
			case 'a', 'A':
				addvalues = field.SlideTo(TO_LEFT)
				break
			case 'w', 'W':
				addvalues = field.SlideTo(TO_UP)
				break
			case 's', 'S':
				addvalues = field.SlideTo(TO_DOWN)
				break
			}
			if addvalues {
				added, hasEmpty := field.AddRandomValues(1)
				if added == 0 || !hasEmpty && !HasAvailableSteps(field) {
					finished = true
				}
			}
		}
	}
}
