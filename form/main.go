package main

import (
	"errors"
	"fmt"

	"github.com/awesome-gocui/gocui"
)

var (
	gui *gocui.Gui
)

type form struct {
	mainView    *gocui.View
	formInput   *gocui.View
	formVisible bool
	nameInput   *gocui.View
	valueInput  *gocui.View
	listView    []string
	currentView int
}

func (app *form) layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	// main window
	if view, err := g.SetView("main", 0, 0, maxX-1, maxY-1, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		//view.BgColor = gocui.GetColor("#ffffff")
		view.BgColor = gocui.ColorCyan
		app.mainView = view
	}

	// form
	if view, err := g.SetView("form", maxX/2-15, maxY/2-4, maxX/2+15, maxY/2+3, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.BgColor = gocui.ColorBlack
		view.Visible = app.formVisible
		view.SetWritePos(0, 1)
		fmt.Fprint(view, "Name:")
		view.SetWritePos(0, 3)
		fmt.Fprint(view, "Value:")
		app.formInput = view
	} else {
		view.Visible = app.formVisible
	}
	// name input
	x0, y0, x1, _ := app.formInput.Dimensions()
	x0 = x0 + 7
	x1--
	y0++
	if view, err := g.SetView("nameInput", x0, y0, x1, y0+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		//view.BgColor = gocui.ColorBlue
		view.Editable = true
		view.Frame = false
		view.Visible = app.formVisible
		app.nameInput = view
		app.listView = append(app.listView, view.Name())
		if app.formVisible {
			g.SetCurrentView(view.Name())
		}
	} else {
		view.Visible = app.formVisible
		if app.formVisible {
			if len(app.mainView.BufferLines()) >= maxY {
				app.mainView.Clear()
			}
			fmt.Fprintln(app.mainView, view.Buffer())

			view.Frame = false
			if app.listView[app.currentView] == view.Name() {
				view.Frame = true
			}
		}
	}
	y0 = y0 + 2
	// value input
	if view, err := g.SetView("valueInput", x0, y0, x1, y0+2, 0); err != nil {
		if !errors.Is(err, gocui.ErrUnknownView) {
			return err
		}
		view.SelBgColor = gocui.Get256Color(238)
		view.Frame = false
		view.Editable = true
		view.Visible = app.formVisible
		app.listView = append(app.listView, view.Name())
		app.valueInput = view
	} else {
		view.Visible = app.formVisible
		if app.formVisible {
			view.Frame = false
			if app.listView[app.currentView] == view.Name() {
				view.Frame = true
			}
		}
	}

	return nil
}

func (app *form) nextView(g *gocui.Gui, v *gocui.View) error {
	// set current bgcolor to back
	// v.BgColor = gocui.ColorBlack
	app.currentView++
	if app.currentView >= len(app.listView) {
		app.currentView = 0
	}
	g.SetCurrentView(app.listView[app.currentView])
	// nv := g.CurrentView()
	// // _, y, _, _ := nv.Dimensions()
	// // app.formInput.SetCursor(0, y)
	// // fmt.Fprintln(app.mainView, y)
	// nv.BgColor = gocui.Get256Color(238)
	return nil
}

func (app *form) showForm(g *gocui.Gui, v *gocui.View) error {
	//app.formVisible = !app.formVisible
	if !app.formVisible {
		app.formVisible = true
		g.SetCurrentView(app.listView[app.currentView])
	} else {
		app.formVisible = false
		g.SetCurrentView(app.mainView.Name())
	}
	return nil
}

func (app *form) keybindings(g *gocui.Gui) error {
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
	if err := g.SetKeybinding("", 'f', gocui.ModNone, app.showForm); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyTab, gocui.ModNone, app.nextView); err != nil {
		return err
	}
	if err := g.SetKeybinding("", gocui.KeyCtrlF, gocui.ModNone, app.showForm); err != nil {
		return err
	}

	return nil
}

func main() {
	var err error
	//gui, err = gocui.NewGui(gocui.OutputNormal, true)
	gui, err = gocui.NewGui(gocui.Output256, true)
	if err != nil {
		panic(err)
	}
	defer gui.Close()
	app := form{
		formVisible: false,
		currentView: 0,
		listView:    make([]string, 0),
	}
	gui.SetManagerFunc(app.layout)

	err = app.keybindings(gui)
	if err != nil {
		panic(err)
	}

	err = gui.MainLoop()
	if err != nil && !errors.Is(err, gocui.ErrQuit) {
		panic(err)
	}
}
