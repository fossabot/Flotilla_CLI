/*
* @Author: Ximidar
* @Date:   2018-11-30 15:43:19
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-11-30 17:08:58
 */

package CommonBlocks

import (
	"errors"
	"fmt"

	"github.com/ximidar/gocui"
)

// Tabs will organize all Tab Objects
type Tabs struct {
	x, y    int
	Name    string
	TabList []string
	Tabs    []*Tab
}

func NewTabs(x, y int, Name string) *Tabs {
	tabs := Tabs{x: x, y: y, Name: Name, Tabs: make([]*Tab, 0), TabList: make([]string, 0)}
	return &tabs
}

// Layout Will Layout tabs as needed
func (tabs *Tabs) Layout(g *gocui.Gui) error {
	// MaxX, _ := g.Size()
	// NumberOfTabs := len(tabs.Tabs)
	spacing := 30 //int(math.Abs(math.Min(30, int(float64(NumberOfTabs)/float64(MaxX)))))

	// iterate over tabs and arrange them
	x := tabs.x
	y := tabs.y
	h := 2

	// Update the layout
	for _, tab := range tabs.Tabs {
		tab.X = x
		tab.Y = y
		tab.W = spacing
		tab.H = h
		g.Update(tab.Layout)

		x = x + spacing + 2

	}
	return nil
}

// AddTab will add a new tab to the tab bar
func (tabs *Tabs) AddTab(Name string, Label string, Handler func(g *gocui.Gui, v *gocui.View) error) {
	tab := NewTab(0, 0, 0, 0, Name, Label, Handler)

	tabs.TabList = append(tabs.TabList, Name)
	tabs.Tabs = append(tabs.Tabs, tab)
}

// Tab will be a button that will bring up a specific screen
type Tab struct {
	X, Y, W, H int
	Name       string
	Label      string
	Handler    func(g *gocui.Gui, v *gocui.View) error
}

// NewTab will Create a new tab
func NewTab(x int, y int, w int, h int, Name, Label string, handler func(g *gocui.Gui, v *gocui.View) error) *Tab {
	tab := Tab{X: x, Y: y, W: w, H: h, Name: Name, Label: Label, Handler: handler}
	return &tab
}

// Layout will lay the tab in the appropriate position
func (tab *Tab) Layout(g *gocui.Gui) error {
	v, err := g.SetView(tab.Name, tab.X, tab.Y, tab.X+tab.W, tab.Y+tab.H)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
	}
	if err := g.SetKeybinding(tab.Name, gocui.KeyEnter, gocui.ModNone, tab.Handler); err != nil {
		//fmt.Println("Couldn't Bind")
		return err
	}

	if err := g.SetKeybinding(tab.Name, gocui.MouseLeft, gocui.ModNone, tab.Handler); err != nil {
		//fmt.Println("Couldn't Bind")
		return err
	}
	if err := tab.centerLabel(v); err != nil {
		//fmt.Println("Couldn't Center View")
		return err
	}

	return nil
}

func (tab *Tab) centerLabel(v *gocui.View) error {
	w, _ := v.Size()
	if len(tab.Label) > w {
		return errors.New("label is bigger than the button")
	}

	offsetSize := (w - len(tab.Label)) / 2
	spaceOffset := ""
	for i := 0; i < offsetSize; i++ {
		spaceOffset = spaceOffset + " "
	}
	fmt.Fprint(v, fmt.Sprintf("%v%v", spaceOffset, tab.Label))
	return nil
}
