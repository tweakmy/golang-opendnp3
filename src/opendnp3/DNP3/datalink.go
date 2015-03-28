package DNP3

import (
	//"fmt"
	"errors"
	"encoding/binary"
	crc "github.com/howeyc/crc16"
)

const MaxFrameSize = 292

//Data link Direction
const (
	MasterToSlave byte = 0x80 //Primary bit
	SlavetoMaster byte = 0x00 //Secondary bit
)

//All the Function Code Key
const (
	P0 byte = 0x40	//Reset Link
	P1 byte = 0x41	//Reset User Process
	P2 byte = 0x42	//Test Link
	P3 byte = 0x43	//User Data - Confirm Expected
	P4 byte = 0x44	//User Data - No Confirm
	P9 byte = 0x49	//Request Link Status
	S0 byte = 0x00	//Confirm - Ack
	S1 byte = 0x01	//Confirm - Nack
	S11 byte = 0x0B	//Link Status
	S14 byte = 0x0E	//Not functioning
	S15 byte = 0x0F	//Function not implemented
)
//type Datalink struct{
//	IsMaster bool				// The master/slave bit set on all messages
//	UseConfirms bool			// If true, the Data link layer will send data requesting confirmation
//	NumRetry	size_t			// The number of retry attempts the link will attempt after the initial try
//	LocalAddr 	uint			// dnp3 address of the local device
//	RemoteAddr	uint			// dnp3 address of the remote device
//	Timeout millis_t			// the response timeout in milliseconds for confirmed requests	
//}

type Datalink struct{
	IsMaster bool				// The master/slave bit set on all messages
	nextFCB bool 				// Remember next FCB bit
	LinkConfig
}

//Decode header to see if it is worthy to follow up message
func (d Datalink) DecodeHeader (DataLinkRecvBuffer []byte)(err error) {
	
	err = nil
	
	//Decode first byte
	if DataLinkRecvBuffer[0] != 0x5 {
		println("Incorrect header [0]:" + string(DataLinkRecvBuffer[0]))
		err = errors.New("Error in header")
	}
	
	//Decode second byte
	if DataLinkRecvBuffer[1] != 0x64 {
		println("Incorrect header [1]:" + string(DataLinkRecvBuffer[1]))
		err = errors.New("Error in header")
	}
	
	//Convert `little eidian` to integer 
	//recvDestAddr := Convert2BytesToInteger(DataLinkRecvBuffer[4:6]) 
	//recvSrcAddr := Convert2BytesToInteger(DataLinkRecvBuffer[6:8]) 
	//Try not use external
	recvDestAddr := uint(binary.LittleEndian.Uint16(DataLinkRecvBuffer[4:6]))
	recvSrcAddr := uint(binary.LittleEndian.Uint16(DataLinkRecvBuffer[6:8]))
//	fmt.Println(recvDestAddr, recvSrcAddr)
//	fmt.Println(d.LinkConfig.LocalAddr, d.LinkConfig.RemoteAddr)
//	fmt.Println(d.LocalAddr, d.RemoteAddr)
	
	//Decode Dest Address
	if recvDestAddr != d.LinkConfig.LocalAddr {
		println("Ignore this Destination Address:" + string(recvDestAddr))
		err = errors.New("Incorrect Address")
	}
	
	//Decode Remote Address
	if recvSrcAddr !=  d.LinkConfig.RemoteAddr {
		println("Ignore this Remote Address:" + string(recvSrcAddr))
		err = errors.New("Incorrect Address")
	}
	
	return 	  
}

//Decode datalink and provide response
func (d Datalink) Decode (DataLinkRecvBuffer []byte, RespChan chan []byte) {
	
	var RespControlByte, RespPri_fc byte
	var FcbRes, FcvRes bool
	
	//Decode the Control byte
	dir,fcb,fcv,dfc,pri_fc := d.decodeControl(DataLinkRecvBuffer[3]) 		 
	
	//Check incoming CRC error
	
	//Decide which Respondant (master or slave)
	if dir == MasterToSlave {
		RespControlByte = 0x00 //Set the Slave direction
		FcbRes, FcvRes, RespPri_fc =  d.SlaveResp(fcv,fcb,pri_fc)
	}else{
		RespControlByte = 0x80 //Set the Master direction
		FcbRes, FcvRes, RespPri_fc = d.MasterResp(dfc,pri_fc)
	}
	
	//Convert Address to Remote Address
	DestAddrinByte := make([]byte,2)
	SrcAddrinByte := make([]byte,2)
	binary.LittleEndian.PutUint16(SrcAddrinByte,uint16(d.LinkConfig.LocalAddr))
	binary.LittleEndian.PutUint16(DestAddrinByte,uint16(d.LinkConfig.RemoteAddr))
	
	//Build the response control byte
	RespControlByte = RespControlByte & RespPri_fc
	if FcbRes == true { 
		RespControlByte = RespControlByte & 0x40 
	}
	if FcvRes == true { 
		RespControlByte = RespControlByte & 0x20 
	}
	
	//Build up the Data Link Message
	resp := make([]byte,0,10)
	resp = append(resp, 
		0x05, 0x64, 						//Header
		0x10,								//Length
		RespControlByte,    				//Control
		SrcAddrinByte[0],SrcAddrinByte[1],	//Src address
		DestAddrinByte[0],DestAddrinByte[1])//Dest Address
	
	//Create CRC
	var dlCRC uint16
	crc.Update(dlCRC,DNPTABLE,resp)
	//binary.LittleEndian.
	CRCinbyte := make([]byte,2)
	binary.LittleEndian.PutUint16(CRCinbyte,dlCRC)
	
	//!!???Datalink Overflow
	//dfc = false
	
	resp = append(resp,CRCinbyte[0],CRCinbyte[1])
	println(resp)
	RespChan <- resp
}

