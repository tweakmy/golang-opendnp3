/*
Copyright (C) 2014 Jo Ee Liew liewjoee@yahoo.com

Licensed to the Apache Software Foundation (ASF) under one
or more contributor license agreements.  See the NOTICE file
distributed with this work for additional information
regarding copyright ownership.  The ASF licenses this file
to you under the Apache License, Version 2.0 (the
"License"); you may not use this file except in compliance
with the License.  You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing,
software distributed under the License is distributed on an
"AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY
KIND, either express or implied.  See the License for the
specific language governing permissions and limitations
under the License.
*/

package APL

import (
	"net"
	"strconv"
	"fmt"
	h "opendnp3/helper"
	//"bytes"
	//"log"
	//"time"
)
//import "github.com/golang/glog"


//Common Attribute share among the TcpClient and TcpServer
type TcpBase struct{
	Port uint64  			//Define the socket number
	Address string 			//Ip address or Hostname
	//logger *Logger
	TCPConn *net.TCPConn
	TCPAddr *net.TCPAddr
	Apbuffer []byte 		//The fixed buffer
}

type TcpClient struct{
	Name string		//Name of the Channel
	TcpBase 		//Implement TCPBase
}

type TcpServer struct{
	Name string		//Name of the Channel
	TcpBase 		//Implement TCPBase
}

type tcpCallBackFunc func([]byte, int) int

func (t *TcpClient) DoOpen(){

	t.ResolveTCP() //Run TCP Resolve, if error it will quit the app

	conn, err := net.DialTCP("tcp", nil, t.TCPAddr)
	if err != nil {
		Logger.Logged(2,LEV_ERROR,"Dial failed:" + err.Error())
		//println("Dial failed:", err.Error())
	}else
	{
		t.TCPConn = conn
	}
}

func (t *TcpServer) DoOpen(){

	t.ResolveTCP() //Run TCP Resolve, if error it will quit the app

  	ln, err := net.ListenTCP("tcp", t.TCPAddr)
	if err != nil {
		Logger.Logged(2,LEV_FATAL,"Could not open port:" + err.Error())
	}else{
		Logger.Logged(2,LEV_INFO,"Listening to " + t.Address +":" + strconv.Itoa(int(t.Port)))
	}

  	//Waiting for client to connect
  	conn, err := ln.AcceptTCP()
  	if err != nil {
		Logger.Logged(2,LEV_ERROR,"Could not accept TCP:" + err.Error())
  	}else{
		Logger.Logged(2,LEV_INFO,"Connected Client at " + t.Address +":" + strconv.Itoa(int(t.Port)))
    	t.TCPConn = conn
  	}
}

//Common to both Server and Client
func (t *TcpBase) ResolveTCP() {
	
	addr, err := net.ResolveTCPAddr("tcp",t.Address + ":" + strconv.Itoa(int(t.Port)))
	if err != nil {
		Logger.Logged(2,LEV_FATAL,"ResolveTCPAddr failed:" + err.Error())
	}else
	{
		t.TCPAddr  = addr
	}
}

func (t *TcpBase) DoAsyncWrite(apBuffer []byte) {
		_, err := t.TCPConn.Write(apBuffer)
		if err != nil {
			Logger.Logged(2,LEV_FATAL,"Write to server failed:" + err.Error())
		}
}

func (t *TcpBase) DoAsyncRead() (pBuffer *h.Buffer) {

		bufferLen , err := t.TCPConn.Read(t.Apbuffer)
		if err != nil {
			Logger.Logged(2,LEV_FATAL,"Read data failed:" + err.Error())
			t.TCPConn.Close()
		}
		
		//Logged the raw in hex
		t.PrintHexRaw(bufferLen)
		
		//Create a byte buffer
		pBuffer = h.NewBuffer(t.Apbuffer[:bufferLen])
		return 
}

func (t *TcpBase) DoClose() {
		Logger.Logged(2,LEV_INFO,"TCP Connection Close")
		t.TCPConn.Close()

}

func (t *TcpBase) PrintHexRaw(bufferLen int) {
		raw_mesg := ""
		for i := 0; i < bufferLen; i++ {
        	raw_mesg += fmt.Sprintf("%02x ", t.Apbuffer[i])
    	}
		Logger.Logged(2,LEV_RAW,raw_mesg)
}