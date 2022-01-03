package main

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

const (
	GridWidth  = 2
	GridHeight = 2

	Left  = "left"
	Right = "right"
	Up    = "up"
	Down  = "down"
)

func gridCoordinatesFromDesktopNumber(desktopNumber int) (x int, y int) {
	x = (desktopNumber - 1) % GridWidth
	y = (desktopNumber - 1) / GridHeight

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

func canSwitch(direction string, desktopNumber int) bool {
	x, y := gridCoordinatesFromDesktopNumber(desktopNumber)

	if direction == Right && x == GridWidth {
		return false
	}

	if direction == Left && x == 0 {
		return false
	}

	if direction == Down && y == GridHeight {
		return false
	}

	if direction == Up && y == 0 {
		return false
	}

	return true
}

func getDesktopModifier(direction string) (int, error) {
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

func switchDesktop(direction string) error {
	desktopNumber, err := getDesktopNumber()
	if err != nil {
		return err
	}

	modifier, err := getDesktopModifier(direction)
	if err != nil {
		return err
	}

	if canSwitch(direction, desktopNumber) {
		newDesktopNumber := desktopNumber + modifier
		err := setDesktopNumber(newDesktopNumber)

		if err != nil {
			return err
		}
	} else {
		x, y := gridCoordinatesFromDesktopNumber(desktopNumber)
		return fmt.Errorf(
			"Can't switch %s from desktop number %d (%d, %d)",
			direction, desktopNumber, x, y,
		)
	}

	return nil
}

func main() {
	direction := os.Args[1]
	err := switchDesktop(direction)

	if err != nil {
		fmt.Println(err)
	}
}
