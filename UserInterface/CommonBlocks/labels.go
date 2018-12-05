/*
* @Author: Ximidar
* @Date:   2018-12-04 15:47:49
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-04 15:52:16
 */

package CommonBlocks

import (
	"fmt"

	"github.com/ximidar/Flotilla/Flotilla_CLI/FlotillaInterface"
	"github.com/ximidar/gocui"
)

// Label is an object for displaying a temporary message on the screen
type Label struct {
	X, Y, W, H int
	Name       string
	Message    string
	FI         *FlotillaInterface.FlotillaInterface
}

// NewLabel will construct a new Filesystem object
func NewLabel(name string, message string, x int, y int, w int, h int) (*Label, error) {

	fs := new(Label)
	fs.Name = name
	fs.Message = message
	fs.X = x
	fs.Y = y
	fs.W = w
	fs.H = h

	return fs, nil

}

// Layout will tell gocui how to layout this widget
func (fs *Label) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	middleX := maxX / 2
	middleY := maxY / 2

	v, err := g.SetView(fs.Name, middleX, middleY, middleX+fs.W, middleY+fs.H)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

	}
	fmt.Fprintln(v, fs.Message)
	g.SetViewOnTop(v.Name())

	return nil
}
