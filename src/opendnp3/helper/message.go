package helper

import (
	//"bytes"
)

//This is provided to give facilities for the Slave or Master to process
type RecvMessage struct{
	Length uint       //Length of the whole message, which does not include CRC ans 
	IsMaster bool
	Fcb bool	
	Fcv bool
	Dfc bool
	FunctionCode byte  //This will be a constant var pass in
	SrcAddr uint
	DestAddr uint	
	
	Transport int
	TpFinFir byte  //??!! Todo: to depecreated in future
	IsInitalBytesACPI bool //This is to indicate the first frame will be have ACPI, in the multi frame, second frame will not have ACPI 
	TpSeq uint
	
	AppFinFir byte
	AppCon bool
	AppUns bool
	AppSeq uint
	//Breakup IIN to the specific flags
	AllStationsMsgReceived bool
	Class01Avail bool
	Class02Avail bool
	Class03Avail bool
	TimesyncRequired bool
	PointInLocal bool
	DeviceTrouble bool
	DeviceRestart bool
	FuncCodeNotImplemented bool
	RequestObjectsUnknown bool
	ParamInvalid bool
	BufferOverflow bool
	OperationAlreadyInProgress bool
	CorruptConfig bool
	Reserved0 bool
	Reserved1 bool	
	AppFuncCode byte
	
	PUserDataBuffer *Buffer //Raw user data will be store (minus the CRC, including the multi-frame or multi-fragment)
	CrcError bool //This will check if this message is to be processed or ignored  
	ChanCrcByte chan []byte  	//Check if it is refering to the correct byte
	ChanIsCorrectCRC chan bool	//Check if it is correct CRC
}