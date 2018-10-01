/*
* @Author: Ximidar
* @Date:   2018-08-25 10:12:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-10-01 02:53:50
 */

package flotilla_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	ms "github.com/ximidar/Flotilla/data_structures"
	"log"
	_ "os"
	"strconv"
	"time"
)

const (
	// address name
	NAME = "commango."

	// reply subs
	LIST_PORTS       = NAME + "list_ports"
	INIT_COMM        = NAME + "init_comm"
	CONNECT_COMM     = NAME + "connect_comm"
	DISTCONNECT_COMM = NAME + "disconnect_comm"
	WRITE_COMM       = NAME + "write_comm"
	GET_STATUS       = NAME + "get_status"

	// pubs
	READ_LINE = NAME + "read_line"
	WRITE_LINE = NAME + "write_line"
	STATUS_UPDATE = NAME + "status_update"
)

// empty []byte for giving an empty payload
var EMPTY []byte

type Mango struct {
	NC        *nats.Conn
	Emit_Line chan string

	Timeout time.Duration
}

func NewMango() (*Mango, error) {
	mgo := new(Mango)
	var err error
	mgo.NC, err = nats.Connect(nats.DefaultURL)
	mgo.Timeout = 100 * time.Millisecond

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
		return nil, err
	}

	return mgo, nil
}

func (mgo *Mango) Make_Request(subject string, payload []byte) ([]byte, error) {

	msg, err := mgo.NC.Request(subject, payload, mgo.Timeout)

	if err != nil {
		panic(err) // TODO make some sort of intelligent way to parse errors
	}

	mgo.Timeout = 100 * time.Millisecond

	return msg.Data, nil

}

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

	if err != nil {
		panic(err)
	}

	//log.Printf("\nInitialize Comm\nSuccess: %v\nResponse: %v\n", response.Success, response.Message)

	return nil

}

func (mgo *Mango) Comm_Get_Status() (*ms.Comm_Status, error){
	call, err := mgo.Make_Request(GET_STATUS, EMPTY)

	if err != nil{
		fmt.Println("Could not get status")
		return nil, err
	}

	return mgo.Deconstruct_status(call)

}

func (mgo *Mango) Deconstruct_status(call []byte) (*ms.Comm_Status, error){
	reply := new(ms.Reply_JSON)
	err := json.Unmarshal(call, reply)

	if err != nil {
		fmt.Println("Could not deconstruct json package")
		panic(err)
	}

	if reply.Success {
		status := new(ms.Comm_Status)
		err = json.Unmarshal(reply.Message, status)
		if err != nil {
			fmt.Println("Could not deconstruct status")
			return nil, err
		}
		return status, nil
	} else {
		return nil, errors.New("Could not get comm status")
	}
}

func (mgo *Mango) Comm_Connect() error {
	mgo.Timeout = 10 * time.Second //Ten Seconds to Connect
	call, err := mgo.Make_Request(CONNECT_COMM, EMPTY)

	if err != nil {
		fmt.Println("Could not connect")
		return err
	}

	reply := new(ms.Reply_String)
	err = json.Unmarshal(call, reply)

	//log.Printf("\nInitialize Comm\nSuccess: %v\nResponse: %v\n", reply.Success, reply.Message)

	if !reply.Success {
		return errors.New(fmt.Sprintf("Could not connect: %v", reply.Message))
	}

	return nil
}

func (mgo *Mango) Comm_Disconnect() error {
	mgo.Timeout = 10 * time.Second //Ten Seconds to disconnect
	call, err := mgo.Make_Request(DISTCONNECT_COMM, EMPTY)

	if err != nil {
		fmt.Println("Could not disconnect")
		return err
	}

	reply := new(ms.Reply_String)
	err = json.Unmarshal(call, reply)

	//log.Printf("\nInitialize Comm\nSuccess: %v\nResponse: %v\n", reply.Success, reply.Message)

	if !reply.Success {
		return errors.New(fmt.Sprintf("Could not disconnect: %v", reply.Message))
	}

	return nil
}

func (mgo *Mango) Comm_Get_Available_Ports() ([]string, error) {

	call, err := mgo.Make_Request(LIST_PORTS, []byte(""))

	if err != nil {
		return nil, err
	}

	// deconstruct the reply
	reply := new(ms.Reply_JSON)
	var ports []string
	err = json.Unmarshal(call, reply)

	if err != nil {
		fmt.Println("Could not deconstruct json package")
		panic(err)
	}

	if reply.Success {
		err = json.Unmarshal(reply.Message, &ports)
		if err != nil {
			fmt.Println("Could not deconstruct ports")
			return nil, err
		}
		//fmt.Println(ports)
	} else {
		return nil, errors.New("Could not get ports")
	}

	return ports, nil

}

func (mgo *Mango) Comm_Write(command string) error {
	expected_bytes := len(command)
	call, err := mgo.Make_Request(WRITE_COMM, []byte(command))

	if err != nil {
		log.Println("Could not Write Comm")
		return err
	}

	// Check if the bytes match and if the call was successful
	reply := new(ms.Reply_String)
	err = json.Unmarshal(call, reply)
	written, _ := strconv.Atoi(reply.Message)
	//log.Printf("\nWrite Comm\nSuccess: %v\nResponse: %v\n", reply.Success, written)

	if !reply.Success {
		return errors.New(fmt.Sprintf("Could not write comm: %v", reply.Message))
	}

	if expected_bytes != written {
		return errors.New(fmt.Sprintf("Expected %v != Written %v", expected_bytes, written))
	}

	return nil
}
