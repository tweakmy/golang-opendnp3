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

package APL

import (
	"os"
	"io"
	"fmt"
	//"fmt"
	//"time"
	"log"
)
//import "golang.org/log" //Modified version of golang/log package

/*
Using the default log package to do the logging
*/

var (
	Logger *DNP3Logger //This is the logger which will be used all over the opendnp3 package
	infoLogger *log.Logger	//This is the Logger defined in the go package log
	errorLogger *log.Logger
)

type FilterLevel uint

const(
	LEV_INFO FilterLevel =	0x01
	LEV_WARNING =		0x02
	LEV_ERROR =	0x04
	LEV_FATAL =	0x08
	LEV_INTERPRET =	0x10
	LEV_RAW =	0x20
	LEV_DEBUG =		0x40
)

type DNP3Logger struct{

}

//Initialize the logger to start logging
func (l *DNP3Logger) Init(){
	l.InitDefault(os.Stdout,os.Stderr)
}

//Define logging medium
func (l *DNP3Logger) InitDefault(iowriter io.Writer,errwriter io.Writer){
	//!!!??? Todo: implement fatal logger
	infoLogger = log.New(iowriter,"",log.Ldate|log.Ltime|log.Lshortfile)
	errorLogger = log.New(errwriter,"",log.Ldate|log.Ltime|log.Lshortfile)
}

//Logged function and filtering
//callLevel refers whether the Logger.logged function was call within the same package or not
//If call function was from the same package then use 2, else use 3
func (l *DNP3Logger) Logged(calldepth int, filterlevel FilterLevel, messageLog string){
	/*
	!!!??? Todo: Setup the filtering
	*/
	
	switch filterlevel {
		case LEV_INFO:
			infoLogger.Output(calldepth,"INFO: " + messageLog)
		case LEV_INTERPRET:
			infoLogger.Output(calldepth,"INTERPRET: " + messageLog)	
		case LEV_RAW:
			infoLogger.Output(calldepth,"RAW: " + messageLog)
		case LEV_ERROR:
			errorLogger.Output(calldepth,"ERROR: " + messageLog)
		case LEV_FATAL:
			errorLogger.Output(calldepth,"FATAL: " + messageLog)
			os.Exit(1) //Exit if it is fatal
	}
}

//Logged f function
func (l *DNP3Logger) Loggedf(calldepth int, filterlevel FilterLevel, messageLog string, v ...interface{}){
			l.Logged(calldepth, filterlevel, fmt.Sprintf( messageLog, v...))
}
