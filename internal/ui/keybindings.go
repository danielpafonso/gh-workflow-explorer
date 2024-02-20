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
	_ = writeTootle
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

func (app *App) debug(g *gocui.Gui, v *gocui.View) error {
	runs := make([]string, 0)
	_, cy := v.Cursor()
	_, oy := v.Origin()
	_, size := v.Size()
	runs = append(runs, fmt.Sprintf("l: %d", len(app.runs)))
	runs = append(runs, fmt.Sprintf("s: %d", size))
	runs = append(runs, fmt.Sprint(cy))
	runs = append(runs, fmt.Sprint(oy))

	g.DeleteView("status")
	app.StatusView(strings.Join(runs, "\n"))
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
	for _, id := range runsToDelete {
		app.api.DeleteWorkflow(id)
		// _ = id
		// time.Sleep(1000 * time.Millisecond)
	}
	// update Main view
	app.WriteMain()
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
	if err := g.SetKeybinding("", 'f', gocui.ModNone, app.debug); err != nil {
		return err
	}
	return nil
}
