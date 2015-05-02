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
	"encoding/binary"
	"fmt"
	//"bytes"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)

//Function Code in the application layer
const (
	//Transfers function
	ReqConfirm_00 = 0x00
	ReqRead_01 = 0x01
	ReqWrite_02 = 0x02
	
	//Control Related Functions
	ReqSelect_03 = 0x03
	ReqOperate_04 = 0x04
	ReqDirectOperate_05 = 0x05
	ReqDirectOperateNoAck_06 = 0x06
	
	//Freeze Related Functions
	ReqImmediateFreeze_07 = 0x07  //Copy specified objects to a freeze buffer
	ReqImmediateFreezeNoAck_08 = 0x08
	ReqFreezeAndClear_09 = 0x09  //Copy specified objects to freeze buffer, and then cleart the object
	ReqFreezeAndClearNoAck_0A = 0x0A
	ReqFreezeWithTime_0B = 0x0B
	ReqFreezeWithTimeNoAck_0C = 0x0C

	//Application Control Functions
	ReqColdRestart_0D = 0x0D
	ReqWarmRestart_0E = 0x0E
	ReqInitaliseDataToDefaults_0F = 0x0F
	ReqInitaliseApplication_10 = 0x10
	ReqStartApplication_11 = 0x11
	ReqStopApplication_12 = 0x12
	
	//Configuration Functions
	ReqSaveConf_13 = 0x13
	ReqEnabledUnsolicited_14 = 0x14
	ReqDisableUnsolicited_15 = 0x15
	ReqAssignClass_16 = 0x16   //Assign Specific data objects to a particular class
	
	//Time Sync Functions
	ReqDelayMeasurement_0x17 = 0x17
	ReqRecordCurrentTime_0x18 = 0x18  //Causes outstation to record the time of its clock on completion of receipt
	
	//File Functions, response should be object group 70 and variant 4
	ReqOpenFile_0x19 = 0x19
	ReqCloseFile_0x1A = 0x1A
	ReqDeleteFile_0x1B = 0x1B
	ReqGetFileInfo_0x1C = 0x1C
	ReqFileAuthenticate_0x1D = 0x1D
	ReqFileAbort_0x1E = 0x1E
	
	//Response Function Code
	RespConfirm_0x00 = 0x00
	RespResponse_0x81 = 0x81
	RespUnsolicited_0x82 = 0x82
)	

//IIN flag; Only a response will have IIN flag which is 2 byte
const(
	//First byte
	AllStationsMsgReceived_15 = 0x8000
	Class01Avail_14 = 0x4000
	Class02Avail_13 = 0x2000
	Class03Avail_12 = 0x1000
	TimesyncRequired_11 = 0x0800
	PointInLocal_10 = 0x0400
	DeviceTrouble_09 = 0x0200
	DeviceRestart_08 = 0x0100
	
	//Second byte
	FuncCodeNotImplemented_07 = 0x0080
	RequestObjectsUnknown_06 = 0x0040
	ParamInvalid_05 = 0x0020
	BufferOverflow_04 = 0x0010
	OperationAlreadyInProgress_03 = 0x0008
	CorruptConfig_02 = 0x0004
	Reserved0_01 = 0x0002
	Reserved1_00 = 0x0001
)
type Application struct {
	Uns Unsolicited
}

//Decode the application part of the message
func (a Application) Decode(pBuffer *h.Buffer, pRecvMesg *h.RecvMessage, isMaster bool, isInitalByteACPI bool){ 
	
	numReadByte := 1 //this is due to the tranport layer byte
		
	//If this is not the multi-frame or this is the first of the multi-frame
	if isInitalByteACPI {
		
		//Capture the inital 2 byte
		ACF_FC := pBuffer.Next(2)
		numReadByte += 2 //Proceed the 2 more bytes
		
		//Read the Application Control Field
		a.DecodeACF(ACF_FC[0], pRecvMesg)
		
		//Read in the Application Function Code if thisi is a response message
		pRecvMesg.AppFuncCode = ACF_FC[1]
		a.PrintFuncCode(pRecvMesg.AppFuncCode, isMaster) //Print application function Code	
		
		//If this is a response message then the next 2 byte is IIN
		if !isMaster{
			iin_bytes := pBuffer.Next(2)
			numReadByte += 2 //Proceed the 2 more bytes
			a.DecodeIIN(uint(binary.LittleEndian.Uint16(iin_bytes)),pRecvMesg)
		}
		
		//Initialize the Message UserData bytes with empty bytes; otherwise it will put the application on panic with nil pointer
		pRecvMesg.PUserDataBuffer = h.NewBuffer([]byte(""))
		
	} 
	
	//Join all the user data together without the CRC
	//Read the Remainder of buffer read + the number of Readbyte is less than 18
	//every 16 bytes of user data is padded with 2 byte CRC
	var userDataLength int
	for {
		bufLen := pBuffer.Len()
		//If there is data still to be read
		if numReadByte + bufLen > 0 {
			if numReadByte + bufLen < 18 {
				userDataLength = bufLen-2
			}else{
				userDataLength = 16
			}
			
			_, err := pRecvMesg.PUserDataBuffer.Write(pBuffer.Next(userDataLength))
			if err != nil{
					//!!?? Error handling
			}
			
			//Check CRC
			crc := pBuffer.Next(2)
			pRecvMesg.ChanCrcByte <- crc
			fmt.Println("stuck here")
			if !<-pRecvMesg.ChanIsCorrectCRC {
				pRecvMesg.CrcError = true
				PrintErrorCRC(crc)
			}		
		}
		
		numReadByte = 0
	} 
}

