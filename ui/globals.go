package ui

import rl "github.com/gen2brain/raylib-go/raylib"

// Offsets and size from Source texture
const (
	GapWidth                  float32 = 1.0 // Gap between elements in the texture
	SmileyW                   float32 = 24.0
	SmileyH                   float32 = 24.0
	SmileyOffsetX             float32 = 2.0
	SmileyOffsetY             float32 = 25.0
	CellW                     float32 = 16.0
	CellH                     float32 = 16.0
	MineNumberOffsetX         float32 = 1.0
	MineNumberOffsetY         float32 = 67.0
	MineBlockOffsetX          float32 = 1.0
	MineBlockOffsetY          float32 = 50.0
	WindowTLCornerOffsetX     float32 = 170.0 // TL = TopLeft
	WindowTLCornerOffsetY     float32 = 0
	WindowTLCornerW           float32 = 76.0
	WindowTLCornerH           float32 = 54.0
	WindowTRCornerOffsetX     float32 = 247.0 // TR = TopRight
	WindowTRCornerOffsetY     float32 = 0
	WindowTRCornerW           float32 = 72.0
	WindowTRCornerH                   = WindowTLCornerH
	WindowTopFillOffsetX      float32 = 320.0 // The pattern used to repeat to fill up the top
	WindowTopFillOffsetY      float32 = 0
	WindowTopFillW                    = CellW
	WindowTopFillH                    = WindowTLCornerH
	WindowLeftBorderOffsetX   float32 = 337.0 // Left Border
	WindowLeftBorderOffsetY   float32 = 0
	WindowLeftBorderW         float32 = 11.0
	WindowLeftBorderH                 = CellH
	WindowRightBorderOffsetX  float32 = 349.0 // Right Border
	WindowRightBorderOffsetY  float32 = 0
	WindowRightBorderW        float32 = 8.0
	WindowRightBorderH                = CellH
	WindowBottomBorderOffsetX float32 = 337.0 // Bottom Border
	WindowBottomBorderOffsetY float32 = 17.0
	WindowBottomBorderW               = CellW
	WindowBottomBorderH       float32 = 8.0
	WindowBLCornerOffsetX     float32 = 337.0 // BL = Bottom Left
	WindowBLCornerOffsetY     float32 = 26.0
	WindowBLCornerW           float32 = 11.0
	WindowBLCornerH           float32 = 8.0
	WindowBRCornerOffsetX     float32 = 350.0 // BR = Bottom Right
	WindowBRCornerOffsetY     float32 = 26.0
	WindowBRCornerW           float32 = 8.0
	WindowBRCornerH                   = WindowBLCornerH
	LCDDigitOffsetX           float32 = 1.0
	LCDDigitOffsetY           float32 = 1.0
	LCDDigitW                 float32 = 13.0
	LCDDigitH                 float32 = 23.0
)

