package ui

import (
	"github.com/awesome-gocui/gocui"
)

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}

func clearDebug(g *gocui.Gui, v *gocui.View) error {
	_ = g.DeleteView("status")
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
	if err := g.SetKeybinding("", gocui.KeySpace, gocui.ModNone, app.debug); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'd', gocui.ModNone, quit); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'r', gocui.ModNone,
		func(g *gocui.Gui, v *gocui.View) error {
			app.WriteMain()
			return nil
		}); err != nil {
		return err
	}
	if err := g.SetKeybinding("", 'f', gocui.ModNone, clearDebug); err != nil {
		return err
	}
	return nil
}
