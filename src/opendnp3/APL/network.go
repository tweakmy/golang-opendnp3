package APL

import (
	"net"
	"strconv"
	"log"
	"fmt"
	//"os"
	//"time"
)


//Common Attribute share among the TcpClient and TcpServer
type TcpBase struct{
	Port uint64  			//Define the socket number
	Address string 			//Ip address or Hostname
	//logger *Logger
	TCPConn *net.TCPConn 
	TCPAddr *net.TCPAddr	
	Apbuffer []byte 		//The fixed buffer
}

type TcpClient struct{
	Name string		//Name of the Channel
	TcpBase 		//Implement TCPBase
}

type TcpServer struct{
	Name string		//Name of the Channel
	TcpBase 		//Implement TCPBase
}

type tcpCallBackFunc func([]byte, int) int

func (t *TcpClient) DoOpen(){
	
	t.ResolveTCP() //Run TCP Resolve, if error it will quit the app
	
	conn, err := net.DialTCP("tcp", nil, t.TCPAddr)
	if err != nil {
		println("Dial failed:", err.Error())
	}else
	{
		t.TCPConn = conn
	}
}	

func (t *TcpServer) DoOpen(){
	
	t.ResolveTCP() //Run TCP Resolve, if error it will quit the app
    
    ln, err := net.ListenTCP("tcp", t.TCPAddr)
	if err != nil {
		//panic(err)
		println("Could not open port:", err.Error())
	}else{
		println("Listening to " + t.Address +":" + strconv.Itoa(int(t.Port)))
	} 
    
    //Waiting for client to connect
    conn, err := ln.AcceptTCP()
    if err != nil {
         println("Could not accept TCP:", err.Error()) 
     }else{
     	println("Connected Client at " + t.Address +":" + strconv.Itoa(int(t.Port)))
     	t.TCPConn = conn
     	//defer t.DoOpen()		
     }
}

//Common to both Server and Client


func (t *TcpBase) ResolveTCP() {
	//println(t.Port)
	addr, err := net.ResolveTCPAddr("tcp",t.Address + ":" + strconv.Itoa(int(t.Port)))
	if err != nil {
		//println("ResolveTCPAddr failed:", err.Error())
		log.Fatal("ResolveTCPAddr failed:", err.Error())
		//os.Exit(1)
	}else
	{
		t.TCPAddr  = addr
	}
}

func (t *TcpBase) DoAsyncWrite(apBuffer []byte) {
		_, err := t.TCPConn.Write(apBuffer)
		if err != nil {
			println("Write to server failed:", err.Error())
			//os.Exit(1)
		}
}	

func (t *TcpBase) DoAsyncRead(apBuffer *[]byte) {
	 	
		bufferLen , err := t.TCPConn.Read(t.Apbuffer)
		if err != nil {
			fmt.Print(apBuffer)
			println("Read data failed:", err.Error())
			//os.Exit(1)
			t.TCPConn.Close()
		}//else{
		
		*apBuffer = make([]byte, 0, bufferLen) //Change the size correctly
		*apBuffer = append(*apBuffer,t.Apbuffer[:bufferLen]...)
		//BufferSize = tempSize
		//	println("make new buffer")
		//	apBuffer := make([]byte,0, bufferSize)
		//	apBuffer = append(apBuffer,tempBuffer[:bufferSize]...)
		//}		
}

func (t *TcpBase) DoClose() {
	 	t.TCPConn.Close()
}		