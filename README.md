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