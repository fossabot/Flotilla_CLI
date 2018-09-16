/*
* @Author: Ximidar
* @Date:   2018-08-25 10:12:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-09-15 19:54:28
 */

package mango_interface

import (
	_"errors"
	_"fmt"
	ms "github.com/ximidar/mango_structures"
	"github.com/nats-io/go-nats"
	_"os"
	"time"
	"log"
	"encoding/json"
)

const(
	// address name
	NAME = "commango."

	// reply subs
	LIST_PORTS = NAME + "list_ports"
	INIT_COMM = NAME + "init_comm"
	CONNECT_COMM = NAME + "connect_comm"
	DISTCONNECT_COMM = NAME + "disconnect_comm"
	WRITE_COMM = NAME + "write_comm"

	// pubs
	READ_LINE = NAME + "read_line"
)

type Mango struct {
	NC *nats.Conn
	
}

/*
Basic structure of requesting different things
subj, payload := args[0], []byte(args[1])

msg, err := nc.Request(subj, []byte(payload), 100*time.Millisecond)
if err != nil {
	if nc.LastError() != nil {
		log.Fatalf("Error in Request: %v\n", nc.LastError())
	}
	log.Fatalf("Error in Request: %v\n", err)
}
log.Printf("Published [%s] : '%s'\n", subj, payload)
log.Printf("Received [%v] : '%s'\n", msg.Subject, string(msg.Data))
*/

func NewMango() (*Mango, error) {
	mgo := new(Mango)
	var err error
	mgo.NC, err = nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
		return nil, err
	}

	return mgo, nil
}

func (mgo *Mango) Make_Request(subject string, payload []byte) ([]byte, error){

	msg, err := mgo.NC.Request(subject, payload, 100*time.Millisecond)

	if err != nil{
		panic(err) // TODO make some sort of intelligent way to parse errors
	}

	return msg.Data, nil

}

// func (mgo *Mango) Get_Comm_Signal() (chan *dbus.Signal, error) {
// 	c := make(chan string, 10)
// 	return c, nil
// }

func (mgo *Mango) Comm_Set_Connection_Options(port string, baud int32) error {

	init_comm := new(ms.Init_Comm)
	init_comm.Port = port
	init_comm.Baud = int(baud)
	payload, _ := json.Marshal(init_comm)
	call, err := mgo.Make_Request(INIT_COMM, payload)

	if err != nil {
		return err
	}

	response := new(ms.Reply_String)
	err = json.Unmarshal(call, response)

	if err != nil{
		panic(err)
	}

	log.Printf("Success: %v\nResponse: %v\n", response.Success, response.Message)

	return nil

}

// func (mgo *Mango) Comm_Connect() error {
// 	call := mgo.Comm_Obj.Call("com.mango_core.commango.Open_Comm", 0)

// 	if call.Err != nil {
// 		return call.Err
// 	}
// 	return nil
// }

// func (mgo *Mango) Comm_Disconnect() error {
// 	call := mgo.Comm_Obj.Call("com.mango_core.commango.Close_Comm", 0)

// 	if call.Err != nil {
// 		return call.Err
// 	}
// 	return nil

// }

// //Call(method string, flags Flags, args ...interface{}) *Call
// func (mgo *Mango) Comm_Get_Available_Ports() ([]string, error) {

// 	call := mgo.Comm_Obj.Call("com.mango_core.commango.Get_Available_Ports", 0)

// 	if call.Err != nil {
// 		return nil, call.Err
// 	}

// 	if len(call.Body) > 0 {
// 		ports, ok := call.Body[0].([]string)
// 		if !ok {
// 			return nil, errors.New("Could not convert body to []string")
// 		}
// 		return ports, nil
// 	}

// 	return []string{}, nil

// }

// func (mgo *Mango) Comm_Write(command string) error {
// 	expected_bytes := len(command)
// 	call := mgo.Comm_Obj.Call("com.mango_core.commango.Write_Comm", 0, command)

// 	if call.Err != nil {
// 		return call.Err
// 	}

// 	if len(call.Body) > 0 {
// 		bytes_written, ok := call.Body[0].(int)

// 		if !ok {
// 			return errors.New("Could not cast body to int")
// 		}

// 		if bytes_written != expected_bytes {
// 			return errors.New("expected_bytes != written bytes")
// 		}

// 	} else {
// 		return errors.New("Call did not return any bytes")
// 	}

// 	return nil
// }
