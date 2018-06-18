/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-06-17 18:17:16
*/

package cmd

import(
	"fmt"
	"log"
	"time"
	"os"
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
	reader_active bool
	RootGUI *gocui.Gui
	Connection_Info string
	Monitor_View string
	Send_View string
}

func New_Cli_Gui() *Cli_Gui {
	cgui := new(Cli_Gui)
	cgui.Printer_Name, _ = os.Hostname() 
	cgui.reader_active = false

	// names
	cgui.Connection_Info = "connection_info"
	cgui.Monitor_View = "monitor_view"
	cgui.Send_View = "send_view"
	return cgui
}

func (gui *Cli_Gui) screen_init() (err error){
	gui.RootGUI, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.RootGUI.Close()

	gui.RootGUI.Cursor = true

	gui.RootGUI.SetManagerFunc(gui.layout)

	if err := gui.RootGUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.RootGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}

	fmt.Println("Exiting!")
	return
}

func (gui *Cli_Gui) layout(g *gocui.Gui) (err error) {
	maxX, maxY := g.Size()
	if v, err := g.SetView(gui.Connection_Info, 0, 0, maxX/5, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title = gui.Printer_Name
		v.Autoscroll = true
		v.Wrap = true
		gui.reader_active = true
		go gui.Reader_fmt()
	}

	return nil
}

func (gui *Cli_Gui) quit(g *gocui.Gui, v *gocui.View) error {
	gui.reader_active = false
	return gocui.ErrQuit
}

func (gui *Cli_Gui) Reader_fmt() {
	
	counter := 0
	for gui.reader_active{
		time.Sleep(500 * time.Millisecond)
		counter += 1
		gui.RootGUI.Update(func(g *gocui.Gui) error {
			v, err := g.View(gui.Connection_Info)
			if err != nil {
				log.Println(err)
				return err
			}
			
			mess := fmt.Sprintf("Hello at count: %v", counter)
			fmt.Fprintln(v, mess)
			return err
		})
			
	}
}


