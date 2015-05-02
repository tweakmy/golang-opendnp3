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

package DNP3

import (
	//"fmt"
	//"bytes"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)

type IntfCommCh interface {
	DoClose()
	DoOpen()
	DoAsyncWrite([]byte)
	DoAsyncRead() *h.Buffer//Return the buffer size
	//Should duck-type to network or serial
}

type Slave struct{
	Name  string   		//Name of the Slave as it might be possible to have multiple slave
	CommCh IntfCommCh	//Interface comm channel
	Config SlaveConfig	//Decoupled from the user config, as user might want to change on the fly
	Dl Datalink
	Tp Transport
	Ap Application
}

//Trying to stay away from state programming and use channel
func (s *Slave) Start() {
	//By default slave unsolicited should be turn off before comms initialization
	s.Ap.Uns.Class1Enabled = false
	s.Ap.Uns.Class2Enabled = false
	s.Ap.Uns.Class3Enabled = false
	//s.CommCh.DoClose()
	s.CommCh.DoOpen() //Attemp to open the channel

	var recvMesg h.RecvMessage  //Create decoded message container
	
	//Create a byte buffer pointer
	pRecvBuffer := s.CommCh.DoAsyncRead()
	
	//Create a channel to pass in
	recvMesg.ChanCrcByte = make(chan []byte)
	recvMesg.ChanIsCorrectCRC = make(chan bool)
	
	//Start Processing the CRC
	go ValidateCRCandLenfromDNP3(pRecvBuffer.Bytes(), &recvMesg)
	
	//Slave will decode the received Message
	s.Decode(pRecvBuffer, &recvMesg)
	
	//Process the user data after merging
	if recvMesg.PUserDataBuffer.Len() > 0 {
		s.Ap.ProcessUserData(s,&recvMesg)
	}
	s.CommCh.DoClose()
}

//Decode incoming message and provide response
func (s *Slave) Decode( pRecvBuffer *h.Buffer, pRecvMesg *h.RecvMessage) {
	
	//Decode Datalink
	if pRecvBuffer.Len() > 9 {
		err := s.Dl.Decode(pRecvBuffer,pRecvMesg)
		if err != nil{
		
		}
	}
		
	//Decode Transport layer
	if pRecvBuffer.Len() > 10 {
		err := s.Tp.Decode(pRecvBuffer,pRecvMesg)
		if err != nil{
			
		}
		
	}
	
	if pRecvBuffer.Len() > 11 {
		s.Ap.Decode(pRecvBuffer, pRecvMesg, pRecvMesg.IsMaster, pRecvMesg.IsInitalBytesACPI)
	}	
	
}

//Slave might attach more then one channel
//Like SEL RTU
//!!Future::Will create each eventbuffer for each Channel
func (s *Slave) AttachTCP(tcpserver *apl.TcpServer){
	s.CommCh = tcpserver
}

//Slave might attach more then one channel
//Like SEL RTU
func (s *Slave) AttachSerial(serial *apl.Serial){

}

//In case, no value was defined by user, then load the default value
//The userconfig will be pass in to the slave instead of writing the configuration on the fly
//This will user to preset a few value before starting the slave
func (s *Slave) Configure(userLinkConfig *LinkConfig) {

	s.Dl = Datalink{ IsMaster:false,
				LinkConfig:LinkConfig{
					UseConfirms:false,
					NumRetry:1,
					LocalAddr:1,
					RemoteAddr:100,
					Timeout:100},
				}
}