//Prove Master Response
func (d Datalink) MasterResp(dfc bool, pri_fc byte) (RespFCV bool, RespFCB bool, RespPri_fc byte) {
	
	return false, false, 0xff
}
//Provide Slave Response
func (d Datalink) SlaveResp(fcv bool,fcb bool, pri_fc byte) (RespRes bool,RespDfc bool, RespPri_fc byte) {
	
	RespDfc = false
	
	switch{
		case pri_fc == P0:
		
			println("Confirm - Ack")
			RespPri_fc = S0
			d.setExpectToggleNxFCB(fcv,fcb)	
		
		case pri_fc == P1:
		
			println("Reset user process function has been discontinued")
			RespPri_fc = S15
		
		case pri_fc == P2:
		
			println("Confirm -Ack")
			RespPri_fc = S0 
			//If expected bit is correct set	
			RespDfc = d.checkExpectFCBandSetNxFCB(fcv,fcb)
		
		case pri_fc == P3:	
			RespPri_fc = S15 //Temporary parse a wrong response to the Master
			//if d.checkExpectFCBandSetNxFCB(fcv,fcb) {
			//	println("Confirm -Ack")
			//	RespPri_fc = S0
			//}else{
			//	//Otherwise, Slave to ignore the message
			//	println("Wrong FCB bit, ignoring message")
			//	RespPri_fc = 0xff
			//}
			//Attached with User Data 
		case pri_fc == P4:
			//There is no data link response
			println("User Data - No confirm")
			RespPri_fc = P4	
			//Attached with User Data
			
		case pri_fc == P9:
			//There is no data link response
			println("Confirm -Ack")
			RespPri_fc = S0	
			
			
		default:
			RespPri_fc = S15			
	}	
	return false, RespDfc, RespPri_fc
}

//Check FCB bit is properly toggled, else set dfc bit
func (d Datalink) checkExpectFCBandSetNxFCB ( fcv bool, fcb bool) bool {
	
	if fcb == true {
		if d.nextFCB == fcb {
			d.setExpectToggleNxFCB(fcv, fcb)
			return false	
		}
		return true
	}
	return false
}

//Expect NextFCB bit to be toggle
func (d Datalink) setExpectToggleNxFCB (fcv bool, fcb bool){
	
	//Store the the fcb bit if the fcv bit is set
	if fcv == true {		
		if d.nextFCB = true; fcb == true {
				d.nextFCB = false
		} 
	}	

}

//Decode datalink Control byte
func (d Datalink) decodeControl (ControlByte byte)(DIR byte,  
	FCB bool, FCV bool, DFC bool, PRI_FC_raw byte) {
	
	//Decode bit 7
	DIR = ControlByte & 0x80

	//Decode bit 6
	if (ControlByte & 0x40) == 0x40 {
		FCB = (ControlByte & 0x30) >> 5 == 1	
		FCV = (ControlByte & 0x30) >> 4 == 1
	}else{
		DFC = (ControlByte & 0x30) >> 4 == 1
	}
	
	//Bit Masking Primary and Function Code 
	PRI_FC_raw = ControlByte & 0x4F
	println("Received:")
	println(ControlByte)
	switch {
		case PRI_FC_raw == P0 :
			println("Reset Link")
		case PRI_FC_raw == P1 :
			println("Reset User Process")
		case PRI_FC_raw == P2 :
			println("Test Link")
		case PRI_FC_raw == P3 : 
			println("User Data - Confirm Expected")
		case PRI_FC_raw == P4 : 
			println("User Data - No confirm")
		case PRI_FC_raw == P9 : 
			println("Request Link Status")				
		case PRI_FC_raw == S0 : 
			println("Confirm - ACK")
		case PRI_FC_raw == S1 : 
			println("Confirm - NACK")
		case PRI_FC_raw == S11 : 
			println("Respond - Link Status")			
		case PRI_FC_raw == S14 : 
			println("Respond - Not Functioning")	
		case PRI_FC_raw == S15 : 
			println("Respond - Link Not Implemented")
		default:
			var FC = string(PRI_FC_raw)
			println("Respond - Unkown Function Code:" + FC)
	}
	
	return	
}