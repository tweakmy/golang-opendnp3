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
	//"binary"
	"fmt"
	"errors"
	//"bytes"
	
	"encoding/binary"
	//crc "github.com/howeyc/crc16"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)

const MaxFrameSize = 292

//Data link Direction
const (
	MasterToSlave byte = 0x80 //Primary bit
	SlavetoMaster byte = 0x00 //Secondary bit
)


//All the Function Code Key
const (
	P0_ResetLink byte = 0x40	//Reset Link
	P1_ResetUserProcess byte = 0x41	//Reset User Process
	P2_TestLink byte = 0x42	//Test Link
	P3_UserDataConfirm byte = 0x43	//User Data - Confirm Expected
	P4_UserDataNoConfirm byte = 0x44	//User Data - No Confirm
	P9_ReqLinkStatus byte = 0x49	//Request Link Status
	S0_ConfirmAck byte = 0x00	//Confirm - Ack
	S1_ConfirmNack byte = 0x01	//Confirm - Nack
	S11_LinkStatus byte = 0x0B	//Link Status
	S14_NotFunc byte = 0x0E	//Not functioning
	S15_FuncNotImpl byte = 0x0F	//Function not implemented
)

type Datalink struct {
	IsMaster bool
	LinkConfig
}

//#start #start #length #control #dest #dest #src #src #crc #crc
//Decode Data link which is 10 byte
func (d Datalink) Decode (pbuffer *h.Buffer, pRecvMesg *h.RecvMessage )(err error) {
	fmt.Println("dummy")
	//??Todo: to try unit test on less message
	//initialize the error to nothing
	err = nil
	
	//Decode start byte
	startHeader := pbuffer.Next(2)
	if startHeader[0] != 0x5 || startHeader[1] != 0x64 {
		apl.Logger.Loggedf(3,apl.LEV_ERROR,"Incorrect header: %02x %02x",startHeader[0],startHeader[1])
		err = errors.New("Incorrect header")	
	}
	//Decode the Length
	length := pbuffer.Next(1)
	pRecvMesg.Length = uint(length[0])
	d.PrintLength(pRecvMesg.Length)

	//Decode the control message
	d.DecodeControl(pbuffer, pRecvMesg)
	
	//Decode the Destination Address
	destAddr := pbuffer.Next(2)	
	pRecvMesg.DestAddr = uint(binary.LittleEndian.Uint16(destAddr))
	d.PrintAddress(destAddr,"DestAddr",pRecvMesg.DestAddr)
	
	//Decode the Source Address
	srcAddr := pbuffer.Next(2)	
	pRecvMesg.SrcAddr = uint(binary.LittleEndian.Uint16(srcAddr))
	d.PrintAddress(srcAddr,"SourceAddr",pRecvMesg.SrcAddr)
	
	//Check CRC
	crc := pbuffer.Next(2)
	pRecvMesg.ChanCrcByte <- crc
	if !<-pRecvMesg.ChanIsCorrectCRC {
		pRecvMesg.CrcError = true
		PrintErrorCRC(crc)
	}
	
	return 	  
}

//Decode Control byte which is of single byte
func (d Datalink) DecodeControl(pbuffer *h.Buffer,pRecvMesg *h.RecvMessage) {
	ControlByte, err  := pbuffer.ReadByte()
	if err !=nil{
		apl.Logger.Logged(3,apl.LEV_ERROR,"Control Error" + err.Error())
	}
	//Decode bit 7 of the control byte the direction bit
	pRecvMesg.IsMaster = ControlByte & 0x80 >> 7 == 1
	
	//Decode bit 6
	if (ControlByte & 0x40) == 0x40 {
			pRecvMesg.Fcb = (ControlByte & 0x30) >> 5 == 1
			pRecvMesg.Fcv = (ControlByte & 0x30) >> 4 == 1
	}else{
			pRecvMesg.Dfc = (ControlByte & 0x30) >> 4 == 1
	}
	
	//Decode bit 5-0 
	pRecvMesg.FunctionCode = ControlByte & 0x4F
	
	//Output the Control byte to logger
	d.PrintControl(ControlByte,pRecvMesg.IsMaster,pRecvMesg.Fcb,pRecvMesg.Fcv,pRecvMesg.Dfc,pRecvMesg.FunctionCode)	
}

//Print total length
func (d Datalink) PrintLength(length uint){
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"Length: %d",length)
}

//Print datalink Function Code byte
func (d Datalink) PrintControl (controlByte byte, isMaster bool, fcb bool, fcv bool, dfc bool, fc byte) {
	var fcString string
	switch {
		case fc == P0_ResetLink :
			fcString = "FC:Reset Link"
		case fc == P1_ResetUserProcess :
			fcString = "Reset User Process"
		case fc == P2_TestLink :
			fcString = "Test Link"
		case fc == P3_UserDataConfirm : 
			fcString = "Confirm Expected"
		case fc == P4_UserDataNoConfirm : 
			fcString = "User Data - No confirm"
		case fc == P9_ReqLinkStatus : 
			fcString = "Test Link"	
		case fc == S0_ConfirmAck : 
			fcString = "Test Link"
		case fc == S1_ConfirmNack : 
			fcString = "Confirm - NACK"
		case fc == S11_LinkStatus : 
			fcString = "Respond - Link Status"		
		case fc == S14_NotFunc : 
			fcString = "Respond - Not Functioning"
		case fc == S15_FuncNotImpl : 
			fcString = "Respond - Link Not Implemented"
		default:
			//var FC = string(PRI_FC_raw)
			fcString ="Unknown Function Code"	
	}
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x IsMaster: %t FCB: %t FCV: %t DFC: %t FCB: %s",
		controlByte,isMaster,fcb,fcv,dfc,fcString)
}

//Print the Datalink Address
func (d Datalink) PrintAddress(addressByte []byte, whichAddr string ,address uint){
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x %02x %s: %d",addressByte[0],addressByte[1],whichAddr,address)
}


//
////Check FCB bit is properly toggled, else set dfc bit
//func (d Datalink) checkExpectFCBandSetNxFCB ( fcv bool, fcb bool) bool {
//	
//	if fcb == true {
//		if d.nextFCB == fcb {
//			d.setExpectToggleNxFCB(fcv, fcb)
//			return false	
//		}
//		return true
//	}
//	return false
//}
//
////Expect NextFCB bit to be toggle
//func (d Datalink) setExpectToggleNxFCB (fcv bool, fcb bool){
//	
//	//Store the the fcb bit if the fcv bit is set
//	if fcv == true {		
//		if d.nextFCB = true; fcb == true {
//				d.nextFCB = false
//		} 
//	}	
//
//}

