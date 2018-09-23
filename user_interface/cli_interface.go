/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-09-22 23:13:55
 */

package user_interface

import (
	"fmt"
	"github.com/ximidar/gocui"
	"github.com/ximidar/mango_cli/mango_interface"
	"github.com/nats-io/go-nats"
	"log"
	"strconv"
	_"time"
)

type Cli_Gui struct {
	reader_active bool
	RootGUI       *gocui.Gui

	Connection_Info   string
	Monitor_View      string
	Monitor           Monitor_Interface
	Send_View         string
	Baud_Button       string
	Port_Button       string
	Connect_Button    string
	Disconnect_Button string
	Info_View         string

	port string
	baud int32

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
	gui.Disconnect_Button = "disconnect_button"
	gui.Info_View = "info_view"

	gui.port = ""
	gui.baud = -1

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

	gui.reader_active = true
	go gui.Comm_Relay()

	if err := gui.RootGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return
}

func (gui *Cli_Gui) quit(g *gocui.Gui, v *gocui.View) error {
	gui.reader_active = false
	return gocui.ErrQuit
}

func (gui *Cli_Gui) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	gui.Monitor = New_Monitor(gui.Monitor_View, 31, 0)
	send_bar := New_Send_Bar(gui.Send_View, 31, maxY-3, gui.write_to_comm)
	gui.connection_info_layout(g)
	port_button := New_Explode_Button(gui.Port_Button, 0, 8, 14, "Port Select", gui.get_ports, gui.port_select)
	baud_button := New_Explode_Button(gui.Baud_Button, 15, 8, 15, "Baud Select", gui.get_bauds, gui.baud_select)
	connect_button := New_Button(gui.Connect_Button, 0, 11, 30, "Connect", gui.connect_comm)
	disconnect_button := New_Button(gui.Disconnect_Button, 0, 14, 30, "Disconnect", gui.disconnect_comm)
	g.Update(gui.Monitor.Layout)
	g.Update(send_bar.Layout)
	g.Update(port_button.Layout)
	g.Update(baud_button.Layout)
	g.Update(connect_button.Layout)
	g.Update(disconnect_button.Layout)
	return nil
}

func (gui *Cli_Gui) write_to_comm(mess string) {
	//gui.Monitor.Write(gui.RootGUI, mess)
	gui.Mango.Comm_Write(mess)
}

func (gui *Cli_Gui) get_bauds() []string {
	return []string{"250000", "230400", "115200", "57600", "38400", "19200", "9600"}
}

func (gui *Cli_Gui) baud_select(selection string) {
	gui.Monitor.Write(gui.RootGUI, fmt.Sprintf("Selection %v ", selection))
	if temp_baud, err := strconv.Atoi(selection); err == nil {
		gui.baud = int32(temp_baud)
	} else {
		gui.Monitor.Write(gui.RootGUI, "default to 115200")
		gui.baud = 115200
	}
}

func (gui *Cli_Gui) connect_comm(g *gocui.Gui, v *gocui.View) error {
	gui.Monitor.Write(g, "connect!")
	gui.Mango.Comm_Set_Connection_Options(gui.port, gui.baud)
	gui.Mango.Comm_Connect()
	return nil
}

func (gui *Cli_Gui) disconnect_comm(g *gocui.Gui, v *gocui.View) error {
	gui.Monitor.Write(g, "disconnect!")
	gui.Mango.Comm_Disconnect()
	return nil
}

func (gui *Cli_Gui) get_ports() []string {
	ports, err := gui.Mango.Comm_Get_Available_Ports()

	if err != nil {
		return []string{"Check commango daemon"}
	}

	if len(ports) == 0 {
		return []string{"no ports available"}
	}

	return ports
}

func (gui *Cli_Gui) port_select(selection string) {
	gui.Monitor.Write(gui.RootGUI, fmt.Sprintf("Selection %v ", selection))
	gui.port = selection
}

func (gui *Cli_Gui) connection_info_layout(g *gocui.Gui) (err error) {
	if v, err := g.SetView(gui.Connection_Info, 0, 0, 30, 7); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(g.Size())
			panic(err)
		}
		v.Title = "Connection Info"
		status, err := gui.Mango.Comm_Get_Status()
		if err != nil{
			fmt.Fprintln(v, err.Error())
		}

		fmt.Fprintln(v, fmt.Sprintf("Port: %v\nBaud: %v\nConnected: %v", status.Port, status.Baud, status.Connected))
		update_status := func(msg *nats.Msg){
			newstatus, err := gui.Mango.Deconstruct_status(msg.Data)
			if err != nil{
				fmt.Fprintln(v, err.Error())
			}
			v.Clear()
			v.SetCursor(v.Origin())
			fmt.Fprintln(v, fmt.Sprintf("Port: %v\nBaud: %v\nConnected: %v", newstatus.Port, newstatus.Baud, newstatus.Connected))
		}

		gui.Mango.NC.Subscribe("commango.status_update", update_status)
	}

	return nil
}

func (gui *Cli_Gui) Comm_Relay() {

	for gui.reader_active {
		select {
		case read := <-gui.Mango.Emit_Line:
			gui.Monitor.Write(gui.RootGUI, read)
		}
		// busy_sleeper := time.Duration(50 * time.Millisecond)
		// time.Sleep(busy_sleeper)

	}
}
