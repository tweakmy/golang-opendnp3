package DNP3

import (
	//"bytes"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)

//Define a list of the objects
const (
	ClassObject_3C = 60
)

type object struct{
	Cov ClassObjectVariant
}

//Process the object and pass in ...callback values
func (o object) ProcessObjectAndVariant(pBuffer *h.Buffer,cbv ...interface{}){
	var objectByte []byte 
	
	//This will continue to read all bytes until all bytes are read
	for pBuffer.Len() !=0 {
		//Read the Object Number 
		objectByte = pBuffer.Next(1)
		
		//Logged the decoded 
		//o.PrintObject(objectByte[0])
		
		//Pass in bytes Buffer holder to read more
		switch objectByte[0]{
			case ClassObject_3C: o.Cov.ProcessVariants(pBuffer,cbv)
		} 	
	}	
}

//Print the object number
func (o object) PrintObject(objectGroup byte){
	var objectNumString string
	switch objectGroup {
		case ClassObject_3C:
			objectNumString = "Class Object"	
	}
	apl.Logger.Loggedf(3,apl.LEV_INFO,"%02x Object %s",objectGroup, objectNumString)
}
