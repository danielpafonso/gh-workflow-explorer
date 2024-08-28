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
		maxX, maxY := app.gui.Size()
		app.status.X0 = maxX/2 - xSize/2
		app.status.Y0 = maxY/2 - ySize/2
		app.status.X1 = app.status.X0 + xSize
		app.status.Y1 = app.status.Y0 + ySize
		app.statusVisible = true
		app.statusView.Clear()
		fmt.Fprint(app.statusView, strings.Join(lines, "\n"))
		return nil
	})
}

func (app *App) FilterWindow(text string) {
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		padding := 2
		if app.filter.visible {
			lines := strings.Split(text, "\n")
			ySize := len(lines) + 1
			xSize := 0
			for i, line := range lines {
				xSize = maxInts(xSize, len(line)+padding*2)
				// add padding
				lines[i] = fmt.Sprintf("%s%s", strings.Repeat(" ", padding), line)
			}
			maxX, maxY := app.gui.Size()
			app.filter.pos.X0 = maxX/2 - xSize/2
			app.filter.pos.Y0 = maxY/2 - ySize/2
			app.filter.pos.X1 = app.filter.pos.X0 + xSize
			app.filter.pos.Y1 = app.filter.pos.Y0 + ySize
			app.filter.visible = true
			app.filter.view.Clear()
			fmt.Fprint(app.filter.view, strings.Join(lines, "\n"))
		}
		return nil
	})
}

func (app *App) scrollMain(g *gocui.Gui, v *gocui.View, dy int) error {
	// disable warning
	_ = g
	if v != nil {
		_, size := v.Size()
		_, cy := v.Cursor()
		_, oy := v.Origin()
		cMove := cy + dy
		overflow := true
		// check if lines overflow
		runs := 0
		for _, run := range app.runs {
			if run.show {
				runs += 1
			}
		}
		//if len(app.runs) < size {
		if runs < size {
			overflow = false
			//size = len(app.runs)
			size = runs
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
		_, oy = v.Origin()
		v.Subtitle = fmt.Sprintf("%d/%d", cMove+oy+1, v.LinesHeight()-1)
	}
	return nil
}

// func (app *App) layout(*gocui.Gui) error {
func (app *App) Layout(*gocui.Gui) error {
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
	if view, err := app.gui.SetView("filterlist", maxX/2, -1, maxX, 8, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Frame = false
		app.filterListView = view
	} else {
		view.Clear()
		view.WriteString(fmt.Sprintf(
			"Name: %s\n\nCommit: %s\n\nStatus: %s",
			app.filter.fields.Name,
			app.filter.fields.Commit,
			app.filter.fields.Status,
		))
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
		helpLine := "<q> exit    <UP/DOWN arrow> nav    <space> toogle    <a> (un)select all    <f> filter    <d> delete    <r> refresh"
		view.SetWritePos(maxX/2-len(helpLine)/2, 0)
		view.WriteString(helpLine)
	}
	// status view
	if view, err := app.gui.SetView("status", app.status.X0, app.status.Y0, app.status.X1, app.status.Y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Visible = app.statusVisible
		app.statusView = view
	} else {
		view.Visible = app.statusVisible
	}
	// filter window view
	app.filter.pos.X0 = maxX/2 - 20
	app.filter.pos.Y0 = maxY/2 - 4
	app.filter.pos.X1 = app.filter.pos.X0 + 40
	app.filter.pos.Y1 = app.filter.pos.Y0 + 8

	if view, err := app.gui.SetView("filter", app.filter.pos.X0, app.filter.pos.Y0, app.filter.pos.X1, app.filter.pos.Y1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Title = " Filter "
		// view.FrameColor = gocui.ColorMagenta
		view.Visible = app.filter.visible
		app.filter.view = view
		view.WriteString("Name:\n\nCommit:\n\nStatus:\n\n\t<enter> confirm\t<Esc> cancel")
	} else {
		view.Visible = app.filter.visible
	}

	if view, err := app.gui.SetView("filter-name", app.filter.pos.X0+9, app.filter.pos.Y0, app.filter.pos.X1-1, app.filter.pos.Y0+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Visible = app.filter.visible
		view.Frame = false
		view.Editable = true
		app.filter.inputs = append(app.filter.inputs, view)
	} else {
		view.Visible = app.filter.visible
	}
	if view, err := app.gui.SetView("filter-commit", app.filter.pos.X0+9, app.filter.pos.Y0+2, app.filter.pos.X1-1, app.filter.pos.Y0+4, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Visible = app.filter.visible
		view.Frame = false
		view.Editable = true
		app.filter.inputs = append(app.filter.inputs, view)
	} else {
		view.Visible = app.filter.visible
	}
	if view, err := app.gui.SetView("filter-status", app.filter.pos.X0+9, app.filter.pos.Y0+4, app.filter.pos.X1-1, app.filter.pos.Y0+6, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.Visible = app.filter.visible
		view.Frame = false
		view.Editable = true
		app.filter.inputs = append(app.filter.inputs, view)
	} else {
		view.Visible = app.filter.visible
	}

	return nil
}
