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

type viewCoord struct {
	X0 int
	X1 int
	Y0 int
	Y1 int
}

type FilterFields struct {
	Name   string
	Commit string
	Status string
	Older  string
}
type filterData struct {
	view    *gocui.View
	pos     viewCoord
	visible bool
	fields  FilterFields
	inputs  []*gocui.View
	focus   int
}

// App creates UI and run workflow explorer
type App struct {
	api            internal.GithubApi
	gui            *gocui.Gui
	repoView       *gocui.View
	filterListView *gocui.View
	columnsView    *gocui.View
	columns        []columnsTable
	mainView       *gocui.View
	statusView     *gocui.View
	statusVisible  bool
	status         viewCoord
	filter         filterData
	showRuns       int
	runs           []workflows
}

func NewAppUI(config internal.GithubApi) *App {
	return &App{
		api: config,
		columns: []columnsTable{
			{0, "Workflow Name", 13},
			{1, "Commit Name", 11},
			{2, "Status", 6},
			{3, "Updated", 7},
		},
		statusVisible: true,
		status:        viewCoord{X1: 2, Y1: 2},
		filter: filterData{
			visible: false,
			pos:     viewCoord{X1: 1, Y1: 1},
			inputs:  make([]*gocui.View, 0),
		},
		runs: make([]workflows, 0),
	}
}

func maxInts(x, y int) int {
	if x < y {
		return y
	}
	return x
}

func (app *App) WriteRepoOnwer() {
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		fmt.Fprintf(app.repoView, " Owner: %s\n\n Repo: %s", app.api.Owner, app.api.Repo)
		return nil
	})
}

func (app *App) WriteColumns() {
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		// clear view
		app.columnsView.Clear()

		// write columns
		app.columnsView.SetWritePos(7, 0)
		fmt.Fprintf(
			app.columnsView,
			"%s%s   %s%s   %s%s  %s\n",
			app.columns[0].text,
			strings.Repeat(" ", app.columns[0].spaces-len(app.columns[0].text)),
			app.columns[1].text,
			strings.Repeat(" ", app.columns[1].spaces-len(app.columns[1].text)),
			app.columns[2].text,
			strings.Repeat(" ", app.columns[2].spaces-len(app.columns[2].text)),
			app.columns[3].text,
		)
		return nil
	})
}

func (app *App) WriteMain(keepPosition ...bool) {
	// default value
	keepPos := false
	if len(keepPosition) > 0 {
		keepPos = keepPosition[0]
	}
	app.gui.UpdateAsync(func(g *gocui.Gui) error {
		// delete loading window
		app.statusVisible = false

		var cx, cy, ox, oy int
		// clear view
		if keepPos {
			cx, cy = app.mainView.Cursor()
			ox, oy = app.mainView.Origin()
		}
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
					"  [%s] %s%s   %s%s   %s%s   %s\n",
					toogle,
					run.run.Name,
					strings.Repeat(" ", app.columns[0].spaces-utf8.RuneCountInString(run.run.Name)),
					run.run.Title,
					strings.Repeat(" ", app.columns[1].spaces-utf8.RuneCountInString(run.run.Title)),
					run.run.Conclusion,
					//
					strings.Repeat(" ", app.columns[2].spaces-utf8.RuneCountInString(run.run.Conclusion)),
					run.run.Updated.Format(time.RFC3339),
				)
			}
		}
		app.mainView.SetCursor(cx, cy)
		app.mainView.SetOrigin(ox, oy)
		if app.showRuns == 0 {
			app.mainView.Subtitle = "0/0"
		} else {
			app.mainView.Subtitle = fmt.Sprintf("%d/%d", cy+oy+1, app.showRuns)
		}
		return nil
	})
}

func (app *App) refreshWorkflows() error {
	workflowsRuns, err := app.api.ListWorkflows()
	if err != nil {
		return err
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
	// update show runs
	app.showRuns = len(app.runs)
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
	app.gui.SetManagerFunc(app.Layout)
	//app.gui.SetManager(app)

	// set Keybings
	if err := app.keybindings(app.gui); err != nil {
		return err
	}

	// write dynamic text
	app.WriteRepoOnwer()
	app.StatusView("Requesting workflows\nPlease Wait...")

	go func() {
		// get workflow list
		err := app.refreshWorkflows()
		if err != nil {
			app.StatusView(fmt.Sprintf("Error when requesting:\n%s", err))
		} else {
			app.WriteMain()
			app.WriteColumns()
		}
	}()

	// enter UI mainloop
	if err := app.gui.MainLoop(); err != nil && !errors.Is(err, gocui.ErrQuit) {
		return err
	}
	return nil
}
