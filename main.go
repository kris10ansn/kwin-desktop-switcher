package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

type GridContext struct {
	Width  int
	Height int
}

type Direction = string

const (
	Left  Direction = "left"
	Right Direction = "right"
	Up    Direction = "up"
	Down  Direction = "down"
)

func gridCoordinatesFromDesktopNumber(desktopNumber int, gridContext *GridContext) (x int, y int) {
	x = (desktopNumber - 1) % gridContext.Width
	y = (desktopNumber - 1) / gridContext.Height

	return (x + 1), (y + 1)
}

func getDesktopNumber() (int, error) {
	output, err := exec.Command("qdbus", "org.kde.KWin", "/KWin", "currentDesktop").Output()

	if err != nil {
		return -1, err
	}

	outputString := string(output)
	desktopNumber, err := strconv.Atoi(strings.Trim(outputString, "\n"))

	if err != nil {
		return -1, err
	}

	return desktopNumber, nil
}

func setDesktopNumber(n int) error {
	_, err := exec.Command("qdbus", "org.kde.KWin", "/KWin", "setCurrentDesktop", fmt.Sprint(n)).Output()

	return err
}

func canSwitch(direction Direction, desktopNumber int, gridContext *GridContext) bool {
	x, y := gridCoordinatesFromDesktopNumber(desktopNumber, gridContext)

	if direction == Right && x == gridContext.Width {
		return false
	}

	if direction == Left && x == 0 {
		return false
	}

	if direction == Down && y == gridContext.Height {
		return false
	}

	if direction == Up && y == 0 {
		return false
	}

	return true
}

func getDesktopModifier(direction Direction) (int, error) {
	switch direction {
	case Right:
		return 1, nil
	case Left:
		return -1, nil
	case Down:
		return 2, nil
	case Up:
		return -2, nil
	default:
		return 0, (fmt.Errorf("Unknown direction %s", direction))
	}
}

func SwitchDesktop(direction Direction, gridContext *GridContext) error {
	desktopNumber, err := getDesktopNumber()
	if err != nil {
		return err
	}

	modifier, err := getDesktopModifier(direction)
	if err != nil {
		return err
	}

	if canSwitch(direction, desktopNumber, gridContext) {
		newDesktopNumber := desktopNumber + modifier
		err := setDesktopNumber(newDesktopNumber)

		if err != nil {
			return err
		}
	} else {
		x, y := gridCoordinatesFromDesktopNumber(desktopNumber, gridContext)
		return fmt.Errorf(
			"Can't switch %s from desktop number %d (%d, %d)",
			direction, desktopNumber, x, y,
		)
	}

	return nil
}

func main() {
	direction := os.Args[1]

	context := &GridContext{Width: 2, Height: 2}
	err := SwitchDesktop(direction, context)

	if err != nil {
		fmt.Println(err)
	}
}
