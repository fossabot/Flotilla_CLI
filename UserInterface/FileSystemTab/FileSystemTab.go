/*
* @Author: Ximidar
* @Date:   2018-12-02 13:26:45
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-04 17:26:39
 */

package FileSystemTab

import (
	"github.com/ximidar/Flotilla/Flotilla_CLI/FlotillaInterface"
	"github.com/ximidar/gocui"
)

const (
	// FileView : Name for the File View
	FileViewName = "FileView"

	// FileInfo : Name for the File Info View
	FileInfoName = "FileInfo"

	// SelectButton : Name of the button to select files
	SelectButtonName = "SelectButton"

	// PathBar : Name for the view in charge of the current path
	PathBarName = "PathBar"
)

// FileSystemTab is an object for displaying the FileSystem
type FileSystemTab struct {
	X, Y int
	Name string
	FI   *FlotillaInterface.FlotillaInterface
}

// NewFileSystemTab will construct a new Filesystem object
func NewFileSystemTab(name string, x int, y int) (*FileSystemTab, error) {

	fs := new(FileSystemTab)
	fs.Name = name
	fs.X = x
	fs.Y = y

	return fs, nil

}

// Layout will tell gocui how to layout this widget
func (fs *FileSystemTab) Layout(g *gocui.Gui) error {
	maxX, maxY := g.Size()

	_, err := g.SetView(fs.Name, maxX+1, maxY+1, maxX+2, maxY+2)

	if err != nil {
		if err != gocui.ErrUnknownView {
			return err
		}
		// Update keybindings
	}

	return nil
}
