/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-26 13:53:29
 */

package user_interface

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"github.com/ximidar/mango_cli/mango_interface"
	"log"
)

type Cli_Gui struct {
	reader_active bool
	RootGUI       *gocui.Gui

	Connection_Info string
	Monitor_View    string
	Monitor         Monitor_Interface
	Send_View       string
	Baud_Button     string
	Port_Button     string
	Connect_Button  string
	Info_View       string

	port string
	baud string

	Mango *mango_interface.Mango
}

func New_Cli_Gui() *Cli_Gui {
	gui := new(Cli_Gui)
	gui.reader_active = false

	// names
	gui.Connection_Info = "connection_info"
	gui.Monitor_View = "monitor_view"
	gui.Send_View = "send_view"
	gui.Baud_Button = "baud_button"
	gui.Port_Button = "port_button"
	gui.Connect_Button = "connect_button"
	gui.Info_View = "info_view"
	var err error
	gui.Mango, err = mango_interface.NewMango()
	if err != nil {
		panic(err)
	}
	return gui
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

func (gui *Cli_Gui) nextView(g *gocui.Gui, v *gocui.View) (err error) {

	if v.Name() == gui.Connection_Info {
		_, err = setCurrentViewOnTop(g, gui.Send_View)
		g.Cursor = true
	} else {
		_, err = setCurrentViewOnTop(g, gui.Connection_Info)
		g.Cursor = false
	}

	return err
}

func (gui *Cli_Gui) Screen_Init() (err error) {
	gui.RootGUI, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.RootGUI.Close()

	gui.RootGUI.Cursor = true
	gui.RootGUI.Mouse = true
	gui.RootGUI.Highlight = true
	gui.RootGUI.SelFgColor = gocui.ColorGreen

	gui.RootGUI.SetManagerFunc(gui.Layout)

	if err := gui.RootGUI.SetKeybinding("", gocui.KeyCtrlC, gocui.ModNone, gui.quit); err != nil {
		log.Panicln(err)
	}

	if err := gui.RootGUI.SetKeybinding("", gocui.KeyTab, gocui.ModNone, gui.nextView); err != nil {
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

func (gui *Cli_Gui) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	gui.Monitor = New_Monitor(gui.Monitor_View, 31, 0)
	send_bar := New_Send_Bar(gui.Send_View, 31, maxY-3, gui.write_to_monitor)
	gui.connection_info_layout(g)
	exb := New_Explode_Button("test", 0, 8, 30, "explode", gui.Info_Loader, gui.Selection_Callback)
	g.Update(gui.Monitor.Layout)
	g.Update(send_bar.Layout)
	g.Update(exb.Layout)

	return nil
}

func (gui *Cli_Gui) write_to_monitor(mess string) {
	gui.Monitor.Write(gui.RootGUI, mess)
}

func (gui *Cli_Gui) Info_Loader() []string {
	return []string{"Hello", "My", "name", "is", "Matt"}
}

func (gui *Cli_Gui) Selection_Callback(selection string) {
	gui.Monitor.Write(gui.RootGUI, fmt.Sprintf("Selection %v ", selection))
}

func (gui *Cli_Gui) connection_info_layout(g *gocui.Gui) (err error) {
	if v, err := g.SetView(gui.Connection_Info, 0, 0, 30, 7); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(g.Size())
			panic(err)
		}
		v.Title = "Connection Info"

	}

	return nil
}

func (gui *Cli_Gui) Reader_fmt() {

	reader, err := gui.Mango.Get_Comm_Signal()
	if err != nil {
		panic(err)
	}
	for gui.reader_active {
		select {
		case read := <-reader:
			mess, ok := read.Body[0].(string)
			if ok {
				gui.Monitor.Write(gui.RootGUI, mess)
			}
		default:
			continue

		}

	}
}
