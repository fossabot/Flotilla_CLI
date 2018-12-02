/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-02 15:27:31
 */

package UserInterface

import (
	"github.com/ximidar/Flotilla/Flotilla_CLI/UserInterface/CommTab"
	"github.com/ximidar/Flotilla/Flotilla_CLI/UserInterface/CommonBlocks"
	"github.com/ximidar/Flotilla/Flotilla_CLI/UserInterface/FileSystemTab"
	"github.com/ximidar/gocui"
)

// CliGui is an object that will instantiate the ui
type CliGui struct {
	TabList        *CommonBlocks.Tabs
	CommTab        *commtab.CommTab
	FileTab        *FileSystemTab.FileSystemTab
	CurrentTabName string
	RootGUI        *gocui.Gui
}

// NewCliGui is the constructor for CliGui
func NewCliGui() (*CliGui, error) {
	cli := new(CliGui)
	cli.TabList = CommonBlocks.NewTabs(0, 0, "Tabs", cli)
	cli.CurrentTabName = "CommTab"

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

	// Make Tabs
	gui.TabList.AddTab("CommTab", "Comm")
	gui.TabList.AddTab("FileTab", "Files")

	gui.RootGUI.SetManagerFunc(gui.Layout)

	if err := gui.RootGUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		return err
	}

	// if err := gui.RootGUI.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.nextView); err != nil {
	// 	return err
	// }

	err = gui.setupCommTab()
	err = gui.setupFileTab()
	if err != nil {
		return err
	}

	if err := gui.RootGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		panic(err)
		//return err
	}
	return

}

func (gui *CliGui) UpdateTab(name string) {
	gui.CurrentTabName = name

	// Delete all views that are not apart of the currently slected tab
	for _, view := range gui.RootGUI.Views() {
		name = view.Name()
		if name != "CommTab" && name != "FileTab" && name != "Tabs" {
			gui.RootGUI.DeleteKeybindings(name)
			gui.RootGUI.DeleteView(name)
		}
	}

}

// CheckSize makes sure the size of the screen is big enough to accomodate the tool
func (gui *CliGui) CheckSize(x, y int) bool {
	if x > 88 || y > 20 {
		return true
	}
	return false
}

// Layout is a function for Gocui to help layout the screen
func (gui *CliGui) Layout(g *gocui.Gui) error {
	x, y := g.Size()
	if !gui.CheckSize(x, y) {
		return nil
	}

	g.Update(gui.TabList.Layout)

	//Update Tab based on Selected Tab
	switch gui.CurrentTabName {
	case "CommTab":
		g.Update(gui.CommTab.Layout)
	case "FileTab":
		g.Update(gui.FileTab.Layout)
	}

	return nil
}

// CommTabHandlder will controll pulling up the CommTab Contents.
func (gui *CliGui) CommTabHandler(g *gocui.Gui, v *gocui.View) error {
	return nil
}

func (gui *CliGui) setupCommTab() error {

	gui.CommTab = commtab.NewCommTab(0, 3, gui.RootGUI)
	gui.CommTab.Name = "CommContents"

	return nil
}

func (gui *CliGui) setupFileTab() error {
	var err error
	gui.FileTab, err = FileSystemTab.NewFileSystemTab("FileContents", 0, 0, 30, 2)
	return err
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

// func (gui *CliGui) nextView(g *gocui.Gui, v *gocui.View) (err error) {

// 	lentabs := len(gui.TabList)

// 	if gui.CurrentTabNumber >= lentabs {
// 		gui.CurrentTabNumber = 0
// 	} else {
// 		gui.CurrentTabNumber++
// 	}

// 	_, err = setCurrentViewOnTop(g, gui.TabList[gui.CurrentTabNumber])
// 	g.Cursor = true

// 	return err
// }

func (gui *CliGui) quit(g *gocui.Gui, v *gocui.View) error {
	// TODO add a function here that will tell all running tabs to quit
	return gocui.ErrQuit
}
