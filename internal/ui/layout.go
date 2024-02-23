package ui

import (
	"errors"
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

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

func (app *App) setStatus(status string) error {
	// Calculate size
	maxX, maxY := app.gui.Size()
	lines := strings.Split(status, "\n")
	auxY := len(lines)
	app.statusPos.y0 = maxY/2 - auxY/2
	app.statusPos.y1 = app.statusPos.y0 + auxY + 1
	auxX := 0
	for _, v := range lines {
		auxX = maxInts(auxX, len(v))
	}
	app.statusPos.x0 = maxX/2 - auxX/2
	app.statusPos.x1 = app.statusPos.x0 + auxX + 1

	// update view
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		app.statusView.Clear()
		app.statusView.Visible = true
		fmt.Fprint(app.statusView, status)

		return nil
	})

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
	// status view
	if view, err := app.gui.SetView("statusWindow", app.statusPos.x0, app.statusPos.y0, app.statusPos.x1, app.statusPos.y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Visible = app.statusVisible
		app.statusView = view
	}
	return nil
}
