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
	"time"
)

type ClassMask struct {
	Class1 bool
	Class2 bool
	Class3 bool
}

//Select which time
type TypeWhichTime struct {
	LocalTime bool
	UTCTime bool
}

type TypeEventMaxConfig struct{
	MaxBinaryEvent size_t
	MaxAnalogEvent size_t
	MaxCounterEvents size_t
	MaxVtoEvents size_t
}

type TypeVtoRouterConfig struct{
	//Not defined yet; I dont use it and not sure how to use it
}

type millis_t time.Duration 
type size_t uint

//To be share with Master and Slave
type BaseConfig struct{
	DisableUnsol bool 			// if true, fully disables unsolicited mode as if the slave didn't support it, a IIN Func code not implemented will be return
	UnsolMask ClassMask 		// controls what unsol classes are enabled
	AllowTimeSync bool 			// if true, the slave will request time synchronization on an interval
	WhichTime TypeWhichTime 	// when defined it will set to use local time, otherwise by default will use UTC time
	TimeSyncPeriod millis_t 	// The period of time sync interval in milliseconds
	UnsolPackDelay millis_t		// The amount of time the slave will wait before sending new unsolicited data ( <= 0 == immediate)
	UnsolRetryDelay millis_t	// How long the slave will wait before retrying an unsuccessful unsol response
	MaxFragSize size_t			// The maximum fragment size the slave will use for data it sends
}

type SlaveConfig struct{
	MaxControls int   			// The maximum number of controls the slave will attempt to process from a single APDU
	VtoWriterQueueSize size_t	// The number of objects to store in the VtoWriter queue.
	EventMaxConfig TypeEventMaxConfig	// Structure that defines the maximum number of events to buffer
	StaticBinary GrpVar			// The default group/variation to use for static binary responses
	StaticAnalog GrpVar			// The default group/variation to use for static analog responses
	StaticCounter GrpVar		// The default group/variation to use for static counter responses
	StaticSetpointStatus GrpVar	// The default group/variation to use for static setpoint status responses
	EventBinary GrpVar			// The default group/variation to use for binary event responses
	EventAnalog GrpVar			// The default group/variation to use for analog event responses
	EventCounter GrpVar			// The default group/variation to use for counter event responses
	EventVto GrpVar				// The default group/variation to use for VTO event responses
	//Observerable to be refractor
	BaseConfig
}

type AppConfig struct{
	RspTimeout millis_t 		// The response/confirm timeout in millisec
	NumRetry size_t				// Number of retries performed for applicable frames
	//FragSize size_t			// Redundant member in the BaseConfig	// The maximum size of received application layer fragments
}

type LinkConfig struct{
	UseConfirms bool			// If true, the Data link layer will send data requesting confirmation
	NumRetry	size_t			// The number of retry attempts the link will attempt after the initial try
	LocalAddr 	uint			// dnp3 address of the local device
	RemoteAddr	uint			// dnp3 address of the remote device
	Timeout millis_t			// the response timeout in milliseconds for confirmed requests
}

type VtoConfig struct{
	VtoRouterConfig []TypeVtoRouterConfig
}

