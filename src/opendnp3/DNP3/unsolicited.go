package DNP3

import (
	//"fmt"
	h "opendnp3/helper"
)

type Unsolicited struct{
	Class1Enabled bool
	Class2Enabled bool
	Class3Enabled bool
}

//Process the User Data on the Slave
func (u Unsolicited) DisableUnsolicited(pSlave *Slave, pRecvMesg *h.RecvMessage){
	
	var obj object
	
	//Expecting a Class 1, Class 2 or Class 3
	var class1, class2, class3 bool
	
	obj.ProcessObjectAndVariant(pRecvMesg.PUserDataBuffer,
		&class1,
		&class2,
		&class3)
	if class1 { pSlave.Ap.Uns.Class1Enabled = false }
	if class2 { pSlave.Ap.Uns.Class2Enabled = false }
	if class3 { pSlave.Ap.Uns.Class2Enabled = false }
}
