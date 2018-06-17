/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-06-16 17:13:08
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
    fmt.Println("Hello, I work!")
    screen_init()
  },
}

type Cli_Gui struct{
	Printer_Name string

}

func New_Cli_Gui() *Cli_Gui {
	cgui := new(Cli_Gui)
	cgui.Printer_Name = "Test this var"
	return cgui
}

func screen_init() {
	g, err := gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer g.Close()

	g.SetManagerFunc(layout)

	if err := g.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, quit); err != nil {
		log.Panicln(err)
	}

	if err := g.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
}

func layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()
	if v, err := g.SetView("hello", maxX/2-7, maxY/2, maxX/2+7, maxY/2+2); err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		fmt.Fprintln(v, "Hello world!")
	}
	return nil
}

func quit(g *gocui.Gui, v *gocui.View) error {
	return gocui.ErrQuit
}
 
