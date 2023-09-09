package ui

import (
	"encoding/base64"
	"fmt"
	rl "github.com/gen2brain/raylib-go/raylib"
	"github.com/jeffreybian/mines/assets"
	"github.com/jeffreybian/mines/game"
	log "github.com/sirupsen/logrus"
	_ "image/png"
)

type GameWindow struct {
	image              *rl.Image
	mainTexture        rl.Texture2D
	field              *game.MineField
	fieldOffsetX       float32
	fieldOffsetY       float32
	fieldW             float32
	fieldH             float32
	startTime          float64
	seconds            int
	width              float32
	smileyIndex        int // The index to the Smiley
	smileyOffsetX      float32
	smileyOffsetY      float32
	mouseDownRow       int // Temporary current row under mouse
	mouseDownCol       int // Temporary current col under mouse
	rightButtonDownRow int
	rightButtonDownCol int
	locked             bool
}

var MainGameWindow *GameWindow

func init() {
	MainGameWindow = &GameWindow{
		fieldOffsetX:  WindowLeftBorderW,
		fieldOffsetY:  WindowTopFillH,
		smileyOffsetY: 16.0,
	}
}

func (g *GameWindow) InitSystem(rows, cols, numberOfMines int) {
	rl.InitWindow(1, 1, "")
	rawTexture, err := base64.StdEncoding.DecodeString(assets.TextureBase64)
	if err != nil {
		panic(err)
	}
	g.image = rl.LoadImageFromMemory(".png", rawTexture, int32(len(rawTexture)))
	if g.image == nil {
		panic("Not able to load image")
	}
	g.mainTexture = rl.LoadTextureFromImage(g.image)

	rl.SetTargetFPS(30)

	m, err := game.NewMineField(rows, cols, numberOfMines)
	if err != nil {
		log.Fatal(err)
	}
	g.ResetGame(m)
}

func (g *GameWindow) ResetGame(field *game.MineField) {
	g.field = field
	var width, height = g.calculateWindowSize()
	title := fmt.Sprintf("%dx%d/%d", field.Cols, field.Rows, field.Mines)
	rl.SetWindowSize(int(width), int(height))
	rl.SetWindowTitle(title)
	g.width = float32(width)
	g.fieldW = float32(field.Cols) * CellW
	g.fieldH = float32(field.Rows) * CellH
	g.smileyOffsetX = (g.width - SmileyW) / 2.0
	g.mouseDownRow = -1
	g.mouseDownCol = -1
	g.rightButtonDownRow = -1
	g.rightButtonDownCol = -1
	g.startTime = -1.0
	g.seconds = 0
	g.locked = false
}

func (g *GameWindow) Loop() {
	for !rl.WindowShouldClose() {
		g.handleInput()
		rl.BeginDrawing()
		rl.ClearBackground(rl.LightGray)
		g.drawWindow()
		g.drawOverlays()
		g.drawAllCells()
		rl.EndDrawing()
	}
}

func (g *GameWindow) CleanUp() {
	rl.UnloadImage(g.image)
	rl.UnloadTexture(g.mainTexture)
	rl.CloseWindow()
}

func (g *GameWindow) calculateWindowSize() (int32, int32) {
	return int32(float32(g.field.Cols)*CellW + WindowLeftBorderW + WindowRightBorderW), int32(float32(g.field.Rows)*CellH + WindowBottomBorderH + WindowTopFillH)
}

func (g *GameWindow) drawCell(blockState game.BlockState, position rl.Vector2) {
	var index = 9
	if blockState&game.MaskRevealed == game.MaskRevealed {
		if blockState&game.MaskMine == game.MaskMine {
			index = 14
		} else if blockState&0b1111 == game.MaskEmpty {
			index = 0
		} else if (blockState & game.MaskNeighboringMines) > 0 {
			index = int(blockState & game.MaskNeighboringMines)
		}
		if blockState&game.MaskMine == game.MaskMine && blockState&game.MaskInvalid == game.MaskInvalid {
			index = 15
		}
	} else if blockState&game.MaskFlagged == game.MaskFlagged {
		if blockState&game.MaskInvalid == game.MaskInvalid {
			index = 16
		} else {
			index = 11
		}
	} else if blockState&game.MaskQuestioned == game.MaskQuestioned {
		index = 12
	}
	rl.DrawTextureRec(g.mainTexture, cellSourceRects[index], position, rl.White)
}

