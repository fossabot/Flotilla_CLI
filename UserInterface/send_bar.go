/*
* @Author: Ximidar
* @Date:   2018-08-25 21:59:56
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-10-17 13:22:39
 */
package UserInterface

import (
	"fmt"

	"github.com/ximidar/gocui"
)

// SendBar is a gui object for sending strings
type SendBar struct {
	name    string
	x, y    int
	handler func(message string)
}

// NewSendBar will construct a SendBar Object
func NewSendBar(name string, x, y int, handler func(message string)) *SendBar {
	return &SendBar{name: name,
		x:       x,
		y:       y,
		handler: handler}
}

// Layout is SendBar's Gocui Layout Function
func (w *SendBar) Layout(g *gocui.Gui) error {
	maxX, _ := g.Size()

	if v, err := g.SetView(w.name, w.x, w.y, maxX-1, w.y+2); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title = "Send"
		v.Editable = true
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.sendViewClear); err != nil {
			return err
		}
	}
	return nil
}

func (w *SendBar) sendViewClear(g *gocui.Gui, v *gocui.View) error {
	contents := v.Buffer()
	v.Clear()
	v.SetCursor(v.Origin())

	// send the contents somewhere
	w.handler(contents)

	return nil
}
