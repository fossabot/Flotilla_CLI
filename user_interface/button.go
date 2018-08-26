/*
* @Author: Ximidar
* @Date:   2018-08-25 21:58:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-26 15:16:58
 */
package user_interface

import (
	"fmt"
	"github.com/jroimartin/gocui"
	"errors"
)

type Button struct {
	name    string
	x, y    int
	w       int
	label   string
	handler func(g *gocui.Gui, v *gocui.View) error
}

func New_Button(name string, x,y,w int, label string, handler func(g *gocui.Gui, v *gocui.View) error) *Button{
	return &Button{name: name,
				   x:x,
				   y:y,
				   w:w,
				   label:label,
				   handler:handler,
	}
}

func (b *Button) Layout(g *gocui.Gui) error {
	v, err := g.SetView(b.name, b.x, b.y, b.x+b.w, b.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.KeyEnter, gocui.ModNone, b.handler); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.MouseLeft, gocui.ModNone, b.handler); err != nil {
			return err
		}
		if err := b.center_label(v); err != nil{
			return err
		}
	}
	return nil
}

func (b *Button) center_label(v *gocui.View) error{
	w, _ := v.Size()
	if len(b.label) > w{
		return errors.New("Label is bigger than the button!")
	}

	offset_size := (w - len(b.label)) / 2 
	space_offset := ""
	for i := 0; i < offset_size ; i++{
		space_offset = space_offset + " "
	}
	fmt.Fprint(v, fmt.Sprintf("%v%v", space_offset, b.label))
	return nil
}

type Explode_Button struct {
	name            string
	x, y            int
	w               int
	label           string
	get_body        func() []string
	select_callback func(selection string)
}

func New_Explode_Button(name string, x, y, w int, label string, get_body func() []string, select_callback func(selection string)) *Explode_Button {
	return &Explode_Button{name: name,
		x:               x,
		y:               y,
		w:               w,
		label:           label,
		get_body:        get_body,
		select_callback: select_callback,
	}
}

func (b *Explode_Button) Layout(g *gocui.Gui) error {
	v, err := g.SetView(b.name, b.x, b.y, b.x+b.w, b.y+2)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.KeyEnter, gocui.ModNone, b.explode); err != nil {
			return err
		}
		if err := g.SetKeybinding(b.name, gocui.MouseLeft, gocui.ModNone, b.explode); err != nil {
			return err
		}
		if err := b.center_label(v); err != nil{
			return err
		}
	}
	return nil
}

func (b *Explode_Button) center_label(v *gocui.View) error{
	w, _ := v.Size()
	if len(b.label) > w{
		return errors.New("Label is bigger than the button!")
	}

	offset_size := (w - len(b.label)) / 2 
	space_offset := ""
	for i := 0; i < offset_size ; i++{
		space_offset = space_offset + " "
	}
	fmt.Fprint(v, fmt.Sprintf("%v%v", space_offset, b.label))
	return nil
}

func (b *Explode_Button) explode(g *gocui.Gui, v *gocui.View) error {
	body := b.get_body()
	midx, midy := g.Size()
	midx = midx / 2
	midy = midy / 2
	name := fmt.Sprintf("%s_explode", b.name)
	explode := NewExplode(name, midx, midy, body, b.select_callback)
	g.Update(explode.Layout)
	return nil
}

type Explode struct {
	name            string
	x, y            int
	w, h            int
	body            []string
	select_callback func(selection string)
}

func NewExplode(name string, x, y int, body []string, select_callback func(selection string)) *Explode {
	w := 0
	for _, l := range body {
		if len(l) > w {
			w = len(l)
		}
	}
	h := len(body) + 1
	w = w + 1

	return &Explode{name: name, x: x, y: y, w: w, h: h, body: body, select_callback: select_callback}
}

func (w *Explode) Layout(g *gocui.Gui) error {
	v, err := g.SetView(w.name, w.x, w.y, w.x+w.w, w.y+w.h)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		g.InputEsc=true
		if err := g.SetKeybinding(w.name, gocui.KeyEsc, gocui.ModNone, w.destroy); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyEnter, gocui.ModNone, w.select_and_destroy); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.MouseLeft, gocui.ModNone, w.select_and_destroy); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowUp, gocui.ModNone, w.move_select_up); err != nil {
			return err
		}
		if err := g.SetKeybinding(w.name, gocui.KeyArrowDown, gocui.ModNone, w.move_select_down); err != nil {
			return err
		}
		for _, line := range w.body {
			fmt.Fprintln(v, line)
		}

		// Make it selected and highlight the first choice
		v.Highlight = true
		v.SelBgColor = gocui.ColorGreen
		v.SelFgColor = gocui.ColorBlack
		if _, err := g.SetCurrentView(v.Name()); err != nil {
			return err
		}

	}
	return nil
}

func (w *Explode) select_and_destroy(g *gocui.Gui, v *gocui.View) error {
	//send selected item
	v.SetCursor(v.Cursor())
	_, cy := v.Cursor()
	l, err := v.Line(cy)
	if err != nil {
		l = ""
		panic(err)
	}
	w.select_callback(l)
	w.destroy(g, v)
	return nil
}

func (w *Explode) destroy(g *gocui.Gui, v *gocui.View) error {
	g.DeleteView(w.name)
	g.DeleteKeybindings(w.name)
	return nil
}

func (w *Explode) move_select_up(g *gocui.Gui, v *gocui.View) error {
	_, cury := v.Cursor()
	orgx, orgy := v.Origin()

	desty := cury - 1

	if desty == orgy-1 {
		desty = orgy
	}

	v.SetCursor(orgx, desty)
	return nil

}

func (w *Explode) move_select_down(g *gocui.Gui, v *gocui.View) error {
	_, cury := v.Cursor()
	orgx, orgy := v.Origin()

	desty := cury + 1

	if desty == (orgy + w.h) {
		desty = (orgy + w.h)
	}

	v.SetCursor(orgx, desty)
	return nil

}
