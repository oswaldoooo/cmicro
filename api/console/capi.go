package console_api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"strconv"
)

/*
the api for console
*/
func init() {
	RegisterMap["origin"] = newbasic_client    //basic client
	RegisterMap["http"] = newhttp_basic_client //basic http client
}

const nullmsg = ""

type actmode int

var BAD_CLIENT = errors.New("client is bad")
var BAD_METHOD = errors.New("bad register method")
var RegisterMap = make(map[string]func(string, int) (ConsoleClient, error)) //register the instance of console client
/*
	typename support tcp/udp,http,socket,

if it's http url,don't need fill port.address is url
*/
type Address struct {
	TypeName string `json:"typename" form:"typename" yaml:"typename"`
	Add      string `json:"address" form:"address" yaml:"address"`
	Port     int    `json:"port" form:"port" yaml:"port"`
}
type ConsoleClient interface {
	RegisterService(serviceName string, instanceId string, act actmode, address *Address) error
	Ping() bool //verify console is still online
	DeregisterService(instanceId string) error
}

// register the new client use target method
func NewClient(name, address string, port int) (ConsoleClient, error) {
	if regfun, ok := RegisterMap[name]; ok {
		con, err := regfun(address, port)
		if err != nil {
			return nil, err
		}
		if con.Ping() {
			return con, nil
		}
		return nil, BAD_CLIENT
	} else {
		return nil, BAD_METHOD
	}
}

/*
basic instances for console client
*/

type basic_client struct {
	con net.Conn
}

func newbasic_client(add string, port int) (ConsoleClient, error) {
	con, err := net.Dial("tcp", add+":"+strconv.Itoa(port))
	if err == nil {
		return &basic_client{con: con}, nil
	}
	return nil, err
}

// actmode
const (
	STOP       = 567
	START      = 220
	RESTART    = 443
	RELOAD     = 789
	Deregister = 999
)

func (s *basic_client) RegisterService(serviceName string, instanceId string, act actmode, address *Address) error {
	ans := make(map[string]string)
	if len(serviceName) > 0 {
		ans["name"] = serviceName
	}
	if len(instanceId) > 0 {
		ans["id"] = instanceId
	}
	if len(address.TypeName) > 0 {
		ans["type"] = address.TypeName
	}
	if len(address.Add) > 0 {
		ans["address"] = address.Add
	}
	if address.Port > 0 {
		ans["port"] = strconv.Itoa(address.Port)
	}
	switch act {
	case STOP, START, RESTART, RELOAD:
		ans["actmode"] = strconv.Itoa(int(act))
	}
	rep, err := s.Send(&request{ArgVal: ans})
	if err == nil {
		if len(rep.Err) == 0 {
			return nil
		}
		return errors.New(rep.Err)
	}
	return err
}

func (s *basic_client) Ping() bool {
	_, err := s.con.Write([]byte(nullmsg))
	return err == nil
}

func (s *basic_client) DeregisterService(instanceId string) error {
	ans := make(map[string]string)
	if len(instanceId) > 0 {
		ans["id"] = instanceId
	}
	ans["actmode"] = strconv.Itoa(Deregister)
	rep, err := s.Send(&request{ArgVal: ans})
	if err == nil {
		if len(rep.Err) == 0 {
			return nil
		}
		return errors.New(rep.Err)
	}
	return err
}

type http_basic_client struct {
	Address string
}

// http edition client
func (s *http_basic_client) RegisterService(serviceName string, instanceId string, act actmode, address *Address) error {
	ans := make(map[string]string)
	if len(serviceName) > 0 {
		ans["name"] = serviceName
	}
	if len(instanceId) > 0 {
		ans["id"] = instanceId
	}
	if len(address.TypeName) > 0 {
		ans["type"] = address.TypeName
	}
	if len(address.Add) > 0 {
		ans["address"] = address.Add
	}
	if address.Port > 0 {
		ans["port"] = strconv.Itoa(address.Port)
	}
	switch act {
	case STOP, START, RESTART, RELOAD:
		ans["actmode"] = strconv.Itoa(int(act))
	}
	content, err := json.Marshal(&request{ArgVal: ans})
	if err == nil {
		var req *http.Request
		req, err = http.NewRequest("POST", s.Address+"register", bytes.NewReader(content))
		if err == nil {
			var r *http.Response
			r, err = http.DefaultClient.Do(req)
			if err == nil {
				var lang int
				buffer := make([]byte, Buffer_Size)
				lang, err = r.Body.Read(buffer)
				var rep response
				rep, err = decode(buffer[:lang], err)
				if err == nil {
					if len(rep.Err) > 0 {
						err = errors.New(rep.Err)
					}
				}
			}
		}
	}
	return err

}

func (s *http_basic_client) Ping() bool {
	_, err := http.DefaultClient.Get(s.Address + "ping")
	return err == nil
}

func (s *http_basic_client) DeregisterService(instanceId string) error {
	ans := make(map[string]string)
	if len(instanceId) > 0 {
		ans["id"] = instanceId
	}
	content, err := encode(&request{ans})
	if err == nil {
		var req *http.Request
		req, err = http.NewRequest("POST", s.Address+"deregister", bytes.NewReader(content))
		if err == nil {
			var r *http.Response
			r, err = http.DefaultClient.Do(req)
			if err == nil {
				var lang int
				buffer := make([]byte, Buffer_Size)
				lang, err = r.Body.Read(buffer)
				var rep response
				rep, err = decode(buffer[:lang], err)
				if err == nil {
					if len(rep.Err) > 0 {
						err = errors.New(rep.Err)
					}
				}
			}
		}
	}
	return err
}

func newhttp_basic_client(add string, port int) (ConsoleClient, error) {
	con := &http_basic_client{Address: "http://" + add + ":" + strconv.Itoa(port) + "/v1/"}
	if con.Ping() {
		return con, nil
	}
	return nil, BAD_CLIENT
}
