package ui

import (
	"fmt"
	"strings"

	"github.com/awesome-gocui/gocui"
)

func (app *App) toogleMain(g *gocui.Gui, v *gocui.View) error {
	xPosition := 3
	_, cy := v.Cursor()
	_, oy := v.Origin()
	// update run
	writeTootle := ' '
	calculatedLine := cy + oy
	for i := range app.runs {
		if app.runs[i].line == calculatedLine {
			app.runs[i].toogle = !app.runs[i].toogle
			if app.runs[i].toogle {
				writeTootle = '*'
			}
			break
		}
	}
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		v.SetCursor(xPosition, cy+oy)
		v.EditDelete(false)
		v.EditWrite(writeTootle)
		// reset position, why? don't know..
		v.SetCursor(xPosition, cy)
		return nil
	})
	return nil
}

func (app *App) toogleAllRuns(g *gocui.Gui, v *gocui.View) error {
	toSelect := make([]int, 0)
	for i := range app.runs {
		if !app.runs[i].show {
			// skip hidden runs
			continue
		}
		if !app.runs[i].toogle {
			toSelect = append(toSelect, i)
		}
	}
	if len(toSelect) > 0 {
		// selected non-selected
		for _, i := range toSelect {
			app.runs[i].toogle = true
		}
	} else {
		// unselect all
		for i := range app.runs {
			if !app.runs[i].show {
				// skip hidden runs
				continue
			}
			app.runs[i].toogle = false
		}
	}
	app.WriteMain(true)
	return nil
}

func (app *App) filterOpen(g *gocui.Gui, v *gocui.View) error {
	// toogle on filter window
	app.filter.visible = true
	app.filter.focus = 0
	g.SetCurrentView(app.filter.inputs[0].Name())
	app.filter.inputs[0].BgColor = gocui.ColorCyan
	// Add current filter?
	return nil
}

func (app *App) filterClose(g *gocui.Gui, v *gocui.View) error {
	v.BgColor = gocui.ColorDefault
	for _, input := range app.filter.inputs {
		input.Clear()
	}
	app.filter.visible = false
	g.SetCurrentView(app.mainView.Name())
	return nil
}

func (app *App) filterFocus(g *gocui.Gui, v *gocui.View) error {
	// update current
	v.BgColor = gocui.ColorDefault
	tmp := v.Buffer()
	v.Clear()
	v.WriteString(tmp)
	// update next
	app.filter.focus = (app.filter.focus + 1) % len(app.filter.inputs)
	next := app.filter.inputs[app.filter.focus]
	g.SetCurrentView(next.Name())
	next.BgColor = gocui.ColorCyan
	next.MoveCursor(len(next.Buffer()), 0)
	return nil
}

func (app *App) filterApply(g *gocui.Gui, v *gocui.View) error {
	// update filter view
	app.filter.fields.Name = app.filter.inputs[0].Buffer()
	app.filter.fields.Commit = app.filter.inputs[1].Buffer()
	app.filter.fields.Status = app.filter.inputs[2].Buffer()

	// update runs list
	app.filterRuns()
	// close
	app.filterClose(g, v)
	app.WriteMain()
	return nil
}

// func (app *App) filterRuns(g *gocui.Gui, v *gocui.View) error {
func (app *App) filterRuns() {
	app.showRuns = 0
	for i, workflow := range app.runs {
		filter := strings.Contains(workflow.run.Name, app.filter.fields.Name) &&
			strings.Contains(workflow.run.Title, app.filter.fields.Commit) &&
			strings.Contains(workflow.run.Conclusion, app.filter.fields.Status)

		// apply filter
		app.runs[i].show = filter

		// update runs to show counter
		if filter {
			app.showRuns += 1
		}
	}
	// return nil
}

func (app *App) refreshMain(g *gocui.Gui, v *gocui.View) error {
	app.StatusView("Requesting workflows\nPlease Wait...")
	go func() {
		// clean runs
		if len(app.runs) > 0 {
			app.runs = []workflows{}
		}
		// clean filters
		app.filter.fields.Name = ""
		app.filter.fields.Commit = ""
		app.filter.fields.Status = ""
		// get data from api
		app.refreshWorkflows()
		// update main
		app.WriteMain()
		app.WriteColumns()
	}()
	return nil
}

func (app *App) deleteRuns(g *gocui.Gui, v *gocui.View) error {
	runsToDelete := make([]int, 0)
	for i := 0; i < len(app.runs); i++ {
		if app.runs[i].toogle {
			// add to delete
			runsToDelete = append(runsToDelete, app.runs[i].run.ID)
			// remove from runs array
			app.runs = append(app.runs[:i], app.runs[i+1:]...)
			i--
			// update runs to show counter
			app.showRuns -= 1
		}
	}
	// delete runs
	go func() {
		for i, id := range runsToDelete {
			fill := int((float32(i) / float32(len(runsToDelete))) * 15)
			app.StatusView(fmt.Sprintf("Deleting Workruns: \n [%s%s]", strings.Repeat("=", fill), strings.Repeat(" ", 15-fill)))
			app.statusView.Subtitle = fmt.Sprintf("%d/%d", i+1, len(runsToDelete))
			app.api.DeleteWorkflow(id)
		}
		// update Main view
		app.WriteMain()
	}()
	return nil
}

func (app *App) keybindings(g *gocui.Gui) error {
	// inline function: exit
	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	// calling functions functions
	if err := g.SetKeybinding("main", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return app.scrollMain(g, v, -1)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return app.scrollMain(g, v, 1)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", gocui.KeySpace, gocui.ModNone, app.toogleMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'd', gocui.ModNone, app.deleteRuns); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'r', gocui.ModNone, app.refreshMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'r', gocui.ModAlt, func(g *gocui.Gui, v *gocui.View) error {
		app.WriteMain()
		return nil
	}); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'a', gocui.ModNone, app.toogleAllRuns); err != nil {
		return err
	}
	if err := g.SetKeybinding("main", 'f', gocui.ModNone, app.filterOpen); err != nil {
		return err
	}

	//   filter Window
	// Add keybings to all inputs
	for _, view := range []string{"filter-name", "filter-commit", "filter-status"} {
		if err := g.SetKeybinding(view, gocui.KeyTab, gocui.ModNone, app.filterFocus); err != nil {
			return err
		}
		if err := g.SetKeybinding(view, gocui.KeyEsc, gocui.ModNone, app.filterClose); err != nil {
			return err
		}
		if err := g.SetKeybinding(view, gocui.KeyEnter, gocui.ModNone, app.filterApply); err != nil {
			return err
		}
	}
	return nil
}
