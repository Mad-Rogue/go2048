package main

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/eiannone/keyboard"
)

func GetKey() (ch rune, err error) {
	if err = keyboard.Open(); err != nil {
		return
	}
	defer keyboard.Close()

	var (
		chChan  = make(chan rune, 1)
		errChan = make(chan error, 1)
	)
	go func(chChan chan<- rune, errChan chan<- error) {
		ch, _, err := keyboard.GetSingleKey()
		if err != nil {
			errChan <- err
		}
		chChan <- ch
	}(chChan, errChan)

	select {
	case ch = <-chChan:
	case err = <-errChan:
	}

	return
}

var clear map[string]func()

func init() {
	clear = make(map[string]func())
	clear["linux"] = func() {
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["darwin"] = func() {
		fmt.Print("\033[H\033[2J")
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls")
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}

func ClearScreen() {
	value, ok := clear[runtime.GOOS]
	if ok {
		value()
	} else { //unsupported platform, last chance
		fmt.Print("\033[H\033[2J")
	}
}

func Max(x, y byte) byte {
	if x < y {
		return y
	}
	return x
}
func Min(x, y byte) byte {
	if y < x {
		return y
	}
	return x
}
func Center(s string, w int) string {
	return fmt.Sprintf("%[1]*s", -w, fmt.Sprintf("%[1]*s", (w+len(s))/2, s))
}
