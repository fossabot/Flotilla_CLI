/*
* @Author: Ximidar
* @Date:   2018-08-25 21:58:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 22:44:25
*/
package user_interface

import(
	"github.com/jroimartin/gocui"
	"fmt"
)

type Button struct{
	name    string
	x, y    int
	w       int
	label   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func (b *Button) Layout(g *gocui.Gui) error{
	v, err := g.SetView(b.name, b.x, b.y, b.x+b.w, b.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(b.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.KeyEnter, gocui.ModNone, b.handler); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.MouseLeft, gocui.ModNone, b.handler); err != nil {
			panic(err)
		}
		fmt.Fprint(v, b.label)
	}
	return nil
}

type Explode_Button struct{
	name    string
	x, y    int
	w       int
	label   string
	info_loader func() []string
}

func New_Explode_Button(name string, x,y,w int, label string, info_loader func() []string) *Explode_Button{
	return &Explode_Button{name:name,
						   x:x,
						   y:y,
						   w:w,
						   label: label,
						   info_loader:info_loader,
						}
}

func (b *Explode_Button) Layout(g *gocui.Gui) error{
	v, err := g.SetView(b.name, b.x, b.y, b.x+b.w, b.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if _, err := g.SetCurrentView(b.name); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.KeyEnter, gocui.ModNone, b.explode); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.MouseLeft, gocui.ModNone, b.explode); err != nil {
			return err
		}
		fmt.Fprint(v, b.label)
	}
	return nil
}

func (b *Explode_Button) explode(g *gocui.Gui, v *gocui.View) error{
	body := b.info_loader()
	midx, midy := g.Size()
	midx = midx / 2
	midy = midy / 2
	name := fmt.Sprintf("%s_explode", b.name)
	explode := NewExplode(name, midx, midy, body)
	g.Update(explode.Layout)
	return nil
}

type Explode struct {
	name string
	x, y int
	w, h int
	body []string
}

func NewExplode(name string, x, y int, body []string) *Explode {
	w := 0
	for _, l := range body {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(body) + 1
	w = w + 1

	return &Explode{name: name, x: x, y: y, w: w, h: h, body: body}
}

func (w *Explode) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.destroy); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.destroy); err != nil {
			return err
		}
		for _, line := range w.body{
			fmt.Fprintln(v, line)
		}
		
	}
	return nil
}

func (w *Explode) destroy(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(w.name)
	return nil
}