//Decode Application Control Field
func (a Application) DecodeACF(acf byte, pRecvMesg *h.RecvMessage){
	
	//Check if this is the first and final frame
	var fir, fin uint
	 
	pRecvMesg.AppFinFir = acf & 0xC0
	if pRecvMesg.AppFinFir == FirstOfMulti_01 {
		fir = 0
		fin = 1
	}else if pRecvMesg.AppFinFir == NotFirstNotLast_00 {
		fir = 0
		fin = 0	
	}else if pRecvMesg.AppFinFir == FinalFrame_10 {
		fir = 1 
		fin = 0
	}else if pRecvMesg.AppFinFir == OneFrame_11 {
		fir = 1 
		fin = 1	
	}
	
	//Check CON, UNS and SEQ number
	pRecvMesg.AppCon = acf & 0x20 >> 5 == 1  	//Application confirm
	pRecvMesg.AppUns = acf & 0x10 >> 4 == 1  	//Application unsolicited
	pRecvMesg.AppSeq = uint(acf & 0x0F)  		//Application sequence
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x Application Fir: %d, Fin: %d, Con: %t, Uns: %t, Seq: %d",
		acf,fir,fin, pRecvMesg.AppCon, pRecvMesg.AppCon, pRecvMesg.AppSeq )

}

//Decode Application Internal Indication Field
func (a Application) DecodeIIN(iin uint, pRecvMesg *h.RecvMessage){
	//From First Byte		
	pRecvMesg.AllStationsMsgReceived = (iin & AllStationsMsgReceived_15 >> 15  == 1) 
	pRecvMesg.Class01Avail 			 = (iin & Class01Avail_14 >> 14  == 1)
	pRecvMesg.Class02Avail 			 = (iin & Class02Avail_13 >> 13 == 1)
	pRecvMesg.Class03Avail 			 = (iin & Class03Avail_12 >> 12 ==1)
	pRecvMesg.TimesyncRequired 		 = (iin & TimesyncRequired_11 >> 11 == 1)
	pRecvMesg.PointInLocal 			 = (iin & PointInLocal_10 >> 10 == 1)
	pRecvMesg.DeviceTrouble 		 = (iin & DeviceTrouble_09 >> 9 == 1)
	pRecvMesg.DeviceRestart 		 = (iin & DeviceRestart_08 >> 8 == 1)
	
	//From Second Byte
 	pRecvMesg.FuncCodeNotImplemented = (iin & FuncCodeNotImplemented_07 >> 7 == 1)
	pRecvMesg.RequestObjectsUnknown  = (iin & RequestObjectsUnknown_06 >> 6 == 1)
	pRecvMesg.ParamInvalid 			 = (iin & ParamInvalid_05 >> 5 == 1)
	pRecvMesg.BufferOverflow 		 = (iin & BufferOverflow_04 >> 4 == 1)
	pRecvMesg.OperationAlreadyInProgress = (iin & OperationAlreadyInProgress_03 >> 3 == 1)
	pRecvMesg.CorruptConfig = (iin & CorruptConfig_02 >> 2 == 1)
	pRecvMesg.Reserved0 			 = (iin & Reserved0_01 >> 1 == 1)
	pRecvMesg.Reserved1				 = (iin & Reserved1_00 == 1)

}

//Process the User Data
func (a Application) ProcessUserData(pSlave *Slave, pRecvMesg *h.RecvMessage){
	switch pRecvMesg.AppFuncCode {
		//Process the Request Disable Unsolicited
		case ReqDisableUnsolicited_15:
			a.Uns.DisableUnsolicited(pSlave, pRecvMesg)
	}
}

