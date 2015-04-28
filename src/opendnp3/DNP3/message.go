package DNP3

import (
)

//This is provided to give facilities for the Slave or Master to process
type RecvMessage struct{
	length uint       //Length of the whole message
	IsMaster bool
	fcb bool	
	fcv bool
	dfc bool
	functionCode byte  //This will be a constant var pass in
	srcAddr uint
	destAddr uint	
	
	transport int
	tpFinFir byte  //??!! Todo: to depecreated in future
	isInitalBytesACPI bool //This is to indicate the first frame will be have ACPI, in the multi frame, second frame will not have ACPI 
	tpSeq uint
	
	appFinFir byte
	appCon bool
	appUns bool
	appSeq uint
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
	appFuncCode byte
	
	crcError bool //This will check if this message is to be processed or ignored  
}