func (g *GameWindow) drawAllCells() {
	offsetX, offsetY := g.fieldOffsetX, g.fieldOffsetY
	rows, cols := g.field.Rows, g.field.Cols
	states := g.field.States
	for r := 0; r < rows; r++ {
		for c := 0; c < cols; c++ {
			state := states[r*cols+c]
			if r == g.mouseDownRow && c == g.mouseDownCol {
				state = game.MaskEmpty | game.MaskRevealed
			}
			g.drawCell(
				state,
				rl.Vector2{
					X: offsetX + float32(c)*CellW,
					Y: offsetY + float32(r)*CellH,
				},
			)
		}
	}
}

func (g *GameWindow) drawWindow() {

	cols := g.field.Cols
	rows := g.field.Rows
	tex := g.mainTexture

	// 1. Draws top part
	// The width of top left corner + top right corner = 8 cells,
	// We need to repeatedly fill in the top fill (`cols` - 8) number of times
	// Then overlay the numbers and the smiley upon.
	// 1.1 Draws the top part container
	tlCornerSourceRect := windowElementSourceRects[0]
	topFillSourceRect := windowElementSourceRects[2]

	topFillOffsetX := tlCornerSourceRect.Width

	rl.DrawTextureRec(tex, tlCornerSourceRect, rl.Vector2{X: 0, Y: 0}, rl.White)

	for i := 0; i < cols-8; i++ {
		rl.DrawTextureRec(tex, topFillSourceRect, rl.Vector2{X: topFillOffsetX, Y: 0}, rl.White)
		topFillOffsetX += CellW
	}

	rl.DrawTextureRec(tex, windowElementSourceRects[1], rl.Vector2{X: topFillOffsetX, Y: 0}, rl.White)

	// 1.2 TODO: Now overlay the Smiley and numbers!

	// 2. Draws other borders and corners around
	leftBorderSourceRect := windowElementSourceRects[3]
	rightBorderSourceRect := windowElementSourceRects[4]
	leftBorderOffsetY := WindowTLCornerH

	bottomBorderSourceRect := windowElementSourceRects[7]

	blCornerSourceRect := windowElementSourceRects[5]
	brCornerSourceRect := windowElementSourceRects[6]
	bottomBorderOffsetX := blCornerSourceRect.Width

	for i := 0; i < rows; i++ {
		rl.DrawTextureRec(tex, leftBorderSourceRect, rl.Vector2{X: 0, Y: leftBorderOffsetY}, rl.White)
		leftBorderOffsetY += CellH
	}
	// BL Corner
	rl.DrawTextureRec(tex, blCornerSourceRect, rl.Vector2{X: 0, Y: leftBorderOffsetY}, rl.White)
	// Continues to bottom border
	for i := 0; i < cols; i++ {
		rl.DrawTextureRec(tex, bottomBorderSourceRect, rl.Vector2{X: bottomBorderOffsetX, Y: leftBorderOffsetY}, rl.White)
		bottomBorderOffsetX += CellW
	}
	// BR Corner
	rl.DrawTextureRec(tex, brCornerSourceRect, rl.Vector2{X: bottomBorderOffsetX, Y: leftBorderOffsetY}, rl.White)

	// Right border
	for i := 0; i < rows; i++ {
		rl.DrawTextureRec(tex, rightBorderSourceRect, rl.Vector2{X: bottomBorderOffsetX, Y: WindowTLCornerH + float32(i)*CellH}, rl.White)
	}
}

// drawOverlays draws the smiley and numbers etc.
func (g *GameWindow) drawOverlays() {
	// Numbers:
	tex := g.mainTexture
	numFlags := g.field.Flags
	g.drawLCDNumber(tex, numFlags, rl.Vector2{X: 18.0, Y: 16.0})
	if g.startTime >= 0 {
		g.seconds = int(rl.GetTime() - g.startTime)
	}
	g.drawLCDNumber(tex, g.seconds, rl.Vector2{X: g.width - 53.0, Y: 16.0})
	rl.DrawTextureRec(tex, smileySourceRects[g.smileyIndex], rl.Vector2{X: g.smileyOffsetX, Y: 16.0}, rl.White)
}

