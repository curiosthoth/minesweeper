package main

import (
	"github.com/jeffreybian/mines/ui"
)

func main() {
	ui.MainGameWindow.InitSystem()
	ui.MainGameWindow.Loop()
	ui.MainGameWindow.CleanUp()
}
