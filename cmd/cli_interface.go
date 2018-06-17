/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-06-16 17:52:16
*/

package cmd

import(
	"fmt"
	"log"
	"github.com/spf13/cobra"
	"github.com/jroimartin/gocui"

)

func init() {
  rootCmd.AddCommand(printerface)
}

var printerface = &cobra.Command{
  Use:   "ui",
  Short: "Show the cli UI for MangoOS",
  Long:  `This will open the cli UI for MangoOS. This has tools for monitoring the command line and starting prints (or it will in the future)`,
  Run: func(cmd *cobra.Command, args []string) {
  	cligui := New_Cli_Gui()
  	cligui.screen_init()
  },
}

type Cli_Gui struct{
	Printer_Name string
	RootGUI *gocui.Gui
	Connection_Info *gocui.View
	Monitor_View *gocui.View
	Send_View *gocui.View
}

func New_Cli_Gui() *Cli_Gui {
	cgui := new(Cli_Gui)
	cgui.Printer_Name = "Test this var"
	return cgui
}

func (gui *Cli_Gui) screen_init() (err error){
	gui.RootGUI, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.RootGUI.Close()

	gui.RootGUI.SetManagerFunc(gui.layout)

	if err := gui.RootGUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.RootGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return
}

func (gui *Cli_Gui) layout(g *gocui.Gui) (err error) {
	maxX, maxY := g.Size()
	if gui.Connection_Info, err = g.SetView(gui.Printer_Name, 0, 0, maxX/5, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(gui.Connection_Info, gui.Printer_Name)
	}

	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
 
