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
	//"bytes"
	apl "opendnp3/APL"
	h "opendnp3/helper"
)

const (
	FirstOfMulti_01 byte = 0x40
	NotFirstNotLast_00 byte = 0x00
	FinalFrame_10 byte = 0x80
	OneFrame_11 byte = 0xC0
)

type Transport struct{}

//Decode the tranport layer of the message
//TPDU can only carry 249 user data
func (t Transport) Decode (pbuffer *h.Buffer, pRecvMesg *h.RecvMessage) (err error) {
		
	buffer := pbuffer.Next(1)
	
	err = nil 
	var fir, fin uint
	//First frame , final frame of transport layer
	pRecvMesg.TpFinFir = buffer[0] & 0xC0 
	if pRecvMesg.TpFinFir == FirstOfMulti_01 {
		fir = 0
		fin = 1
		pRecvMesg.IsInitalBytesACPI = false
	}else if pRecvMesg.TpFinFir == NotFirstNotLast_00 {
		fir = 0
		fin = 0	
		pRecvMesg.IsInitalBytesACPI = false	
	}else if pRecvMesg.TpFinFir == FinalFrame_10 {
		fir = 1
		fin = 0
		pRecvMesg.IsInitalBytesACPI = true	//This will trigger application layer to read the first 2 byte or 4 byte as ACPI
	}else if pRecvMesg.TpFinFir == OneFrame_11 {
		fir = 1
		fin = 1
		pRecvMesg.IsInitalBytesACPI = true	//This will trigger application layer to read the first 2 byte or 4 byte as ACPI	
		//Initialize the Message UserData bytes
		pRecvMesg.PUserDataBuffer = h.NewBuffer([]byte(""))
	}
		
	//Sequence Number for transport 0-63 allowed
	pRecvMesg.TpSeq = uint(buffer[0] & 0x3F)
	apl.Logger.Loggedf(3,apl.LEV_INTERPRET,"%02x Transport Fir: %d Fin: %d Transport Seq: %d",buffer,fir,fin,pRecvMesg.TpSeq)
	return 
}
