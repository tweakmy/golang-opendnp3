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
	"bytes"
	"encoding/binary"
	//crc "github.com/howeyc/crc16"
	apl "opendnp3/APL"
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
func (d Datalink) Decode (pbuffer *bytes.Buffer, pRecvMesg *RecvMessage )(err error) {
	fmt.Println("dummy")
	//??Todo: to try unit test on less message
	//??Todo: potential of program crashing because it might be less than 10 bytes
	buffer := make([]byte,10)
	
	//??Todo: implement error catch on the buffer.read function
	/*Read reads the next len(p) bytes from the buffer or until the buffer is drained. 
	The return value n is the number of bytes read. 
	If the buffer has no data to return, err is io.EOF (unless len(p) is zero); otherwise it is nil.
	*/ 
	pbuffer.Read(buffer) //Read it to the buffer
	
	//??Todo: to use bufferlen
	//bufferLen := pbuffer.Len
	
	//initialize the error to nothing
	err = nil
	
	//Decode start byte
	if buffer[0] != 0x5 || buffer[1] != 0x64 {
		apl.Logger.Loggedf(3,apl.LEV_ERROR,"Incorrect header: %02x %02x",buffer[0],buffer[1])
		err = errors.New("Incorrect header")		
	}
	
	pRecvMesg.length = uint(buffer[2])
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"Length: %d",pRecvMesg.length)
	
	//Decode the control message
	ControlByte := buffer[3]
	
	//Decode bit 7 of the control byte the direction bit
	pRecvMesg.IsMaster = ControlByte & 0x80 >> 7 == 1
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"IsMaster: %t",pRecvMesg.IsMaster)
		//Decode bit 6
		if (ControlByte & 0x40) == 0x40 {
			pRecvMesg.fcb = (ControlByte & 0x30) >> 5 == 1
			apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"FCB: %t", pRecvMesg.fcb)	
			
			pRecvMesg.fcv = (ControlByte & 0x30) >> 4 == 1
			apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"FCV: %t", pRecvMesg.fcv)
		}else{
			pRecvMesg.dfc = (ControlByte & 0x30) >> 4 == 1
			apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"DFC: %t", pRecvMesg.dfc)
		}
	
	//Bit Masking Primary and Function Code 
	pRecvMesg.functionCode = ControlByte & 0x4F
	d.PrintControl(pRecvMesg.functionCode)
	
	
	//Decode the Destination Address	
	pRecvMesg.destAddr = uint(binary.LittleEndian.Uint16(buffer[4:6]))
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x %02x DestAddr: %d",buffer[4],buffer[5],pRecvMesg.destAddr)
	
	//Decode the Source Address
	pRecvMesg.srcAddr = uint(binary.LittleEndian.Uint16(buffer[6:8]))
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x %02x SrcAddr: %d", buffer[6],buffer[7],pRecvMesg.srcAddr)
	
	//Check CRC
	crc := uint16(binary.LittleEndian.Uint16(buffer[8:10]))
	if crc != ValidateCRC(buffer[0:8]){
		apl.Logger.Loggedf(3,apl.LEV_ERROR,"Incorrect CRC: %02x, %02x ",buffer[8],buffer[9])
		pRecvMesg.crcError = true
	}
	
	return 	  
}

//Decode datalink Control byte
func (d Datalink) PrintControl (fc byte) {
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
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x FC: %s",fc,fcString)
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

