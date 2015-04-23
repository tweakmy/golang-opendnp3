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