//Print Application Function Code
func (a Application) PrintFuncCode(funcCode byte, isMaster bool){
	var funcCodeString string
	if isMaster {
		switch funcCode {
			case ReqConfirm_00: funcCodeString = "Request Confirm"
			case ReqRead_01: funcCodeString = "Request Read"
			case ReqWrite_02: funcCodeString = "Request Write"
			case ReqSelect_03: funcCodeString = "Select" 
			case ReqOperate_04: funcCodeString = "Operate"
			case ReqDirectOperate_05: funcCodeString = "Direct Operate"
			case ReqDirectOperateNoAck_06: funcCodeString = "Direct Operate NoAck"
			case ReqImmediateFreeze_07: funcCodeString = "Immediate Freeze"
			case ReqImmediateFreezeNoAck_08: funcCodeString = "Immediate Freeze NoAck"
			case ReqFreezeAndClear_09: funcCodeString = "Freeze and Clear"
			case ReqFreezeAndClearNoAck_0A: funcCodeString = "Freeze and Clear NoAck"
			case ReqFreezeWithTime_0B: funcCodeString = "Freeze with Time"
			case ReqFreezeWithTimeNoAck_0C: funcCodeString = "Freeze with Time NoAck"
			case ReqColdRestart_0D: funcCodeString = "Cold Restart"
			case ReqWarmRestart_0E: funcCodeString = "Warm Restart"
			case ReqInitaliseDataToDefaults_0F: funcCodeString = "Initalise Data to Defaults"
			case ReqInitaliseApplication_10: funcCodeString = "Initalise Application"
			case ReqStartApplication_11: funcCodeString = "Start Application"
			case ReqStopApplication_12: funcCodeString = "Stop Application"
			case ReqSaveConf_13: funcCodeString = "Save Configuration"
			case ReqEnabledUnsolicited_14: funcCodeString = "Enabled Unsolicited"
			case ReqDisableUnsolicited_15: funcCodeString = "Disable Unsolicited"
			case ReqAssignClass_16: funcCodeString = "Assign Class"
			case ReqDelayMeasurement_0x17: funcCodeString = "Delay Measurement"
			case ReqRecordCurrentTime_0x18: funcCodeString = "Record Current Time"
			case ReqOpenFile_0x19: funcCodeString = "Open File"
			case ReqCloseFile_0x1A: funcCodeString = "Close File"
			case ReqDeleteFile_0x1B: funcCodeString = "Delete File"
			case ReqGetFileInfo_0x1C: funcCodeString = "Get File Info"
			case ReqFileAuthenticate_0x1D: funcCodeString = "File Authenticate"
			case ReqFileAbort_0x1E: funcCodeString = "File Abort"
			default: funcCodeString = "Function Code Unknown"
			}
	}else{ 
		//if this is a response
		switch funcCode {
			case RespConfirm_0x00: funcCodeString = "Response Confirm"
			case RespResponse_0x81: funcCodeString = "Response"
			case RespUnsolicited_0x82: funcCodeString = "Response Unsolicited"
			default: funcCodeString = "Function Code Unknown"
		}
	}
	
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET, "%02x App Func Code: %s",funcCode, funcCodeString)
}

//Print the IIN bit
func (a Application) PrintIIN(iin uint){
		
		//Only if there is no flag is set
		if iin & 0xFFFF == 0x0000{
			apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%16b IIN: 0x0000",iin) 	
		}else{
			apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%16b IIN:",iin)
		}
		//From First Byte		
		if iin & AllStationsMsgReceived_15 >> 15  == 1{
			 apl.Logger.Logged(3,apl.LEV_INTERPRET,"AllStationsMsgReceived: 1")
		}	 
		if iin & Class01Avail_14 >> 14  == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"Class01Avail: 1")
		}	
		if iin & Class02Avail_13 >> 13 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"Class02Avail: 1")
		}	
		if iin & Class03Avail_12 >> 12 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"Class03Avail: 1")
		}	
		if iin & TimesyncRequired_11 >> 11 == 1 {
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"TimesyncRequired: 1")
		}	
	 	if iin & PointInLocal_10 >> 10 == 1{ 
	 		apl.Logger.Logged(3,apl.LEV_INTERPRET,"PointInLocal: 1")
	 	}	
		if iin & DeviceTrouble_09 >> 9 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"DeviceTrouble: 1")
		}	
		if iin & DeviceRestart_08 >> 8 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"DeviceRestart: 1")
		}
		//From Second Byte
		if iin & FuncCodeNotImplemented_07 >> 7 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"FuncCodeNotImplemented: 1")
		}	
		if iin & RequestObjectsUnknown_06 >> 6 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"RequestObjectsUnknown: 1")
		}	
		if iin & ParamInvalid_05 >> 5 == 1{
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"ParamInvalid: 1")
		}
		if iin & BufferOverflow_04 >> 4 == 1{
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"BufferOverflow: 1")
		}	
		if iin & OperationAlreadyInProgress_03 >> 3 == 1{
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"OperationAlreadyInProgress: 1")
		}	
		if iin & CorruptConfig_02 >> 2 == 1{
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"CorruptConfig: 1")
		}	
		if iin & Reserved0_01 >> 1 == 1{ 
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"Reserved0: 1")
		}	
		if iin & Reserved1_00 == 1{
			apl.Logger.Logged(3,apl.LEV_INTERPRET,"Reserved1: 1")
		}	
}
