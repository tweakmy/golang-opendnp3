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
	crc "github.com/howeyc/crc16"
	h "opendnp3/helper"
)

//Create the dnp polynomial table
var DNPTABLE = crc.MakeTable(0xA6BC)

func GenerateCRC(buffer []byte) uint16{
	return crc.Checksum(buffer,DNPTABLE)
}

//Calculate the CRC.first 10 byte as datalink byte, followed by next 18 byte or maybe less than 18
//Utilize go channel to synchronize the output
func ValidateCRCandLenfromDNP3(dnp3MsgByte []byte, pRecvMesg *h.RecvMessage){ 
	
	pRecvMesg.CrcError = false	

	//Setup buffer to read consequtively
	buffer := h.NewBuffer(dnp3MsgByte)
	
	//If buffer is not empty
	//Calculate CRC for first 8 byte data + 2 byte crc
	//fmt.Println(buffer.Len())
	if buffer.Len() >= 10 {
		//Try to validate the first 10 bytes
		isCorrectCRC := ValidateCRC(buffer.Next(10))
		
		//Check if we are talking about the same CRC in the message
		if binary.LittleEndian.Uint16(buffer.ReadBackward(2)) == binary.LittleEndian.Uint16(<- pRecvMesg.ChanCrcByte) {
			pRecvMesg.ChanIsCorrectCRC <- isCorrectCRC
		}
	}
	
	//Continue reading until end of the message
	//Calculate CRC for first 16 byte data + 2 byte crc
	for {
		if buffer.Len() > 0 {
			//Try to validate the first 10 bytes
			isCorrectCRC := ValidateCRC(buffer.Next(18))
			
			
			//Check if we are talking about the same CRC in the message
			if binary.LittleEndian.Uint16(buffer.ReadBackward(2)) == binary.LittleEndian.Uint16(<- pRecvMesg.ChanCrcByte) {
				pRecvMesg.ChanIsCorrectCRC <- isCorrectCRC
				fmt.Println("end up here")
			}else{
				apl.Logger.Logged(3,apl.LEV_ERROR,"incorrect CRC")
			}
			
		}
	}
}

//Accept the Compute CRC whole set of data bytes and last 2 byte as the CRC 
func ValidateCRC(buffer []byte) bool {
	crc := uint16(binary.LittleEndian.Uint16(buffer[(len(buffer)-2):len(buffer)]))
//	fmt.Println(crc)
//	fmt.Println(buffer[0:(len(buffer)-2)])
	return crc == GenerateCRC(buffer[0:(len(buffer)-2)])
}

//Print CRC
func PrintErrorCRC(crcByte []byte){
	apl.Logger.Loggedf(3,apl.LEV_ERROR,"Incorrect CRC: %02x %02x ",crcByte[0],crcByte[1])
}
