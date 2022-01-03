package main

import (
	"fmt"
	"os"

	"github.com/kris10ansn/kwin-desktop-switcher/pkg/switcher"
)

func main() {
	direction := os.Args[1]

	context := &switcher.GridContext{Width: 2, Height: 2}
	err := switcher.SwitchDesktop(direction, context)

	if err != nil {
		fmt.Println(err)
	}
}
