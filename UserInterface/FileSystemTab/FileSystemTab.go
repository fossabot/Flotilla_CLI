/*
* @Author: Ximidar
* @Date:   2018-12-02 13:26:45
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-05 15:48:53
 */

package FileSystemTab

import (
	"github.com/ximidar/Flotilla/Flotilla_CLI/FlotillaInterface"
	"github.com/ximidar/Flotilla/Flotilla_File_Manager/Files"
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
	X, Y      int
	Name      string
	FI        *FlotillaInterface.FlotillaInterface
	Structure map[string]*Files.File

	// Views
	FileView FileViewInterface

	// File Manipulation
	CurrentlySelectedFile *Files.File
	CurrentDirectory      *FolderNode
}

// NewFileSystemTab will construct a new Filesystem object
func NewFileSystemTab(name string, x int, y int) (*FileSystemTab, error) {

	fs := new(FileSystemTab)
	fs.Name = name
	fs.X = x
	fs.Y = y

	// Set up the flotillainterface
	var err error
	fs.FI, err = FlotillaInterface.NewFlotillaInterface()
	if err != nil {
		return nil, err
	}
	fs.Structure, err = fs.FI.GetFileStructure()
	if err != nil {
		return nil, err
	}
	fs.initNode()

	// Set up the fileview
	fs.FileView = NewFileView(FileViewName, fs.X, fs.Y, fs.SelectFile)
	fs.FileView.AddFileNames(fs.CurrentDirectory.GetFileList()...)

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
	g.Update(fs.FileView.Layout)

	return nil
}

func (fs *FileSystemTab) SelectFile(file string) {

}

// initNode will apply the structure to the folder node and instantiate the root folder
func (fs *FileSystemTab) initNode() {
	var contents map[string]*Files.File
	if fs.Structure["root"].Contents != nil {
		contents = fs.Structure["root"].Contents
	}

	fs.CurrentDirectory = &FolderNode{PreviousNode: nil, Contents: contents, Path: "root", Info: fs.Structure["root"]}
}

// FolderNode is a linked list that will serve to keep our directory history intact
type FolderNode struct {
	PreviousNode *FolderNode
	Contents     map[string]*Files.File
	Path         string
	Info         *Files.File
}

// GetFileList will gather all the file names in Contents and return them as a list
func (fn *FolderNode) GetFileList() []string {
	var files []string

	if fn.PreviousNode != nil {
		files = append(files, "..")
	}

	if fn.Contents == nil {
		return files
	}

	for _, value := range fn.Contents {
		files = append(files, value.Name)
	}
	return files
}
