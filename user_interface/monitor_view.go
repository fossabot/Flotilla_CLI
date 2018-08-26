/*
* @Author: Ximidar
* @Date:   2018-08-25 22:00:52
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-26 13:54:37
 */
package user_interface

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"os"
)

type Monitor_Interface interface {
	Write(g *gocui.Gui, mess string)
	Read(g *gocui.Gui) string
	Clear(g *gocui.Gui)
	Layout(g *gocui.Gui) error
}

type Monitor struct {
	name string
	x, y int
}

func New_Monitor(name string, x, y int) Monitor_Interface {
	return &Monitor{name: name, x: x, y: y}
}

func (w *Monitor) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	if v, err := g.SetView(w.name, w.x, w.y, maxX-1, maxY-4); err != nil {
		if err != gocui.ErrUnknownView {
			fmt.Println(err)
			return err
		}
		v.Title, err = os.Hostname()
		if err != nil {
			return err
		}
		v.Autoscroll = true
		v.Wrap = true
	}
	return nil
}

func (w *Monitor) Write(g *gocui.Gui, mess string) {

	g.Update(func(g *gocui.Gui) error {
		v, err := g.View(w.name)
		if err != nil {
			return err
		}
		fmt.Fprintln(v, mess)
		return err

	})

}

func (w *Monitor) Clear(g *gocui.Gui) {
	v, err := g.View(w.name)
	if err != nil {
		return
	}
	v.Clear()
	v.SetCursor(v.Origin())
}

func (w *Monitor) Read(g *gocui.Gui) string {
	v, err := g.View(w.name)
	if err != nil {
		return ""
	}
	return v.Buffer()
}
