package console_api

import "encoding/json"

/*
the encode and decode request and response for console api
*/

type request struct {
	ArgVal map[string]string `json:"argval"`
}
type response struct {
	Err string `json:"error"`
}

func encode(origin *request) ([]byte, error) {
	return json.Marshal(origin)
}

func decode(msg []byte, err error) (response, error) {
	var res response
	if err != nil {
		return res, err
	}
	err = json.Unmarshal(msg, &res)
	return res, err
}

type Site interface {
	Send(req *request) (response, error)
}

var Buffer_Size = 10 << 10 //exchange buffer size
// this error just transport error,not logic error,logic error should put in response
func (s *basic_client) Send(req *request) (response, error) {
	content, err := encode(req)
	buffer := make([]byte, Buffer_Size)
	var resp response
	if err == nil {
		_, err = s.con.Write(content)
		if err == nil {
			var lang int
			lang, err = s.con.Read(buffer)
			resp, err = decode(buffer[:lang], err)
		}
	}
	return resp, err
}