var cellSourceRects = []rl.Rectangle{
	{MineBlockOffsetX + CellW + GapWidth, MineBlockOffsetY, CellW, CellH},     // 0 (empty, clean slot)
	{MineNumberOffsetX, MineNumberOffsetY, CellW, CellH},                      // 1
	{MineNumberOffsetX + CellW + GapWidth, MineNumberOffsetY, CellW, CellH},   // 2
	{MineNumberOffsetX + (CellW+GapWidth)*2, MineNumberOffsetY, CellW, CellH}, // 3
	{MineNumberOffsetX + (CellW+GapWidth)*3, MineNumberOffsetY, CellW, CellH}, // 4
	{MineNumberOffsetX + (CellW+GapWidth)*4, MineNumberOffsetY, CellW, CellH}, // 5
	{MineNumberOffsetX + (CellW+GapWidth)*5, MineNumberOffsetY, CellW, CellH}, // 6
	{MineNumberOffsetX + (CellW+GapWidth)*6, MineNumberOffsetY, CellW, CellH}, // 7
	{MineNumberOffsetX + (CellW+GapWidth)*7, MineNumberOffsetY, CellW, CellH}, // 8
	{MineBlockOffsetX, MineBlockOffsetY, CellW, CellH},                        // block, unrevealed - 9
	{MineBlockOffsetX + CellW + GapWidth, MineBlockOffsetY, CellW, CellH},     // empty clean slot, revealed - 10, same as 0
	{MineBlockOffsetX + (CellW+GapWidth)*2, MineBlockOffsetY, CellW, CellH},   // flagged, unrevealed - 11
	{MineBlockOffsetX + (CellW+GapWidth)*3, MineBlockOffsetY, CellW, CellH},   // questioned, unrevealed - 12
	{MineBlockOffsetX + (CellW+GapWidth)*4, MineBlockOffsetY, CellW, CellH},   // questioned, revealed - 13
	{MineBlockOffsetX + (CellW+GapWidth)*5, MineBlockOffsetY, CellW, CellH},   // mine!, revealed - 14
	{MineBlockOffsetX + (CellW+GapWidth)*6, MineBlockOffsetY, CellW, CellH},   // mine stepped on!, revealed - 15
	{MineBlockOffsetX + (CellW+GapWidth)*7, MineBlockOffsetY, CellW, CellH},   // wrongly flagged mine, revealed - 16
}

var smileySourceRects = []rl.Rectangle{
	{SmileyOffsetX, SmileyOffsetY, SmileyW, SmileyH},                        // Normal Smiley - 0
	{SmileyOffsetX + SmileyW + GapWidth, SmileyOffsetY, SmileyW, SmileyH},   // Pressed Down Smiley - 1
	{SmileyOffsetX + (SmileyW+GapWidth)*2, SmileyOffsetY, SmileyW, SmileyH}, // "oh" Smiley - 2
	{SmileyOffsetX + (SmileyW+GapWidth)*3, SmileyOffsetY, SmileyW, SmileyH}, // Win Smiley -3
	{SmileyOffsetX + (SmileyW+GapWidth)*4, SmileyOffsetY, SmileyW, SmileyH}, // Failed Smiley - 4
}

var windowElementSourceRects = []rl.Rectangle{
	{WindowTLCornerOffsetX, WindowTLCornerOffsetY, WindowTLCornerW, WindowTLCornerH},                 // TL Corner 0
	{WindowTRCornerOffsetX, WindowTRCornerOffsetY, WindowTRCornerW, WindowTRCornerH},                 // TR Corner
	{WindowTopFillOffsetX, WindowTopFillOffsetY, WindowTopFillW, WindowTopFillH},                     // Top fill
	{WindowLeftBorderOffsetX, WindowLeftBorderOffsetY, WindowLeftBorderW, WindowLeftBorderH},         // Left Border 3
	{WindowRightBorderOffsetX, WindowRightBorderOffsetY, WindowRightBorderW, WindowRightBorderH},     // Right Border
	{WindowBLCornerOffsetX, WindowBLCornerOffsetY, WindowBLCornerW, WindowBLCornerH},                 // BL Corner
	{WindowBRCornerOffsetX, WindowBRCornerOffsetY, WindowBRCornerW, WindowBRCornerH},                 // BR Corner 6
	{WindowBottomBorderOffsetX, WindowBottomBorderOffsetY, WindowBottomBorderW, WindowBottomBorderH}, // Bottom border
}

var lcdNumberSourceRects = []rl.Rectangle{
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*9, LCDDigitOffsetY, LCDDigitW, LCDDigitH}, // 0
	{LCDDigitOffsetX, LCDDigitOffsetY, LCDDigitW, LCDDigitH},                          // 1
	{LCDDigitOffsetX + LCDDigitW + GapWidth, LCDDigitOffsetY, LCDDigitW, LCDDigitH},   // 2
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*2, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*3, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*4, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*5, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*6, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*7, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
	{LCDDigitOffsetX + (LCDDigitW+GapWidth)*8, LCDDigitOffsetY, LCDDigitW, LCDDigitH},
}
