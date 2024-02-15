package internal

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

// App creates UI and run workflow explorer
type App struct {
	api        GithubApi
	gui        *gocui.Gui
	repoView   *gocui.View
	filterView *gocui.View
	mainView   *gocui.View
}

func NewAppUI(config GithubApi) *App {
	return &App{
		api: config,
	}
}

func (app *App) Layout() func(*gocui.Gui) error {
	return func(g *gocui.Gui) error {
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
		// Main view
		if view, err := app.gui.SetView("main", -1, 8, maxX, maxY-2, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			view.Wrap = true
			app.gui.SetCurrentView("main")
			app.mainView = view
			view.WriteString("asdasdasdasd\nasd")
		}
		// help view
		if view, err := app.gui.SetView("help", -1, maxY-2, maxX, maxY, 0); err != nil {
			if !errors.Is(err, gocui.ErrUnknownView) {
				return err
			}
			helpLine := "<q> exit    <arrows> nav    <space> toogle    <f> filter    <d> delete    <r> refresh"
			view.SetWritePos(maxX/2-len(helpLine)/2, 0)
			view.WriteString(helpLine)
		}
		return nil
	}
}

func (app *App) WriteRepoOnwer() {
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		fmt.Fprintf(app.repoView, "Owner: %s\n\n Repo:%s", app.api.Owner, app.api.Repo)
		return nil
	})
}

func (app *App) WriteFilter() {
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		fmt.Fprintf(app.filterView, "Owner: %s\n\n Repo:%s", app.api.Owner, app.api.Repo)
		return nil
	})
}

func (app *App) WriteMain(text string, clearView ...bool) {
	clear := false
	if len(clearView) > 0 {
		clear = clearView[0]
	}
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		if clear {
			app.mainView.Clear()
		}
		fmt.Fprint(app.mainView, text)
		return nil
	})
}

func (app *App) StartUI() error {
	var err error
	// create terminal GUI
	app.gui, err = gocui.NewGui(gocui.OutputNormal, true)
	if err != nil {
		return err
	}
	defer app.gui.Close()
	// set graphical manager
	app.gui.SetManagerFunc(app.Layout())

	// write dynamic text
	app.WriteRepoOnwer()
	app.WriteFilter()
	app.WriteMain(fmt.Sprintln(app.api.Auth))

	// set Keybings
	if err := app.gui.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'q', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeyArrowUp, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeyArrowDown, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeyArrowLeft, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeyArrowRight, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", gocui.KeySpace, gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'f', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'd', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	if err := app.gui.SetKeybinding("", 'r', gocui.ModNone, func(g *gocui.Gui, v *gocui.View) error {
		return gocui.ErrQuit
	}); err != nil {
		return err
	}
	// enter UI mainloop
	if err := app.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
