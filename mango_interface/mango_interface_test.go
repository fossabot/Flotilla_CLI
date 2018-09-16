/*
* @Author: Ximidar
* @Date:   2018-08-25 10:51:03
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-09-15 19:55:00
 */
package mango_interface_test

import (
	_"fmt"
	"github.com/ximidar/mango_cli/mango_interface"
	"testing"
	_"time"
)

func Test_Get_Available_Ports(t *testing.T) {
	// fmt.Println("Testing Get Available Ports")
	// mgo, err := mango_interface.NewMango()

	// if err != nil {
	// 	t.Fatal(err)
	// }
	// ports, err := mgo.Comm_Get_Available_Ports()
	// if err != nil {
	// 	fmt.Println("Could not get available ports", err)
	// 	t.Fatal(err)
	// }

	// fmt.Println(ports)
}

func Test_Comm_set_up_and_write(t *testing.T) {
	mgo, err := mango_interface.NewMango()

	if err != nil {
		t.Fatal(err)
	}

	mgo.Comm_Set_Connection_Options("/dev/ttyACM0", 115200)
	// mgo.Comm_Connect()

	// duration := time.Duration(5 * time.Second)
	// time.Sleep(duration)

	// read_chan, err := mgo.Get_Comm_Signal()
	// stop_reading := false

	// read_func := func() {
	// 	for read := range read_chan {
	// 		if stop_reading {
	// 			break
	// 		}

	// 		fmt.Printf("%s", read.Body[0])
	// 	}
	// }

	// go read_func()

	// pause_dur := time.Duration(100 * time.Millisecond)
	// writes := []string{"Hello!", "My", "Name", "Is", "Matt"}

	// for _, write := range writes {
	// 	mgo.Comm_Write(write)
	// 	time.Sleep(pause_dur)
	// }

	// stop_reading = true

	// mgo.Comm_Disconnect()
}
