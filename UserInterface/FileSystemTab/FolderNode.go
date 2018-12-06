/*
* @Author: Ximidar
* @Date:   2018-12-06 09:22:16
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-12-06 09:44:08
 */

package FileSystemTab

import (
	"errors"

	"github.com/ximidar/Flotilla/Flotilla_File_Manager/Files"
)

var (
	// ErrFileNotFound will return if the folder node cannot find a specified file
	ErrFileNotFound = errors.New("File not found in folder node")

	// ErrNotADirectory will return if the move command specifies a file
	ErrNotADirectory = errors.New("File is not a Directory")
)

// FolderNode is a linked list that will serve to keep our directory history intact
type FolderNode struct {
	PreviousNode *FolderNode
	Contents     map[string]*Files.File
	Info         *Files.File
}

// NewFolderNode will create a new Folder Node
func NewFolderNode(previousNode *FolderNode, contents map[string]*Files.File, info *Files.File) *FolderNode {
	return &FolderNode{PreviousNode: previousNode, Contents: contents, Info: info}
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

// ReturnFileByName will return the file info based on the name of the file
func (fn *FolderNode) ReturnFileByName(name string) (*Files.File, error) {
	file, ok := fn.Contents[name]
	if !ok {
		return nil, ErrFileNotFound
	}

	return file, nil
}

// MoveToFolder will return a new node with the named folder
func (fn *FolderNode) MoveToFolder(name string) (*FolderNode, error) {

	file, err := fn.ReturnFileByName(name)
	if err != nil {
		return nil, err
	}

	if file.IsDir {
		newNode := NewFolderNode(fn, file.Contents, file)
		return newNode, nil
	}

	return nil, ErrNotADirectory

}
