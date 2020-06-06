package main

import (
	"encoding/hex"
	"fmt"

	"github.com/gizak/termui/v3"
)

func toChar(b byte) byte {
	if b < 32 || b > 126 {
		return '.'
	}
	return b
}

func ColorHex(byte_array []byte, from, to, offset int) []termui.Cell {
	var cells []termui.Cell
	style_default := termui.Style{termui.ColorClear, termui.ColorClear, termui.ModifierClear}
	style_highlight := termui.Style{15, 21, termui.ModifierClear}

	for j := 0; j < len(byte_array)/16; j++ {
		_cells := termui.RunesToStyledCells([]rune(fmt.Sprintf(" %#010x | ", j*16+offset)), style_default)
		cells = append(cells, _cells...)
		end := (j + 1) * 16
		fill_to_end := end
		if len(byte_array) < end {
			end = len(byte_array) - 1
		}
		ascii_cells := []termui.Cell{}
		for i, b := range byte_array[j*16 : end] {
			bytes := []byte{b}
			if j*16+i >= from && j*16+i < to {
				h := fmt.Sprintf("%s", hex.EncodeToString(bytes))
				cells = append(cells, termui.RunesToStyledCells([]rune(h), style_highlight)...)
				cells = append(cells, termui.RunesToStyledCells([]rune(" "), style_default)...)
				ascii_cells = append(ascii_cells, termui.RunesToStyledCells([]rune(string(toChar(b))), style_highlight)...)
			} else {
				h := fmt.Sprintf("%s ", hex.EncodeToString(bytes))
				cells = append(cells, termui.RunesToStyledCells([]rune(h), style_default)...)
				ascii_cells = append(ascii_cells, termui.RunesToStyledCells([]rune(string(toChar(b))), style_default)...)
			}

			if i == 7 {
				cells = append(cells, termui.RunesToStyledCells([]rune(" "), style_default)...)
			}
			if fill_to_end != end {
				i := 0
				for i = 0; i < fill_to_end-end; i++ {
					cells = append(cells, termui.RunesToStyledCells([]rune("   "), style_default)...)
				}
				if i > 9 {
					cells = append(cells, termui.RunesToStyledCells([]rune(" "), style_default)...)
				}
			}
		}
		cells = append(cells, termui.RunesToStyledCells([]rune("| "), style_default)...)
		cells = append(cells, ascii_cells...)
		cells = append(cells, termui.RunesToStyledCells([]rune("\n"), style_default)...)
	}

	return cells
}

type HexView struct {
	Content       []byte
	Start, Stop   int
	ContentOffset int
	ContentSize   int
	ColoredString string
	View          *CellParagraph
	Rendered      bool
}

func (hexview *HexView) Render() error {
	if hexview.Rendered == true {
		return nil
	}
	hexview.Rendered = true
	hexview.View.Cells =
		ColorHex(hexview.Content, hexview.Start, hexview.Stop, hexview.ContentOffset)
	return nil
}

func (hexview *HexView) Selection() []byte {
	return hexview.Content[hexview.Start:hexview.Stop]
}

func (hexview *HexView) HandleInput(input string) bool {
	switch input {
	case "h":
		if hexview.Start > 0 {
			hexview.Start -= 1
			hexview.Stop -= 1
		}
		break
	case "l":
		if hexview.Stop < hexview.ContentSize-1 {
			hexview.Start += 1
			hexview.Stop += 1
		}
		break
	case "H":
		if hexview.Stop > hexview.Start+1 {
			hexview.Stop -= 1
		}
		break

	case "j":
		if hexview.Stop+16 < hexview.ContentSize {
			hexview.Stop += 16
			hexview.Start += 16
		} else {
			diff := hexview.Stop - hexview.Start
			hexview.Stop = hexview.ContentSize - 1
			hexview.Start = hexview.Stop - diff
		}
		break
	case "L":
		if hexview.Stop < hexview.ContentSize-1 {
			hexview.Stop += 1
		}
		break
	case "k":
		if hexview.Start-16 < 0 {
			hexview.Start = 0
			hexview.Stop -= hexview.Start
		} else {
			hexview.Start -= 16
			hexview.Stop -= 16
		}
		break
	default:
		return false
	}
	hexview.Rendered = false
	return true
}
