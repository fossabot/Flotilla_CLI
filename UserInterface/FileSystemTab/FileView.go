/*
* @Author: Ximidar
* @Date:   2018-12-04 17:25:32
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-04 17:35:52
 */

package FileSystemTab

import (
	"github.com/ximidar/Flotilla/Flotilla_File_Manager/Files"
	"github.com/ximidar/gocui"
)

// FileView is the UI element that will show the file tree of the flotilla system
type FileView struct {
	Name      string
	X, Y, W   int
	Structure map[string]*Files.File
}

// NewFileView will create a new FileView Object
func NewFileView(name string, x int, y int, w int) (*FileView, error) {
	fv := new(FileView)
	fv.Name = name
	fv.X = x
	fv.Y = y
	fv.W = w

	return fv, nil
}

// Layout tells the gocui system how to lay out the ui element
func (fv *FileView) Layout(g *gocui.Gui) error {
	_, maxY := g.Size()
	_, err := g.SetView(fv.Name, fv.X, fv.Y, fv.X+fv.W, maxY-1)
	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// Update KeyBindings
	}
	return nil
}
