package kits

import (
	"encoding/json"
	"io"
	"net"
)

type Conn interface {
	Decode([]byte) error
	NeedSave() bool      //when decode failed,is need save here
	Save([]byte)         //save function,and never failed
	Do() error           //you should do reponse or bad_response insid it after did it's work
	IsClose() bool       //decide is close the session
	Response() error     //send response to client
	Bad_Response() error //if there is error here
}

var (
	Buffer_Size = 5 << 10 //outsid can change the buffer size
)

func ListenAndServe(socket_method, address string, conreg func(net.Conn) Conn, errchan chan<- error) {
	listener, err := net.Listen(socket_method, address)
	if err != nil {
		errchan <- err
		return
	}
	for {
		con, err := listener.Accept()
		if err == nil {
			tcon := conreg(con)
			buffer := make([]byte, Buffer_Size)
			var lang int
			go func() {
				for !tcon.IsClose() {
					lang, err = con.Read(buffer)
					if err == nil {
						err = tcon.Decode(buffer[:lang])
						if err == nil {
							err = tcon.Do()
						} else if tcon.NeedSave() {
							//decode failed and need save
							tcon.Save(buffer[:lang])
						}
					}
					if err != nil && err != io.EOF {
						errchan <- err
					}
				}
			}()
		} else if err != io.EOF {
			errchan <- err
		}
	}
}

/*
client zone start

	|	|	|
	|	|	|
	V	V	V
*/
type Client interface {
	Register(net.Conn)    //just register the client
	IsClose() bool        //need close connection?
	NeedWaitReturn() bool //wait server response,then do next
	GetBack() error       //get server response
}

func Dial(socket_method, address string, tcl Client, msgchan <-chan any, errchan chan<- error) {
	con, err := net.Dial(socket_method, address)
	if err != nil {
		errchan <- err
		return
	}
	tcl.Register(con)
	for !tcl.IsClose() {
		content, err := json.Marshal(<-msgchan)
		if err == nil {
			_, err = con.Write(content)
			if err == nil {
				if tcl.NeedWaitReturn() { //wait response
					err = tcl.GetBack()
				} //just do next
			}
		}
		if err != nil && err != io.EOF {
			errchan <- err
		}
	}
}
