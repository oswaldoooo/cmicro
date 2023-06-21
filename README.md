## **CMIRCO**
![](console.png)
### **console_api**
```go
// example
cl, err := console_api.NewClient("origin", "localhost", 9000)
if err != nil {
    fmt.Println("connect to console failed")
    os.Exit(1)
}
err = client.RegisterService("serviceName", "serviceId", console_api.RELOAD, &console_api.Address{Add: "localhost", Port: 8000, TypeName: "tcp"})
if err != nil && err != io.EOF {
    fmt.Println(err.Error())
    os.Exit(1)
}


//deregister service from console
client.DeregisterService("serviceId")
```
### **enhance socket**
**server example**
```go
/* 
you need implement kits.Conn,and a function for register your kits.Conn
*/
import (
    "github.com/oswaldoooo/cmirco/kits"
    "json"
    "net"
)
type request struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}
type response struct {
	Ans string `json:"ans"`
	Err string `json:"error"`
}
type mirco_con struct {
	Request request
	Resp    response
	close   bool
	net.Conn
}

func register_conn(con net.Conn) kits.Conn {//register your kits.Conn
	return &mirco_con{Conn: con}
}


func (s *mirco_con) Decode(data []byte) error {//decode from bytes to your struct and then store it in connection
	return json.Unmarshal(data, &s.Request)
}

func (s *mirco_con) Do() error {//main method after decode,and you should return client response if you want
	var err error
	if len(s.Request.FirstName) == 0 || len(s.Request.LastName) == 0 {
		s.Resp.Err = "name is not completed"
		err = s.Bad_Response()
	} else {
		s.Resp.Ans = "hello " + s.Request.FirstName + "." + s.Request.LastName
		err = s.Response()
	}
	// flush the tcon data
	s.Request = request{}
	s.Resp = response{}
	return err
}

func (s *mirco_con) IsClose() bool {//is close connection
	return s.close
}

func (s *mirco_con) Response() error {//return noraml response
	content, err := json.Marshal(&s.Resp)
	if err == nil {
		_, err = s.Write(content)
	}
	return err
}

func (s *mirco_con) Bad_Response() error {//return a exception response
	content, err := json.Marshal(&s.Resp)
	if err == nil {
		_, err = s.Write(content)
	}
	return err
}

func main() {
	errchan := make(chan error)
	go kits.ListenAndServe("tcp", "localhost:8000", register_conn, errchan)
	for {
		select {
		case err := <-errchan:
			common.OuptutWithPrefix(&common.Prefix{Prefix: "error", Time: true}, err.Error())
		}
	}
}

```
**client example**
```go
import (
    "github.com/oswaldoooo/cmirco/kits"
    "net"
    "json"
)
type show_client struct {
	net.Conn
	Close bool
}

func (s *show_client) Register(con net.Conn) {//register your dial
	s.Conn = con
}

func (s *show_client) IsClose() bool {//is close the connection
	return s.Close
}

func (s *show_client) NeedWaitReturn() bool {//need wait server's response
	return true
}

var repchan = make(chan *response)

func (s *show_client) GetBack() error {//get server's response
	buffer := make([]byte, 5<<10)
	lang, err := s.Read(buffer)
	if err == nil {
		var rep response
		err = json.Unmarshal(buffer[:lang], &rep)
		if err == nil {
			repchan <- &rep
		}
	}
	return err
}

func main() {
	msgpip := make(chan any)
	errchan := make(chan error)
	cli := show_client{}
	go kits.Dial("tcp", "localhost:8000", &cli, msgpip, errchan)
	var req request
	reader := bufio.NewReader(os.Stdin)
	var msarr []string
	go func() {
		var err error
		err = common.SetRelease("error.log", true, 0600)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(1)
		}
		for {
			err = <-errchan
			common.OuptutWithPrefix(&common.Default_Err_Prefix, err.Error())
		}
	}()
	for {
		msg, _ := reader.ReadString('\n')
		msg = strings.TrimSpace(msg)
		msarr = strings.Split(msg, " ")
		req.FirstName = msarr[0]
		req.LastName = msarr[1]
		msgpip <- &req
		fmt.Println(<-repchan)
	}
}
```
### **common output**
```go
import "github.com/oswaldoooo/cmirco/api/common"
var erroroutput=common.Prefix{Prefix:"error",Time:true}
var warnoutput=common.Prefix{Prefix:"warning",Time:true}
func main(){
    // common.SetDebug()//direct output to terminal
    // common.SetRelease("running.log",false,0600)//output to running.log
    common.Output("without any prefix")
    common.OutputWithPrefix(&erroroutput,"error there")
    common.OutputWithPrefix(&warnoutput,"warn there")
}
```