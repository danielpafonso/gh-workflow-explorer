package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

func (app *App) StatusView(text string) {
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		padding := 2
		// calculate view Size
		lines := strings.Split(text, "\n")
		ySize := len(lines) + 1
		xSize := 0
		for i, line := range lines {
			xSize = maxInts(xSize, len(line)+padding*2)
			// add padding
			lines[i] = fmt.Sprintf("%s%s", strings.Repeat(" ", padding), line)
		}
		// create view
		maxX, _ := app.gui.Size()
		// maxX, maxY := app.gui.Size()
		//if view, err := app.gui.SetView("status", maxX/2-xSize/2, maxY/2, maxX/2+xSize/2, maxY/2+ySize, 0); err != nil {
		if view, err := app.gui.SetView("status", maxX-xSize, 5, maxX, 5+ySize, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			fmt.Fprintln(view, strings.Join(lines, "\n"))
		}
		return nil
	})
}

func (app *App) scrollMain(g *gocui.Gui, v *gocui.View, dy int) error {
	if v != nil {
		_, size := v.Size()
		_, cy := v.Cursor()
		_, oy := v.Origin()
		cMove := cy + dy
		overflow := true
		// check if lines overflow
		if len(app.runs) < size {
			overflow = false
			size = len(app.runs)
		}
		if dy < 0 {
			// scroll up
			if cy+dy < 0 {
				// jump to end
				if overflow {
					v.SetOrigin(0, v.LinesHeight()-size-1)
				}
				cMove = size - 1
			} else if cy+dy < 2 {
				if oy+dy >= 0 {
					v.SetOrigin(0, oy+dy)
					cMove = cy
				}
			}
		} else {
			// scroll down
			if cy+dy == size {
				// jump to start
				if overflow {
					v.SetOrigin(0, 0)
				}
				cMove = 0
			} else if cy+dy >= size-2 {
				if oy+dy < v.LinesHeight()-size {
					// move origin down
					v.SetOrigin(0, oy+dy)
					// keep cursor
					cMove = cy
				}
			}
		}
		// move cursor
		v.SetCursor(0, cMove)
	}
	return nil
}

func (app *App) layout(*gocui.Gui) error {
	maxX, maxY := app.gui.Size()
	// repo view
	if view, err := app.gui.SetView("repo", -1, -1, maxX/2, 8, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Frame = false
		app.repoView = view
	}
	// filter view
	if view, err := app.gui.SetView("filter", maxX/2, -1, maxX, 8, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Frame = false
		app.filterView = view
	}
	// columns name view
	if view, err := app.gui.SetView("column", -1, 8, maxX, 10, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		app.columnsView = view
	}
	// Main view
	if view, err := app.gui.SetView("main", -1, 10, maxX, maxY-2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Highlight = true
		view.SelBgColor = gocui.ColorCyan
		app.gui.SetCurrentView("main")
		app.mainView = view
	}
	// help view
	if view, err := app.gui.SetView("help", -1, maxY-2, maxX, maxY, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		helpLine := "<q> exit    <UP/DOWN arrow> nav    <space> toogle    <f> filter    <d> delete    <r> refresh"
		view.SetWritePos(maxX/2-len(helpLine)/2, 0)
		view.WriteString(helpLine)
	}
	return nil
}
