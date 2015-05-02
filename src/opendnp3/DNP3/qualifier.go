package DNP3

import (
	//"fmt"
	//"bytes"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)
const (
	Q0_8bit_Index = 0
	Q1_16bit_Index = 1
	Q2_32bit_Index = 2
	Q3_8bit_RangeAbsolute = 3
	Q4_16bit_RangeAbsolute = 4
	Q5_32bit_RangeAbsolute = 5
	Q6_AllMode = 6 //For this the index should be 0
	Q7_8bit_NonRangeMode = 7
	Q8_16bit_NonRangeMode = 8
	Q9_32bit_NonRangeMode = 9
	QB_AsForIsize_ObjIdenMode = 11
)

type Qualifier struct{}

//This function will tell how the user data and range should be read
func (q Qualifier) ProcessQualifier(pBuffer *h.Buffer) (indexSize uint, QualifierCode uint) {
	qualifierByte := pBuffer.Next(1)
	indexSize = uint(qualifierByte[0] & 0x70 >> 4) 
	QualifierCode = uint(qualifierByte[0] & 0x0F) 
	q.PrintQualifier(qualifierByte[0], indexSize , QualifierCode)
	return
}


func (q Qualifier) PrintQualifier(qualifierByte byte, indexSize uint, QualifierCode uint){
	var qualifierCodeString string
	switch QualifierCode{
		case Q0_8bit_Index: qualifierCodeString = "8 bit Range Index"
		case Q1_16bit_Index: qualifierCodeString = "16 bit Range Index"
		case Q2_32bit_Index: qualifierCodeString = "32 bit Range Index"
		case Q3_8bit_RangeAbsolute: qualifierCodeString = "8 bit Range Absolute"
		case Q4_16bit_RangeAbsolute: qualifierCodeString = "16 bit Range Absolute"
		case Q5_32bit_RangeAbsolute: qualifierCodeString = "32 bit Range Absolute"
		case Q6_AllMode: qualifierCodeString = "All Mode"
		case Q7_8bit_NonRangeMode : qualifierCodeString = "8 bit Non Range"
		case Q8_16bit_NonRangeMode : qualifierCodeString = "16 bit Non Range"
		case Q9_32bit_NonRangeMode : qualifierCodeString = "32 bit Non Range"
		case QB_AsForIsize_ObjIdenMode : qualifierCodeString = "As For I size"
	}
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x Qindex: %d Qcode: %s",qualifierByte, indexSize,  qualifierCodeString)
}