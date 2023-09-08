package game

import "github.com/veandco/go-sdl2/sdl"

const DefaultFps = 30
const LowestFps = 10

type App struct {
	fps              uint32
	windows          map[*RenderWindow]RenderListener
	timeStepMs       uint64
	renderFrameDelta float32
	quit             bool
	vsync            bool
	appListener      AppListener
	initialized      bool
}

func NewApp(appListener AppListener) *App {
	return &App{
		fps:              DefaultFps,
		timeStepMs:       1000 / DefaultFps,
		renderFrameDelta: 1 / float32(DefaultFps),
		quit:             false,
		vsync:            false,
		appListener:      appListener,
		initialized:      false,
	}
}

func (a *App) Initialize() {
	a.appListener.OnStartUp(a)
	a.initialized = true
}

func (a *App) SetFps(fps uint32) {
	a.fps = fps
	a.timeStepMs = uint64(1000 / fps)
	a.renderFrameDelta = 1 / float32(fps)
}

func (a *App) GetFps() uint32 {
	return a.fps
}

func (a *App) AddNewWindow(title string, x, y, width, height int32, flags uint32, bgColor sdl.Color) *RenderWindow {
	w := NewRenderWindow(title, x, y, width, height, flags, bgColor)
	a.windows[w] = nil
	return w
}
func (a *App) AssociateListener(w *RenderWindow, listener RenderListener) {
	a.windows[w] = listener
}

func (a *App) Loop() {
	if !a.initialized {
		panic("The 'Initialize' method is not called.")
	}
	var nextGameStep = sdl.GetTicks64()
	var timeStepMs = a.timeStepMs
	for !a.quit {
		var now = sdl.GetTicks64()
		if nextGameStep <= now || a.vsync {
			lowestFps := LowestFps
			for nextGameStep <= now && lowestFps > 0 {
				a.appListener.OnGameStep(a, float32(now-nextGameStep)/1000)
				nextGameStep += timeStepMs
				lowestFps -= 1
			}

			for window, listener := range a.windows {
				window.preRenderOneFrame()
				if listener != nil {
					listener.OnRender(window, a.renderFrameDelta)
				}
				window.postRenderOneFrame()
			}
		} else {
			sdl.Delay(uint32(nextGameStep - now))
		}
	}
}

func (a *App) RequestGameQuit() {
	a.quit = true
}
func (a *App) Terminate() {
	if !a.initialized {
		panic("The 'Initialize' method is not called.")
	}
	a.appListener.OnShutDown(a)
	for w, _ := range a.windows {
		w.Destroy()
	}
	sdl.Quit()
}
