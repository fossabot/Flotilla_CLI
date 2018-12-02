/*
* @Author: Ximidar
* @Date:   2018-12-02 13:26:45
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-02 15:18:31
 */

package FileSystemTab

import (
	"fmt"

	"github.com/ximidar/Flotilla/Flotilla_CLI/FlotillaInterface"
	"github.com/ximidar/gocui"
)

// FileSystemTab is an object for displaying the FileSystem
type FileSystemTab struct {
	X, Y, W, H int
	Name       string
	FI         *FlotillaInterface.FlotillaInterface
}

// NewFileSystemTab will construct a new Filesystem object
func NewFileSystemTab(name string, x int, y int, w int, h int) (*FileSystemTab, error) {

	fs := new(FileSystemTab)
	fs.Name = name
	fs.X = x
	fs.Y = y
	fs.W = w
	fs.H = h

	return fs, nil

}

// Layout will tell gocui how to layout this widget
func (fs *FileSystemTab) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	middleX := maxX / 2
	middleY := maxY / 2

	Message := "New Tab!"

	v, err := g.SetView(fs.Name, middleX, middleY, middleX+fs.W, middleY+fs.H)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}

	}
	fmt.Fprintln(v, Message)
	g.SetViewOnTop(v.Name())

	return nil
}
