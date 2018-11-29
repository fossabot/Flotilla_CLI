/*
* @Author: Ximidar
* @Date:   2018-11-29 13:14:25
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-11-29 13:27:44
 */

// Package commtab is the user interface for connecting and monitoring
// the serial line for flotilla
package commtab

import (
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/nats-io/go-nats"
	"github.com/ximidar/Flotilla/Flotilla_CLI/FlotillaInterface"
	"github.com/ximidar/Flotilla/Flotilla_CLI/UserInterface/CommonBlocks"
	"github.com/ximidar/gocui"
)

// CommTab Creates a Command Line GUI
type CommTab struct {
	readerActive bool
	RootGUI      *gocui.Gui

	ConnectionInfo   string
	MonitorView      string
	Monitor          MonitorInterface
	SendView         string
	BaudButton       string
	PortButton       string
	ConnectButton    string
	DisconnectButton string
	InfoView         string

	port string
	baud int32

	FlotillaInterface *FlotillaInterface.FlotillaInterface
}

// NewCommTab will Create a CommTab object
func NewCommTab() *CommTab {
	gui := new(CommTab)
	gui.readerActive = false

	// names
	gui.ConnectionInfo = "connection_info"
	gui.MonitorView = "monitor_view"
	gui.SendView = "send_view"
	gui.BaudButton = "baud_button"
	gui.PortButton = "port_button"
	gui.ConnectButton = "connect_button"
	gui.DisconnectButton = "disconnect_button"
	gui.InfoView = "info_view"

	gui.port = ""
	gui.baud = -1

	var err error
	gui.FlotillaInterface, err = FlotillaInterface.NewFlotillaInterface()
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

func (gui *CommTab) nextView(g *gocui.Gui, v *gocui.View) (err error) {

	_, err = setCurrentViewOnTop(g, gui.SendView)
	g.Cursor = true

	return err
}

// ScreenInit will initialize the screen
func (gui *CommTab) ScreenInit() (err error) {
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

	gui.readerActive = true
	gui.CommRelay()

	if err := gui.RootGUI.MainLoop(); err != nil && err != gocui.ErrQuit {
		log.Panicln(err)
	}
	return
}

func (gui *CommTab) quit(g *gocui.Gui, v *gocui.View) error {
	gui.readerActive = false
	return gocui.ErrQuit
}

// Layout is CommTab's gocui Layout Function
func (gui *CommTab) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	gui.Monitor = NewMonitor(gui.MonitorView, 31, 0)
	sendBar := NewSendBar(gui.SendView, 31, maxY-3, gui.writeToComm)
	gui.connectionInfoLayout(g)
	portButton := CommonBlocks.NewExplodeButton(gui.PortButton, 0, 8, 14, "Port Select", gui.getPorts, gui.portSelect)
	baudButton := CommonBlocks.NewExplodeButton(gui.BaudButton, 15, 8, 15, "Baud Select", gui.getBauds, gui.baudSelect)
	connectButton := CommonBlocks.NewButton(gui.ConnectButton, 0, 11, 30, "Connect", gui.connectComm)
	disconnectButton := CommonBlocks.NewButton(gui.DisconnectButton, 0, 14, 30, "Disconnect", gui.disconnectComm)
	g.Update(gui.Monitor.Layout)
	g.Update(sendBar.Layout)
	g.Update(portButton.Layout)
	g.Update(baudButton.Layout)
	g.Update(connectButton.Layout)
	g.Update(disconnectButton.Layout)
	return nil
}

func (gui *CommTab) writeToComm(mess string) {
	//gui.Monitor.Write(gui.RootGUI, mess)
	gui.FlotillaInterface.CommWrite(mess)
}

func (gui *CommTab) getBauds() []string {
	return []string{"250000", "230400", "115200", "57600", "38400", "19200", "9600"}
}

func (gui *CommTab) baudSelect(selection string) {
	gui.Monitor.Write(gui.RootGUI, fmt.Sprintf("Selection %v ", selection))
	if tempBaud, err := strconv.Atoi(selection); err == nil {
		gui.baud = int32(tempBaud)
	} else {
		gui.Monitor.Write(gui.RootGUI, "default to 115200")
		gui.baud = 115200
	}
}

func (gui *CommTab) connectComm(g *gocui.Gui, v *gocui.View) error {
	gui.Monitor.Write(g, "connect!")
	gui.FlotillaInterface.CommSetConnectionOptions(gui.port, gui.baud)
	gui.FlotillaInterface.CommConnect()
	return nil
}

func (gui *CommTab) disconnectComm(g *gocui.Gui, v *gocui.View) error {
	gui.Monitor.Write(g, "disconnect!")
	gui.FlotillaInterface.CommDisconnect()
	return nil
}

func (gui *CommTab) getPorts() []string {
	ports, err := gui.FlotillaInterface.CommGetAvailablePorts()

	if err != nil {
		return []string{"Check comFlotillaInterface daemon"}
	}

	if len(ports) == 0 {
		return []string{"no ports available"}
	}

	return ports
}

func (gui *CommTab) portSelect(selection string) {
	gui.Monitor.Write(gui.RootGUI, fmt.Sprintf("Selection %v ", selection))
	gui.port = selection
}

func (gui *CommTab) connectionInfoLayout(g *gocui.Gui) (err error) {
	if v, err := g.SetView(gui.ConnectionInfo, 0, 0, 30, 7); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(g.Size())
			panic(err)
		}
		v.Title = "Connection Info"
		status, err := gui.FlotillaInterface.CommGetStatus()
		if err != nil {
			fmt.Fprintln(v, err.Error())
		}

		fmt.Fprintln(v, fmt.Sprintf("Port: %v\nBaud: %v\nConnected: %v", status.Port, status.Baud, status.Connected))
		updateStatus := func(msg *nats.Msg) {
			newStatus, err := gui.FlotillaInterface.DeconstructStatus(msg.Data)
			if err != nil {
				fmt.Fprintln(v, err.Error())
			}
			v.Clear()
			v.SetCursor(v.Origin())
			fmt.Fprintln(v, fmt.Sprintf("Port: %v\nBaud: %v\nConnected: %v", newStatus.Port, newStatus.Baud, newStatus.Connected))
		}

		gui.FlotillaInterface.NC.Subscribe("comFlotillaInterface.status_update", updateStatus)
	}

	return nil
}

// CommRelay will subscribes functions to incoming data from Nats
func (gui *CommTab) CommRelay() {

	gui.FlotillaInterface.NC.Subscribe("comFlotillaInterface.read_line", gui.CommReadSub)
	gui.FlotillaInterface.NC.Subscribe("comFlotillaInterface.write_line", gui.CommWriteSub)

}

// CommReadSub will reveive Comm Messages from the Nats Server
func (gui *CommTab) CommReadSub(msg *nats.Msg) {

	data := string(msg.Data)
	data = fmt.Sprintf("\u001b[46mRECV:\u001b[0m \n%v", data)
	data = strings.Replace(data, "\n", "\n\t", -1)
	data = strings.Replace(data, "echo:", "", -1)
	gui.Monitor.Write(gui.RootGUI, data)
}

// CommWriteSub will Revieve updates from the Nats server on Written Messages
func (gui *CommTab) CommWriteSub(msg *nats.Msg) {
	data := string(msg.Data)
	data = strings.Replace(data, "\n", "", -1)
	data = fmt.Sprintf("\u001b[44mSENT: %v\u001b[0m", data)
	gui.Monitor.Write(gui.RootGUI, data)
}
