/*
* @Author: Ximidar
* @Date:   2018-08-25 10:12:08
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-10-17 14:13:22
 */

package FlotillaInterface

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	_ "os"
	"strconv"
	"time"

	"github.com/nats-io/go-nats"
	DS "github.com/ximidar/Flotilla/DataStructures"
	CS "github.com/ximidar/Flotilla/DataStructures/CommStructures"
)

// EMPTY []byte for giving an empty payload
var EMPTY []byte

// FlotillaInterface is an interface to the Nats server
type FlotillaInterface struct {
	NC       *nats.Conn
	EmitLine chan string

	Timeout time.Duration
}

// NewFlotillaInterface will construct the FlotillaInterface
func NewFlotillaInterface() (*FlotillaInterface, error) {
	fi := new(FlotillaInterface)
	var err error
	fi.NC, err = nats.Connect(nats.DefaultURL)
	fi.Timeout = 100 * time.Millisecond

	if err != nil {
		log.Fatalf("Can't connect: %v\n", err)
		return nil, err
	}

	return fi, nil
}

// MakeRequest will construct a Nats Request and send it
func (fi *FlotillaInterface) MakeRequest(subject string, payload []byte) ([]byte, error) {

	msg, err := fi.NC.Request(subject, payload, fi.Timeout)

	if err != nil {
		panic(err) // TODO make some sort of intelligent way to parse errors
	}

	fi.Timeout = 100 * time.Millisecond

	return msg.Data, nil

}

// CommSetConnectionOptions will set the connection options to the Comm Object
func (fi *FlotillaInterface) CommSetConnectionOptions(port string, baud int32) error {

	initComm := new(CS.InitComm)
	initComm.Port = port
	initComm.Baud = int(baud)
	payload, _ := json.Marshal(initComm)
	call, err := fi.MakeRequest(CS.InitializeComm, payload)

	if err != nil {
		return err
	}

	response := new(DS.ReplyString)
	err = json.Unmarshal(call, response)

	if err != nil {
		panic(err)
	}

	//log.Printf("\nInitialize Comm\nSuccess: %v\nResponse: %v\n", response.Success, response.Message)

	return nil

}

// CommGetStatus will get the status of the Comm Object
func (fi *FlotillaInterface) CommGetStatus() (*CS.CommStatus, error) {
	call, err := fi.MakeRequest(CS.GetStatus, EMPTY)

	if err != nil {
		fmt.Println("Could not get status")
		return nil, err
	}

	return fi.DeconstructStatus(call)

}

// DeconstructStatus will figure out if the call succeeded or not
func (fi *FlotillaInterface) DeconstructStatus(call []byte) (*CS.CommStatus, error) {
	reply := new(DS.ReplyJSON)
	err := json.Unmarshal(call, reply)

	if err != nil {
		fmt.Println("Could not deconstruct json package")
		panic(err)
	}

	if reply.Success {
		status := new(CS.CommStatus)
		err = json.Unmarshal(reply.Message, status)
		if err != nil {
			fmt.Println("Could not deconstruct status")
			return nil, err
		}
		return status, nil
	}
	return nil, errors.New("Could not get comm status")

}

// CommConnect will query the Nats object to connect the Comm Object
func (fi *FlotillaInterface) CommConnect() error {
	fi.Timeout = 10 * time.Second //Ten Seconds to Connect
	call, err := fi.MakeRequest(CS.ConnectComm, EMPTY)

	if err != nil {
		fmt.Println("Could not connect")
		return err
	}

	reply := new(DS.ReplyString)
	err = json.Unmarshal(call, reply)
	if err != nil {
		return err
	}

	//log.Printf("\nInitialize Comm\nSuccess: %v\nResponse: %v\n", reply.Success, reply.Message)

	if !reply.Success {
		return fmt.Errorf("Could not connect: %v", reply.Message)
	}

	return nil
}

// CommDisconnect will query the Nats Server to Disconnect the Comm Object
func (fi *FlotillaInterface) CommDisconnect() error {
	fi.Timeout = 10 * time.Second //Ten Seconds to disconnect
	call, err := fi.MakeRequest(CS.DisconnectComm, EMPTY)

	if err != nil {
		fmt.Println("Could not disconnect")
		return err
	}

	reply := new(DS.ReplyString)
	err = json.Unmarshal(call, reply)
	if err != nil {
		return err
	}

	//log.Printf("\nInitialize Comm\nSuccess: %v\nResponse: %v\n", reply.Success, reply.Message)

	if !reply.Success {
		return fmt.Errorf("Could not disconnect: %v", reply.Message)
	}

	return nil
}

// CommGetAvailablePorts will query the Nats Server for all available ports
func (fi *FlotillaInterface) CommGetAvailablePorts() ([]string, error) {

	call, err := fi.MakeRequest(CS.ListPorts, []byte(""))

	if err != nil {
		return nil, err
	}

	// deconstruct the reply
	reply := new(DS.ReplyJSON)
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

// CommWrite will write a message to the Comm Object over Nats
func (fi *FlotillaInterface) CommWrite(command string) error {
	expectedBytes := len(command)
	call, err := fi.MakeRequest(CS.WriteComm, []byte(command))

	if err != nil {
		log.Println("Could not Write Comm")
		return err
	}

	// Check if the bytes match and if the call was successful
	reply := new(DS.ReplyString)
	err = json.Unmarshal(call, reply)
	if err != nil {
		return err
	}
	written, _ := strconv.Atoi(reply.Message)
	//log.Printf("\nWrite Comm\nSuccess: %v\nResponse: %v\n", reply.Success, written)

	if !reply.Success {
		return fmt.Errorf("Could not write comm: %v", reply.Message)
	}

	if expectedBytes != written {
		return fmt.Errorf("Expected %v != Written %v", expectedBytes, written)
	}

	return nil
}
