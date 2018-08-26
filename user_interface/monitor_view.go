/*
* @Author: Ximidar
* @Date:   2018-08-25 22:00:52
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-26 13:26:15
*/
package user_interface

import(
	"github.com/jroimartin/gocui"
	"fmt"
	"os"
)

type Monitor_Interface interface{
	Write(mess string)
	Read() string
	Clear()
	Layout(g *gocui.Gui) error
}

type Monitor struct{
	name    string
	x, y int
	g *gocui.Gui
	v *gocui.View
}

func New_Monitor(name string, x,y int) Monitor_Interface{
	return &Monitor{name:name, x:x, y:y}
}

func (w *Monitor) Layout(g *gocui.Gui) error{
	maxX, maxY := g.Size()	
	w.g = g

	if v, err := g.SetView(w.name, w.x, w.y, maxX - 1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title, err = os.Hostname() 
		if err != nil{
			return err
		}
		v.Autoscroll = true
		v.Wrap = true
		w.v = v
	}
	return nil
}

func (w *Monitor) Write(mess string){

	w.g.Update(func(g *gocui.Gui) error {
		v, err := g.View(w.name)
		if err != nil {
			return err
		}		
		fmt.Fprintln(v, mess)
		return err
		
	})

}

func (w *Monitor) Clear() {
	w.v.Clear()
	w.v.SetCursor(w.v.Origin())
}

func (w *Monitor) Read() string{
	return w.v.Buffer()
}