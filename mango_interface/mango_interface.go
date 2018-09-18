/*
* @Author: Ximidar
* @Date:   2018-08-25 10:12:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-09-16 16:12:54
 */

package mango_interface

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/nats-io/go-nats"
	ms "github.com/ximidar/mango_structures"
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

	// pubs
	READ_LINE = NAME + "read_line"
	WRITE_LINE = NAME + "write_line"
)

// empty []byte for giving an empty payload
var EMPTY []byte

type Mango struct {
	NC        *nats.Conn
	Emit_Line chan string
}

func NewMango() (*Mango, error) {
	mgo := new(Mango)
	var err error
	mgo.NC, err = nats.Connect(nats.DefaultURL)

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
		return nil, err
	}

	// Subscribe to Read_Line
	mgo.Emit_Line = make(chan string, 20)
	mgo.NC.Subscribe(READ_LINE, mgo.emit_readline_msg)

	return mgo, nil
}

func (mgo *Mango) Make_Request(subject string, payload []byte) ([]byte, error) {

	msg, err := mgo.NC.Request(subject, payload, 100*time.Millisecond)

	if err != nil {
		panic(err) // TODO make some sort of intelligent way to parse errors
	}

	return msg.Data, nil

}

func (mgo *Mango) emit_readline_msg(msg *nats.Msg) {
	mgo.Emit_Line <- string(msg.Data)
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

func (mgo *Mango) Comm_Connect() error {
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
