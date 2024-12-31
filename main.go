package main

import (
	"flag"
	"github.com/jeffreybian/mines/ui"
)

func parseCliArgs() (int, int, int) {
	rows, cols, numberOfMines := 0, 0, 0
	flag.IntVar(&rows, "r", 9, "Number of rows")
	flag.IntVar(&cols, "c", 9, "Number of columns")
	flag.IntVar(&numberOfMines, "m", 10, "Number of mines")
	flag.Parse()
	return rows, cols, numberOfMines
}

func main() {
	rows, cols, numberOfMines := parseCliArgs()
	ui.MainGameWindow.InitSystem(rows, cols, numberOfMines)
	ui.MainGameWindow.Loop()
	ui.MainGameWindow.CleanUp()
}
