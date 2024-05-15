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
	app.refreshMain(g, v)
	return nil
}

func (app *App) filterRuns(g *gocui.Gui, v *gocui.View) error {
	// toogle on/off filter window
	app.filterVisible = !app.filterVisible
	data := fmt.Sprintf("name:\nvisible: %v\nOK   Cancel", app.filterVisible)
	app.FilterWindow(data)
	//app.StatusView(fmt.Sprintf("name:\nvisible: %v\n OK   Cancel", app.filterVisible))
	return nil
}

func (app *App) refreshMain(g *gocui.Gui, v *gocui.View) error {
	app.WriteMain()

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
	if err := g.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return app.scrollMain(g, v, -1)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			return app.scrollMain(g, v, 1)
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, app.toogleMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gocui.ModNone, app.deleteRuns); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'r', gocui.ModNone, app.refreshMain); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'a', gocui.ModNone, app.toogleAllRuns); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'f', gocui.ModNone, app.filterRuns); err != nil {
		return err
	}
	return nil
}
