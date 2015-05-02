package DNP3

import (
	//"bytes"
	//"fmt"
	"reflect"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)

type ClassObjectVariant struct{
	Q Qualifier
}

//Read the One (Variants and Qualifier) at a time
//Each Classs Object is of #1 byte Variants and 1 byte Qualifier
func (co ClassObjectVariant) ProcessVariants(pBuffer *h.Buffer,cbv []interface{}){

	//There should be set of set the true
	variant := pBuffer.Next(1)
	switch variant[0] {
		//Checkout http://play.golang.org/p/mawvuoGUfP for more explaination
		case 2: reflect.ValueOf(cbv[0]).Elem().Set(reflect.ValueOf(true))
		case 3: reflect.ValueOf(cbv[1]).Elem().Set(reflect.ValueOf(true))
		case 4: reflect.ValueOf(cbv[2]).Elem().Set(reflect.ValueOf(true))
	}
	
	//Print the decoded Object Variation
	co.PrintVariant(variant[0])
	
	//Process the qualifier
	co.Q.ProcessQualifier(pBuffer)
}

//Print the object variant
func (co ClassObjectVariant) PrintVariant(variant byte ){
	var whichClass int
	switch variant {
		case 2: whichClass = 1
		case 3: whichClass = 2
		case 4: whichClass = 3		
	}	
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x Object Variant: Class %d",variant, whichClass)
}