/*
* @Author: Ximidar
* @Date:   2018-08-25 10:12:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 11:44:33
 */

package mango_interface

import (
	"errors"
	"fmt"
	"github.com/godbus/dbus"
	"os"
)

type Mango struct {
	Conn     *dbus.Conn
	Comm_Obj dbus.BusObject
}

func NewMango() (*Mango, error) {
	mgo := new(Mango)
	var err error
	mgo.Conn, err = dbus.SessionBus()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to connect to session bus:", err)
		return nil, err
	}

	mgo.Comm_Obj = mgo.Conn.Object("com.mango_core.commango", "/com/mango_core/commango")

	return mgo, nil
}

func (mgo *Mango) Get_Comm_Signal() (chan *dbus.Signal, error) {
	matchstr := "type='signal',sender='com.mango_core.commango',interface='com.mango_core.commango',path='/com/mango_core/commango'"
	mgo.Conn.BusObject().Call("org.freedesktop.DBus.AddMatch", 0, matchstr)
	c := make(chan *dbus.Signal, 10)
	mgo.Conn.Signal(c)

	return c, nil
}

func (mgo *Mango) Comm_Set_Connection_Options(port string, baud int32) error {
	call := mgo.Comm_Obj.Call("com.mango_core.commango.Init_Comm", 0, port, baud)

	if call.Err != nil {
		return call.Err
	}
	return nil

}

func (mgo *Mango) Comm_Connect() error {
	call := mgo.Comm_Obj.Call("com.mango_core.commango.Open_Comm", 0)

	if call.Err != nil {
		return call.Err
	}
	return nil
}

func (mgo *Mango) Comm_Disconnect() error {
	call := mgo.Comm_Obj.Call("com.mango_core.commango.Close_Comm", 0)

	if call.Err != nil {
		return call.Err
	}
	return nil

}

//Call(method string, flags Flags, args ...interface{}) *Call
func (mgo *Mango) Comm_Get_Available_Ports() ([]string, error) {

	call := mgo.Comm_Obj.Call("com.mango_core.commango.Get_Available_Ports", 0)
	fmt.Println(call.Body[0])

	if call.Err != nil {
		return nil, call.Err
	}

	if len(call.Body) > 0 {
		ports, ok := call.Body[0].([]string)
		if !ok {
			return nil, errors.New("Could not convert body to []string")
		}
		return ports, nil
	}

	return []string{""}, nil

}

func (mgo *Mango) Comm_Write(command string) error {
	expected_bytes := len(command)
	call := mgo.Comm_Obj.Call("com.mango_core.commango.Write_Comm", 0, command)

	if call.Err != nil {
		return call.Err
	}

	if len(call.Body) > 0 {
		bytes_written, ok := call.Body[0].(int)

		if !ok {
			return errors.New("Could not cast body to int")
		}

		if bytes_written != expected_bytes {
			return errors.New("expected_bytes != written bytes")
		}

	} else {
		return errors.New("Call did not return any bytes")
	}

	return nil
}
