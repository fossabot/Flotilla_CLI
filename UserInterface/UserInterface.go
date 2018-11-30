/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-11-30 14:41:06
 */

package UserInterface

import (
	"github.com/ximidar/Flotilla/Flotilla_CLI/UserInterface/CommTab"
	"github.com/ximidar/gocui"
)

// CliGui is an object that will instantiate the ui
type CliGui struct {
	TabList          []string
	CurrentTabNumber int
	RootGUI          *gocui.Gui
}

// NewCliGui is the constructor for CliGui
func NewCliGui() (*CliGui, error) {
	cli := new(CliGui)
	cli.TabList = append(cli.TabList, "CommTab")
	cli.CurrentTabNumber = 0

	return cli, nil
}

// ScreenInit will initialize the screen for the gui
func (gui *CliGui) ScreenInit() (err error) {
	gui.RootGUI, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		return err
	}
	defer gui.RootGUI.Close()

	gui.RootGUI.Cursor = true
	gui.RootGUI.Mouse = true
	gui.RootGUI.Highlight = true
	gui.RootGUI.SelFgColor = gocui.ColorGreen

	gui.RootGUI.SetManagerFunc(gui.Layout)

	if err := gui.RootGUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		return err
	}

	// if err := gui.RootGUI.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.nextView); err != nil {
	// 	return err
	// }

	if err := gui.RootGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		return err
	}
	return

}

// Layout is a function for Gocui to help layout the screen
func (gui *CliGui) Layout(g *gocui.Gui) error {
	err := gui.setupCommTab(g)

	if err != nil {
		return err
	}
	return nil
}

func (gui *CliGui) setupCommTab(g *gocui.Gui) error {
	CommTab := commtab.NewCommTab(g)
	CommTab.Name = "CommTab"
	g.Update(CommTab.Layout)
	return nil
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	view, err := g.SetViewOnTop(name)
	if err != nil {
		view.SetCursor(view.Origin())
	}

	return view, err
}

func (gui *CliGui) nextView(g *gocui.Gui, v *gocui.View) (err error) {

	lentabs := len(gui.TabList)

	if gui.CurrentTabNumber >= lentabs {
		gui.CurrentTabNumber = 0
	} else {
		gui.CurrentTabNumber++
	}

	_, err = setCurrentViewOnTop(g, gui.TabList[gui.CurrentTabNumber])
	g.Cursor = true

	return err
}

func (gui *CliGui) quit(g *gocui.Gui, v *gocui.View) error {
	// TODO add a function here that will tell all running tabs to quit
	return gocui.ErrQuit
}