func (g *GameWindow) drawLCDNumber(texture rl.Texture2D, num int, position rl.Vector2) {
	// Extract the digits, maximum 3 digits
	digits := [3]int{
		num / 100,
		num / 10,
		num % 10,
	}
	rl.DrawTextureRec(texture, lcdNumberSourceRects[digits[0]], position, rl.White)
	rl.DrawTextureRec(texture, lcdNumberSourceRects[digits[1]], rl.Vector2{X: position.X + LCDDigitW, Y: position.Y}, rl.White)
	rl.DrawTextureRec(texture, lcdNumberSourceRects[digits[2]], rl.Vector2{X: position.X + LCDDigitW*2, Y: position.Y}, rl.White)
}

func (g *GameWindow) handleInput() {
	row, col, isOverSmiley := g.mouseTest()
	if g.locked && !isOverSmiley {
		return
	}
	if rl.IsMouseButtonDown(rl.MouseLeftButton) {
		if isOverSmiley {
			g.smileyIndex = 1
		} else {
			g.smileyIndex = 2
		}
		if row >= 0 && col >= 0 {
			state := g.field.Get(row, col)
			if state&game.MaskRevealed == game.MaskRevealed {
				return
			}
			if state&game.MaskFlagged != game.MaskFlagged {
				// Saves this, need to recover later if necessary
				g.mouseDownRow, g.mouseDownCol = row, col
			}
		}

		if rl.IsMouseButtonDown(rl.MouseRightButton) {
			// TODO: Later finish the 9-block reveal
		}
	} else if rl.IsMouseButtonDown(rl.MouseRightButton) {
		// Flag
		if row >= 0 && col >= 0 && (g.rightButtonDownRow != row || g.rightButtonDownCol != col) {
			state := g.field.Get(row, col)
			if state&game.MaskRevealed == game.MaskRevealed {
				return
			}
			g.field.Flag(row, col)
			g.rightButtonDownRow, g.rightButtonDownCol = row, col
		}
	}
	if rl.IsMouseButtonReleased(rl.MouseLeftButton) {
		g.smileyIndex = 0

		if isOverSmiley {
			m, err := game.NewMineField(g.field.Rows, g.field.Cols, g.field.Mines)
			if err != nil {
				log.Fatal(err)
			}
			g.ResetGame(m)
		}

		if row >= 0 && col >= 0 {
			winOrLose := g.field.RevealAndCheckWinOrLose(row, col)
			if winOrLose != 0 {
				g.startTime = -1.0
				g.locked = true
				if winOrLose == 1 {
					// Win
					g.smileyIndex = 3
				} else if winOrLose == -1 {
					// Lose
					g.smileyIndex = 4
				}
			} else {
				if g.startTime < 0 {
					g.startTime = rl.GetTime()
				}
			}
		}
		g.mouseDownRow, g.mouseDownCol = -1, -1
	}
	if rl.IsMouseButtonReleased(rl.MouseRightButton) {
		g.rightButtonDownRow, g.rightButtonDownCol = -1, -1
	}
}

// mouseTest returns tuple of (row, col, isOverSmiley)
func (g *GameWindow) mouseTest() (int, int, bool) {
	mousePosition := rl.GetMousePosition()
	mx, my := mousePosition.X, mousePosition.Y
	sx, sy := g.smileyOffsetX, g.smileyOffsetY
	fx, fy := g.fieldOffsetX, g.fieldOffsetY
	rows, cols := float32(g.field.Rows), float32(g.field.Cols)
	r, c := -1, -1
	isOverSmiley := false
	if mx >= sx && mx <= sx+SmileyW && my >= sy && my <= sy+SmileyH {
		isOverSmiley = true
	}
	if mx >= fx && mx <= fx+cols*CellW && my >= fy && my <= fy+rows*CellH {
		// in mine area
		c = int((mx - fx) / CellW)
		r = int((my - fy) / CellH)
	}

	return r, c, isOverSmiley
}
