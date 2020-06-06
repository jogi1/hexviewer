package main

import (
	"flag"
	"fmt"
	ui "github.com/gizak/termui/v3"
	"io/ioutil"
	"os"
)

type Viewer struct {
	HexView               HexView
	TypeView              TypeView
	Grid                  *ui.Grid
	TermWidth, TermHeight int
	Filename              string
	File                  []byte
	Input                 string
}

func (viewer *Viewer) Render() {
	viewer.HexView.Render()
	viewer.TypeView.Render()
}

func (viewer *Viewer) HandleInput(input string) {
	viewer.HexView.View.Title = input
	if input == "<F12>" {
		viewer.TypeView.Rendered = false
		viewer.TypeView.BigEndian = !viewer.TypeView.BigEndian
		return
	}
	if input == "G" {
		for {
			viewer.HexView.ContentOffset += 16 * 16
			if viewer.HexView.ContentOffset > len(viewer.File) {
				viewer.HexView.ContentOffset -= 16 * 16
				break
			}
		}
		viewer.HexView.Content = viewer.File[viewer.HexView.ContentOffset : viewer.HexView.ContentOffset+viewer.HexView.ContentSize]
		viewer.HexView.Rendered = false
		return
	}
	if input == "<C-j>" || input == "N" {
		if input == "<C-j>" {
			viewer.HexView.ContentOffset += 16
		}
		if input == "N" {
			viewer.HexView.ContentOffset += 128
		}
		if viewer.HexView.ContentOffset+16*16 > len(viewer.File) {
			viewer.HexView.ContentOffset += len(viewer.File) - viewer.HexView.ContentOffset
		}

		viewer.HexView.Content = viewer.File[viewer.HexView.ContentOffset : viewer.HexView.ContentOffset+viewer.HexView.ContentSize]
		viewer.HexView.Rendered = false
		return
	}
	if input == "<C-k>" || input == "P" {
		if input == "<C-k>" {
			viewer.HexView.ContentOffset -= 16
		}
		if input == "P" {
			viewer.HexView.ContentOffset -= 128
		}
		if viewer.HexView.ContentOffset < 0 {
			viewer.HexView.ContentOffset = 0
		}

		viewer.HexView.Content = viewer.File[viewer.HexView.ContentOffset : viewer.HexView.ContentOffset+viewer.HexView.ContentSize]
		viewer.HexView.Rendered = false
		return
	}
	if viewer.Input == "hex" {
		if viewer.HexView.HandleInput(input) {
			viewer.TypeView.SetType(viewer.HexView.Selection())
		}
	}
}

func (viewer *Viewer) Init(filename string) error {
	var err error
	viewer.Input = "hex"

	viewer.File, err = ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	viewer.TermWidth, viewer.TermHeight = ui.TerminalDimensions()
	viewer.HexView.View = NewCellParagraph()
	viewer.HexView.View.Title = "Hex View"
	viewer.HexView.Stop = 1
	viewer.HexView.ContentSize = 16 * 16
	if viewer.HexView.ContentSize > len(viewer.File) {
		viewer.HexView.ContentSize = len(viewer.File)
	}
	viewer.HexView.Content = viewer.File[viewer.HexView.ContentOffset : viewer.HexView.ContentOffset+viewer.HexView.ContentSize]
	if err = viewer.HexView.Render(); err != nil {
		return err
	}

	viewer.TypeView.Init()

	viewer.Grid = ui.NewGrid()
	viewer.Grid.SetRect(0, 0, viewer.TermWidth, viewer.TermHeight)
	viewer.Grid.Set(
		ui.NewRow(2.0/3.0,
			ui.NewCol(1.0, viewer.HexView.View),
		),
		ui.NewRow(1.0/3.0,
			ui.NewCol(1.0, viewer.TypeView.View),
		),
	)

	return nil
}

func main() {
	flag.Parse()

	var viewer Viewer

	if len(flag.Args()) != 1 {
		fmt.Printf("need to supply a file")
		os.Exit(1)
	}

	if err := ui.Init(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	defer ui.Close()

	err := viewer.Init(flag.Args()[0])
	if err != nil {
		fmt.Printf("%#v", err)
		os.Exit(1)
	}

	quit := false

	ui.Render(viewer.Grid)

	for e := range ui.PollEvents() {
		if e.Type == ui.KeyboardEvent {
			switch e.ID {
			case "Q":
				quit = true
				break
			default:
				viewer.HandleInput(e.ID)
				break
			}

			if quit {
				break
			}

			ui.Clear()
			viewer.Render()
			ui.Render(viewer.Grid)
		}
	}
	os.Exit(0)
}
