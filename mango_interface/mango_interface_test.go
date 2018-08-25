/*
* @Author: Ximidar
* @Date:   2018-08-25 10:51:03
* @Last Modified by:   Ximidar
* @Last Modified time: 2018-08-25 11:12:48
*/
package mango_interface_test

import (
	"testing"
	"fmt"
	"github.com/ximidar/mango_cli/mango_interface"
)

func Test_Get_Available_Ports(t *testing.T){
	fmt.Println("Testing Get Available Ports")	
	mgo, err := mango_interface.NewMango()

	if err != nil{
		t.Fatal(err)
	}
	ports, err := mgo.Comm_Get_Available_Ports()
	if err != nil{
		fmt.Println("Could not get available ports", err)
		t.Fatal(err)
	}

	fmt.Println(ports)
}