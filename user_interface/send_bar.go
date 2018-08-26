/*
* @Author: Ximidar
* @Date:   2018-08-25 21:59:56
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-26 13:26:07
*/
package user_interface

import(
	"github.com/jroimartin/gocui"
	"fmt"
)

type Send_Bar struct{
	name    string
	x, y    int
	handler func(message string)
}

func New_Send_Bar(name string, x,y int, handler func(message string)) *Send_Bar{
	return &Send_Bar{name:name,
					 x:x,
					 y:y,
					 handler:handler}
}

func (w *Send_Bar) Layout(g *gocui.Gui) error{
	maxX, _ := g.Size()

	if v, err := g.SetView(w.name, w.x, w.y, maxX -1 , w.y + 2); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title = "Send"
		v.Editable = true
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.send_view_clear); err != nil {
			return err
		}
	}
	return nil
}

func (w *Send_Bar) send_view_clear(g *gocui.Gui, v *gocui.View) error {
	contents := v.Buffer()
	v.Clear()
	v.SetCursor(v.Origin())

	// send the contents somewhere
	w.handler(contents)

	return nil
}
