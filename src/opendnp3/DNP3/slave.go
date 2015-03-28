package DNP3

import (
	"fmt"
)
import apl "opendnp3/APL"

//const DNP3Header = [2]byte{0x5,0x64}

type IntfCommCh interface {
	DoClose()
	DoOpen()
	DoAsyncWrite([]byte)
	DoAsyncRead(*[]byte) //Return the buffer size
	//Should duck-type to network or serial
}

type Slave struct{
	Name  string   		//Name of the Slave as it might be possible to have multiple slave
	CommCh IntfCommCh	//Interface comm channel
	Config SlaveConfig	//Decoupled from the user config, as user might want to change on the fly 
	Dl Datalink 
}

//Trying to stay away from state programming and use channel
func (s *Slave) Start() {
	//s.CommCh.DoClose()
	s.CommCh.DoOpen() //Attemp to open the channel
	
	//Initaiting data link
	var RecvBuffer []byte
	s.CommCh.DoAsyncRead(&RecvBuffer)
	s.Decode(RecvBuffer)
}

//Decode incoming message and provide response
func (s *Slave) Decode( recvBuffer []byte) {

	bufferSize := len(recvBuffer)
	dlLayerChan := make(chan []byte) //Create Datalink Response Channel

	//Wrong frame size
	if bufferSize < 10 {	
		//Informed the received buffer is wrong size
		println("Incorrect frame size: %d",bufferSize )
		return
	}

	datalinkBuffer := make([]byte,0,10)
	datalinkBuffer = append( datalinkBuffer , recvBuffer[:10]... ) //Copy only 10 byte

	//Make case to ignore the message or not
	err := s.Dl.DecodeHeader(datalinkBuffer)
	if err != nil {
		println(err.Error())
		return
	}

//	tpLayerChan := make(chan []byte) //Create Tranport Layer Response Channel
//	appLayerChan := make(chan []byte) //Create Tranport Layer Response Channel
	
	if bufferSize >= 11 {
//		tpLayerBuffer := make([]byte,0,10)
//		tpLayerBuffer = append( tpLayerBuffer , recvBuffer[10:11]... ) //Copy only 10 byte
//		go s.DecodeTranportLayer( &datalinkBuffer , dlChan)
	}
	
	//Decode and provide response for the Datalink	
	if bufferSize >= 10 {
		//Create a extract 10 byte message	
		//go s.Dl.Decode(datalinkBuffer , dlLayerChan)
		s.Dl.Decode(datalinkBuffer , dlLayerChan)
		datalinkResp := <- dlLayerChan
		fmt.Println(datalinkResp) 	 
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