package ui

import (
	"errors"
	"fmt"
	"strings"
	"time"
	"unicode/utf8"

	"github-workflow-explorer/internal"

	"github.com/awesome-gocui/gocui"
)

// workflows stores runs data and interface data, if is toogle and filtered
type workflows struct {
	line   int
	show   bool
	toogle bool
	run    internal.WorkflowRun
}

type columnsTable struct {
	index  int
	text   string
	spaces int
}

type viewPosition struct {
	x0 int
	y0 int
	x1 int
	y1 int
}

// App creates UI and run workflow explorer
type App struct {
	api           internal.GithubApi
	gui           *gocui.Gui
	repoView      *gocui.View
	filterView    *gocui.View
	columnsView   *gocui.View
	statusView    *gocui.View
	statusPos     viewPosition
	statusVisible bool
	columns       []columnsTable
	mainView      *gocui.View
	runs          []workflows
}

func NewAppUI(config internal.GithubApi) *App {
	return &App{
		api: config,
		columns: []columnsTable{
			{0, "Workflow Name", 13},
			{1, "Commit Name", 11},
			{2, "Status", 6},
		},
		runs: make([]workflows, 0),
		statusPos: viewPosition{
			x0: 0,
			x1: 0,
			y0: 1,
			y1: 1,
		},
		statusVisible: false,
	}
}

func maxInts(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (app *App) WriteRepoOnwer() {
	// app.gui.UpdateAsync(func(g *gocui.Gui) error {
	app.gui.Update(func(g *gocui.Gui) error {
		fmt.Fprintf(app.repoView, "Owner: %s\n\n Repo:%s", app.api.Owner, app.api.Repo)
		return nil
	})
}

func (app *App) WriteFilter() {
	// app.gui.UpdateAsync(func(g *gocui.Gui) error {
	app.gui.Update(func(g *gocui.Gui) error {
		// change this writing
		fmt.Fprintf(app.filterView, "Owner: %s\n\n Repo:%s", app.api.Owner, app.api.Repo)
		return nil
	})
}

func (app *App) WriteColumns() {
	// app.gui.UpdateAsync(func(g *gocui.Gui) error {
	app.gui.Update(func(g *gocui.Gui) error {
		// clear view
		app.columnsView.Clear()

		// write columns
		app.columnsView.SetWritePos(7, 0)
		fmt.Fprintf(
			app.columnsView,
			"%s%s   %s%s   %s\n",
			app.columns[0].text,
			strings.Repeat(" ", app.columns[0].spaces-len(app.columns[0].text)),
			app.columns[1].text,
			strings.Repeat(" ", app.columns[1].spaces-len(app.columns[1].text)),
			app.columns[2].text,
		)
		return nil
	})
}

func (app *App) WriteMain() {
	// app.gui.UpdateAsync(func(g *gocui.Gui) error {
	app.gui.Update(func(g *gocui.Gui) error {
		// hide status window
		app.statusView.Visible = false

		// clear view
		cx, cy := app.mainView.Cursor()
		ox, oy := app.mainView.Origin()
		app.mainView.Clear()

		// write lines
		// assuming one line por item, change this when smart print
		for i, run := range app.runs {
			// update line
			app.runs[i].line = i
			if run.show {
				// write line
				toogle := " "
				if run.toogle {
					toogle = "*"
				}
				fmt.Fprintf(
					app.mainView,
					"  [%s] %s%s   %s%s   %s%d %d\n",
					toogle,
					run.run.Name,
					strings.Repeat(" ", app.columns[0].spaces-utf8.RuneCountInString(run.run.Name)),
					run.run.Title,
					strings.Repeat(" ", app.columns[1].spaces-utf8.RuneCountInString(run.run.Title)),
					run.run.Conclusion,
					i,
					run.run.ID,
				)
			}
		}
		app.mainView.SetCursor(cx, cy)
		app.mainView.SetOrigin(ox, oy)
		return nil
	})
}

func (app *App) refreshWorkflows() error {
	// workflowsRuns, err := app.api.ListWorkflows()
	// if err != nil {
	// 	return err
	// }
	workflowsRuns := make([]internal.WorkflowRun, 20)
	for i := range workflowsRuns {
		workflowsRuns[i] = internal.WorkflowRun{
			Name:       fmt.Sprintf("NAME  %d  NAME", i),
			Title:      "Some title to fill space",
			Conclusion: "done",
			ID:         i + 1,
		}
	}
	for i, workflowRun := range workflowsRuns {
		// calculates columns size
		app.columns[0].spaces = maxInts(app.columns[0].spaces, utf8.RuneCountInString(workflowRun.Name))
		app.columns[1].spaces = maxInts(app.columns[1].spaces, utf8.RuneCountInString(workflowRun.Title))
		app.columns[2].spaces = maxInts(app.columns[2].spaces, utf8.RuneCountInString(workflowRun.Conclusion))
		app.runs = append(app.runs, workflows{
			line:   i,
			show:   true,
			toogle: false,
			run:    workflowRun,
		})
	}
	return nil
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
	app.gui.SetManagerFunc(app.layout)

	// set Keybings
	if err := app.keybindings(app.gui); err != nil {
		return err
	}

	// write dynamic text
	app.WriteRepoOnwer()
	app.WriteFilter()
	// app.StatusView("Requesting workflows\nPlease Wait...")
	app.setStatus("Requesting workflows\nPlease Wait...")

	go func() {
		time.Sleep(time.Second * 1)
		// app.statusView.Visible = false
		// get workflow list
		app.refreshWorkflows()

		app.WriteMain()
		app.WriteColumns()
	}()

	// enter UI mainloop
	if err := app.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
