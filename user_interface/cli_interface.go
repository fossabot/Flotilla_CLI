/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 10:09:28
*/

package user_interface

import(
	"fmt"
	"log"
	"time"
	"os"
	"github.com/jroimartin/gocui"
)

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

func (gui *Cli_Gui) Screen_Init() (err error){
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

func (gui *Cli_Gui) quit(g *gocui.Gui, v *gocui.View) error {
	gui.reader_active = false
	return gocui.ErrQuit
}

func (gui *Cli_Gui) layout(g *gocui.Gui) error{
	gui.send_view_layout(g)
	gui.monitor_view_layout(g)
	gui.connection_info_layout(g)
	
	return nil
}

func (gui *Cli_Gui) connection_info_layout(g *gocui.Gui) (err error) {
	maxX, maxY := g.Size()
	if v, err := g.SetView(gui.Connection_Info, 0, 0, maxX/5, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title = "Connection Info"
	}

	return nil
}

func (gui *Cli_Gui) send_view_layout(g *gocui.Gui) (err error){
	maxX, maxY := g.Size()

	if v, err := g.SetView(gui.Send_View, maxX/5 + 1, maxY - maxY/10, maxX - 1, maxY-1); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title = "Send"
	}

	return nil

}

func (gui *Cli_Gui) monitor_view_layout(g *gocui.Gui) (err error){
	maxX, maxY := g.Size()

	if v, err := g.SetView(gui.Monitor_View, maxX/5 + 1, 0, maxX - 1, (maxY - maxY/10) - 1); err != nil {
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

func (gui *Cli_Gui) Reader_fmt() {
	
	counter := 0
	for gui.reader_active{
		time.Sleep(10 * time.Millisecond)
		counter += 1
		gui.RootGUI.Update(func(g *gocui.Gui) error {
			v, err := g.View(gui.Monitor_View)
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


