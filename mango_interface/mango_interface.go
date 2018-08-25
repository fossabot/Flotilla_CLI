/*
* @Author: Ximidar
* @Date:   2018-08-25 10:12:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 11:14:45
*/

package mango_interface

import (
	"fmt"
	"github.com/godbus/dbus"
	"os"
	"errors"
)


type Mango struct {
	Conn *dbus.Conn

}

func NewMango() (*Mango, error){
	mgo := new(Mango)
	var err error
	mgo.Conn, err = dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		return nil, err
	}

	return mgo, nil
}

func (mgo *Mango) Get_Comm_Signal() (chan *dbus.Signal, error){
	matchstr := "type='signal',sender='com.mango_core.commango',interface='com.mango_core.commango',path='/com/mango_core/commango'"
	mgo.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchstr)
	c := make(chan *dbus.Signal, 10)
	mgo.Conn.Signal(c)

	return c, nil
}

func (mgo *Mango) Get_Comm_Service() (error){
	return nil

}

func (mgo *Mango) Comm_Set_Connection_Options() {

}

func (mgo *Mango) Comm_Connect() {

}

func (mgo *Mango) Comm_Disconnect() {

}

//Call(method string, flags Flags, args ...interface{}) *Call
func (mgo *Mango) Comm_Get_Available_Ports() ([]string, error){

	obj := mgo.Conn.Object("com.mango_core.commango", "/com/mango_core/commango")
	call := obj.Call("com.mango_core.commango.Get_Available_Ports", 0)
	fmt.Println(call.Body[0])

	if call.Err != nil{
		return nil, call.Err
	}

	if len(call.Body) > 0{
		ports, ok := call.Body[0].([]string)
		if !ok{
			return nil, errors.New("Could not convert body to []string")
		} 
		return ports, nil
	}

	return []string{""}, nil
	

}

func (mgo *Mango) Comm_Write(command string) {

}