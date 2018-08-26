/*
* @Author: Ximidar
* @Date:   2018-06-16 16:39:58
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 22:43:56
*/

package user_interface

import(
	"fmt"
	"log"
	_"time"
	"os"
	"github.com/jroimartin/gocui"
	"github.com/ximidar/mango_cli/mango_interface"
)

type Cli_Gui struct{
	Printer_Name string
	reader_active bool
	RootGUI *gocui.Gui

	Connection_Info string
	Monitor_View string
	Send_View string
	Baud_Button string
	Port_Button string
	Connect_Button string
	Info_View string

	port string
	baud string

	Mango *mango_interface.Mango
}

func New_Cli_Gui() *Cli_Gui {
	gui := new(Cli_Gui)
	gui.Printer_Name, _ = os.Hostname() 
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
	if err != nil{
		panic(err)
	}
	return gui
}

func setCurrentViewOnTop(g *gocui.Gui, name string) (*gocui.View, error) {
	if _, err := g.SetCurrentView(name); err != nil {
		return nil, err
	}
	view, err := g.SetViewOnTop(name)
	if err != nil{
		view.SetCursor(view.Origin())
	}

	return view, err
}


func (gui *Cli_Gui) nextView(g *gocui.Gui, v *gocui.View) (err error){
	
	if v.Name() == gui.Connection_Info{
		_, err = setCurrentViewOnTop(g, gui.Send_View)
		g.Cursor = true
	} else {
		_, err = setCurrentViewOnTop(g, gui.Connection_Info)
		g.Cursor = false
	}

	return err
}

func (gui *Cli_Gui) Screen_Init() (err error){
	gui.RootGUI, err = gocui.NewGui(gocui.OutputNormal)
	if err != nil {
		log.Panicln(err)
	}
	defer gui.RootGUI.Close()

	gui.RootGUI.Cursor = true
	gui.RootGUI.Mouse = true
	gui.RootGUI.Highlight = true
	gui.RootGUI.SelFgColor = gocui.ColorGreen

	gui.RootGUI.SetManagerFunc(gui.layout)

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

func (gui *Cli_Gui) layout(g *gocui.Gui) error{
	gui.send_view_layout(g)
	gui.monitor_view_layout(g)
	gui.connection_info_layout(g)
	exb := New_Explode_Button("test", 0, 8, 30, "explode", return_strings)
	g.Update(exb.Layout)
	
	return nil
}

func return_strings() []string{
	return []string{"Hello", "My", "name", "is", "Matt"}
}

func (gui *Cli_Gui) connection_info_layout(g *gocui.Gui) (err error) {
	//maxX, _ := g.Size()
	if v, err := g.SetView(gui.Connection_Info, 0, 0, 30, 7); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(g.Size())
			panic(err)
		}
		v.Title = "Connection Info"
		
	}

	return nil
}

func msg_me(g *gocui.Gui, v *gocui.View) error{
	view, err := g.View("monitor_view")
	if err != nil {
		log.Println(err)
		return err
	}
	x,y := g.Size()
	fmt.Fprintln(view, fmt.Sprintf("button press %v %v", x, y))
	
	return err
}

func (gui *Cli_Gui) send_view_layout(g *gocui.Gui) (err error){
	maxX, maxY := g.Size()

	if v, err := g.SetView(gui.Send_View, 30 + 1, maxY - maxY/10, maxX - 1, (maxY - maxY/10) + 2); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title = "Send"
		v.Editable = true
		if _, err = setCurrentViewOnTop(g, gui.Send_View); err != nil {
			return err
		}
	}
	gui.RootGUI.SetKeybinding(gui.Send_View, gocui.KeyEnter, gocui.ModNone, gui.send_view_clear)
	return nil

}

func (gui *Cli_Gui) send_view_clear(g *gocui.Gui, v *gocui.View) error {
	contents := v.Buffer()
	gui.Mango.Comm_Write(contents)
	v.Clear()
	v.SetCursor(v.Origin())

	return nil
}

func (gui *Cli_Gui) monitor_view_layout(g *gocui.Gui) (err error){
	maxX, maxY := g.Size()

	if v, err := g.SetView(gui.Monitor_View, 30 + 1, 0, maxX - 1, (maxY - maxY/10) - 1); err != nil {
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
	
	reader, err := gui.Mango.Get_Comm_Signal()
	if err != nil{
		panic(err)
	}
	for gui.reader_active{
		select{
		case read := <- reader:
			gui.RootGUI.Update(func(g *gocui.Gui) error {
				v, err := g.View(gui.Monitor_View)
				if err != nil {
					log.Println(err)
					return err
				}
				
				mess, ok := read.Body[0].(string)
				if ok{
					fmt.Fprintln(v, mess)
				}
				return err
			})
		default:
			continue

		}
		
			
	}
}


