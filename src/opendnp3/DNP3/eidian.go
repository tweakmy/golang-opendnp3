package DNP3

import (
	"bytes"
	"encoding/binary"
	//"fmt"
)

func Convert2BytesToInteger(buffer []byte ) uint {
	
	var toInterger uint16
	
	buf := bytes.NewReader(buffer)
	err := binary.Read(buf, binary.LittleEndian, &toInterger)
	if err != nil {
		println("binary.Read failed:", err)
	}
	//fmt.Print(toInterger)
	
	return uint(toInterger)
}

